package main

import (
	"fmt"
	"os"

	"github.com/nexbit/tnotify/telegram"

	"github.com/jaffee/commandeer"
)

func main() {
	err := commandeer.Run(telegram.NewTelegram())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
