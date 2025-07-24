package routers

import (
	"github.com/gin-gonic/gin"

	"subscriptions/rest-service/internal/api/handlers"
)

func subscriptionRouter(router *gin.RouterGroup, handler handlers.SubHandler) {
	subsRouter := router.Group("/subs")
	{
		subsRouter.GET("/", handler.GetAllSubscriptions)
		subsRouter.GET("/:id", handler.GetSubscriptionByID)
		subsRouter.POST("/", handler.CreateSubscription)
		subsRouter.PUT("/:id", handler.FullUpdateSubscription)
		subsRouter.PATCH("/:id", handler.PatchUpdateSubscription)
		subsRouter.DELETE("/:id", handler.DeleteSubscription)
	}
}
