package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "github.com/almaraz333/finance-tracker-proto-files/expense"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Expenst struct {
	Category  string
	Amount    float64
	CreatedAt time.Time
}

func createExpense(w http.ResponseWriter, r *http.Request, c pb.ExpenseClient) {
	var expense = Expenst{}

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

func main() {
	PORT := 8080

	mux := http.NewServeMux()

	fmt.Printf("listening on port: %v \n", PORT)

	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	c := pb.NewExpenseClient(conn)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	mux.HandleFunc("POST /api/expenses", func(w http.ResponseWriter, r *http.Request) {
		createExpense(w, r, c)
	})

	if err := http.ListenAndServe("127.0.0.1:"+fmt.Sprint(PORT), mux); err != nil {
		log.Fatal(err.Error())
	}

}
