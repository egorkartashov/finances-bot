package users

import (
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"
)

type userStorage interface {
	Save(user User) error
	Get(id int64) (u User, ok bool, err error)
}

type Usecase struct {
	storage userStorage
}

func NewUsecase(storage userStorage) *Usecase {
	return &Usecase{storage}
}

func (u *Usecase) GetOrRegister(userID int64) (user User, err error) {
	user, ok, err := u.storage.Get(userID)
	if err != nil {
		return
	}

	if !ok {
		user = User{
			id:       userID,
			Currency: currency.RUB,
		}
		err = u.storage.Save(user)
		if err != nil {
			return
		}
	}

	return user, nil
}

func (u *Usecase) SetCurrency(userID int64, curr currency.Currency) error {
	user, err := u.GetOrRegister(userID)
	if err != nil {
		return err
	}

	user.Currency = curr
	return u.storage.Save(user)
}
