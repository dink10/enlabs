package payments

import "time"

// Account is a user model.
type SourceType struct {
	ID    int    `json:"id" pg:",pk"`
	Value string `json:"value"`
}

// Account is a user model.
type Account struct {
	ID      int `pg:",pk"`
	Balance float64
}

// Account is a user model.
type Payment struct {
	ID            int `pg:",pk"`
	CreatedAt     time.Time
	AccountID     int
	TransactionID string
	State         string
	Amount        float64
	SourceType    int
	Processed     bool
}
