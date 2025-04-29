package services

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

func peparedJob(job *models.AnsibleJob, ansible *models.Ansible)(string, error) {
	job.Step = "Pre-job"
	var jobDir = filepath.Join(viper.GetString("job.baseDir"), job.ID.Hex())
	var vaultURL = viper.GetString("vault.uri")
	vaultToken, err := utils.Decrypt(viper.GetString("vault.token"))
	
	if err != nil {
		return "Decrypt vault token failed", err
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
		job.Status = models.Failed
		job.Error = err.Error()
		return "", err
	}
	defer file.Close()
	
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	err = encoder.Encode(extraVarsData)
	if err != nil {
		job.Status = models.Failed
		job.Error = err.Error()
		return "", err
	}

	ansible.Args = append(ansible.Args, fmt.Sprintf("--extra-vars=@%s", extraVarsFile))

	var preJob = models.Playbook{Name: "pre_job_run.yml", ProjectDir: viper.GetString("playbook.baseDir")}
	job.Status = models.Preparing
	job.Step = "Pre-job"
	o, e := preJob.Run(ansible)
	if e != nil {
		job.Status = models.Failed
		job.Error = e.Error()
		return o, e
	} 
	job.Status = models.Prepared
	return o, nil
	
	
}

func doJob(job *models.AnsibleJob, ansible *models.Ansible)(string, error) {
	//reset extra vars, remove --extra-vars from ansible.Args
	for i, arg := range ansible.Args {
		if strings.HasPrefix(arg, "--extra-vars=") {
			ansible.Args = append(ansible.Args[:i], ansible.Args[i+1:]...)
			break
		}
	}
	var playbook = models.Playbook{Name: job.Playbook, ProjectDir: viper.GetString("playbook.baseDir")}
	job.Status = models.Running
	job.Step = "Job"
	job.Playbook = playbook.Name

	var signedCert = filepath.Join(viper.GetString("job.baseDir"), job.ID.Hex(), "p4s_signed_cert")
	ansible.Args = append(ansible.Args, fmt.Sprintf("--private-key=%s", signedCert))
	if len(job.OutputPlugin) > 0 {
		ansible.Env = outputPluginMap[job.OutputPlugin]
	}
	o, e := playbook.Run(ansible)
	if e != nil {
		job.Status = models.Failed
		job.Output = o
		job.Error = e.Error()
		return o, e
	} 
	job.Status = models.Completed
	job.Output = o
	return o, nil
}

func jobHandler(c chan *models.AnsibleJob, ansible *models.Ansible, updateStatusCallback func (*models.AnsibleJob) error) {	
	outloop:
	for{
		select {
		case  job := <- c:
			utils.Log.Info(fmt.Sprintf("Status %s is arrived", job.Status))
			if job.Status == models.Scheduled {
				if output, err := peparedJob(job, ansible); err == nil {
					utils.Log.Info(fmt.Sprintf("Pre-job run is ok and output %s", output))
				} 
				utils.Log.Info("pre-job is done and callback to update job status")
				updateStatusCallback(job)
				c <- job
			} else if job.Status == models.Prepared {
				if output, err := doJob(job, ansible); err == nil {
					utils.Log.Info(fmt.Sprintf("Job run is ok and output %s", output))
				} 
				utils.Log.Info("job is done and callback to update job status")
				updateStatusCallback(job)
				c <- job
			} else {
                job.EndTime = time.Now()
				utils.Log.Info("update end time of job and callback to update job status")
                updateStatusCallback(job)
				// Run Clean up job
				if err := os.RemoveAll(filepath.Join(viper.GetString("job.baseDir"), job.ID.Hex())); err != nil {
					utils.Log.Error(err.Error())
				}
				break outloop
			}
            
		}
	}
}

// RunPlaybook launches the specified playbook.
func RunPlaybook(j *models.AnsibleJob, i *models.Inventory, updateJobFunc func(*models.AnsibleJob) error) (*models.AnsibleJob, error) {
	utils.Log.Info(fmt.Sprintf("job id is %s and inventory id is %s", j.ID.Hex(), i.ID.Hex()))
	
	j.StartTime = time.Now()
	i.AddLocalhost()
	utils.Log.Info(fmt.Sprintf("Add localhost %s to inventory", i.Localhost.Name))

	var b = viper.GetString("job.baseDir")	

	// create a directory for the job as basis
	if err := os.MkdirAll(filepath.Join(b, j.ID.Hex()), os.ModePerm); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	// create inventory file for job
	var inventoryFile = filepath.Join(b, j.ID.Hex(), "hosts.yml")
	if err := i.ToYaml(inventoryFile); err != nil {
		return nil, err
	}

	var ansible = models.NewAnsible("playbook")
	ansible.Args = append(ansible.Args, fmt.Sprintf("--inventory-file=%s", inventoryFile))
	done := make(chan *models.AnsibleJob, 1)
	done <- j
	go jobHandler(done, ansible, updateJobFunc)

	return j, nil
}
