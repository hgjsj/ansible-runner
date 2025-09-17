package db

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type PostgreSQL struct {
	db *gorm.DB
	DBService
}

type AnsibleJobSQL struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	Playbook     string    `gorm:"type:text"`
	Status       string    `gorm:"type:text"`
	Output       string    `gorm:"type:text"`
	Error        string    `gorm:"type:text"`
	Step         string    `gorm:"type:text"`
	StartTime    time.Time `gorm:"type:timestamptz"`
	EndTime      time.Time `gorm:"type:timestamptz"`
	InventoryID  uint64    `gorm:"type:bigint"`
	OutputPlugin string    `gorm:"type:text"`
}

type inventorySQL struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"type:text;uniqueIndex"`
	Property string `gorm:"type:jsonb"`
}


var PostgreSQLClient *PostgreSQL

// NewPostgreSQL creates a new PostgreSQL configuration.impletemet with singleton pattern
func newPostgreSQL() *PostgreSQL {
	once.Do(func() {
		fmt.Print("initialize PostgreSQL client")

		password, err := utils.Decrypt(viper.GetString("postgresql.password"))
		if err != nil {
			utils.Log.Error(err.Error())
			panic(err)
		}

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			viper.GetString("postgresql.host"),
			viper.GetString("postgresql.user"),
			password,
			viper.GetString("postgresql.database"),
			viper.GetString("postgresql.port"),
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			utils.Log.Error(err.Error())
			panic(err)
		}

		// if err := db.AutoMigrate(&AnsibleJobSQL{}, &inventorySQL{}); err != nil {
		//  	utils.Log.Error(err.Error())
		//  	panic(err)
		// }
		PostgreSQLClient = &PostgreSQL{db: db}

	})

	return PostgreSQLClient
}

func (p *PostgreSQL) get(ctx context.Context, obj any) error {
	res := p.db.WithContext(ctx).Find(obj)
	if res.Error != nil {
		utils.Log.Error(res.Error.Error())
		return res.Error
	}
	return nil
}

func (p *PostgreSQL) create(ctx context.Context, data any) error {
	// insert a record to the table
	res := p.db.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(data)
	if res.Error != nil {
		utils.Log.Error(res.Error.Error())
		return res.Error
	}
	return nil
}

func (p *PostgreSQL) list(ctx context.Context, objs any) error {
	res := p.db.WithContext(ctx).Find(objs)
	if res.Error != nil {
		utils.Log.Error(res.Error.Error())
		return res.Error
	}
	return nil
}
func (p *PostgreSQL) delete(ctx context.Context, obj any) error {
	res := p.db.WithContext(ctx).Delete(obj)
	if res.Error != nil {
		utils.Log.Error(res.Error.Error())
		return res.Error
	}
	return nil
}

func (p *PostgreSQL) close() error {
	db, err := p.db.DB()
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	db.Close()
	return nil
}

func (p *PostgreSQL) GetInventory(id string) (*models.Inventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	iSQL := inventorySQL{ID: n}
	if err := p.get(ctx, &iSQL); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	var objmap map[string]json.RawMessage
	if err = json.Unmarshal([]byte(iSQL.Property), &objmap); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	var i models.Inventory
	i.ID = strconv.FormatUint(iSQL.ID, 10)
	i.Name = iSQL.Name
	if err = json.Unmarshal(objmap["groups"], &i.Groups); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	return &i, nil
}

func (p *PostgreSQL) CreateInventory(i *models.Inventory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objMap := make(map[string]any)
	objMap["groups"] = i.Groups
	property, err := json.Marshal(objMap)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	inv := inventorySQL{
		Name:     i.Name,
		Property: string(property),
	}
	if i.ID != "" {
		inv.ID, err = strconv.ParseUint(i.ID, 10, 64)
		if err != nil {
			utils.Log.Error(err.Error())
			return err
		}
	}
	if err := p.create(ctx, &inv); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	i.ID = strconv.FormatUint(inv.ID, 10)
	return nil
}

func (p *PostgreSQL) DeleteInventory(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	iSQL := inventorySQL{ID: n}

	if err := p.delete(ctx, &iSQL); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

func (p *PostgreSQL) ListInventory() ([]models.Inventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var invs []inventorySQL
	if err := p.list(ctx, &invs); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	var result []models.Inventory
	for _, iSQL := range invs {
		var objmap map[string]json.RawMessage
		if err := json.Unmarshal([]byte(iSQL.Property), &objmap); err != nil {
			utils.Log.Error(err.Error())
			return nil, err
		}
		var i models.Inventory
		i.ID = strconv.FormatUint(iSQL.ID, 10)
		i.Name = string(iSQL.Name)
		if err := json.Unmarshal(objmap["groups"], &i.Groups); err != nil {
			utils.Log.Error(err.Error())
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}

func (p *PostgreSQL) GetAnsibleJob(id string) (*models.AnsibleJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	jSQL := AnsibleJobSQL{ID: n}
	if err := p.get(ctx, &jSQL); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	j := models.AnsibleJob{
		ID:           strconv.FormatUint(jSQL.ID, 10),
		Playbook:     jSQL.Playbook,
		Status:       jSQL.Status,
		Output:       jSQL.Output,
		Error:        jSQL.Error,
		Step:         jSQL.Step,
		StartTime:    jSQL.StartTime,
		EndTime:      jSQL.EndTime,
		InventoryID:  strconv.FormatUint(jSQL.InventoryID, 10),
		OutputPlugin: jSQL.OutputPlugin,
	}

	return &j, nil
}

func (p *PostgreSQL) CreateAnsibleJob(j *models.AnsibleJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := strconv.ParseUint(j.InventoryID, 10, 64)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}

	jSQL := &AnsibleJobSQL{
		Playbook:     j.Playbook,
		Status:       j.Status,
		Output:       j.Output,
		Error:        j.Error,
		Step:         j.Step,
		StartTime:    j.StartTime,
		EndTime:      j.EndTime,
		InventoryID:  n,
		OutputPlugin: j.OutputPlugin,
	}
	if j.ID != "" {
		jSQL.ID, err = strconv.ParseUint(j.ID, 10, 64)
		if err != nil {
			utils.Log.Error(err.Error())
			return err
		}
	}

	if err := p.create(ctx, jSQL); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	j.ID = strconv.FormatUint(jSQL.ID, 10)
	return nil
}

func (p *PostgreSQL) DeleteAnsibleJob(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}

	jSQL := AnsibleJobSQL{ID: n}
	if err := p.delete(ctx, &jSQL); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

func (p *PostgreSQL) ListAnsibleJob() ([]models.AnsibleJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var jobsSQL []AnsibleJobSQL
	if err := p.list(ctx, &jobsSQL); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	var jobs []models.AnsibleJob
	for _, jSQL := range jobsSQL {
		j := models.AnsibleJob{
			ID:           strconv.FormatUint(jSQL.ID, 10),
			Playbook:     jSQL.Playbook,
			Status:       jSQL.Status,
			Output:       jSQL.Output,
			Error:        jSQL.Error,
			Step:         jSQL.Step,
			StartTime:    jSQL.StartTime,
			EndTime:      jSQL.EndTime,
			InventoryID:  strconv.FormatUint(jSQL.InventoryID, 10),
			OutputPlugin: jSQL.OutputPlugin,
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (p *PostgreSQL) CloseService() {
	p.close()
}
