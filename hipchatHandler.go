package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	var out, stderr bytes.Buffer
	var err error
	switch r.Method {
	case "POST":
		tokenArgs := []string{"--token", handler.tokenService.Token}
		args := parseMessage(getMessage(w, r), tokenArgs)
		pipeIndex := SliceIndex(len(args), func(i int) bool { return args[i] == "|" })
		if pipeIndex == -1 {
			out, stderr, err = standardCommand(args)
		} else {
			out, stderr, err = pipedCommand(args, pipeIndex)
		}
		if err != nil {
			sendMessage(stderr.Bytes(), w)
			return
		}
		resp := prepareResponse(out)
		respBody, err := json.Marshal(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendMessage(respBody, w)
	default:
		sendMessage([]byte("Bad Message"), w)
		http.Error(w, fmt.Sprintf("Method %s not supported", r.Method), 404)
	}
}

func pipedCommand(args []string, pipeIndex int) (bytes.Buffer, bytes.Buffer, error) {
	var out, stderr2 bytes.Buffer
	c1 := exec.Command("oc", args[0:pipeIndex]...)
	c2 := exec.Command(args[pipeIndex+1], args[pipeIndex+2:]...)
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r
	c2.Stderr = &stderr2
	c2.Stdout = &out
	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	err := c2.Wait()
	if err != nil {
		return out, stderr2, err
	}
	return out, stderr2, nil
}

func standardCommand(args []string) (bytes.Buffer, bytes.Buffer, error) {
	cmd := exec.Command("oc", args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	go func() {
		time.Sleep(3000)
		cmd.Process.Kill()
	}()
	return out, stderr, err
}
func sendMessage(body []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
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

//SliceIndex get the index of an element in a slice
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
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
