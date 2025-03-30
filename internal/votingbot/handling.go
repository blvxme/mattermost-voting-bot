package votingbot

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (bot *VotingBot) Run() {
	var err error

	for {
		bot.wsClient, err = model.NewWebSocketClient4("ws://"+bot.config.server.Host+bot.config.server.Path, bot.client.AuthToken)
		if err != nil {
			logrus.Warnf("Failed to connect to WS server: %+v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		logrus.Debugf("Connected\n")

		bot.wsClient.Listen()

		for event := range bot.wsClient.EventChannel {
			go bot.handleEvent(event)
		}
	}
}

func (bot *VotingBot) handleEvent(event *model.WebSocketEvent) {
	if event.EventType() != model.WebsocketEventPosted {
		return
	}
	logrus.Debugf("New %s event\n", event.EventType())

	post := &model.Post{}
	err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)
	if err != nil {
		logrus.Errorf("Failed to unmarshal post: %+v\n", err)
		return
	}
	logrus.Debugf("Post: %+v\n", post)

	if post.UserId == bot.user.Id {
		return
	}

	bot.handlePost(post)
}

func (bot *VotingBot) handlePost(post *model.Post) {
	if strings.HasPrefix(post.Message, "!voting-create") {
		logrus.Infof("New create command\n")
		bot.createVoting(post)
	} else if strings.HasPrefix(post.Message, "!voting-vote") {
		logrus.Infof("New vote command\n")
		bot.registerVote(post)
	} else if strings.HasPrefix(post.Message, "!voting-results") {
		logrus.Infof("New results command\n")
		bot.sendResults(post)
	} else if strings.HasPrefix(post.Message, "!voting-end") {
		logrus.Infof("New end command\n")
		bot.endVoting(post)
	} else if strings.HasPrefix(post.Message, "!voting-delete") {
		logrus.Infof("New delete command\n")
		bot.deleteVoting(post)
	}
}
