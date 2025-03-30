package votingbot

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/sirupsen/logrus"
	"regexp"
)

func (bot *VotingBot) sendMessage(channelId, message string, rootId string) {
	post := &model.Post{ChannelId: channelId, Message: message, RootId: rootId}

	_, _, err := bot.client.CreatePost(post)
	if err != nil {
		logrus.Errorf("Failed to create post: %+v\n", err)
	}
}

func getCommandArgs(command string) (args []string) {
	args = make([]string, 0)

	matches := regexp.MustCompile(`"([^"]+)"`).FindAllStringSubmatch(command, -1)

	for _, match := range matches {
		if len(match) > 1 {
			args = append(args, match[1])
		}
	}

	return
}
