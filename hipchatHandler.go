package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
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
		args := []string{"--token", handler.tokenService.Token, "get", "pods"}
		cmd := exec.Command("oc", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := stringResponse{Message: out.String()}
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

type stringResponse struct {
	Message string `json:"message"`
}
