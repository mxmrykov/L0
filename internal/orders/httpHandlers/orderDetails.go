package httpHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/mxmrykov/L0/internal/models"
	"net/http"
)

func OrderHandler(order *models.Order, w http.ResponseWriter, r *http.Request) {

	response, err := json.MarshalIndent(order, "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling: %v", err)
	}

	_, err = w.Write(response)

	if err != nil {
		fmt.Printf("Error at reponding order: %v", err)
	}
}
