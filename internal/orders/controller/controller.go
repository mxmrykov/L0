package controller

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mxmrykov/L0/internal/caches"
	"net/http"
)

type OrderCacheHandler struct {
	oc *caches.OrderCache
}

func NewOrderController(oc *caches.OrderCache) *OrderCacheHandler {
	return &OrderCacheHandler{
		oc: oc,
	}
}
func (och *OrderCacheHandler) GetOrderController(c echo.Context) error {

	order := och.oc.GetOrderByUid(c.Param("order"))
	or, err := json.MarshalIndent(order, "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling respond: %v", err)
	}
	return c.JSONBlob(http.StatusOK, or)
}

func (och *OrderCacheHandler) GetAllOrders(c echo.Context) error {
	order := och.oc.GetOrders()

	or, err := json.MarshalIndent(order, "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling respond: %v", err)
	}

	return c.JSONBlob(http.StatusOK, or)
}
