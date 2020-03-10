package burgers

import (
	"context"
	"errors"
	"github.com/AlisherFozilov/crud/pkg/crud/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BurgersSvc struct {
	pool *pgxpool.Pool // dependency
}

func NewBurgersSvc(pool *pgxpool.Pool) *BurgersSvc {
	if pool == nil {
		panic(errors.New("pool can't be nil"))
	}
	return &BurgersSvc{pool: pool}
}

func (service *BurgersSvc) InitDB() error {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return NewDbError(err)
	}
	_, err = conn.Query(context.Background(), burgersDDL)
	if err != nil {
		return NewQueryError(burgersDDL, err)
	}
	return nil
}

func (service *BurgersSvc) BurgersList() (list []models.Burger, err error) {
	list = make([]models.Burger, 0) // TODO: for REST API
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return nil, NewDbError(err) // TODO: wrap to specific error
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), getIdNamePriceSQL)
	if err != nil {
		return nil, NewQueryError(getIdNamePriceSQL, err) // TODO: wrap to specific error
	}
	defer rows.Close()

	for rows.Next() {
		item := models.Burger{}
		err := rows.Scan(&item.Id, &item.Name, &item.Price)
		if err != nil {
			return nil, NewDbError(err) // TODO: wrap to specific error
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, NewDbError(err)
	}

	return list, nil
}

func (service *BurgersSvc) Save(model models.Burger) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return NewDbError(err) // TODO: wrap to specific error
	}
	defer conn.Release()


	_, err = conn.Exec(context.Background(), insertNamePriceSQL, model.Name, model.Price)
	if err != nil {
		return NewQueryError(insertNamePriceSQL, err)
	}

	return nil
}

func (service *BurgersSvc) RemoveById(id int64) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return NewDbError(err) // TODO: wrap to specific error
	}
	defer conn.Release()


	_, err = conn.Exec(context.Background(), setRemovedTrueByIdSQL, id)
	if err != nil {
		return NewQueryError(setRemovedTrueByIdSQL, err)
	}

	return nil
}
