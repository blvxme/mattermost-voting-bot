package votingbot

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/sirupsen/logrus"
)

type VotingBot struct {
	config   *votingBotConfig
	client   *model.Client4
	wsClient *model.WebSocketClient
	user     *model.User
	team     *model.Team
	channel  *model.Channel
}

func NewBot() (bot *VotingBot) {
	config, err := getConfig()
	if err != nil {
		logrus.Panicf("Failed to get configuration: %+v\n", err)
		return
	}
	logrus.Debugf("Using configuration: %+v\n", config)

	client := model.NewAPIv4Client(config.server.String())
	client.SetToken(config.token)

	user, resp, err := client.GetUser("me", "")
	if err != nil {
		logrus.Panicf("Failed to get user: %+v\n", err)
		return
	}
	logrus.Debugf("User: %+v, reponse: %+v\n", user, resp)

	team, resp, err := client.GetTeamByName(config.team, "")
	if err != nil {
		logrus.Panicf("Failed to get team: %+v\n", err)
		return
	}
	logrus.Debugf("Team: %+v, reponse: %+v\n", team, resp)

	bot = &VotingBot{
		config:   config,
		client:   client,
		wsClient: nil,
		user:     user,
		team:     team,
	}

	return
}
