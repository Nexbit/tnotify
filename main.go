// Copyright (c) 2019, Nexbit di Paolo Furini.
// You may use, distribute and modify this code under the
// terms of the MIT license.
// You should have received a copy of the MIT license with
// this file.

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
