package repository

import (
	"bot/domain"
	"context"
	"errors"
)

var InsertError = errors.New("failed to insert event")

const insertEventQuery = `INSERT INTO orders (
							symbol, 
                          	side, 
							type, 
							order_price, 
							order_size,
							actual_price, 
							actual_amount,
							timestamp
						) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

func (repo *OrderEventsStorage) StoreEvent(ctx context.Context, event domain.OrderEvent) error {
	commandTag, err := repo.pool.Exec(ctx, insertEventQuery,
		event.ExecOrder.Symbol,
		event.ExecOrder.Side,
		event.ExecOrder.Type,
		event.ExecOrder.LimitPrice,
		event.ExecOrder.Quantity,
		event.Price,
		event.Amount,
		event.ExecOrder.TS)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return InsertError
	}
	return nil
}
