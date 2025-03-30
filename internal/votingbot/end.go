package votingbot

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/sirupsen/logrus"
	"mattermost-voting-bot/internal/data/repository"
)

func (bot *VotingBot) endVoting(post *model.Post) {
	args := getCommandArgs(post.Message)
	if len(args) != 1 {
		bot.sendMessage(post.ChannelId, `
Ты, кажется, неверно отправил команду. Требуется ровно 1 аргумент: ID голосования.
Вот, как надо: !voting-end "ID голосования"
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
			fmt.Sprintf("Голосование %s может завершить только его создатель!", result.Title),
			post.Id,
		)

		return
	}

	if result.IsEnded {
		bot.sendMessage(post.ChannelId, "Ошибка: голосование уже завершено!", post.Id)
		return
	}

	result.IsEnded = true

	ended, err := repo.Update(result)
	if err != nil {
		logrus.Errorf("Failed to update voting: %+v\n", err)
		bot.sendMessage(post.ChannelId, "Что-то пошло не так :( Попробуй снова", post.Id)

		return
	}

	winner, winnerVotes := "", -1
	optionsStr := ""
	for option, votes := range ended[0].Options {
		optionsStr += fmt.Sprintf("* %s: %d голосов (%d%%)\n", option, votes, (100*votes)/len(ended[0].Users))
		if votes > winnerVotes {
			winner = option
			winnerVotes = votes
		} else if votes == winnerVotes && winner != "" {
			winner += ", " + option
		}
	}

	bot.sendMessage(
		post.ChannelId,
		fmt.Sprintf("Голосование %s завершено. Победитель: %s. Результаты голосования:\n%s", ended[0].Title, winner, optionsStr),
		post.Id,
	)
}
