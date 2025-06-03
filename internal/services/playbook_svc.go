package services

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/looplab/fsm"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var outputPluginMap = map[string][]string{
	models.TINYJSON: {"ANSIBLE_CALLBACK_WHITELIST=tiny_json","ANSIBLE_STDOUT_CALLBACK=tiny_json"},
	models.JSON:    {"ANSIBLE_CALLBACK_WHITELIST=json","ANSIBLE_STDOUT_CALLBACK=json"},
}

// GetPlaybooks returns a list of playbook file names with .yml extension from the specified base directory.
func GetPlaybooks(baseDir string) ([]string, error) {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	var playbooks []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yml" {
			playbooks = append(playbooks, file.Name())
		}
	}
	return playbooks, nil

}


// before_prepare is a callback function that runs before the prepare event in the FSM.
func before_prepare(c context.Context, e *fsm.Event) {
	j := e.Args[0].(*models.AnsibleJob)
	i := e.Args[1].(*models.Inventory)
	a := e.Args[2].(*models.Ansible)
	i.AddLocalhost()
	var b = viper.GetString("job.baseDir")
	var jobDir = filepath.Join(b, j.ID.Hex())
	j.StartTime = time.Now()
	j.Step = "Pre-job"
		// create a directory for the job as basis
	if err := os.MkdirAll(jobDir, os.ModePerm); err != nil {
		utils.Log.Error(err.Error())
		j.Status = models.Failed
		return
	}

	// create inventory file for job
	var inventoryFile = filepath.Join(jobDir, "hosts.yml")
	if err := i.ToYaml(inventoryFile); err != nil {
		utils.Log.Error(err.Error())
		j.Status = models.Failed
		e.FSM.Event(context.Background(), "cleanup", j)
		return
	}
	a.Args = append(a.Args, fmt.Sprintf("--inventory-file=%s", inventoryFile))

	
	var vaultURL = viper.GetString("vault.uri")
	vaultToken, err := utils.Decrypt(viper.GetString("vault.token"))
	
	if err != nil {
		utils.Log.Error(err.Error())
		j.Status = models.Failed
		e.FSM.Event(context.Background(), "cleanup", j)
		return
	}

	// create a yaml file extra-vars.yaml for ansible
	extraVarsFile := filepath.Join(jobDir, "extra-vars.yaml")
	extraVarsData := map[string]string{
		"job_dir":       jobDir,
		"vault_address": vaultURL,
		"vault_token":   vaultToken,
	}
	file, err := os.Create(extraVarsFile)
	if err != nil {
		utils.Log.Error(err.Error())
		j.Status = models.Failed
		e.FSM.Event(context.Background(), "cleanup", j)
		return
	}
	defer file.Close()
	
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	err = encoder.Encode(extraVarsData)
	if err != nil {
		utils.Log.Error(err.Error())
		j.Status = models.Failed
		e.FSM.Event(context.Background(), "cleanup", j)
		return 
	}
	a.Args = append(a.Args, fmt.Sprintf("--extra-vars=@%s", extraVarsFile))
	j.Status = models.Preparing

	var preJob = models.Playbook{Name: "pre_job_run.yml", ProjectDir: viper.GetString("playbook.baseDir")}
	if _, err := preJob.Run(a); err != nil {
		utils.Log.Error(fmt.Sprintf("Pre-job run failed: %s", err.Error()))
		j.Status = models.Failed
		e.FSM.Event(context.Background(), "cleanup", j)
		return 
	}
	
	j.Status = models.Prepared
}

func after_prepare(c context.Context, e *fsm.Event) {
	j := e.Args[0].(*models.AnsibleJob)
	a := e.Args[2].(*models.Ansible)
	
	UpdateJobStatus(j)
	if e.FSM.Current() == models.Prepared {
		e.FSM.Event(context.Background(), "run", j, a)
	}
}

func before_run(c context.Context, e *fsm.Event) {
	j := e.Args[0].(*models.AnsibleJob)
	a := e.Args[1].(*models.Ansible)

	//reset extra vars, remove --extra-vars from ansible.Args
	for i, arg := range a.Args {
		if strings.HasPrefix(arg, "--extra-vars=") {
			a.Args = append(a.Args[:i], a.Args[i+1:]...)
			break
		}
	}

		var playbook = models.Playbook{Name: j.Playbook, ProjectDir: viper.GetString("playbook.baseDir")}
	j.Status = models.Running
	j.Step = "Job"
	j.Playbook = playbook.Name

	var signedCert = filepath.Join(viper.GetString("job.baseDir"), j.ID.Hex(), "p4s_signed_cert")
	a.Args = append(a.Args, fmt.Sprintf("--private-key=%s", signedCert))
	if len(j.OutputPlugin) > 0 {
		a.Env = outputPluginMap[j.OutputPlugin]
	}
	o, err := playbook.Run(a)
	if e != nil {
		j.Status = models.Failed
		j.Output = o
		j.Error = err.Error()
		
	} else {
		j.Status = models.Completed
		j.Output = o
	}

	e.FSM.Event(context.Background(), "cleanup", j)
}

func after_run(c context.Context, e *fsm.Event) {
	j := e.Args[0].(*models.AnsibleJob)
	UpdateJobStatus(j)
	if e.FSM.Current() == models.Completed {
		e.FSM.Event(context.Background(), "cleanup", j)
	}
}


func cleanupJob(c context.Context, e *fsm.Event) {
	j := e.Args[0].(*models.AnsibleJob)

	if err := os.RemoveAll(filepath.Join(viper.GetString("job.baseDir"), j.ID.Hex())); err != nil {
		utils.Log.Error(err.Error())
	}
}
	

// InitFSM initializes the finite state machine (FSM) for AnsibleJob.
func InitFSM() *fsm.FSM {
	// Initialize the FSM for AnsibleJob
	return fsm.NewFSM(
		models.Scheduled,
		fsm.Events{
			{Name: "prepare", Src: []string{models.Scheduled}, Dst: models.Prepared},
			{Name: "run", Src: []string{models.Prepared}, Dst: models.Completed},
			{Name: "cleanup", Src: []string{models.Prepared, models.Completed, models.Failed}, Dst: models.Done},
		},
		fsm.Callbacks{
			"before_prepare": before_prepare,
			"after_prepare":  after_prepare,
			"before_run":     before_run,
			"after_run":      after_run,
			"before_cleanup": cleanupJob,
		},
	)
}

// RunPlaybook launches the specified playbook.
func RunPlaybook(j *models.AnsibleJob, i *models.Inventory, updateJobFunc func(*models.AnsibleJob) error) (*models.AnsibleJob, error) {
	utils.Log.Info(fmt.Sprintf("job id is %s and inventory id is %s", j.ID.Hex(), i.ID.Hex()))


	fsm := InitFSM()
	utils.Log.Info(fmt.Sprintf("Add localhost %s to inventory", i.Localhost.Name))

	var ansible = models.NewAnsible("playbook")

	go func ()  {
		fsm.Event(context.Background(),"prepare", j, i, ansible)
	}()

	return j, nil
}
