package karmabot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SlackConnResp struct {
	OK  bool
	URL string
}

func GetConn(logger *zap.SugaredLogger) (*websocket.Conn, error) {
	client := &http.Client{}
	var conn *websocket.Conn
	token := fmt.Sprintf("Bearer %s", os.Getenv("SLACK_APP_TOKEN"))

	req, _ := http.NewRequest(http.MethodPost, SlackConnURL, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return conn, fmt.Errorf("error getting response: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return conn, fmt.Errorf("error reading body: %w", err)
	}

	logger.Debugf("connection URL response: %v", string(body))

	slackConnResp := SlackConnResp{}
	err = json.Unmarshal(body, &slackConnResp)
	if err != nil {
		return conn, fmt.Errorf("error unmarshaling slack connection response: %w", err)
	}

	conn, _, err = websocket.DefaultDialer.Dial(slackConnResp.URL, nil)
	if err != nil {
		return conn, fmt.Errorf("error obtaining connection: %w", err)
	}

	logger.Debug("connection obtained")
	return conn, nil
}

func WriteAcknowledgement(conn *websocket.Conn, logger *zap.SugaredLogger, messageType int, p []byte) error {
	e := &Envelope{}
	err := json.Unmarshal(p, e)

	payload, err := json.Marshal(e)
	if err != nil {
		fmt.Errorf("error marshaling acknowledgement payload: %w", err)
	}
	err = conn.WriteMessage(messageType, payload)
	if err != nil {
		fmt.Errorf("error sending message acknowledgement: %w", err)
	}
	logger.Debug("message acknowledged")
	return nil
}
