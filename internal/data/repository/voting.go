package repository

import (
	"context"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
	"time"
)

type VotingRepository struct {
	connection *tarantool.Connection
	config     *votingRepositoryConfig
}

func NewVotingRepository() (repo *VotingRepository, err error) {
	config, err := getConfig()
	if err != nil {
		err = fmt.Errorf("failed to get voting repository configuration: %s", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  config.host + ":" + config.port,
		User:     config.user,
		Password: config.password,
	}

	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		err = fmt.Errorf("failed to connect to tarantool: %s", err.Error())
		return
	}

	repo = &VotingRepository{
		connection: conn,
		config:     config,
	}

	return
}

func DestroyVotingRepository(repo *VotingRepository) (err error) {
	err = repo.connection.CloseGraceful()
	return
}
