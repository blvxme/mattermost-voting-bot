package repository

import (
	"fmt"
	"os"
)

type votingRepositoryConfig struct {
	host string
	port string

	user     string
	password string
}

func getConfig() (config *votingRepositoryConfig, err error) {
	host := os.Getenv("TARANTOOL_HOST")
	if host == "" {
		err = fmt.Errorf("TARANTOOL_HOST environment variable not set")
		return
	}

	port := os.Getenv("TARANTOOL_PORT")
	if port == "" {
		err = fmt.Errorf("TARANTOOL_PORT environment variable not set")
		return
	}

	user := os.Getenv("TARANTOOL_USER")
	if user == "" {
		err = fmt.Errorf("TARANTOOL_USER environment variable not set")
		return
	}

	password := os.Getenv("TARANTOOL_PASSWORD")
	if password == "" {
		err = fmt.Errorf("TARANTOOL_PASSWORD environment variable not set")
		return
	}

	config = &votingRepositoryConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}

	return
}
