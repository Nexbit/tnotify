// Copyright (c) 2019, Nexbit di Paolo Furini.
// You may use, distribute and modify this code under the
// terms of the MIT license.
// You should have received a copy of the MIT license with
// this file.

package telegram

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Telegram is the struct that receives command line arguments.
type Telegram struct {
	User  string `help:"Recipient User or Channel ID"`
	Key   string `help:"API Key of your Telegram bot"`
	Text  string `help:"Text of the message"`
	Icon  string `help:"(optional) Icon before title or message text (UTF code)"`
	Title string `help:"(optional) Title displayed in bold between the icon (if provided) and the message text"`
}

const sendMessageURL = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&parse_mode=markdown&text=%s"

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

	if m.Title != "" {
		text = "*" + m.Title + "*%0A%0A" + text
	}
	if m.Icon != "" {
		if code, err := strconv.ParseInt(m.Icon, 16, 32); err == nil {
			text = string(code) + " " + text
		}
	}

	resp, err := http.Get(fmt.Sprintf(sendMessageURL, m.Key, m.User, text))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
