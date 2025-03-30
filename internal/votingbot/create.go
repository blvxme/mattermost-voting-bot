package votingbot

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/sirupsen/logrus"
	"mattermost-voting-bot/internal/data/entity"
	"mattermost-voting-bot/internal/data/repository"
)

func (bot *VotingBot) createVoting(post *model.Post) {
	args := getCommandArgs(post.Message)
	if len(args) < 3 {
		bot.sendMessage(post.ChannelId, `
Ты, кажется, неверно отправил команду. Требуется как минимум 3 аргумента: заголовок голосования и два варианта ответа.
Вот, как надо: !voting-create "Заголовок голосования" "Вариант ответа 1" "Вариант ответа 2" ... "Вариант ответа N"
`, post.Id)

		return
	}
	logrus.Infof("Command args: %+v\n", args)

	repo, err := repository.NewVotingRepository()
	if err != nil {
		logrus.Errorf("Failed to create voting repository: %+v\n", err)
		bot.sendMessage(post.ChannelId, "Что-то пошло не так :( Попробуй снова", post.Id)
		return
	}

	defer func() {
		if err = repository.DestroyVotingRepository(repo); err != nil {
			logrus.Errorf("Failed to destroy voting repository: %+v\n", err)
		}
	}()

	options := make(map[string]int)
	for i := 1; i < len(args); i++ {
		options[args[i]] = 0
	}

	created, err := repo.Create(entity.VotingEntity{
		Id:        uuid.NewString(),
		CreatorId: post.UserId,
		IsEnded:   false,
		Title:     args[0],
		Options:   options,
		Users:     nil,
	})
	if err != nil {
		logrus.Errorf("Failed to create voting: %+v\n", err)
		bot.sendMessage(post.ChannelId, "Что-то пошло не так :( Попробуй снова", post.Id)
		return
	}
	logrus.Debugf("Created voting: %+v\n", created)

	var optionsStr string
	for option, count := range options {
		optionsStr += fmt.Sprintf("* %s (проголосовавших: %d)\n", option, count)
	}

	bot.sendMessage(
		post.ChannelId,
		fmt.Sprintf("Голосование успешно создано! ID: %s. Вы можете проголосовать за следующие варианты:\n%s", created[0].Id, optionsStr),
		post.Id,
	)
}
