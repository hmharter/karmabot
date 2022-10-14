package karmabot

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type Envelope struct {
	EnvelopeID string `json:"envelope_id"`
	Type       string `json:"type"`
	Reason     string `json:"reason"`
	Payload    string `json:"-"`
}

type MessageEnvelope struct {
	EnvelopeID string         `json:"envelope_id"`
	Payload    MessagePayload `json:"payload"`
	Type       string         `json:"type"`
}

type MessagePayload struct {
	Token       string   `json:"token"`
	TeamID      string   `json:"team_id"`
	ApiAppID    string   `json:"api_app_id"`
	Event       Event    `json:"event"`
	Type        string   `json:"type"`
	AuthedTeams []string `json:"authed_teams"`
	EventID     string   `json:"event_id"`
	EventTime   int      `json:"event_time"`
}

type Event struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
	TS      string `json:"ts"`
}

type SCEnvelope struct {
	EnvelopeID string    `json:"envelope_id"`
	Payload    SCPayload `json:"payload"`
	Type       string    `json:"type"`
}

type SCPayload struct {
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	APIAppID    string `json:"api_app_id"`
}

var BotUser = os.Getenv("SLACK_BOT_USER")

func Process(logger *zap.SugaredLogger, p []byte) error {
	logger.Debug("start processor")

	// Unmarshal received bytes into Envelope
	e := Envelope{}
	err := json.Unmarshal(p, &e)
	if err != nil {
		return fmt.Errorf("error unmarshalling into envelope: %w", err)
	}

	// Get the event type and check it
	logger.Debugf("Envelope Type: %s", e.Type)
	var msg string
	var msgs []string
	switch e.Type {
	case EventsAPI:
		me := MessageEnvelope{}
		err = json.Unmarshal(p, &me)
		if err != nil {
			return fmt.Errorf("error unmarshalling envelope into message envelope: %w", err)
		}
		event := me.Payload.Event
		if event.Type == Message {
			logger.Debugf("Event is a message!")
			if event.User == BotUser {
				logger.Debugf("Bot user message, not processing")
				return nil
			}

			msgs, err = ProcessMessageEvent(logger, event)
			if err != nil {
				return fmt.Errorf("error processing the message: %w", err)
			}
			for _, msg = range msgs {
				logger.Debug("looping through callouts")
				if len(msg) > 0 {
					logger.Debugf("Callout: %v", msg)
					err = PostChatMessage(logger, msg, event.Channel)
					if err != nil {
						return fmt.Errorf("error posting chat message: %w", err)
					}
				}
			}
		}
	case SlashCommands:
		sce := SCEnvelope{}
		err = json.Unmarshal(p, &sce)
		if err != nil {
			return fmt.Errorf("error unmarshalling envelope into slashs command envelope: %w", err)
		}

		logger.Debug("Event is a slash command!")
		msg, err = ProcessSlashCommand(logger, sce)
		if err != nil {
			return fmt.Errorf("error processing the slash command: %w", err)
		}
		if len(msg) > 0 {
			logger.Debugf("Message: %v", msg)
			err = PostEphemeralChatMessage(logger, msg, sce.Payload.ChannelID, sce.Payload.UserID)
			if err != nil {
				return fmt.Errorf("error posting chat message: %w", err)
			}
		}
	}

	return nil
}
