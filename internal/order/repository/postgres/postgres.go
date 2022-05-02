package postgres

import (
	"context"
	"time"

	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
)

type OrderPostgresStorage struct {
	db *pgx.Conn
}

func NewOrderPostgresStorage(dbURI string) order.OrderRepository {
	log.Debug().Msg("new order postgres storage")
	//MakeMigrations(dbURI)

	log.Debug().Msg("pgx connect")
	conn, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &OrderPostgresStorage{
		db: conn,
	}
}

func (ops *OrderPostgresStorage) GetOrderByOrderUID(ctx context.Context, orderNumber string) (*models.Order, error) {
	var result = new(models.Order)
	var accrual int32
	var timeValue pgtype.Timestamp

	err := ops.db.QueryRow(ctx,
		"SELECT user_id, order_id, order_status, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE (debet IS TRUE) AND (order_id = $1);", orderNumber).
		Scan(&result.UserName, &result.Number, &result.Status, &accrual, &timeValue)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	result.Accrual = float32(accrual) / 100
	result.Uploaded = timeValue.Time.Format(time.RFC3339)
	return result, nil

}

func (ops *OrderPostgresStorage) InsertOrder(ctx context.Context, userID string, orderNumber string) error {
	orderItem, err := ops.GetOrderByOrderUID(ctx, orderNumber)
	if orderItem != nil {
		if orderItem.UserName == userID {
			return order.ErrOrderAlreadyInsertedByUser
		} else {
			return order.ErrOrderAlreadyInsertedByOtherUser
		}
	}
	if err != nil {
		return err
	}

	_, err = ops.db.Exec(ctx,
		"INSERT INTO orders "+
			"(user_id, order_id, debet, order_status, accrual) "+
			"VALUES ($1, $2, TRUE, $3, $4);",
		userID, orderNumber, models.OrderStatusNew, 0)

	return err
}

func (ops *OrderPostgresStorage) GetOrdersListByUserID(ctx context.Context, userID string) ([]*models.Order, error) {
	result := make([]*models.Order, 0)

	rows, err := ops.db.Query(ctx,
		"SELECT user_id, order_id, order_status, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE (debet IS TRUE) AND (user_id = $1) "+
			"ORDER BY uploaded_at ASC;", userID)

	if err != nil {
		log.Debug().Err(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Order
		var timeValue pgtype.Timestamp
		var accrual int32
		err := rows.Scan(&item.UserName, &item.Number, &item.Status, &accrual, &timeValue)
		item.Accrual = float32(accrual) / 100
		item.Uploaded = timeValue.Time.Format(time.RFC3339)
		if err != nil {
			log.Debug().Err(err)
			return nil, err
		}
		result = append(result, &item)
	}

	return result, nil
}

func (ops *OrderPostgresStorage) GetBalanceByUserID(ctx context.Context, userID string) (*models.Balance, error) {
	var accrual int32
	err := ops.db.QueryRow(ctx,
		"SELECT COALESCE(SUM(accrual), 0) "+
			"FROM orders "+
			"WHERE (user_id = $1);", userID).Scan(&accrual)

	if err != nil {
		return nil, err
	}

	var withdrawn int32
	err = ops.db.QueryRow(ctx,
		"SELECT COALESCE(SUM(accrual), 0) "+
			"FROM orders "+
			"WHERE (debet IS FALSE) AND (user_id = $1);", userID).Scan(&withdrawn)

	if err != nil {
		return nil, err
	}

	withdrawn = -1 * withdrawn

	return &models.Balance{
		Current:   float32(accrual) / 100,
		Withdrawn: float32(withdrawn) / 100,
	}, nil
}

func (ops *OrderPostgresStorage) WithdrawBalance(ctx context.Context, userID string, bw *models.BalanceWithdraw) error {
	orderItem, err := ops.GetOrderByOrderUID(ctx, bw.OrderID)
	if orderItem != nil {
		return order.ErrOrderBadNumber
	}
	if err != nil {
		return nil
	}

	var balance int32

	err = ops.db.QueryRow(ctx,
		"SELECT COALESCE(SUM(accrual), 0) "+
			"FROM orders "+
			"WHERE user_id = $1 ;", userID).Scan(&balance)

	if err != nil {
		return nil
	}

	if float32(balance)/100 < bw.Sum {
		return order.ErrNotEnougthBalance
	}

	_, err = ops.db.Exec(ctx,
		"INSERT INTO orders "+
			"(user_id, order_id, debet, order_status, accrual) "+
			"VALUES ($1, $2, FALSE, $3, $4);",
		userID, bw.OrderID, models.OrderStatusWithDrawn, -100*bw.Sum)

	return err

}

func (ops *OrderPostgresStorage) Withdrawals(ctx context.Context, userID string) ([]*models.Withdrawals, error) {
	result := make([]*models.Withdrawals, 0)

	rows, err := ops.db.Query(ctx,
		"SELECT order_id, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE (debet IS FALSE) AND (user_id = $1) AND (order_status = $2) "+
			"ORDER BY uploaded_at ASC;", userID, models.OrderStatusWithDrawn)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Withdrawals
		var sum int32
		var processedAt pgtype.Timestamp
		err := rows.Scan(&item.OrderID, &sum, &processedAt)
		if err != nil {
			return nil, nil
		}
		item.Sum = float32(sum) / 100
		item.ProcessedAt = processedAt.Time.Format(time.RFC3339)
		result = append(result, &item)
	}

	return result, nil
}

func (ops *OrderPostgresStorage) UpdateOrder(ctx context.Context, orderNumber string, orderStatus string, orderAccrual int32) error {
	_, err := ops.db.Exec(ctx,
		"UPDATE orders "+
			"SET order_status = $1 , accrual = $2 "+
			"WHERE order_id = $3 ;", orderStatus, orderAccrual, orderNumber)
	return err
}

func (ops *OrderPostgresStorage) GetNotFinnalizedOrdersListByUserID(ctx context.Context, userID string) ([]*models.Order, error) {
	result := make([]*models.Order, 0)

	rows, err := ops.db.Query(ctx,
		"SELECT user_id, order_id, order_status, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE (debet IS TRUE) AND (user_id = $1) AND order_status NOT IN ($2, $3)"+
			"ORDER BY uploaded_at ASC;", userID, models.OrderStatusProcessed, models.OrderStatusInvalid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Order
		var accrual int32
		var uploaded pgtype.Timestamp
		err := rows.Scan(&item.UserName, &item.Number, &item.Status, &accrual, &uploaded)
		if err != nil {
			log.Debug().Err(err)
			return nil, nil
		}
		item.Accrual = float32(accrual) / 100
		item.Uploaded = uploaded.Time.Format(time.RFC3339)
		result = append(result, &item)
	}

	return result, nil
}

func (ops *OrderPostgresStorage) GetNotFinnalizedOrdersList(ctx context.Context) ([]*models.Order, error) {
	result := make([]*models.Order, 0)

	rows, err := ops.db.Query(ctx,
		"SELECT user_id, order_id, order_status, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE (debet IS TRUE) AND order_status NOT IN ($1, $2)"+
			"ORDER BY uploaded_at ASC;", models.OrderStatusProcessed, models.OrderStatusInvalid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Order
		var accrual int32
		var uploaded pgtype.Timestamp
		err := rows.Scan(&item.UserName, &item.Number, &item.Status, &accrual, &uploaded)
		if err != nil {
			log.Debug().Err(err)
			return nil, nil
		}
		item.Accrual = float32(accrual) / 100
		item.Uploaded = uploaded.Time.Format(time.RFC3339)
		result = append(result, &item)
	}

	return result, nil
}
