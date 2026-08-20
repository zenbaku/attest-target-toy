package budget

import "testing"

// TotalSpent charges only processed items: the capped item's cost must not
// appear in it, and invalid items are never charged.
func TestSummarizeTotalSpent(t *testing.T) {
	results := Apply([]Item{{Key: "a", Cost: 8}, {Key: "a", Cost: 5}, {Key: "b", Cost: 1}, {Key: "", Cost: 1}}, 10)
	s := Summarize(results)

	if s.TotalSpent != 9 {
		t.Errorf("got total_spent %d, want 9 — only processed items are charged", s.TotalSpent)
	}
}

// CappedCost totals the cost of capped items only; processed and invalid items
// must not contribute.
func TestSummarizeCappedCost(t *testing.T) {
	results := Apply([]Item{{Key: "a", Cost: 8}, {Key: "a", Cost: 5}, {Key: "a", Cost: 4}, {Key: "b", Cost: -1}}, 10)
	s := Summarize(results)

	if s.CappedCost != 9 {
		t.Errorf("got capped_cost %d, want 9 — capped items cost 5 and 4", s.CappedCost)
	}
	if s.TotalSpent != 8 {
		t.Errorf("got total_spent %d, want 8 — capped costs must not be charged", s.TotalSpent)
	}
}

// Mirrors testdata/items.json at a budget of 10, the sample the run summary
// line is specified against: total_spent=29 and capped_cost=14.
func TestSummarizeSampleInput(t *testing.T) {
	results := Apply([]Item{
		{Key: "alpha", Cost: 8},
		{Key: "alpha", Cost: 5},
		{Key: "alpha", Cost: 2},
		{Key: "beta", Cost: 9},
		{Key: "beta", Cost: 9},
		{Key: "gamma", Cost: 10},
		{Key: "", Cost: 1},
		{Key: "gamma", Cost: -3},
	}, 10)
	s := Summarize(results)

	if s.TotalSpent != 29 {
		t.Errorf("got total_spent %d, want 29", s.TotalSpent)
	}
	if s.CappedCost != 14 {
		t.Errorf("got capped_cost %d, want 14", s.CappedCost)
	}
	if s.Capped != 2 {
		t.Errorf("got %d capped, want 2", s.Capped)
	}
}
