package models

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Ansible represents an Ansible configuration.
type Ansible struct {
	Cmd  string
	Args []string
}

// NewAnsible creates a new Ansible inventory.
func NewAnsible(exe string) *Ansible {
	var cmd = viper.GetString(fmt.Sprintf("ansible.%s", exe))
	if viper.GetString("ansible.baseDir") != "" {
		// set abslute path for ansible
		cmd = filepath.Join(viper.GetString("ansible.baseDir"), viper.GetString("ansible.playbook"))
	}
	return &Ansible{Cmd: cmd, Args: viper.GetStringSlice("ansible.args")}
}

// AnsibleJob represents a job that runs a playbook.
type AnsibleJob struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Output string `json:"output"`
	Error  string
}

// AnsibleRunner is an interface for running Ansible playbooks.
type AnsibleRunner interface {
	Run(*Ansible) (string, error)
}

// Host represents an Ansible host.
type Host struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// Group represents an Ansible group.
type Group struct {
	Name  string            `json:"name"`
	Hosts []Host            `json:"hosts"`
	Vars  map[string]string `json:"vars"`
}

// Inventory represents an Ansible inventory.
type Inventory struct {
	Localhost Host    `json:"localhost"`
	Groups    []Group `json:"groups"`
}

// ToYaml adds a host to the inventory.
func (i Inventory) ToYaml(inventoryFile string) error {
	data, err := yaml.Marshal(i)
	if err != nil {
		return err
	}
	err = os.WriteFile(inventoryFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// NewInventory creates a new Ansible inventory.
func NewInventory() *Inventory {
	return &Inventory{
		Localhost: Host{
			Name: "localhost",
			Vars: map[string]string{"ansible_connection": "local"},
		},
	}
}
