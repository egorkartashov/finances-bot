package expenses

import (
	"sync"
	"time"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
)

type InMemStorage struct {
	mutex       *sync.Mutex
	expensesMap map[int64][]expenses.Expense
}

func NewInMem() *InMemStorage {
	return &InMemStorage{
		mutex:       new(sync.Mutex),
		expensesMap: make(map[int64][]expenses.Expense),
	}
}

func (s *InMemStorage) AddExpense(userID int64, exp expenses.Expense) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.expensesMap[userID] = append(s.expensesMap[userID], exp)
}

func (s *InMemStorage) GetCategoriesTotals(userID int64, minTime time.Time) []expenses.CategoryTotal {
	s.mutex.Lock()
	expensesCopy, ok := s.expensesMap[userID]
	s.mutex.Unlock()

	if !ok {
		return nil
	}

	totalCounters := make(map[string]int32)
	for _, expense := range expensesCopy {
		if expense.Date.Before(minTime) {
			continue
		}
		if _, ok := totalCounters[expense.Category]; !ok {
			totalCounters[expense.Category] = expense.SumRub
		} else {
			totalCounters[expense.Category] += expense.SumRub
		}
	}

	categoriesTotals := make([]expenses.CategoryTotal, len(totalCounters))
	i := 0
	for cat, total := range totalCounters {
		categoriesTotals[i] = expenses.CategoryTotal{
			Category:    cat,
			TotalSumRub: total,
		}
		i += 1
	}

	return categoriesTotals
}
