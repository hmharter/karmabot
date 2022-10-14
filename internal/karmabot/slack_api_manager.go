package karmabot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"go.uber.org/zap"
)

const (
	UserProfileURL = "https://slack.com/api/users.profile.get"
)

type UserProfileWrapper struct {
	OK      bool        `json:"ok"`
	Profile UserProfile `json:"profile"`
	Error   string      `json:"error"`
}

type UserProfile struct {
	RealName string `json:"real_name"`
	Email    string `json:"email"`
}

type ChatMessage struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type EphemeralChatMessage struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`
}

func GetUserInfo(userID string) (*UserKarma, bool, error) {
	var userInfo *UserKarma
	customUser := false

	// If no existing user was found, prepare to get user info from Slack
	userProfile := &UserProfileWrapper{}
	newUserInfo := &UserKarma{
		UserID: userID,
	}

	data := url.Values{}
	data.Add("user", userID)

	token := fmt.Sprintf("Bearer %s", os.Getenv("SLACK_BOT_TOKEN"))

	req, err := http.NewRequest("POST", UserProfileURL, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Authorization", token)
	if err != nil {
		return userInfo, customUser, fmt.Errorf("error creating http request object: %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return userInfo, customUser, fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return userInfo, customUser, fmt.Errorf("response error: %v", resp.Status)
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return userInfo, customUser, fmt.Errorf("error reading response body: %w", err)
	}

	err = json.Unmarshal(r, &userProfile)
	if err != nil {
		return userInfo, customUser, fmt.Errorf("error unmarshalling user profile wrapper: %w", err)
	}
	if !userProfile.OK {
		// This is commented out because it was not working as intended!
		// We need to refine how we detect a custom user, possibly by requiring custom users to be configured
		// via slash command before you can assign karma to them.

		//if userProfile.Error == "user_not_found" {
		//	newUserInfo = &UserKarma{
		//		UserID: userID,
		//	}
		//	customUser = true
		//} else {
		return userInfo, customUser, fmt.Errorf("error getting user: %v", userProfile.Error)
		//}
	}

	newUserInfo.Email = userProfile.Profile.Email
	newUserInfo.RealName = userProfile.Profile.RealName

	Karma[userID] = newUserInfo

	return newUserInfo, customUser, nil
}

func PostChatMessage(logger *zap.SugaredLogger, text string, channel string) error {
	postURL := SlackPostMessageURL
	token := fmt.Sprintf("Bearer %s", os.Getenv("SLACK_BOT_TOKEN"))

	message := ChatMessage{
		Token:   token,
		Channel: channel,
		Text:    text,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling json post body: %w", err)
	}
	logger.Debugf("Post Body: %v", string(body))

	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func PostEphemeralChatMessage(logger *zap.SugaredLogger, text string, channel string, user string) error {
	postURL := SlackPostEphemeralMessageURL
	token := fmt.Sprintf("Bearer %s", os.Getenv("SLACK_BOT_TOKEN"))

	message := EphemeralChatMessage{
		Token:   token,
		Channel: channel,
		Text:    text,
		User:    user,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling json post body: %w", err)
	}
	logger.Debugf("Post Body: %v", string(body))

	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
