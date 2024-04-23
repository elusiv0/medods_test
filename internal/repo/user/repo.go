package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	mapper "github.com/elusiv0/medods_test/internal/mapper/user"
	userDto "github.com/elusiv0/medods_test/internal/model/user"
	"github.com/elusiv0/medods_test/internal/repo"
	userModel "github.com/elusiv0/medods_test/internal/repo/user/model"
	mongoClient "github.com/elusiv0/medods_test/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

const (
	collectionName = "users"
)

var _ repo.UserRepo = (*UserRepo)(nil)

func New(
	client *mongoClient.MongoClient,
	log *slog.Logger,
) *UserRepo {
	collection := client.MongoDatabase.Collection(collectionName)

	return &UserRepo{
		collection: collection,
		logger:     log,
	}
}

func (repo *UserRepo) GetUserByUUID(ctx context.Context, uuid string) (userDto.User, error) {
	userModel := userModel.User{}

	filter := bson.D{{"_id", uuid}}
	if err := repo.collection.FindOne(ctx, filter).Decode(&userModel); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = userDto.ErrUserNotFound
		}
		return userDto.User{}, fmt.Errorf("UserRepo - GetUserByUUID - FindOne: %w", err)
	}

	return mapper.ModelToUser(userModel), nil
}

func (repo *UserRepo) InsertUser(ctx context.Context, user userDto.CreateUser) (string, error) {
	userModel := mapper.CreateUserToUserModel(user)
	result, err := repo.collection.InsertOne(ctx, userModel)

	if err != nil {
		return "", fmt.Errorf("UserRepo - InsertUser - InsertOne: %w", err)
	}

	uuid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("UserRepo - InsertUser - Get Inserted ID: %w", err)
	}

	return uuid.Hex(), nil
}
