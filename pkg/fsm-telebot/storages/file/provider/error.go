package provider

import "github.com/gesemaya/sniper/pkg/fsm-telebot/storages/file"

func newError(provider string, op string, err error) error {
	if err == nil {
		return nil
	}
	return &file.ProviderError{ProviderType: provider, Operation: op, Err: err}
}
