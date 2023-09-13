package middleware

import (
	fsm "github.com/gesemaya/tele/pkg/fsm-telebot"
	tele "github.com/gesemaya/tele/pkg/telebot"
)

// ContextKey is key for telebot.Context storage what uses in middleware.
var ContextKey = "fsm"

// FSMContextMiddleware save FSM FSMContext in telebot.Context.
// Recommend use without manager.
func FSMContextMiddleware(storage fsm.Storage) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			c.Set(ContextKey, fsm.NewFSMContext(c, storage))
			return next(c)
		}
	}
}
