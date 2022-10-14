package karmabot

import (
	"fmt"

	"go.uber.org/zap"
)

func ProcessMessageEvent(logger *zap.SugaredLogger, e Event) ([]string, error) {
	logger.Debug("processing message event")
	// Check to see if the message is assigning Karma

	// We are looking for instances of @username ++++ mentions in messages.
	// In the API, user mentions are formatted like this: <@W6RT3G6Z> where
	// the user id starts with "<@" and is followed by either a capital 'W'
	// or a capital 'U' character and then a collection of upper case letters
	// and numbers and then closed by '>'.

	// We are also looking for at least one '+' (plus) or '-' (minus) character
	// after a space (or ascii non-breaking space) following the closing '>' of the user id.
	// We aren't going to limit the number of '+' (plus) or '-' (minus) characters here,
	// we will trap that further in the processing.
	var messages []string
	var err error

	callouts := ParseCallouts(logger, e.Text)
	if len(callouts) > 0 {
		// This is karma, process it
		logger.Debug("this message totally has karma!")
		messages, err = AssignKarma(logger, callouts, e.User)
		if err != nil {
			return messages, fmt.Errorf("error assigning karma: %w", err)
		}
	}
	logger.Debugf("Messages: %v", messages)
	return messages, nil
}
