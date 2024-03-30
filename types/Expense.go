package types

import "time"

type Expense struct {
	Category  string
	Amount    float64
	CreatedAt time.Time
	Id        int32
}
