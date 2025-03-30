package votingbot

import (
	"fmt"
	"net/url"
	"os"
)

type votingBotConfig struct {
	server *url.URL
	team   string
	user   string
	token  string
}

func getConfig() (config *votingBotConfig, err error) {
	host := os.Getenv("MATTERMOST_HOST")
	if host == "" {
		err = fmt.Errorf("MATTERMOST_HOST environment variable not set")
		return
	}

	port := os.Getenv("MATTERMOST_PORT")
	if port == "" {
		err = fmt.Errorf("MATTERMOST_PORT environment variable not set")
		return
	}

	server, err := url.Parse("http://" + host + ":" + port)
	if err != nil {
		return
	}

	team := os.Getenv("MATTERMOST_TEAM")
	if team == "" {
		err = fmt.Errorf("MATTERMOST_TEAM environment variable not set")
		return
	}

	user := os.Getenv("BOT_USER")
	if user == "" {
		err = fmt.Errorf("BOT_USER environment variable not set")
		return
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		err = fmt.Errorf("BOT_TOKEN environment variable not set")
		return
	}

	config = &votingBotConfig{
		server: server,
		team:   team,
		user:   user,
		token:  token,
	}

	return
}
