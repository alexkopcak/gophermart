package handlers

import (
	"net/http"

	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/alexkopcak/gophermart/internal/order/integration"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

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

				as.UpdateData(ctx, item.Number)
			}
		}()

		log.Debug().Str("package", "handlers").Str("func", "accuralservice").Msg("exit")
		ctx.Next()
	}
}
