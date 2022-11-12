package websocket

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

func OpenAppsyncSubscription(t *testing.T, url string, query string) (func(), chan []byte, error) {

	var err error
	var req *http.Request

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		t.Error(err)
		return nil, nil, err
	}

	req.Header.Set("Sec-WebSocket-Protocol", "graphql-ws")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "SGVsbG8sIHdvcmxkIQ==")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")

	var resp *http.Response

	if resp, err = http.DefaultClient.Do(req); err != nil {
		t.Error(err)
		return nil, nil, err
	}

	if resp.StatusCode != 101 {
		t.Error(resp)
		return nil, nil, err
	}

	var conn *websocket.Conn
	var messages = make(chan []byte)

	if conn, _, err = websocket.DefaultDialer.Dial(url, nil); err != nil {
		t.Error(err)
		return nil, nil, err
	}

	var teardown = func() {
		conn.Close()
		close(messages)
	}

	var teardownWithError = func(err error) error {
		t.Error(err)
		teardown()
		return err
	}

	// make a channel that reads all messages from the websocket

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					t.Error(err)
				}
				return
			}
			messages <- message
		}
	}()

	var payload []byte

	if payload, err = json.Marshal(map[string]interface{}{
		"type": "connection_init",
	}); err != nil {
		return nil, nil, teardownWithError(err)
	}

	if err = conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return nil, nil, teardownWithError(err)
	}

	payload = <-messages

	var msg map[string]interface{}

	if err = json.Unmarshal(payload, &msg); err != nil {
		return nil, nil, teardownWithError(err)
	}

	if msg["type"] != "connection_ack" {
		t.Fatal(msg)
	}

	if payload, err = json.Marshal(map[string]interface{}{
		"type": "start",
		"payload": map[string]interface{}{
			"query": query,
		},
	}); err != nil {
		return nil, nil, teardownWithError(err)
	}

	if err = conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return nil, nil, teardownWithError(err)
	}

	payload = <-messages

	if err = json.Unmarshal(payload, &msg); err != nil {
		return nil, nil, teardownWithError(err)
	}

	if msg["type"] != "data" {
		t.Fatal(msg)
	}

	return teardown, messages, nil
}
