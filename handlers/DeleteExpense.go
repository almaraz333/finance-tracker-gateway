package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/almaraz333/finance-tracker-proto-files/expense"
)

func DelteExpense(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, ok := r.Context().Value("clientKey").(pb.ExpenseClient)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	id := r.PathValue("id")
	idInt, idErr := strconv.ParseInt(id, 10, 32)

	_, err := c.DeleteExpense(ctx, &pb.DeleteExpenseRequest{
		Id: int32(idInt),
	})

	if idErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request, could not parse id"))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted Expense with the id: " + id))
}
