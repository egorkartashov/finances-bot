package users

import "gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/currency"

type User struct {
	id       int64
	Currency currency.Currency
}

func (u *User) GetID() int64 {
	return u.id
}
