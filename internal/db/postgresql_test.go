package db

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"context"
	"testing"
	"time"
)

func init() {
	utils.LoadConfig("../../configs/config")
	utils.InitLogger()
}

func TestNewClient(t *testing.T) {
	NewPostgreSQL()
	t.Logf("new client succesfully successful")

}

func TestCreateInventory(t *testing.T) {
	cli := NewPostgreSQL()

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
	defer cli.Close()
	inv := models.Inventory{
		Name:   "pcm-p4s-postinstall",
		Groups: []models.Group{group},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	if err := cli.CreateInventory(ctx, &inv); err != nil {
		t.Fatalf("create inventory failed with err %s", err.Error())
	} else {
		t.Logf("create inventory %s successful", inv.Name)
	}

}

func TestGetInventory(t *testing.T) {
	cli := NewPostgreSQL()
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if i, err := cli.GetInventory(ctx, "3"); err != nil {
		t.Fatalf("get inventory failed with err %s", err.Error())
	} else {
		t.Logf("get inventory %s with group name %s successful", i.Name, i.Groups[0].Name)
	}
}

func TestListInventory(t *testing.T) {
	cli := NewPostgreSQL()
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if is, err := cli.ListInventory(ctx); err != nil {
		t.Fatalf("list inventory failed with err %s", err.Error())
	} else {
		for _, inv := range is {
			t.Logf("list inventory %s with group name %s successful", inv.Name, inv.Groups[0].Name)
		}
	}
}

func TestDeleteInventory(t *testing.T) {
	cli := NewPostgreSQL()
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cli.DeleteInventory(ctx, "2"); err != nil {
		t.Fatalf("delete inventory failed with err %s", err.Error())
	} else {
		t.Logf("delete inventory successful")
	}
}

func TestCreateAnsibleJob(t *testing.T) {
	cli := NewPostgreSQL()
	defer cli.Close()

	j := models.AnsibleJob{
		Playbook: "test.yml",
		Status:   models.Scheduled,
		//InventoryID:  primitive.NewObjectID(),
		OutputPlugin: models.TINYJSON,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cli.CreateAnsibleJob(ctx, &j); err != nil {
		t.Fatalf("create ansible job failed with err %s", err.Error())
	} else {
		t.Logf("create ansible job %s successful", j.Playbook)
	}
}

func TestGetAnsibleJob(t *testing.T) {
	cli := NewPostgreSQL()
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if j, err := cli.GetAnsibleJob(ctx, "1"); err != nil {
		t.Fatalf("get ansible job failed with err %s", err.Error())
	} else {
		t.Logf("get ansible job %s with status %s successful", j.Playbook, j.Status)
	}
}

func TestListAnsibleJob(t *testing.T) {
	cli := NewPostgreSQL()
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if js, err := cli.ListAnsibleJob(ctx); err != nil {
		t.Fatalf("list ansible job failed with err %s", err.Error())
	} else {
		for _, j := range js {
			t.Logf("list ansible job %s with status %s successful", j.Playbook, j.Status)
		}
	}
}

func TestDeleteAnsibleJob(t *testing.T) {
	cli := NewPostgreSQL()
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cli.DeleteAnsibleJob(ctx, "1"); err != nil {
		t.Fatalf("delete ansible job failed with err %s", err.Error())
	} else {
		t.Logf("delete ansible job successful")
	}
}
