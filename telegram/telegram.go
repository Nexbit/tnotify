// Copyright (c) 2019, Nexbit di Paolo Furini.
// You may use, distribute and modify this code under the
// terms of the MIT license.
// You should have received a copy of the MIT license with
// this file.

package telegram

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Telegram is the struct that receives command line arguments.
type Telegram struct {
	User     string `help:"Recipient User or Channel ID"`
	Key      string `help:"API Key of your Telegram bot"`
	Text     string `help:"Text of the message"`
	Icon     string `help:"(optional) Icon before title or message text (UTF code)"`
	Title    string `help:"(optional) Title displayed in bold between the icon (if provided) and the message text"`
	Html     bool   `help:"(optional) Use html instead of markdown in the message"`
	Success  bool   `help:"(optional) Predefined success icon (overrides -icon argument)"`
	Warning  bool   `help:"(optional) Predefined warning icon (overrides -icon argument)"`
	Error    bool   `help:"(optional) Predefined error icon (overrides -icon argument)"`
	Question bool   `help:"(optional) Predefined question mark icon (overrides -icon argument)"`
}

const sendMessageURL = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&parse_mode=%s&text=%s"

// NewTelegram returns a new Telegram struct.
func NewTelegram() *Telegram { return &Telegram{} }

// Run is called by commandeer package after parsing the command line.
func (m *Telegram) Run() error {
	if m.User == "" {
		return fmt.Errorf("-user is mandatory (user or channel ID)")
	}
	if m.Key == "" {
		return fmt.Errorf("-key is mandatory")
	}
	if m.Text == "" {
		return fmt.Errorf("-text is mandatory")
	}

	// replace \n to %0A (LF encoded)
	text := strings.ReplaceAll(m.Text, "\\n", "%0A")
	icon := m.Icon
	parseMode := "markdown"
	boldStart := "*"
	boldEnd := "*"
	if m.Html {
		parseMode = "html"
		boldStart = "<b>"
		boldEnd = "</b>"
	}

	if m.Title != "" {
		text = boldStart + m.Title + boldEnd + "%0A%0A" + text
	}
	if m.Success {
		icon = "2705"
	}
	if m.Warning {
		icon = "26A0"
	}
	if m.Error {
		icon = "1F6A8"
	}
	if m.Question {
		icon = "2753"
	}
	if icon != "" {
		code, err := strconv.ParseInt(icon, 16, 32)
		if err == nil {
			text = string(code) + " " + text
		} else {
			fmt.Fprintf(os.Stderr, "Error parsing UTF code %q to int - sending message without icon\n", icon)
		}
	}

	resp, err := http.Get(fmt.Sprintf(sendMessageURL, m.Key, m.User, parseMode, text))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
