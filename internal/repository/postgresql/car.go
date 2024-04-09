package postgresql

import (
	"carinfo/internal/entity"
	"carinfo/pkg/logger"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"strconv"
)

type CarRepo struct {
	db Client
}

func NewCarRepo(db Client) *CarRepo {
	return &CarRepo{db: db}
}

const (
	carsTable       = "public.cars"
	carId           = "id"
	model           = "model"
	mark            = "mark"
	year            = "year"
	regNum          = "reg_num"
	ownerName       = "owner_name"
	ownerSurname    = "owner_surname"
	ownerPatronymic = "owner_patronymic"
)

func (r CarRepo) Create(ctx context.Context, car entity.Car) (string, error) {
	logger.Log.Debug("creating car at db starting")
	pgq := `INSERT INTO public.cars (mark, model, year, reg_num, owner_name, owner_surname, owner_patronymic)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`
	err := r.db.QueryRow(ctx, pgq, car.Mark, car.Model, car.Year, car.RegNum, car.Owner.Name, car.Owner.Surname, car.Owner.Patronymic).Scan(&car.ID)
	if err != nil {
		logger.Log.Error("cannot add car to db", "error:", err.Error(), "!")
		return "", err
	}
	logger.Log.Debug("creating car at db ended by id", car.ID)
	return car.ID, nil
}

func (r CarRepo) ReadById(ctx context.Context, id string) (entity.Car, error) {
	logger.Log.Debug("ReadById car from db starting", "id", id)
	car := entity.NewCar()
	pgq := `SELECT id, mark, model, year, reg_num, owner_name, owner_surname, owner_patronymic
			FROM public.cars
			WHERE id = $1`
	rows, err := r.db.Query(ctx, pgq, id)
	if err != nil {
		logger.Log.Error("readById query exec error:", err, "!")
		return entity.Car{}, err
	}
	for rows.Next() {
		if err = rows.Scan(&car.ID, &car.Mark, &car.Model, &car.Year, &car.RegNum, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic); err != nil {
			logger.Log.Error("readById rows iterating error:", err.Error(), "!")
			return entity.Car{}, err
		}
		logger.Log.Debug("car getting from db successful", "id", id)
		return car, nil
	}
	logger.Log.Debug("ReadById car from db ended", "id", id)
	return entity.Car{}, err
}

func (r CarRepo) ReadByRegNum(ctx context.Context, num string) (string, error) {
	logger.Log.Debug("ReadByRegNum car from db starting", "num", num)
	var newId string
	pgq := `SELECT id
			FROM public.cars
			WHERE reg_num = $1`
	rows, err := r.db.Query(ctx, pgq, num)
	if err != nil {
		logger.Log.Error("ReadByRegNum query exec error:", err, "!")
		return newId, err
	}
	for rows.Next() {
		if err = rows.Scan(&newId); err != nil {
			logger.Log.Error("ReadByRegNum rows iterating error:", err.Error(), "!")
			return newId, err
		}
		logger.Log.Debug("car getting from db successful", "num", num)
		return newId, nil
	}
	logger.Log.Debug("ReadByRegNum car from db ended", "num", num)
	return newId, err
}

func (r CarRepo) ReadAll(ctx context.Context, opts map[string]string) ([]entity.Car, error) {
	logger.Log.Debug("ReadAll cars from db starting")
	for k, v := range opts {
		logger.Log.Debug("ReadAll got advanced opts", k, v)
	}
	q := sq.Select(carId, mark, model, year, regNum, ownerName, ownerSurname, ownerPatronymic).
		From(carsTable)
	page := 1
	perPage := 10
	if value, ok := opts["page"]; ok {
		page, _ = strconv.Atoi(value)
	}
	if value, ok := opts["per_page"]; ok {
		perPage, _ = strconv.Atoi(value)
		q = q.Limit(uint64(perPage)).Offset(uint64((page - 1) * perPage))
	}
	for k, v := range opts {
		if k != "page" && k != "per_page" {
			q = q.Where(fmt.Sprintf("%s = '%s'", k, v))
		}
	}
	sql, args, err := q.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logger.Log.Error("findAll query exec error:", err, "!")
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		logger.Log.Error("readAll cars query exec error:", err.Error(), "!")
		return nil, err
	}
	var cars = make([]entity.Car, 0)
	for rows.Next() {
		var car entity.Car
		if err = rows.Scan(&car.ID, &car.Mark, &car.Model, &car.Year, &car.RegNum, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic); err != nil {
			logger.Log.Error("readAll rows iterating error:", err.Error(), "!")
			return nil, err
		}
		cars = append(cars, car)
	}
	logger.Log.Info("allCars getting from db successful")
	logger.Log.Debug("ReadAll cars from db ended")
	return cars, nil
}

func (r CarRepo) Update(ctx context.Context, id string, car entity.Car) (entity.Car, error) {
	logger.Log.Debug("Update car at db staring", id, car)
	q := sq.Update(carsTable).
		Where(sq.Eq{carId: id}).
		PlaceholderFormat(sq.Dollar)
	if car.Mark != "" {
		q = q.Set(mark, car.Mark)
	}
	if car.Model != "" {
		q = q.Set(model, car.Model)
	}
	if car.Year != 0 {
		q = q.Set(year, car.Year)
	}
	if car.RegNum != "" {
		q = q.Set(regNum, car.RegNum)
	}
	if car.Owner.Name != "" {
		q = q.Set(ownerName, car.Owner.Name)
	}
	if car.Owner.Surname != "" {
		q = q.Set(ownerSurname, car.Owner.Surname)
	}
	if car.Owner.Patronymic != "" {
		q = q.Set(ownerPatronymic, car.Owner.Patronymic)
	}
	sql, args, err := q.ToSql()
	if err != nil {
		logger.Log.Error("car info updating query building error", err.Error(), "!")
		return entity.Car{}, err
	}
	res, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		logger.Log.Error("updating data error", err.Error(), "!")
		return entity.Car{}, err
	}
	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		logger.Log.Error("no rows affected")
		return entity.Car{}, fmt.Errorf("no rows affected")
	}
	logger.Log.Info("car editing successful")
	logger.Log.Debug("Update car at db ended", id, car)
	return r.ReadById(ctx, id)
}

func (r CarRepo) Delete(ctx context.Context, id string) error {
	logger.Log.Debug("Delete car staring", "id", id)
	pgq := "DELETE FROM public.cars WHERE id = $1"
	_, err := r.db.Exec(ctx, pgq, id)
	if err != nil {
		logger.Log.Error("car deleting exec error:", err.Error(), "!")
		return err
	}
	logger.Log.Info("car deleted", "id", id)
	logger.Log.Debug("Delete car ended", "id", id)
	return nil
}
