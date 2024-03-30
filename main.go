package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/almaraz333/finance-tracker-gateway/handlers"

	pb "github.com/almaraz333/finance-tracker-proto-files/expense"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const clientKey string = "clientKey"

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

		MiddlewareChain(w, r, Logging, handlers.CreateExpense)
	})

	mux.HandleFunc("GET /api/expenses", func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)

		MiddlewareChain(w, r, Logging, handlers.GetExpenses)
	})

	mux.HandleFunc("DELETE /api/expenses/{id}", func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)

		MiddlewareChain(w, r, Logging, handlers.DelteExpense)
	})

	mux.HandleFunc("PUT /api/expenses", func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), clientKey, c)
		r = r.WithContext(newCtx)

		MiddlewareChain(w, r, Logging, handlers.UpdateExpense)
	})

	if err := http.ListenAndServe("0.0.0.0:"+fmt.Sprint(PORT), mux); err != nil {
		log.Fatal(err.Error())
	}
}
