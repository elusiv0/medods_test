package repo

import (
	"context"

	userDto "github.com/elusiv0/medods_test/internal/model/user"
	tokenModel "github.com/elusiv0/medods_test/internal/repo/token/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo interface {
	GetUserByUUID(ctx context.Context, uuid string) (userDto.User, error)
	InsertUser(ctx context.Context, user userDto.CreateUser) (string, error)
}

type TokenRepo interface {
	GetTokenInfo(ctx context.Context, token string) (tokenModel.Token, error)
	InsertToken(ctx context.Context, token, uuid string) (primitive.ObjectID, error)
	DeleteToken(ctx context.Context, id primitive.ObjectID) error
}
