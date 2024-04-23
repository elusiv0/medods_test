package user

import (
	userDto "github.com/elusiv0/medods_test/internal/model/user"
	userRepo "github.com/elusiv0/medods_test/internal/repo/user/model"
	uuidUtil "github.com/google/uuid"
)

func ModelToUser(userModel userRepo.User) userDto.User {
	return userDto.User{
		UUID: userModel.UUID,
		Name: userModel.Name,
	}
}

func CreateUserToUserModel(userCreate userDto.CreateUser) userRepo.User {
	uuid := uuidUtil.New()

	return userRepo.User{
		UUID: uuid.String(),
		Name: userCreate.Name,
	}
}
