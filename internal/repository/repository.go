package repository

import (
	"carinfo/internal/entity"
	"context"
)

type Cars interface {
	Create(context.Context, entity.Car) (string, error)
	ReadById(context.Context, string) (entity.Car, error)
	ReadByRegNum(context.Context, string) (string, error)
	ReadAll(context.Context, map[string]string) ([]entity.Car, error)
	Update(context.Context, string, entity.Car) (entity.Car, error)
	Delete(context.Context, string) error
}
type Repositories struct {
	Cars Cars
}
