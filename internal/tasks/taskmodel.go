package tasks

import (
	"ansible-runner/internal/models"
)

// JobTask represents a job task.
type JobTask struct {
	Job       *models.AnsibleJob `json:"job,omitempty"`
	Inventory *models.Inventory  `json:"inventory,omitempty"`
}
