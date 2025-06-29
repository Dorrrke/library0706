package inmemory

import (
	"fmt"

	"github.com/Dorrrke/library0706/internal/domain/errors"
	"github.com/Dorrrke/library0706/internal/domain/models"
)

type UserStrage struct {
	userDB map[string]models.User
}

func NewUserStorage() *UserStrage {
	return &UserStrage{
		userDB: make(map[string]models.User),
	}
}

func (us *UserStrage) SaveUser(user models.User) error {
	for _, dbUser := range us.userDB {
		if dbUser.Email == user.Email {
			return fmt.Errorf("error in db: %w", errors.ErrUserAlredyExist)
		}
	}

	us.userDB[user.UID] = user
	return nil
}

func (us *UserStrage) GetUser(user models.UserLogin) (models.User, error) {
	for _, dbUser := range us.userDB {
		if dbUser.Email == user.Email {
			return dbUser, nil
		}
	}
	return models.User{}, errors.ErrIvalidCreds
}
