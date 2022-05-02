package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/alexkopcak/gophermart/internal/order/integration"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AccuralServiceBackground(ouc order.UseCase, as *integration.AccurualService) {
	go func() {
		for {
			orders, err := ouc.GetNotFinnalizedOrdersList(context.Background())
			log.Debug().Err(err)
			for _, item := range orders {
				log.Debug().Str("package", "handlers").Str("func", "AccuralServiceHandler").Str("order", item.Number).Msg("Update order")

				err = as.UpdateData(context.Background(), item.Number)
				log.Debug().Err(err)
			}
			time.Sleep(time.Second * 10)
		}
	}()
}

func AccuralServiceHandler(ouc order.UseCase, as *integration.AccurualService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debug().Str("package", "handlers").Str("func", "accuralservice").Msg("enter")

		userID, err := getUserID(ctx)

		if err != nil {
			ctx.String(http.StatusInternalServerError, "getUserID error")
			ctx.Abort()
			return
		}

		log.Debug().Str("package", "handlers").Str("func", "AccuralServiceHandler").Str("userID", userID).Msg("Get userID")
		orders, err := ouc.GetNotFinnalizedOrdersListByUserID(ctx.Request.Context(), userID)

		if err != nil {
			ctx.String(http.StatusInternalServerError, "GetNotFinnalizedOrdersListByUserID")
			ctx.Abort()
			return
		}

		go func() {
			for _, item := range orders {
				log.Debug().Str("package", "handlers").Str("func", "AccuralServiceHandler").Str("order", item.Number).Msg("Update order")

				err = as.UpdateData(context.Background(), item.Number)
				log.Debug().Err(err)
			}
		}()

		log.Debug().Str("package", "handlers").Str("func", "accuralservice").Msg("exit")
		ctx.Next()
	}
}
