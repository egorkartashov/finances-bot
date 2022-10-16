package storage

import (
	"sync"
	"time"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
)

type Expenses struct {
	mutex       *sync.Mutex
	expensesMap map[int64][]expenses.Expense
}

func NewExpenses() *Expenses {
	return &Expenses{
		mutex:       new(sync.Mutex),
		expensesMap: make(map[int64][]expenses.Expense),
	}
}

func (s *Expenses) AddExpense(userID int64, exp expenses.Expense) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.expensesMap[userID] = append(s.expensesMap[userID], exp)
	return nil
}

func (s *Expenses) GetExpenses(userID int64, minTime time.Time) ([]expenses.Expense, error) {
	s.mutex.Lock()
	expensesCopy, ok := s.expensesMap[userID]
	s.mutex.Unlock()

	if !ok {
		return nil, nil
	}

	filteredExp := make([]expenses.Expense, 0, len(expensesCopy))
	for _, e := range expensesCopy {
		if e.Date.Before(minTime) {
			continue
		}
		filteredExp = append(filteredExp, e)
	}
	return filteredExp, nil
}
