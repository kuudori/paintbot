package validators

import (
	"PaintBackend/internal/config"
	"encoding/json"
	"errors"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"net/url"
)

func ValidateWebAppData(auth string) (*gotgbot.User, error) {
	cfg := config.GetConfig()
	authQuery, err := url.ParseQuery(auth)
	if err != nil {
		return nil, errors.New("failed to parse auth query")
	}

	ok, err := ext.ValidateWebAppQuery(authQuery, cfg.BotToken)
	if err != nil {
		return nil, errors.New("failed to validate data")
	}
	if !ok {
		return nil, errors.New("untrusted data")
	}

	var u gotgbot.User
	err = json.Unmarshal([]byte(authQuery.Get("user")), &u)
	if err != nil {
		return nil, errors.New("failed to get user")
	}

	return &u, nil
}
