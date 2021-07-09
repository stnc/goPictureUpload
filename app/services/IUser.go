package services

import (
	"stncCms/app/domain/entity"
)

//UserAppInterface interface
type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUser(uint64) (*entity.User, error)
	GetUsers() ([]entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
	GetUserByEmailAndPassword2(email string, password string) (*entity.User, map[string]string)
}

type userApp struct {
	request UserAppInterface
}

//UserApp implements the UserAppInterface
var _ UserAppInterface = &userApp{}

func (u *userApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.request.SaveUser(user)
}

func (u *userApp) GetUser(userID uint64) (*entity.User, error) {
	return u.request.GetUser(userID)
}

func (u *userApp) GetUsers() ([]entity.User, error) {
	return u.request.GetUsers()
}

func (u *userApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.request.GetUserByEmailAndPassword(user)
}
func (u *userApp) GetUserByEmailAndPassword2(email string, password string) (*entity.User, map[string]string) {
	return u.request.GetUserByEmailAndPassword2(email, password)
}
