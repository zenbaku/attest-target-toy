// Package budget applies a per-key spending cap to a stream of items.
package budget

// Item is one unit of work charged against its key's budget.
type Item struct {
	Key  string `json:"key"`
	Cost int    `json:"cost"`
}

// Status is the outcome of attempting one Item.
type Status string

const (
	// StatusProcessed means the item fit within its key's remaining budget.
	StatusProcessed Status = "processed"
	// StatusCapped means the item would have exceeded its key's budget.
	StatusCapped Status = "capped"
	// StatusInvalid means the item was malformed and was not charged.
	StatusInvalid Status = "invalid"
)

// Result is the outcome of one Item, carrying the key's spend after the attempt.
type Result struct {
	Item   Item
	Status Status
	Spent  int
}

// Apply charges each item against its key's budget, in order.
//
// An item is processed when its cost fits in the key's remaining budget, and the
// key's spend increases by that cost. An item that would exceed the budget is
// capped and the key's spend is left unchanged — a capped item must not partially
// consume the budget, or a later cheaper item for the same key would be denied
// capacity it should still have.
//
// Items with an empty key or a negative cost are invalid. They are neither
// processed nor charged.
func Apply(items []Item, perKeyBudget int) []Result {
	spent := make(map[string]int, len(items))
	results := make([]Result, 0, len(items))

	for _, item := range items {
		if item.Key == "" || item.Cost < 0 {
			results = append(results, Result{Item: item, Status: StatusInvalid, Spent: spent[item.Key]})
			continue
		}

		if spent[item.Key]+item.Cost > perKeyBudget {
			results = append(results, Result{Item: item, Status: StatusCapped, Spent: spent[item.Key]})
			continue
		}

		spent[item.Key] += item.Cost
		results = append(results, Result{Item: item, Status: StatusProcessed, Spent: spent[item.Key]})
	}

	return results
}

// Summary counts results by status and totals the spend across every key.
type Summary struct {
	Keys      int
	Processed int
	Capped    int
	Invalid   int
	// TotalSpent is the sum of every key's final spend.
	TotalSpent int
}

// Summarize tallies results across all keys.
//
// Every Result carries its key's spend after the attempt, so a key's final spend
// is the one on its last result — including when that last result is capped or
// invalid, neither of which changes the spend it reports.
func Summarize(results []Result) Summary {
	s := Summary{}
	finalSpend := map[string]int{}

	for _, r := range results {
		if r.Item.Key != "" {
			finalSpend[r.Item.Key] = r.Spent
		}
		switch r.Status {
		case StatusProcessed:
			s.Processed++
		case StatusCapped:
			s.Capped++
		case StatusInvalid:
			s.Invalid++
		}
	}

	for _, spent := range finalSpend {
		s.TotalSpent += spent
	}

	s.Keys = len(finalSpend)
	return s
}
