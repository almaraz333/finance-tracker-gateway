package handlers

import (
	"context"
	"encoding/json"
	"github.com/almaraz333/finance-tracker-gateway/types"
	pb "github.com/almaraz333/finance-tracker-proto-files/expense"
	"net/http"
	"time"
)

func CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expense = types.Expense{}

	c, ok := r.Context().Value("clientKey").(pb.ExpenseClient)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request - reading context"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&expense)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request - parsing JSON"))
		return
	}

	res, err := c.CreateExpense(ctx, &pb.CreateExenseRequest{
		Category:  expense.Category,
		CreatedAt: time.Now().UTC().String(),
		Amount:    expense.Amount,
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request - service"))
		return
	}

	resBytes, jsonErr := json.Marshal(res)

	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request - parsing JSON for res"))
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
