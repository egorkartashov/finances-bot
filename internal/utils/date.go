package utils

import "time"

func GetDate(t time.Time) time.Time {
	y, m, d := t.Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	return date
}

func GetTomorrow(t time.Time) time.Time {
	return GetDate(t).AddDate(0, 1, 0)
}
