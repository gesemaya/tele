# fsm-telebot

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vitaliy-ukiru/fsm-telebot?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/gesemaya/sniper/pkg/fsm-telebot.svg)](https://pkg.go.dev/github.com/gesemaya/sniper/pkg/fsm-telebot)
[![Go](https://github.com/gesemaya/sniper/pkg/fsm-telebot/actions/workflows/go.yml/badge.svg?branch=master&style=flat-square)](https://github.com/gesemaya/sniper/pkg/fsm-telebot/actions/workflows/go.yml)
[![golangci-lint](https://github.com/gesemaya/sniper/pkg/fsm-telebot/actions/workflows/golangci-lint.yml/badge.svg?branch=master)](https://github.com/gesemaya/sniper/pkg/fsm-telebot/actions/workflows/golangci-lint.yml)

Finite State Machine for [telebot](https://github.com/gesemaya/sniper). 
Based on [aiogram](https://github.com/aiogram/aiogram) FSM version.

It not a full implementation FSM. It just states manager for telegram bots.

## Install:
```
go get -u github.com/gesemaya/sniper/pkg/fsm-telebot@v1.2.0
```


## Examples
<details>
<summary>simple configuration</summary>

```go
package main

import (
	"os"
	"time"

	"github.com/gesemaya/sniper/pkg/fsm-telebot"
	"github.com/gesemaya/sniper/pkg/fsm-telebot/storages/memory"
	tele "github.com/gesemaya/sniper/pkg/telebot"
)

func main() {
	bot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 3 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	// for example using memory storage
	// but prefer will use redis or file storage.
	storage := memory.NewStorage()
	manager := fsm.NewManager(
		bot,     // tele.Bot
		nil,     // handlers will setups to this group. Default: creates new
		storage, // storage for states and data
		nil,     // context maker. Default: NewFSMContext
	)
	manager.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.Context) error {
		userState, err := state.State()
		if err != nil {
			return c.Send("error: " + err.Error())
		}

		return c.Send(userState.GoString())
	})

}

```

</details>

Many complex examples in directory [examples](./examples).

