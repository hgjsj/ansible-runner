package services

import (
	"ansible-runner/db"
	"ansible-runner/models"
	"ansible-runner/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Playbook represents a playbook configuration.
type Playbook struct {
	Name       string `json:"name"`
	ProjectDir string `json:"project_dir"`
	ExtraVars  []string
}

// Run executes the specified playbook.
func (p *Playbook) Run(a *models.Ansible) (string, error) {
	args := append(a.Args, p.Name)
	output, err := utils.RunCommand(p.ProjectDir, a.Cmd, args...)
	if err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}
	return output, nil
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
				var preJob = Playbook{Name: "pre_job_run.yml", ProjectDir: viper.GetString("playbook.baseDir")}
				job.Status = models.Preparing
				job.Step = "Pre-job"
				db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
				o, e := preJob.Run(ansible)
				if e != nil {
					job.Status = models.Failed
					job.Error = e.Error()
				} else {
					job.Status = models.Prepared
					job.Output = o
				}
				db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
				
				c <- job
			} else if job.Status == models.Prepared {
				var playbook = Playbook{Name: job.Playbook, ProjectDir: viper.GetString("playbook.baseDir")}
				job.Status = models.Running
				job.Step = "Job"
				job.Playbook = playbook.Name
				var signedCert = filepath.Join(viper.GetString("job.baseDir"), job.ID.Hex(), "p4s_signed_cert")
				ansible.Args = append(ansible.Args, fmt.Sprintf("--private-key=%s", signedCert))
				db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
				o, e := playbook.Run(ansible)
				if e != nil {
					utils.Log.Error(e.Error())
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
func RunPlaybook(p string, j string, i *models.Inventory) (*models.AnsibleJob, error) {
	var ansible = models.NewAnsible("playbook")
	job := models.NewAnsibleJob()
	job.Status = "scheduled"
	job.Playbook = p
    job.StartTime = time.Now()

	var inventory = models.NewInventory(true)
	if ! i.ID.IsZero() {
		err := db.NewMongoClient().Retrieve(db.NewMongoClient().InventoryCollection, i.ID.Hex(), i)
		
		if err != nil {
			utils.Log.Error(err.Error())
			return nil, err
		}
		utils.Log.Info(fmt.Sprintf("Get inventory name %s from mongodb", i.Name))
		inventory.Groups = i.Groups
	}
	inventory.Groups = i.Groups
	utils.Log.Info(fmt.Sprintf("Add localhost %s to inventory", inventory.Localhost.Name))

	var b = viper.GetString("job.baseDir")	
	if err := os.MkdirAll(filepath.Join(b, job.ID.Hex()), os.ModePerm); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	 var inventoryFile = filepath.Join(b, job.ID.Hex(), "hosts.yml")
	if err := inventory.ToYaml(inventoryFile); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	ansible.Args = append(ansible.Args, fmt.Sprintf("--inventory-file=%s", inventoryFile))
	ansible.Args = append(ansible.Args, fmt.Sprintf("--extra-vars=job_dir=%s", filepath.Join(b, job.ID.Hex())))	
	done := make(chan *models.AnsibleJob, 1)
	done <- job
	go jobHandler(done, ansible)

	return job, nil
}
