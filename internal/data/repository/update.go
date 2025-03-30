package repository

import (
	"github.com/tarantool/go-tarantool/v2"
	"mattermost-voting-bot/internal/data/entity"
)

func (repo *VotingRepository) Update(entity entity.VotingEntity) (result []entity.VotingEntity, err error) {
	err = repo.connection.Do(
		tarantool.NewUpdateRequest("votings").
			Key(tarantool.StringKey{S: entity.Id}).
			Operations(tarantool.NewOperations().
				Assign(1, entity.CreatorId).
				Assign(2, entity.IsEnded).
				Assign(3, entity.Title).
				Assign(4, entity.Options).
				Assign(5, entity.Users),
			),
	).GetTyped(&result)

	return
}
