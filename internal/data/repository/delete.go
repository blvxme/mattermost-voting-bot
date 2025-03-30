package repository

import (
	"github.com/tarantool/go-tarantool/v2"
	"mattermost-voting-bot/internal/data/entity"
)

func (repo *VotingRepository) DeleteById(id string) (result []entity.VotingEntity, err error) {
	err = repo.connection.Do(
		tarantool.NewDeleteRequest("votings").Key([]interface{}{id}),
	).GetTyped(&result)

	return
}
