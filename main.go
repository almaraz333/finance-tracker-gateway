package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	pb "github.com/almaraz333/finance-tracker-proto-files/expense"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Expense struct {
	Category  string
	Amount    float64
	CreatedAt time.Time
	Id        int32
}

const clientKey int = iota

func delteExpense(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, ok := r.Context().Value(clientKey).(pb.ExpenseClient)

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
		fmt.Println(err.Error())
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted Expense with the id: " + id))
}

func getExpenses(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, ok := r.Context().Value(clientKey).(pb.ExpenseClient)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	res, err := c.GetExpenses(ctx, nil)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		fmt.Println(err.Error())
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

func createExpense(w http.ResponseWriter, r *http.Request) {
	var expense = Expense{}

	c, ok := r.Context().Value(clientKey).(pb.ExpenseClient)

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
		fmt.Println(err.Error())
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

func updateExpense(w http.ResponseWriter, r *http.Request) {
	var expense = Expense{}

	c, ok := r.Context().Value(clientKey).(pb.ExpenseClient)

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
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
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

func Logging(w http.ResponseWriter, r *http.Request) {
	slog.Info("", r.Method, r.URL.Path)
}

func MiddlewareChain(w http.ResponseWriter, r *http.Request, middleware ...func(w http.ResponseWriter, r *http.Request)) {
	for _, mw := range middleware {
		mw(w, r)
	}
}

func main() {
	PORT := 8080

	mux := http.NewServeMux()

	fmt.Printf("listening on port: %v \n", PORT)

	conn, err := grpc.Dial("expense-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	c := pb.NewExpenseClient(conn)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	mux.HandleFunc("POST /api/expenses", func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)

		MiddlewareChain(w, r, Logging, createExpense)
	})

	mux.HandleFunc("GET /api/expenses", func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)

		MiddlewareChain(w, r, Logging, getExpenses)
	})

	mux.HandleFunc("DELETE /api/expenses/{id}", func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)

		MiddlewareChain(w, r, Logging, delteExpense)
	})

	mux.HandleFunc("PUT /api/expenses", func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)

		MiddlewareChain(w, r, Logging, updateExpense)
	})

	if err := http.ListenAndServe("0.0.0.0:"+fmt.Sprint(PORT), mux); err != nil {
		log.Fatal(err.Error())
	}
}
