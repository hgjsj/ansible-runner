package db

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"testing"
)

func init() {
	utils.LoadConfig("../../configs/config")
	utils.InitLogger()
}

func TestNewClient(t *testing.T) {
	cli := newPostgreSQL()
	defer cli.CloseService()	
	t.Logf("new client succesfully successful")
	
}

func TestCreateInventory(t *testing.T) {
	cli := newPostgreSQL()
    defer cli.CloseService()
	host := models.Host{
		Name: "db",
		Vars: map[string]string{
			"os_type":      "linux",
			"ansible_host": "100.68.0.16",
		},
	}
	group := models.Group{
		Name:  "servers",
		Hosts: []models.Host{host},
		Vars: map[string]string{
			"sid":    "RL1",
			"dbtype": "oracle",
		}}
	inv := models.Inventory{
		Name:   "pcm-p4s-postinstall",
		Groups: []models.Group{group},
	}

	if err := cli.CreateInventory( &inv); err != nil {
		t.Fatalf("create inventory failed with err %s", err.Error())
	} else {
		t.Logf("create inventory %s with ID %s successful", inv.Name, inv.ID)
	}

}

func TestGetInventory(t *testing.T) {
	cli := newPostgreSQL()
	defer cli.CloseService()

	if i, err := cli.GetInventory( "1"); err != nil {
		t.Fatalf("get inventory failed with err %s", err.Error())
	} else {
		t.Logf("get inventory %s with group name %s successful", i.Name, i.Groups[0].Name)
	}
}

func TestListInventory(t *testing.T) {
	cli := newPostgreSQL()
	defer cli.CloseService()


	if is, err := cli.ListInventory(); err != nil {
		t.Fatalf("list inventory failed with err %s", err.Error())
	} else {
		for _, inv := range is {
			t.Logf("list inventory %s with group name %s successful", inv.Name, inv.Groups[0].Name)
		}
	}
}

func TestDeleteInventory(t *testing.T) {
	cli := newPostgreSQL()
	defer cli.CloseService()

	if err := cli.DeleteInventory( "2"); err != nil {
		t.Fatalf("delete inventory failed with err %s", err.Error())
	} else {
		t.Logf("delete inventory successful")
	}
}

func TestCreateAnsibleJob(t *testing.T) {
	cli := newPostgreSQL()
	defer cli.CloseService()

	j := models.AnsibleJob{
		Playbook: "test_new1.yml",
		Status:   models.Scheduled,
		InventoryID: "2",
		OutputPlugin: models.TINYJSON,
	}

	if err := cli.CreateAnsibleJob(&j); err != nil {
		t.Fatalf("create ansible job failed with err %s", err.Error())
	} else {
		t.Logf("create ansible job %s successful", j.ID)
	}
}

func TestGetAnsibleJob(t *testing.T) {
	cli := newPostgreSQL()
	defer cli.CloseService()

	if j, err := cli.GetAnsibleJob( "1"); err != nil {
		t.Fatalf("get ansible job failed with err %s", err.Error())
	} else {
		t.Logf("get ansible job %s with status %s successful", j.ID, j.Status)
	}
}

func TestListAnsibleJob(t *testing.T) {
	cli := newPostgreSQL()
	defer cli.CloseService()

	if js, err := cli.ListAnsibleJob(); err != nil {
		t.Fatalf("list ansible job failed with err %s", err.Error())
	} else {
		for _, j := range js {
			t.Logf("list ansible job %s with status %s successful", j.ID, j.Status)
		}
	}
}

func TestDeleteAnsibleJob(t *testing.T) {
	cli := newPostgreSQL()
	defer cli.CloseService()
	if err := cli.DeleteAnsibleJob( "2"); err != nil {
		t.Fatalf("delete ansible job failed with err %s", err.Error())
	} else {
		t.Logf("delete ansible job successful")
	}
}
