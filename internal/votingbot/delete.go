package votingbot

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/sirupsen/logrus"
	"mattermost-voting-bot/internal/data/repository"
)

func (bot *VotingBot) deleteVoting(post *model.Post) {
	args := getCommandArgs(post.Message)
	if len(args) != 1 {
		bot.sendMessage(post.ChannelId, `
Ты, кажется, неверно отправил команду. Требуется ровно 1 аргумент: ID голосования.
Вот, как надо: !voting-delete "ID голосования"
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

	if result.CreatorId != post.UserId {
		bot.sendMessage(
			post.ChannelId,
			fmt.Sprintf("Голосование %s может удалить только его создатель!", result.Title),
			post.Id,
		)

		return
	}

	if !result.IsEnded {
		bot.sendMessage(
			post.ChannelId,
			fmt.Sprintf("Голосование %s не завершено! Сначала нужно завершить голосование.", result.Title),
			post.Id,
		)

		return
	}

	deleted, err := repo.DeleteById(result.Id)
	if err != nil {
		logrus.Errorf("Failed to delete voting: %+v\n", err)
		bot.sendMessage(post.ChannelId, "Что-то пошло не так :( Попробуй снова", post.Id)
		return
	}

	bot.sendMessage(post.ChannelId, fmt.Sprintf("Голосование %s успешно удалено!", deleted[0].Title), post.Id)
}
