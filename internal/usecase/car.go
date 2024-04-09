package usecase

import (
	"carinfo/internal/entity"
	"carinfo/internal/repository"
	"carinfo/internal/repository/postgresql"
	"carinfo/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Cars interface {
	Create(ctx context.Context, cars []entity.Car) ([]entity.Car, error)
	ReadById(ctx context.Context, id string) (entity.Car, error)
	ReadAll(ctx context.Context, options map[string]string) ([]entity.Car, error)
	UpdateInfo(ctx context.Context, id string, car entity.Car) (entity.Car, error)
	Delete(ctx context.Context, id string) error
}

type CarUsecase struct {
	CarRepos repository.Cars
}

func NewCarUsecase(db *pgxpool.Pool) CarUsecase {
	repos := postgresql.NewCarRepo(db)
	return CarUsecase{
		CarRepos: repos,
	}
}

func (c CarUsecase) Create(ctx context.Context, cars []entity.Car) ([]entity.Car, error) {
	var carList []entity.Car
	for _, car := range cars {
		existId, err := c.CarRepos.ReadByRegNum(ctx, car.RegNum)
		if err != nil {
			return nil, err
		}
		if existId != "" {
			logger.Log.Debug("car regNum already exists", car.RegNum)
			continue
		}
		id, err := c.CarRepos.Create(ctx, car)
		if err != nil {
			return nil, err
		}
		newCar, err := c.CarRepos.ReadById(ctx, id)
		if err != nil {
			return nil, err
		}
		carList = append(carList, newCar)
	}
	return carList, nil
}

func (c CarUsecase) ReadById(ctx context.Context, id string) (entity.Car, error) {
	return c.CarRepos.ReadById(ctx, id)
}

func (c CarUsecase) ReadAll(ctx context.Context, options map[string]string) ([]entity.Car, error) {
	return c.CarRepos.ReadAll(ctx, options)
}

func (c CarUsecase) UpdateInfo(ctx context.Context, id string, car entity.Car) (entity.Car, error) {
	return c.CarRepos.Update(ctx, id, car)
}

func (c CarUsecase) Delete(ctx context.Context, id string) error {
	return c.CarRepos.Delete(ctx, id)
}
