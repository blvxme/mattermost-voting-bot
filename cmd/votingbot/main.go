package main

import (
	"github.com/sirupsen/logrus"
	"mattermost-voting-bot/internal/votingbot"
	"os"
)

func main() {
	configureLogging()

	bot := votingbot.NewBot()
	bot.Run()
}

func configureLogging() {
	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
}
