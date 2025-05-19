package appcontext

import "finance-bot/service"

type AppContext struct {
	TransactionService service.TransactionService
	UserService        service.UserService
	UserStateStore     map[int64]string
}
