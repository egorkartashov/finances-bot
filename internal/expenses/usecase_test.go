package expenses_test

//
//import (
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/golang/mock/gomock"
//	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses"
//	expenses_mocks "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/expenses/mocks"
//)
//
//func TestModel_AddExpense_ShouldDelegateToStorage(t *testing.T) {
//	tests := []struct {
//		userID   int64
//		category string
//		sum      int32
//		date     time.Time
//	}{
//		{
//			userID:   1,
//			category: "продукты",
//			sum:      100,
//			date:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
//		},
//		{
//			userID:   101010,
//			category: "1232313131",
//			sum:      -371837931,
//			date:     time.Date(2000, 10, 2, 0, 0, 0, 0, time.UTC),
//		},
//	}
//
//	for _, tt := range tests {
//		tc := tt
//
//		name := fmt.Sprintf("sending expense %s %v %s by user %v", tc.category, tc.sum, tc.date, tc.userID)
//		t.Run(name, func(t *testing.T) {
//			t.Parallel()
//
//			ctrl := gomock.NewController(t)
//
//			expense := expenses.Expense{
//				Category: tc.category,
//				Sum:   tc.sum,
//				Date:     tc.date,
//			}
//
//			storageMock := expenses_mocks.NewMockStorage(ctrl)
//			storageMock.EXPECT().AddExpense(tc.userID, expense)
//
//			expensesModel := expenses.NewUsecase(storageMock)
//			expensesModel.AddExpense(tc.userID, expense)
//		})
//	}
//}
