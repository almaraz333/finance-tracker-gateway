package handlers

import (
	"context"
	"encoding/json"
	pb "github.com/almaraz333/finance-tracker-proto-files/expense"
	"net/http"
	"time"
)

func GetExpenses(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, ok := r.Context().Value("clientKey").(pb.ExpenseClient)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	res, err := c.GetExpenses(ctx, nil)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	resBytes, jsonErr := json.Marshal(res)

	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(resBytes)

	if writeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request - failed to write response"))
		return
	}
}
