package httpHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mxmrykov/L0/internal/caches"
	"net/http"
)

func OrderHandler(c echo.Context, os *caches.OrderCache) error {
	response, err := json.MarshalIndent(os.GetOrderByUid(c.Param("uid")), "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling: %v", err)
		return err
	}

	return c.JSONBlob(http.StatusOK, response)
}
