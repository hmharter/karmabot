package karmabot

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type Flag struct {
	Key   string
	Value string
}

func ProcessSlashCommand(logger *zap.SugaredLogger, sce SCEnvelope) (string, error) {
	var err error
	p := sce.Payload
	message := ""
	userID := ParseUserID(logger, p.Text)
	flags := ParseSlashFlags(p.Text)
	logger.Debugf("UserID: %s, Flag 1 key: %s, Flag 1 value: %s", userID, flags[0].Key, flags[0].Value)

	// This is totally overkill because there's only one flag right now, but you never know.
	for _, flag := range flags {
		switch flag.Key {
		case Name:
			userKarma, ok := canSetName(p, userID)
			if !ok {
				message += "You can only change the real name for yourself or custom users that already have karma.\n"
				continue
			}
			origName := userKarma.RealName
			userKarma.RealName = flag.Value
			uk, _ := json.Marshal(userKarma)
			logger.Debugf("Writing Karma: %s", string(uk))
			err = WriteKarma()
			if err != nil {
				message += "There was an error writing karma.\n"
				continue
			}
			message = fmt.Sprintf("Thanks! %s will henceforth be known as %s.", origName, userKarma.RealName)
		}
	}
	return message, nil
}

// You can only set the name for yourself, or a custom user who already exists.
// We can identify custom users because they have no email address.
func canSetName(p SCPayload, userID string) (*UserKarma, bool) {
	allowed := false
	custom := false
	self := (p.UserID == userID)
	userKarma, ok := Karma[userID]
	if ok {
		custom = (userKarma.Email == "")
	}

	if ok && (self || custom) {
		allowed = true
	}
	return userKarma, allowed
}
