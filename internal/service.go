package internal

import (
	"github.com/fasthttp/router"
	"receipt-processor-challenge/internal/handler"
)

func NewProcessorService(handler *handler.Handler) *router.Router {
	r := router.New()

	r.GET("/receipts/{id}/points", handler.GetPoints)
	r.POST("/receipts/process", handler.ProcessReceipt)
	r.GET("/health", handler.Health)

	return r
}
