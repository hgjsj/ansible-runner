package db

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	db *gorm.DB
	DBService
}

type jobSQL struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Playbook     string
	Status       string
	Output       string `gorm:"type:text"`
	Error        string `gorm:"type:text"`
	Step         string
	StartTime    time.Time
	EndTime      time.Time
	InventoryID  string
	OutputPlugin string
}

type inventorySQL struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	Name     string
	Property string `gorm:"type:jsonb"`
}

var once sync.Once
var PostgreSQLClient *PostgreSQL

// NewPostgreSQL creates a new PostgreSQL configuration.impletemet with singleton pattern
func NewPostgreSQL() *PostgreSQL {
	once.Do(func() {
		if PostgreSQLClient == nil {

			if viper.GetStringMap("postgresql") == nil {
				utils.Log.Error("PostgreSQL config not found and quitting")
				panic("PostgreSQL config not found")
			}

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

			// if err := db.AutoMigrate(&inventorySQL{}, &jobSQL{}); err != nil {
			// 	utils.Log.Error(err.Error())
			// 	panic(err)
			// }
			PostgreSQLClient = &PostgreSQL{db: db}
		}
	})

	return PostgreSQLClient
}

func (p *PostgreSQL) Get(ctx context.Context, obj any) error {
	res := p.db.WithContext(ctx).Find(obj)
	if res.Error != nil {
		utils.Log.Error(res.Error.Error())
		return res.Error
	}
	return nil
}

func (p *PostgreSQL) Create(ctx context.Context, data any) error {
	// insert a record to the table
	res := p.db.WithContext(ctx).Create(data)
	if res.Error != nil {
		utils.Log.Error(res.Error.Error())
		return res.Error
	}
	return nil
}

func (p *PostgreSQL) List(ctx context.Context, objs any) error {
	res := p.db.WithContext(ctx).Find(objs)
	if res.Error != nil {
		utils.Log.Error(res.Error.Error())
		return res.Error
	}
	return nil
}
func (p *PostgreSQL) Delete(ctx context.Context, obj any) error {
	res := p.db.WithContext(ctx).Delete(obj)
	if res.Error != nil {
		utils.Log.Error(res.Error.Error())
		return res.Error
	}
	return nil
}

func (p *PostgreSQL) Close() error {
	db, err := p.db.DB()
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	db.Close()
	return nil
}

func (p *PostgreSQL) GetInventory(ctx context.Context, id string) (*models.Inventory, error) {
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	iSQL := inventorySQL{ID: n}
	if err := p.Get(ctx, &iSQL); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	var objmap map[string]json.RawMessage
	if err = json.Unmarshal([]byte(iSQL.Property), &objmap); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	var i models.Inventory
	i.Name = string(iSQL.Name)
	if err = json.Unmarshal(objmap["groups"], &i.Groups); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	return &i, nil
}

func (p *PostgreSQL) CreateInventory(ctx context.Context, i *models.Inventory) error {
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
	if err := p.Create(ctx, &inv); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

func (p *PostgreSQL) DeleteInventory(ctx context.Context, id string) error {
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	iSQL := inventorySQL{ID: n}
	if err := p.Delete(ctx, &iSQL); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

func (p *PostgreSQL) ListInventory(ctx context.Context) ([]models.Inventory, error) {
	var invs []inventorySQL
	if err := p.List(ctx, &invs); err != nil {
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
		i.Name = string(iSQL.Name)
		if err := json.Unmarshal(objmap["groups"], &i.Groups); err != nil {
			utils.Log.Error(err.Error())
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}

func (p *PostgreSQL) GetAnsibleJob(ctx context.Context, id string) (*models.AnsibleJob, error) {
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	jSQL := jobSQL{ID: n}
	if err := p.Get(ctx, &jSQL); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	// invID, err := primitive.ObjectIDFromHex(jSQL.InventoryID)
	// if err != nil {
	// 	utils.Log.Error(err.Error())
	// 	return nil, err
	// }

	j := models.AnsibleJob{
		Playbook:  jSQL.Playbook,
		Status:    jSQL.Status,
		Output:    jSQL.Output,
		Error:     jSQL.Error,
		Step:      jSQL.Step,
		StartTime: jSQL.StartTime,
		EndTime:   jSQL.EndTime,
		//InventoryID:  invID,
		OutputPlugin: jSQL.OutputPlugin,
	}

	return &j, nil
}

func (p *PostgreSQL) CreateAnsibleJob(ctx context.Context, j *models.AnsibleJob) error {
	jSQL := jobSQL{
		Playbook:  j.Playbook,
		Status:    j.Status,
		Output:    j.Output,
		Error:     j.Error,
		Step:      j.Step,
		StartTime: j.StartTime,
		EndTime:   j.EndTime,
		//InventoryID:  j.InventoryID.Hex(),
		OutputPlugin: j.OutputPlugin,
	}
	if err := p.Create(ctx, &jSQL); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

func (p *PostgreSQL) DeleteAnsibleJob(ctx context.Context, id string) error {
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	jSQL := jobSQL{ID: n}
	if err := p.Delete(ctx, &jSQL); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

func (p *PostgreSQL) ListAnsibleJob(ctx context.Context) ([]models.AnsibleJob, error) {
	var jobsSQL []jobSQL
	if err := p.List(ctx, &jobsSQL); err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	var result []models.AnsibleJob
	for _, jSQL := range jobsSQL {
		// invID, err := primitive.ObjectIDFromHex(jSQL.InventoryID)
		// if err != nil {
		// 	utils.Log.Error(err.Error())
		// 	return nil, err
		// }
		j := models.AnsibleJob{
			Playbook:  jSQL.Playbook,
			Status:    jSQL.Status,
			Output:    jSQL.Output,
			Error:     jSQL.Error,
			Step:      jSQL.Step,
			StartTime: jSQL.StartTime,
			EndTime:   jSQL.EndTime,
			//InventoryID:  invID,
			OutputPlugin: jSQL.OutputPlugin,
		}
		result = append(result, j)
	}
	return result, nil
}
