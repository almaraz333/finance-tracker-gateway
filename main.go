package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func delteExpense(w http.ResponseWriter, r *http.Request, c pb.ExpenseClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

func getExpenses(w http.ResponseWriter, _ *http.Request, c pb.ExpenseClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.GetExpenses(ctx, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	resBytes, jsonErr := json.Marshal(res)

	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(resBytes)

	if writeErr != nil {
		log.Printf("Failed to write response: %v\n", writeErr)
	}
}

func createExpense(w http.ResponseWriter, r *http.Request, c pb.ExpenseClient) {
	var expense = Expense{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&expense)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		fmt.Println(err.Error())
		return
	}

	res, err := c.CreateExpense(ctx, &pb.CreateExenseRequest{
		Category:  expense.Category,
		CreatedAt: time.Now().UTC().String(),
		Amount:    expense.Amount,
	})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Created Expense with the amount: %v", res.GetAmount())
}

func updateExpense(w http.ResponseWriter, r *http.Request, c pb.ExpenseClient) {
	var expense = Expense{}

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
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Updated Expense with id: %v", res.Id)
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
		createExpense(w, r, c)
	})

	mux.HandleFunc("GET /api/expenses", func(w http.ResponseWriter, r *http.Request) {
		getExpenses(w, r, c)
	})

	mux.HandleFunc("DELETE /api/expenses/{id}", func(w http.ResponseWriter, r *http.Request) {
		delteExpense(w, r, c)
	})

	mux.HandleFunc("PUT /api/expenses", func(w http.ResponseWriter, r *http.Request) {
		updateExpense(w, r, c)
	})

	if err := http.ListenAndServe("0.0.0.0:"+fmt.Sprint(PORT), mux); err != nil {
		log.Fatal(err.Error())
	}

}
