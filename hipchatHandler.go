package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var (
	sampleData = `{
		event: 'room_message',
		item: {
			message: {
				date: '2015-01-20T22:45:06.662545+00:00',
				from: {
					id: 1661743,
					mention_name: 'Blinky',
					name: 'Blinky the Three Eyed Fish'
				},
				id: '00a3eb7f-fac5-496a-8d64-a9050c712ca1',
				mentions: [],
				message: '/oc get pods',
				type: 'message'
			},
			room: {
				id: 1147567,
				name: 'The Weather Channel'
			}
		},
		webhook_id: 578829
	}`
)

//HipchatHandler Handle hipchat POST messages from slash command
type HipchatHandler struct {
	tokenService *TokenService
}

func (handler HipchatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		tokenArgs := []string{"--token", handler.tokenService.Token}
		args := parseMessage(getMessage(w, r), tokenArgs)
		cmd := exec.Command("oc", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		go func() {
			time.Sleep(3000)
			cmd.Process.Kill()
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := prepareResponse(out)
		respBody, err := json.Marshal(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
	default:
		http.Error(w, fmt.Sprintf("Method %s not supported", r.Method), 404)
	}
}

func getMessage(rw http.ResponseWriter, req *http.Request) string {
	decoder := json.NewDecoder(req.Body)
	var t jsonRequest
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	return t.Item.Message.Message
}

func parseMessage(message string, tokenArgs []string) []string {
	ocArgs := strings.Replace(message, "/oc ", "", -1)
	args := append(tokenArgs, strings.Split(ocArgs, " ")...)
	return args
}

func prepareResponse(out bytes.Buffer) stringResponse {
	var buffer bytes.Buffer
	buffer.WriteString("<pre>")
	buffer.WriteString(strings.Replace(out.String(), "\n", "<br>", -1))
	buffer.WriteString("</pre>")
	return stringResponse{Message: buffer.String(), MessageFormat: "html", Color: "green"}
}

type jsonRequest struct {
	Event string `json:"event"`
	Item  item   `json:"item"`
}

type item struct {
	Message message `json:"message"`
}

type message struct {
	Message string `json:"message"`
}

type stringResponse struct {
	Message       string `json:"message"`
	MessageFormat string `json:"message_format"`
	Color         string `json:"color"`
}
