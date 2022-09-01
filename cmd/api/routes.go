package api

import (
	"os"

	"github.com/gin-gonic/gin"
)

func (server *Server) Router() error {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/api/initalize-transaction", server.InitiatePayment)

	router.GET("/api/payment-callback", server.VerifyPayment)
	if err := router.Run(os.Getenv("SERVER_ADDRESS")); err != nil {
		return err
	}
	return nil
}
