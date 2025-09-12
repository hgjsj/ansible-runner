package db

import (
	"ansible-runner/internal/models"
	"context"
)

type DBOperator interface {
	Get(ctx context.Context, data any) error
	List(ctx context.Context, data []any) error
	Create(ctx context.Context, data any) error
	Delete(ctx context.Context, data any) error
	Close() error
}

type DBConvertor interface {
	ConvertInventory(i models.Inventory) any
	ConvertAnsibleJob(j models.AnsibleJob) any
}

type DBService interface {
	GetInventory(ctx context.Context, id string) (*models.Inventory, error)
	GetAnsibleJob(ctx context.Context, id string) (*models.AnsibleJob, error)
	ListInventory(ctx context.Context) ([]models.AnsibleJob, error)
	ListAnsibleJob(ctx context.Context) ([]models.AnsibleJob, error)
	CreateInventory(ctx context.Context, i *models.Inventory) error
	CreateAnsibleJob(ctx context.Context, j *models.AnsibleJob) error
	DeleteInventory(ctx context.Context, id string) error
	DeleteAnsibleJob(ctx context.Context, id string) error
}
