package budget

import "testing"

func TestApplyProcessesWithinBudget(t *testing.T) {
	results := Apply([]Item{{Key: "a", Cost: 3}, {Key: "a", Cost: 4}}, 10)

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	for i, r := range results {
		if r.Status != StatusProcessed {
			t.Errorf("result %d: got status %q, want %q", i, r.Status, StatusProcessed)
		}
	}
	if got := results[1].Spent; got != 7 {
		t.Errorf("got spent %d, want 7", got)
	}
}

func TestApplyCapsWhenBudgetExceeded(t *testing.T) {
	results := Apply([]Item{{Key: "a", Cost: 8}, {Key: "a", Cost: 5}}, 10)

	if results[0].Status != StatusProcessed {
		t.Errorf("first item: got %q, want %q", results[0].Status, StatusProcessed)
	}
	if results[1].Status != StatusCapped {
		t.Errorf("second item: got %q, want %q", results[1].Status, StatusCapped)
	}
}

// A capped item must not consume budget. If it did, the third item here would be
// denied capacity it should still have.
func TestCappedItemDoesNotConsumeBudget(t *testing.T) {
	results := Apply([]Item{{Key: "a", Cost: 8}, {Key: "a", Cost: 5}, {Key: "a", Cost: 2}}, 10)

	if got := results[1].Spent; got != 8 {
		t.Errorf("capped item: got spent %d, want 8", got)
	}
	if results[2].Status != StatusProcessed {
		t.Errorf("item after cap: got %q, want %q", results[2].Status, StatusProcessed)
	}
	if got := results[2].Spent; got != 10 {
		t.Errorf("item after cap: got spent %d, want 10", got)
	}
}

func TestBudgetIsPerKey(t *testing.T) {
	results := Apply([]Item{{Key: "a", Cost: 9}, {Key: "b", Cost: 9}}, 10)

	for i, r := range results {
		if r.Status != StatusProcessed {
			t.Errorf("result %d: got %q, want %q — budgets must not be shared across keys", i, r.Status, StatusProcessed)
		}
	}
}

func TestCostEqualToRemainingBudgetIsProcessed(t *testing.T) {
	results := Apply([]Item{{Key: "a", Cost: 10}}, 10)

	if results[0].Status != StatusProcessed {
		t.Errorf("got %q, want %q — the cap is inclusive", results[0].Status, StatusProcessed)
	}
}

func TestInvalidItems(t *testing.T) {
	results := Apply([]Item{{Key: "", Cost: 1}, {Key: "a", Cost: -1}, {Key: "a", Cost: 2}}, 10)

	if results[0].Status != StatusInvalid {
		t.Errorf("empty key: got %q, want %q", results[0].Status, StatusInvalid)
	}
	if results[1].Status != StatusInvalid {
		t.Errorf("negative cost: got %q, want %q", results[1].Status, StatusInvalid)
	}
	if got := results[2].Spent; got != 2 {
		t.Errorf("invalid items must not be charged: got spent %d, want 2", got)
	}
}

func TestSummarize(t *testing.T) {
	results := Apply([]Item{{Key: "a", Cost: 8}, {Key: "a", Cost: 5}, {Key: "b", Cost: 1}, {Key: "", Cost: 1}}, 10)
	s := Summarize(results)

	if s.Keys != 2 {
		t.Errorf("got %d keys, want 2", s.Keys)
	}
	if s.Processed != 2 {
		t.Errorf("got %d processed, want 2", s.Processed)
	}
	if s.Capped != 1 {
		t.Errorf("got %d capped, want 1", s.Capped)
	}
	if s.Invalid != 1 {
		t.Errorf("got %d invalid, want 1", s.Invalid)
	}
}
