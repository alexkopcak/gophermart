package order

import "errors"

var (
	ErrOrderAlreadyInsertedByUser      = errors.New("номер заказа уже был загружен этим пользователем")
	ErrOrderAcceptToHandle             = errors.New("новый номер заказа принят в обработку")
	ErrOrderBadQueryFormat             = errors.New("неверный формат запроса")
	ErrUserNotAuthtorised              = errors.New("пользователь не аутентифицирован")
	ErrOrderAlreadyInsertedByOtherUser = errors.New("номер заказа уже был загружен другим пользователем")
	ErrOrderBadFormat                  = errors.New("неверный формат номера заказа")
	ErrInternalServer                  = errors.New("внутренняя ошибка сервера")

	ErrNotEnougthBalance = errors.New("на счету недостаточно средств")
)
