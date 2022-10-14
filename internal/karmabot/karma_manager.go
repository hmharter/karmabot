package karmabot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"

	"go.uber.org/zap"
)

type UserKarma struct {
	UserID   string `json:"userid"`
	Karma    int    `json:"karma"`
	RealName string `json:"realname"`
	Email    string `json:"email"`
}

var Karma map[string]*UserKarma

func LoadKarma() (map[string]*UserKarma, error) {
	Karma = make(map[string]*UserKarma)
	karmafile := os.Getenv("KARMAFILE")
	if fileExists(karmafile) {
		b, err := ioutil.ReadFile(karmafile)
		if err != nil {
			return Karma, fmt.Errorf("error reading karmafile: %w", err)
		}
		err = json.Unmarshal(b, &Karma)
		if err != nil {
			return Karma, fmt.Errorf("error unmarshalling karma: %w", err)
		}
	}
	return Karma, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func AssignKarma(logger *zap.SugaredLogger, callouts []string, authorID string) ([]string, error) {
	logger.Debug("begin assigning karma")
	var messages []string
	var message string
	var err error

	// Loop through the karma callouts found in the message
	for _, callout := range callouts {
		logger.Debugf("callout: %s", callout)
		customUser := false

		// Get the UserID of the user assigned karma
		userID := ParseUserID(logger, callout)
		logger.Debugf("UserID: %v", userID)
		if len(userID) < 1 {
			return messages, fmt.Errorf("no user ID found")
		}

		// Check for existing Karma user, or get info for new user
		userKarma, ok := Karma[userID]
		if !ok {
			userKarma, customUser, err = GetUserInfo(userID)
			if err != nil {
				return messages, fmt.Errorf("error getting user info: %w", err)
			}
		}
		uk, _ := json.Marshal(userKarma)
		logger.Debugf("Got userKarma: %v", string(uk))

		// If the user assigned karma is the post author, they get no karma.
		if userID == authorID {
			logger.Debug("User ID == Author ID")
			message = fmt.Sprintf("> Nice try, %s.", userKarma.RealName)
			return messages, nil
		}

		// Get the count of karma, and bool for buzzkill
		count, buzzkill := getKarmaCount(logger, callout)
		logger.Debugf("Count: %d, Buzzkill: %v", count, buzzkill)

		// Adjust user karma by count
		userKarma.Karma += count

		// Write new karma to Karmafile
		err = WriteKarma()
		if err != nil {
			return messages, fmt.Errorf("Error writing to karmafile: %w", err)
		}

		realName := userKarma.RealName
		if len(userKarma.RealName) < 1 {
			realName = fmt.Sprintf("@%s", userKarma.UserID)
		}

		if count > 0 {
			message = fmt.Sprintf("> %s %s; their karma has increased to %d.", realName, funPositiveMessage[rand.Intn(len(funPositiveMessage))], userKarma.Karma)
		} else {
			if len(userKarma.RealName) > 1 && realName[len(realName)-1:] == "s" {
				message = fmt.Sprintf("> %s' karma has decreased to %d; %s.", realName, userKarma.Karma, funNegativeMessage[rand.Intn(len(funNegativeMessage))])
			} else {
				message = fmt.Sprintf("> %s's karma has decreased to %d; %s.", realName, userKarma.Karma, funNegativeMessage[rand.Intn(len(funNegativeMessage))])
			}
		}

		if buzzkill {
			message += " :bee: Buzzkill Mode™ has enforced a maximum change of 5 points."
		}

		if customUser { // and they don't have a real name already
			message += fmt.Sprintf("\n>\n>Hey, it looks like %s is a new custom Karma recipient! You can use the ", realName) +
				"slash command `/karmabot @username --name=\"Real Name\"` to set their display name!"
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func getKarmaCount(logger *zap.SugaredLogger, s string) (int, bool) {
	logger.Debug("getting karma count")
	karma := ParseKarma(logger, s)
	direction := Positive
	if string(karma[0]) == "-" {
		direction = Negative
	}
	count := len(karma)
	buzzkill := false
	if math.Abs(float64(count)) > MaxKarma {
		buzzkill = true
		count = MaxKarma
	}
	if direction == Negative {
		count = count * Negative
	}

	return count, buzzkill
}
