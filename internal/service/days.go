package service

import (
	"github.com/Diana172710/-_-_1/mint/internal/models"
	"sort"
)

// DayTotal is a daily subtotal calculated from the current user's filtered rows.
type DayTotal struct {
	Date   string
	Amount int64
	Count  int
}

func DailyTotals(entries []models.Expense) []DayTotal {
	byDate := make(map[string]DayTotal)
	for _, entry := range entries {
		day := byDate[entry.Date]
		day.Date = entry.Date
		day.Amount += entry.Amount
		day.Count++
		byDate[entry.Date] = day
	}
	result := make([]DayTotal, 0, len(byDate))
	for _, day := range byDate {
		result = append(result, day)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Date > result[j].Date })
	return result
}
