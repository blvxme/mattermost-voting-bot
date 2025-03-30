package votingbot

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/sirupsen/logrus"
	"mattermost-voting-bot/internal/data/repository"
)

func (bot *VotingBot) registerVote(post *model.Post) {
	args := getCommandArgs(post.Message)
	if len(args) != 2 {
		bot.sendMessage(post.ChannelId, `
Ты, кажется, неверно отправил команду. Требуется ровно 2 аргумента: ID голосования и выбранный тобой вариант ответа.
Вот, как надо: !voting-vote "ID голосования" "Вариант ответа"
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

	results, err := repo.GetById(args[0])
	if err != nil {
		logrus.Errorf("Failed to get voting: %+v\n", err)
		bot.sendMessage(post.ChannelId, "Что-то пошло не так :( Попробуй снова", post.Id)
		return
	}

	if len(results) == 0 {
		bot.sendMessage(post.ChannelId, "Кажется, голосования с таким ID нет. Попробуй снова.", post.Id)
		return
	}

	result := results[0]

	if result.IsEnded {
		bot.sendMessage(
			post.ChannelId,
			fmt.Sprintf("Голосование %s завершено. Проголосовать уже не получится", result.Title),
			post.Id,
		)

		return
	}

	_, exists := result.Options[args[1]]
	if !exists {
		optionsStr := ""
		for option := range result.Options {
			optionsStr += "* " + option + "\n"
		}

		bot.sendMessage(
			post.ChannelId,
			fmt.Sprintf("Кажется, такого варианта нет. Доступны следующие варианты для голосования:\n%s", optionsStr),
			post.Id,
		)

		return
	}

	for _, user := range result.Users {
		if user == post.UserId {
			bot.sendMessage(post.ChannelId, "Второй раз голосовать запрещено!", post.Id)
			return
		}
	}

	result.Options[args[1]]++
	result.Users = append(result.Users, post.UserId)

	_, err = repo.Update(result)
	if err != nil {
		logrus.Errorf("Failed to update voting: %+v\n", err)
		bot.sendMessage(post.ChannelId, "Что-то пошло не так :( Попробуй снова", post.Id)
		return
	}

	bot.sendMessage(post.ChannelId, fmt.Sprintf("Вы успешно проголосовали за вариант %s!", args[1]), post.Id)
}
