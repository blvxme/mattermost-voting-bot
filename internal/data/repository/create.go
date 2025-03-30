package repository

import (
	"github.com/tarantool/go-tarantool/v2"
	"mattermost-voting-bot/internal/data/entity"
)

func (repo *VotingRepository) Create(entity entity.VotingEntity) (result []entity.VotingEntity, err error) {
	var options []interface{}
	for k, v := range entity.Options {
		options = append(options, map[string]interface{}{k: v})
	}

	tuple := []interface{}{
		entity.Id,
		entity.CreatorId,
		entity.IsEnded,
		entity.Title,
		entity.Options,
		[]interface{}{},
	}

	request := tarantool.NewInsertRequest("votings").Tuple(tuple)
	future := repo.connection.Do(request)

	err = future.GetTyped(&result)

	return
}
