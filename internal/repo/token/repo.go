package token

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	tokenDto "github.com/elusiv0/medods_test/internal/model/token"
	"github.com/elusiv0/medods_test/internal/repo"
	tokenModel "github.com/elusiv0/medods_test/internal/repo/token/model"
	mongoClient "github.com/elusiv0/medods_test/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenRepo struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

const (
	collectionName = "tokens"
)

var _ repo.TokenRepo = (*TokenRepo)(nil)

func New(
	client *mongoClient.MongoClient,
	log *slog.Logger,
) *TokenRepo {
	collection := client.MongoDatabase.Collection(collectionName)

	return &TokenRepo{
		collection: collection,
		logger:     log,
	}
}

func (repo *TokenRepo) GetTokenInfo(ctx context.Context, token string) (tokenModel.Token, error) {
	tokenModel := tokenModel.Token{}

	filter := bson.M{"token": token}

	if err := repo.collection.FindOne(ctx, filter).Decode(&tokenModel); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = tokenDto.ErrRefreshTokenNotRegistered
		}
		return tokenModel, fmt.Errorf("TokenRepo - GetTokenInfo - FindOne: %w", err)
	}

	return tokenModel, nil
}

func (repo *TokenRepo) InsertToken(ctx context.Context, token, uuid string) (primitive.ObjectID, error) {
	tokenModel := tokenModel.Token{
		Token:    token,
		UserUUID: uuid,
	}

	result, err := repo.collection.InsertOne(ctx, tokenModel)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("TokenRepository - InsertToken: %w", err)
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (repo *TokenRepo) DeleteToken(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("TokenRepository - DeleteToken: %w", err)
	}
	if result.DeletedCount == 0 {
		return tokenDto.ErrRefreshTokenNotRegistered
	}

	return nil
}
