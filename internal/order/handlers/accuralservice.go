package handlers

import (
	"fmt"
	"net/http"

	"github.com/alexkopcak/gophermart/internal/order"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AccuralServiceHandler(ouc order.UseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debug().Str("package", "handlers").Str("func", "accuralservice").Msg("enter")

		userID, err := getUserID(ctx)

		if err != nil {
			return
		}

		orders, err := ouc.GetNotFinnalizedOrdersListByUserID(ctx.Request.Context(), userID)

		fmt.Printf("!!!\n%v\n!!!", orders)

		if err != nil {
			ctx.String(http.StatusInternalServerError, "GetNotFinnalizedOrdersListByUserID")
			return
		}

		go func() {
			for _, item := range orders {
				ouc.UpdateOrder(ctx, item)
			}
		}()

		log.Debug().Str("package", "handlers").Str("func", "accuralservice").Msg("exit")
		ctx.Next()
	}
}
