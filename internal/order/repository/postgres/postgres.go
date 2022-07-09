package postgres

import (
	"context"

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
	conn, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &OrderPostgresStorage{
		db: conn,
	}
}

func (ops *OrderPostgresStorage) GetOrderByOrderUID(ctx context.Context, orderNumber string) (*models.Order, error) {
	logger := log.With().Str("package", "postgres").Str("func", "GetOrderByOrderUID").Logger()

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	var result = new(models.Order)
	var accrual int32
	var timeValue pgtype.Timestamp

	logger.Debug().Str("orderNumber", orderNumber).Msg("try to get order")
	err := ops.db.QueryRow(ctx,
		"SELECT user_id, order_id, order_status, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE order_id = $1 ", orderNumber).
		Scan(&result.UserName, &result.Number, &result.Status, &accrual, &timeValue)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	result.Accrual = float32(accrual) / 100
	result.Uploaded = timeValue

	logger.Debug().Int32("order.user", result.UserName).
		Str("order.id", result.Number).
		Str("order.status", result.Status).
		Float32("order.accrual", result.Accrual).
		Time("order.time", result.Uploaded.Time).
		Msg("GetOrderByUID result")

	return result, nil
}

func (ops *OrderPostgresStorage) InsertOrder(ctx context.Context, userID int32, orderNumber string) error {
	logger := log.With().Str("package", "postgres").Str("func", "InsertOrder").Logger()

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	logger.Debug().Int32("userID", userID).Str("orderNumber", orderNumber).Msg("try to add new order")

	cTag, err := ops.db.Exec(ctx,
		"INSERT INTO orders "+
			"(user_id, order_id, debet, order_status, accrual) "+
			"VALUES ($1, $2, TRUE, $3, $4) "+
			"ON CONFLICT (order_id) DO NOTHING",
		userID, orderNumber, models.OrderStatusNew, 0)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return err
	}

	if cTag.RowsAffected() == 0 {
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
	}
	return err
}

func (ops *OrderPostgresStorage) GetOrdersListByUserID(ctx context.Context, userID int32) ([]models.Order, error) {
	logger := log.With().Str("package", "postgres").Str("func", "GetOrdersListByUserID").Logger()

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	result := make([]models.Order, 0)

	logger.Debug().Int32("userID", userID).Msg("try to get order list by user id")
	rows, err := ops.db.Query(ctx,
		"SELECT user_id, order_id, order_status, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE (debet IS TRUE) AND (user_id = $1) "+
			"ORDER BY uploaded_at ASC", userID)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Order
		var timeValue pgtype.Timestamp
		var accrual int32
		err := rows.Scan(&item.UserName, &item.Number, &item.Status, &accrual, &timeValue)
		item.Accrual = float32(accrual) / 100
		item.Uploaded = timeValue
		if err != nil {
			logger.Debug().Err(err).Msg("exit with error")
			return nil, err
		}
		logger.Debug().Int32("order.user", item.UserName).
			Str("order.id", item.Number).
			Str("order.status", item.Status).
			Float32("order.accrual", item.Accrual).
			Time("order.time", item.Uploaded.Time).
			Msg("getOrderListByUID result item")

		result = append(result, item)
	}

	return result, nil
}

func (ops *OrderPostgresStorage) GetBalanceByUserID(ctx context.Context, userID int32) (*models.Balance, error) {
	logger := log.With().Str("package", "postgres").Str("func", "GetBalanceByUserID").Logger()

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	logger.Debug().Int32("userID", userID).Msg("try to get balance by userID")
	var accrual int32
	err := ops.db.QueryRow(ctx,
		"SELECT COALESCE(SUM(accrual), 0) "+
			"FROM orders "+
			"WHERE (user_id = $1);", userID).Scan(&accrual)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return nil, err
	}

	var withdrawn int32
	err = ops.db.QueryRow(ctx,
		"SELECT COALESCE(SUM(accrual), 0) "+
			"FROM orders "+
			"WHERE (debet IS FALSE) AND (user_id = $1);", userID).Scan(&withdrawn)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return nil, err
	}

	withdrawn = -1 * withdrawn

	return &models.Balance{
		Current:   float32(accrual) / 100,
		Withdrawn: float32(withdrawn) / 100,
	}, nil
}

func (ops *OrderPostgresStorage) WithdrawBalance(ctx context.Context, userID int32, bw *models.BalanceWithdraw) error {
	logger := log.With().Str("package", "postgres").Str("func", "WithdrawBalance").Logger()

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	logger.Debug().Int32("userID", userID).Str("orderNumber", bw.OrderID).Float32("sum", bw.Sum).Msg("try to withdraw balance")
	orderItem, err := ops.GetOrderByOrderUID(ctx, bw.OrderID)
	if orderItem != nil {
		logger.Debug().Msg("Bad order number")
		return order.ErrOrderBadNumber
	}
	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return nil
	}

	var balance int32

	err = ops.db.QueryRow(ctx,
		"SELECT COALESCE(SUM(accrual), 0) "+
			"FROM orders "+
			"WHERE user_id = $1 ;", userID).Scan(&balance)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return nil
	}

	if float32(balance)/100 < bw.Sum {
		logger.Debug().Msg("not enougth balance")
		return order.ErrNotEnougthBalance
	}

	_, err = ops.db.Exec(ctx,
		"INSERT INTO orders "+
			"(user_id, order_id, debet, order_status, accrual) "+
			"VALUES ($1, $2, FALSE, $3, $4);",
		userID, bw.OrderID, models.OrderStatusWithDrawn, -100*bw.Sum)

	logger.Debug().Err(err).Msg("exit with error")
	return err
}

func (ops *OrderPostgresStorage) Withdrawals(ctx context.Context, userID int32) ([]*models.Withdrawals, error) {
	logger := log.With().Str("package", "postgres").Str("func", "Withdrawals").Logger()

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	result := make([]*models.Withdrawals, 0)

	rows, err := ops.db.Query(ctx,
		"SELECT order_id, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE (debet IS FALSE) AND (user_id = $1) AND (order_status = $2) "+
			"ORDER BY uploaded_at ASC;", userID, models.OrderStatusWithDrawn)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
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
		item.Sum = float32(sum*-1) / 100
		item.ProcessedAt = processedAt.Time
		result = append(result, &item)
	}

	return result, nil
}

func (ops *OrderPostgresStorage) UpdateOrder(ctx context.Context, orderNumber string, orderStatus string, orderAccrual int32) error {
	logger := log.With().Str("package", "postgres").Str("func", "UpdateOrder").Logger()

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	logger.Debug().Str("orderStatus", orderStatus).Str("orderNumber", orderNumber).Int32("orderAccurual", orderAccrual).Msg("before query")
	comTag, err := ops.db.Exec(context.Background(),
		"UPDATE orders SET order_status = $1 , accrual = $2 WHERE order_id = $3 ;",
		orderStatus, orderAccrual, orderNumber)

	logger.Debug().Int64("Count", comTag.RowsAffected()).Msg("Rows affected")
	logger.Debug().Err(err).Msg("exit with error")
	return err
}

func (ops *OrderPostgresStorage) GetNotFinnalizedOrdersListByUserID(ctx context.Context, userID int32) ([]*models.Order, error) {
	logger := log.With().Str("package", "postgres").Str("func", "GetNotFinnalizedOrdersListByUserID").Logger()

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	result := make([]*models.Order, 0)

	rows, err := ops.db.Query(ctx,
		"SELECT user_id, order_id, order_status, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE (debet IS TRUE) AND (user_id = $1) AND order_status NOT IN ($2, $3)"+
			"ORDER BY uploaded_at ASC;", userID, models.OrderStatusProcessed, models.OrderStatusInvalid)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Order
		var accrual int32
		var uploaded pgtype.Timestamp
		err := rows.Scan(&item.UserName, &item.Number, &item.Status, &accrual, &uploaded)
		if err != nil {
			logger.Debug().Err(err)
			return nil, nil
		}
		item.Accrual = float32(accrual) / 100
		item.Uploaded = uploaded
		result = append(result, &item)
	}

	return result, nil
}

func (ops *OrderPostgresStorage) GetNotFinnalizedOrdersList(ctx context.Context) ([]*models.Order, error) {
	logger := log.With().Str("package", "postgres").Str("func", "GetNotFinnalizedOrdersList").Logger()

	logger.Debug().Msg("enter")
	defer logger.Debug().Msg("exit")

	result := make([]*models.Order, 0)

	rows, err := ops.db.Query(ctx,
		"SELECT user_id, order_id, order_status, accrual, uploaded_at "+
			"FROM orders "+
			"WHERE (debet IS TRUE) AND order_status NOT IN ($1, $2)"+
			"ORDER BY uploaded_at ASC;", models.OrderStatusProcessed, models.OrderStatusInvalid)

	if err != nil {
		logger.Debug().Err(err).Msg("exit with error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Order
		var accrual int32
		var uploaded pgtype.Timestamp
		err := rows.Scan(&item.UserName, &item.Number, &item.Status, &accrual, &uploaded)
		if err != nil {
			logger.Debug().Err(err)
			return nil, nil
		}
		item.Accrual = float32(accrual) / 100
		item.Uploaded = uploaded
		result = append(result, &item)
	}

	return result, nil
}
