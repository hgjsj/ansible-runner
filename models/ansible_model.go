package models

import (
	"ansible-runner/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v3"
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
	return &Ansible{Cmd: cmd, Args: viper.GetStringSlice("ansible.options")}
}

// AnsibleJob represents a job that runs a playbook.
type AnsibleJob struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Playbook string `json:"playbook,omitempty"`
	Status string `json:"status,omitempty"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
	Step   string `json:"step,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime time.Time `json:"end_time,omitempty"`
}
// enum for job status
const (
	Scheduled = "scheduled"
	Preparing = "preparing"
	Prepared = "prepared"
	Running = "running"
	Completed = "completed"
	Failed = "failed"
)
// AnsibleRunner is an interface for running Ansible playbooks.
type AnsibleRunner interface {
	Run(*Ansible) (string, error)
}

// Host represents an Ansible host.
type Host struct {
	Name string            `json:"name,omitempty"`
	Vars map[string]string `json:"vars,omitempty"`
}

// Group represents an Ansible group.
type Group struct {
    Name  string          `json:"name,omitempty"`
	Hosts []Host   `json:"hosts,omitempty"`
	Vars  map[string]string `json:"vars,omitempty"`
}

// Inventory represents an Ansible inventory.
type Inventory struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Localhost *Host    `json:"localhost,omitempty"`
	Groups    []Group `json:"groups,omitempty"`
}


// ToYaml adds a host to the inventory.
func (h *Host) ToYaml() map[string]any {  
    var hostYaml = map[string]any{}  
    if h.Name == "localhost" {
        var varsYaml = map[string]any{"vars": h.Vars}
        hostYaml = map[string]any{h.Name: varsYaml}
    } else{
		if osType, ok := h.Vars["os_type"]; ok && osType == "windows" {
			h.Vars["ansible_connection"] = "winrm"
			h.Vars["ansible_winrm_transport"] = "credssp"
			h.Vars["ansible_port"] = "5986"
			h.Vars["ansible_winrm_server_cert_validation"] = "ignore"
			if tenant, ok := h.Vars["tenant"]; !ok {
				utils.Log.Warn(fmt.Sprintf("Tenant not set for Windows host %s", h.Name))
			} else {
				h.Vars["ansible_password"] = fmt.Sprintf("/%s/%s/%s", viper.GetString("legacy_db"), tenant, h.Name)
			}
		}
        hostYaml = map[string]any{h.Name: h.Vars}
    }
    
    return hostYaml
}

// ToYaml adds a group to the inventory.    
func (g *Group) ToYaml() map[string]any {    
    var groupYaml = map[string]any{}
    var hostsYaml = map[string]any{}
    for _, host := range g.Hosts {
        hostsYaml[host.Name] = host.ToYaml()[host.Name]
    }

    groupYaml["hosts"] = hostsYaml
	if len(g.Vars) > 0 {
    	groupYaml["vars"] = g.Vars
	}
    
    return groupYaml
}

// ToYaml adds a host to the inventory.
func (i *Inventory) ToYaml(inventoryFile string) error {
	var inventoryYaml = i.Localhost.ToYaml()

    for _, group := range i.Groups {
        inventoryYaml[group.Name] = group.ToYaml()
    }

    file, err := os.Create(inventoryFile)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	defer file.Close()

    encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	err = encoder.Encode(inventoryYaml)

	if err != nil {
	 	return err
	}
	return nil
}

// NewInventory creates a new Ansible inventory.
func NewAnsibleJob(initialID bool) *AnsibleJob { 
	j := new(AnsibleJob)
	if initialID{
		j.ID = primitive.NewObjectID()
	}
	
	return j
}

// NewInventory creates a new Ansible inventory.
func NewInventory(localhost bool) *Inventory { 
	i := new(Inventory)
	i.ID = primitive.NewObjectID()
	if localhost {
		h := new(Host)
		h.Name = "localhost"
		h.Vars =map[string]string{"ansible_connection": "local"}
		i.Localhost = h
	}
	return i
}
