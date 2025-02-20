package models

import (
	"ansible-runner/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
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
	ID     string `json:"id"`
	Playbook string `json:"playbook"`
	Status string `json:"status"`
	Output string `json:"output"`
	Error  string `json:"error"`
	Step   string `json:"step"`
	Inventory Inventory `json:"inventory"`
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
    Name  string          `json:"name"`
	Hosts []Host   `json:"hosts"`
	Vars  map[string]string `json:"vars"`
}

// Inventory represents an Ansible inventory.
type Inventory struct {
	Localhost Host    `json:"localhost"`
	Groups    []Group `json:"groups"`
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
    groupYaml["vars"] = g.Vars
    
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
		return err
	}
	defer file.Close()

    encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	err = encoder.Encode(inventoryYaml)

    // data, err := yaml.Marshal(inventoryYaml)

    // if err != nil {
	// 	return err
	// } 
     
	// err = os.WriteFile(inventoryFile, data, 0644)
	if err != nil {
	 	return err
	}
	return nil
}

// NewInventory creates a new Ansible inventory.
func NewInventory() *Inventory {
    

	// host1 := Host{Name: "ss1", Vars: map[string]string{"os_type": "linux"}}
    // host2 := Host{Name: "ss2", Vars: map[string]string{"os_type": "linux"}}
    // host3 := Host{Name: "ss3", Vars: map[string]string{"os_type": "linux"}}
    // vars := map[string]string{"https_proxy": "http://proxy.example.com:8080"}
    // var group = Group{Name: "servers", Hosts: []Host{host1, host2, host3}, Vars: vars}

    // groups := []Group{group}   
    
	return &Inventory{
		Localhost: Host{
			Name: "localhost",
			Vars: map[string]string{"ansible_connection": "local"},
		},
        // Groups: groups,
	}
}
