package handlers

import (
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/add_expense"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/get_currency_options"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/get_report"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/remove_limit"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/set_currency"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/set_limit"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/start"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/unknown_command"
)

var NewStart = start.New
var NewSetCurrency = set_currency.New
var NewAddExpense = add_expense.New
var NewGetReport = get_report.New
var NewSetLimit = set_limit.New
var NewRemoveLimit = remove_limit.New
var NewGetCurrencyOptions = get_currency_options.New
var NewUnknownCommand = unknown_command.New
