package repository

import (
	"faceit/database"
	"faceit/model"
)

type UsersRepository[T model.User] struct {
	Repository[T]
}

func NewUsersRepository() (*UsersRepository[model.User], error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	return &UsersRepository[model.User]{
		Repository: Repository[model.User]{DB: connection},
	}, nil
}
