package service

import (
	"carinfo/internal/entity"
	"carinfo/internal/usecase"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Cars interface {
	Create(ctx context.Context, cars []entity.Car) ([]entity.Car, error)
	ReadById(ctx context.Context, id string) (entity.Car, error)
	ReadAll(ctx context.Context, options map[string]string) ([]entity.Car, error)
	Update(ctx context.Context, id string, car entity.Car) (entity.Car, error)
	Delete(ctx context.Context, id string) error
}

type CarService struct {
	Usecase usecase.Cars
}

func NewCarService(db *pgxpool.Pool) CarService {
	carUsecase := usecase.NewCarUsecase(db)
	return CarService{
		Usecase: carUsecase,
	}
}

func (c CarService) Create(ctx context.Context, cars []entity.Car) ([]entity.Car, error) {
	return c.Usecase.Create(ctx, cars)
}

func (c CarService) ReadById(ctx context.Context, id string) (entity.Car, error) {
	return c.Usecase.ReadById(ctx, id)
}

func (c CarService) ReadAll(ctx context.Context, options map[string]string) ([]entity.Car, error) {
	return c.Usecase.ReadAll(ctx, options)
}

func (c CarService) Update(ctx context.Context, id string, car entity.Car) (entity.Car, error) {
	return c.Usecase.UpdateInfo(ctx, id, car)
}

func (c CarService) Delete(ctx context.Context, id string) error {
	return c.Usecase.Delete(ctx, id)
}
