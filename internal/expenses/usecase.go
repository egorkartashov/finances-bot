package expenses

type Usecase struct {
	cfg               cfg
	tx                tx
	expenseStorage    expenseStorage
	currencyConverter currencyConverter
	userStorage       userStorage
	limitUc           limitUc
}

func NewUsecase(
	cfg cfg, tx tx, expenseStorage expenseStorage, userUc userStorage, currencyConverter currencyConverter,
	limitUc limitUc,
) *Usecase {
	return &Usecase{
		cfg:               cfg,
		expenseStorage:    expenseStorage,
		userStorage:       userUc,
		currencyConverter: currencyConverter,
		limitUc:           limitUc,
		tx:                tx,
	}
}
