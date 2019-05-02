// Copyright (c) 2019, Nexbit di Paolo Furini.
// You may use, distribute and modify this code under the
// terms of the MIT license.
// You should have received a copy of the MIT license with
// this file.

package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	HTML     bool   `help:"(optional) Use html instead of markdown in the message"`
	Success  bool   `help:"(optional) Predefined success icon (overrides -icon argument)"`
	Warning  bool   `help:"(optional) Predefined warning icon (overrides -icon argument)"`
	Error    bool   `help:"(optional) Predefined error icon (overrides -icon argument)"`
	Question bool   `help:"(optional) Predefined question mark icon (overrides -icon argument)"`
	Silent   bool   `help:"(optional) Send message in silent mode (no user notification on the client)"`
	Log      bool   `help:"(optional) Print the API response to stdout on success"`
}

type apiResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

const apiURL = "https://api.telegram.org"
const sendMessageRes = "/bot%s/sendMessage"

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

	// replace \n to actual new-line
	text := strings.ReplaceAll(m.Text, "\\n", "\n")
	icon := m.Icon
	parseMode := "markdown"
	boldStart := "*"
	boldEnd := "*"
	if m.HTML {
		parseMode = "html"
		boldStart = "<b>"
		boldEnd = "</b>"
	}

	if m.Title != "" {
		text = boldStart + m.Title + boldEnd + "\n\n" + text
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

	// make a POST request to telegram API
	data := url.Values{}
	data.Set("parse_mode", parseMode)
	data.Set("chat_id", m.User)
	data.Set("text", text)
	data.Set("disable_notification", strconv.FormatBool(m.Silent))

	destURL, _ := url.ParseRequestURI(apiURL)
	destURL.Path = fmt.Sprintf(sendMessageRes, m.Key)
	destURLStr := destURL.String()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", destURLStr, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	// Let's check if the work actually is done
	// We have seen inconsistencies even when we get 200 OK response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Couldn't parse response body. %+v", err)
	}

	if m.Log {
		log.Println("Response:", string(body))
	}

	var apiResp = new(apiResponse)
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return err
	}
	if !apiResp.Ok {
		description := apiResp.Description
		if description == "" {
			description = "Unknown error"
		}
		return fmt.Errorf("Telegram API error: " + description)
	}

	return nil
}
