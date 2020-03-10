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

func (service *BurgersSvc) BurgersList() (list []models.Burger, err error) {
	list = make([]models.Burger, 0) // TODO: for REST API
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return nil, NewDbError(err) // TODO: wrap to specific error
	}
	defer conn.Release()
	const query = "SELECT id, name, price FROM burgers WHERE removed = FALSE"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, NewQueryError(query, err) // TODO: wrap to specific error
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

	const query = `INSERT INTO burgers (name, price) VALUES ($1, $2);`
	_, err = conn.Exec(context.Background(), query, model.Name, model.Price)
	if err != nil {
		return NewQueryError(query, err)
	}

	return nil
}

func (service *BurgersSvc) RemoveById(id int64) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return NewDbError(err) // TODO: wrap to specific error
	}
	defer conn.Release()

	const query = `UPDATE burgers SET removed = true WHERE id = $1;`
	_, err = conn.Exec(context.Background(), query, id)
	if err != nil {
		return NewQueryError(query, err)
	}

	return nil
}
