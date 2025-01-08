package domain

// Order represents a basic order entity in the system.
type Order struct {
	ID     string   `json:"id"`
	Items  []string `json:"items"`
	Status string   `json:"status"`
}
