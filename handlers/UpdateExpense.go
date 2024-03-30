package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/almaraz333/finance-tracker-gateway/types"
	pb "github.com/almaraz333/finance-tracker-proto-files/expense"
)

func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	var expense = types.Expense{}

	c, ok := r.Context().Value("clientKey").(pb.ExpenseClient)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&expense)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	res, err := c.UpdateExpense(ctx, &pb.UpdateExpenseRequest{
		Category: expense.Category,
		Amount:   expense.Amount,
		Id:       expense.Id,
	})

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
