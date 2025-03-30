package repository

import (
	"github.com/tarantool/go-tarantool/v2"
	"mattermost-voting-bot/internal/data/entity"
)

func (repo *VotingRepository) GetById(id string) (result []entity.VotingEntity, err error) {
	err = repo.connection.Do(
		tarantool.NewSelectRequest("votings").Iterator(tarantool.IterEq).Key([]interface{}{id}),
	).GetTyped(&result)

	return
}
