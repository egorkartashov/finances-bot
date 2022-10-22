package expenses

type Usecase struct {
	cfg               cfg
	tx                tx
	expenseStorage    expenseStorage
	currencyConverter currencyConverter
	userUc            userUc
	limitUc           limitUc
}

func NewUsecase(
	cfg cfg, tx tx, expenseStorage expenseStorage, userUc userUc, currencyProvider currencyConverter, limitUc limitUc,
) *Usecase {
	return &Usecase{
		cfg:               cfg,
		expenseStorage:    expenseStorage,
		userUc:            userUc,
		currencyConverter: currencyProvider,
		limitUc:           limitUc,
		tx:                tx,
	}
}
