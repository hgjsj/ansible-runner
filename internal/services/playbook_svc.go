package services

import (
	"ansible-runner/internal/db"
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
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

func jobHandler(c chan *models.AnsibleJob, ansible *models.Ansible) {	
	outloop:
	for{
		select {
		case  job := <- c:
			utils.Log.Info(fmt.Sprintf("Status %s is arrived", job.Status))
			if job.Status == models.Scheduled {
				var preJob = models.Playbook{Name: "pre_job_run.yml", ProjectDir: viper.GetString("playbook.baseDir")}
				job.Status = models.Preparing
				job.Step = "Pre-job"
				db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
				o, e := preJob.Run(ansible)
				if e != nil {
					job.Status = models.Failed
					job.Error = e.Error()
				} else {
					job.Status = models.Prepared
					utils.Log.Info(fmt.Sprintf("Pre-job run is ok and output %s", o))
				}
				db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
				
				c <- job
			} else if job.Status == models.Prepared {
				var playbook = models.Playbook{Name: job.Playbook, ProjectDir: viper.GetString("playbook.baseDir")}
				job.Status = models.Running
				job.Step = "Job"
				job.Playbook = playbook.Name
				var signedCert = filepath.Join(viper.GetString("job.baseDir"), job.ID.Hex(), "p4s_signed_cert")
				ansible.Args = append(ansible.Args, fmt.Sprintf("--private-key=%s", signedCert))
				if len(job.OutputPlugin) > 0 {
				    ansible.Env = outputPluginMap[job.OutputPlugin]
				}
				db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
				o, e := playbook.Run(ansible)
				if e != nil {
					job.Status = models.Failed
					job.Error = e.Error()
				} else {
					job.Status = models.Completed
					job.Output = o
				}
				db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
				
				c <- job
			} else {
                job.EndTime = time.Now()
                db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
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
func RunPlaybook(j string, i string, o string) (*models.AnsibleJob, error) {
	utils.Log.Info(fmt.Sprintf("job id is %s and inventory id is %s", j, i))
	job, err := db.Retrieve[models.AnsibleJob](db.NewMongoClient(),db.NewMongoClient().JobCollection, j)
	if err != nil{
		return nil, err
	} 
	utils.Log.Info(fmt.Sprintf("Get job name %s from mongodb", job.ID))
    job.StartTime = time.Now()

	inv, err := db.Retrieve[models.Inventory](db.NewMongoClient(),db.NewMongoClient().InventoryCollection, i)
	
	if err != nil {
		return nil, err
	}
	
	utils.Log.Info(fmt.Sprintf("Get inventory name %s from mongodb", inv.Name))

	job.InventoryID = inv.ID
	job.OutputPlugin = o
	inv.AddLocalhost()

	utils.Log.Info(fmt.Sprintf("Add localhost %s to inventory", inv.Localhost.Name))

	var b = viper.GetString("job.baseDir")	
	if err := os.MkdirAll(filepath.Join(b, job.ID.Hex()), os.ModePerm); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	var inventoryFile = filepath.Join(b, job.ID.Hex(), "hosts.yml")
	if err := inv.ToYaml(inventoryFile); err != nil {
		return nil, err
	}

	var ansible = models.NewAnsible("playbook")
	ansible.Args = append(ansible.Args, fmt.Sprintf("--inventory-file=%s", inventoryFile))
	ansible.Args = append(ansible.Args, fmt.Sprintf("--extra-vars=job_dir=%s", filepath.Join(b, job.ID.Hex())))	
	done := make(chan *models.AnsibleJob, 1)
	done <- job
	go jobHandler(done, ansible)

	return job, nil
}
