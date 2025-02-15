package services

import (
	"ansible-runner/models"
	"ansible-runner/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// Playbook represents a playbook configuration.
type Playbook struct {
	Name       string `json:"name"`
	ProjectDir string `json:"project_dir"`
}

// Run executes the specified playbook.
func (p *Playbook) Run(a *models.Ansible) (string, error) {
	args := append(a.Args, p.Name)
	output, err := utils.RunCommand(p.ProjectDir, a.Cmd, args...)
	if err != nil {
		return "", err
	}
	return output, nil
}

// GetPlaybooks returns a list of playbook file names with .yml extension from the specified base directory.
func GetPlaybooks(baseDir string) ([]string, error) {
	files, err := os.ReadDir(baseDir)
	if err != nil {
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

// RunPlaybook launches the specified playbook.
func RunPlaybook(p string, i models.Inventory) (*models.AnsibleJob, error) {
	var playbook = Playbook{Name: p, ProjectDir: viper.GetString("playbook.baseDir")}
	var pre_job = Playbook{Name: "pre_job_run.yml", ProjectDir: viper.GetString("playbook.baseDir")}
	var inventory = models.NewInventory()
	inventory.Groups = i.Groups
	var ansible = models.NewAnsible("playbook")
	job := new(models.AnsibleJob)
	job.ID = uuid.New().String()
	job.Status = "scheduled"
	jobList[job.ID] = job

	var b = viper.GetString("job.baseDir")	
	if err := os.MkdirAll(filepath.Join(b, job.ID), os.ModePerm); err != nil {
		return nil, err
	}

	var inventoryFile = filepath.Join(b, job.ID, "hosts.yml")
	if err := inventory.ToYaml(inventoryFile); err != nil {
		return nil, err
	}

	
	ansible.Args = append(ansible.Args, fmt.Sprintf("--inventory-file=%s", inventoryFile))
	done := make(chan *models.AnsibleJob)

	go func(job *models.AnsibleJob) {
		job.Status = "preparing"
		job.Playbook = playbook.Name
		o, e := pre_job.Run(ansible)
		if e != nil {
			job.Status = "failed"
			job.Error = e.Error()
			done <- job
		} else {
			job.Status = "prepared"
			job.Output = o
		}
		done <- job
	}(job)

	go func(job *models.AnsibleJob) {
		<-done
		if job.Status == "prepared" {
			job.Status = "running"
			job.Playbook = playbook.Name
			o, e := playbook.Run(ansible)
			if e != nil {
				job.Status = "failed"
				job.Error = e.Error()
			} else {
				job.Status = "completed"
				job.Output = o
			}
		} 
	}(job)

	return job, nil
}
