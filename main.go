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

func delteExpense(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, ok := r.Context().Value("context").(pb.ExpenseClient)

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

	c, ok := r.Context().Value("context").(pb.ExpenseClient)

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

	c, ok := r.Context().Value("context").(pb.ExpenseClient)

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

	res, err := c.CreateExpense(ctx, &pb.CreateExenseRequest{
		Category:  expense.Category,
		CreatedAt: time.Now().UTC().String(),
		Amount:    expense.Amount,
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

func updateExpense(w http.ResponseWriter, r *http.Request) {
	var expense = Expense{}

	c, ok := r.Context().Value("context").(pb.ExpenseClient)

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

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

const clientKey int = iota

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

	mux.HandleFunc("POST /api/expenses", Logging(func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)
		createExpense(w, r)
	}))

	mux.HandleFunc("GET /api/expenses", Logging(func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)
		getExpenses(w, r)
	}))

	mux.HandleFunc("DELETE /api/expenses/{id}", Logging(func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)
		delteExpense(w, r)
	}))

	mux.HandleFunc("PUT /api/expenses", Logging(func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)
		updateExpense(w, r)
	}))

	if err := http.ListenAndServe("0.0.0.0:"+fmt.Sprint(PORT), mux); err != nil {
		log.Fatal(err.Error())
	}
}
