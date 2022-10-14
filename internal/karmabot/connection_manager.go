package karmabot

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func OpenConnection(c chan int, logger *zap.SugaredLogger) {
	logger.Debug("get connection")
	conn, err := GetConn(logger)
	if err != nil {
		logger.Error(fmt.Errorf("error getting connection: %w", err))
	}
	defer conn.Close()

	if conn != nil {
		logger.Debug("connection obtained")

		for {
			// Get incoming message from Slack
			messageType, p, readErr := conn.ReadMessage()
			if readErr != nil {
				if websocket.IsUnexpectedCloseError(readErr, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Error(fmt.Errorf("unexpected close error: %w", readErr))
					break
				}

				logger.Error(fmt.Errorf("error reading message: %w", readErr))
				continue
			}
			logger.Debugf("type: %d, msg: %s", messageType, p)

			// Unmarshal message into Envelope. If the message is a disconnect notice, break.
			e := Envelope{}
			err = json.Unmarshal(p, &e)
			if err != nil {
				logger.Error(fmt.Errorf("error unmarshalling into envelope: %w", err))
				continue
			}
			if e.Type == Disconnect {
				logger.Error(fmt.Errorf("closing connection: %v", e.Reason))
				break
			}

			// Send message back as acknowledgement
			err = WriteAcknowledgement(conn, logger, messageType, p)
			if err != nil {
				logger.Error(fmt.Errorf("error processing message received: %w", err))
				continue
			}

			// Send envelope to processor to get response text
			err = Process(logger, p)
			if err != nil {
				logger.Error(fmt.Errorf("error processing message received: %w", err))
				continue
			}
		}
	}

	// When we break out of the for loop, the function ends and the connection is closed.
	// Put an int back on the channel to allow a new connection to be created.
	<-c
}
