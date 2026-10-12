package service

import (
	"github.com/Diana172710/-_-_1/mint/internal/models"
	"testing"
)

func TestDailyTotals(t *testing.T) {
	entries := []models.Expense{{Date: "2026-01-01", Amount: 105}, {Date: "2026-01-03", Amount: 300}, {Date: "2026-01-01", Amount: 95}}
	days := DailyTotals(entries)
	if len(days) != 2 || days[0].Date != "2026-01-03" || days[1].Amount != 200 || days[1].Count != 2 {
		t.Fatalf("unexpected days: %+v", days)
	}
	if entries[0].Amount != 105 {
		t.Fatal("input changed")
	}
	if len(DailyTotals(nil)) != 0 {
		t.Fatal("nonempty report for empty input")
	}
}
