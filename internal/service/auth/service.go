package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/elusiv0/medods_test/internal/model/api"
	tokenDto "github.com/elusiv0/medods_test/internal/model/token"
	"github.com/elusiv0/medods_test/internal/repo"
	tokenRepository "github.com/elusiv0/medods_test/internal/repo/token"
	userRepository "github.com/elusiv0/medods_test/internal/repo/user"
	hashing "github.com/elusiv0/medods_test/internal/util/hash"
	tokenManager "github.com/elusiv0/medods_test/internal/util/token"
)

type AuthService struct {
	userRepo     repo.UserRepo
	tokenRepo    repo.TokenRepo
	logger       *slog.Logger
	tokenManager *tokenManager.TokenManager
}

func New(
	userRepo *userRepository.UserRepo,
	tokenRepo *tokenRepository.TokenRepo,
	log *slog.Logger,
	tokenManager *tokenManager.TokenManager,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		logger:       log,
		tokenManager: tokenManager,
	}
}

func (authService *AuthService) SignIn(ctx context.Context, uuid string) (tokenDto.TokenResponse, error) {
	_, err := authService.userRepo.GetUserByUUID(ctx, uuid)
	if err != nil {
		return tokenDto.TokenResponse{}, fmt.Errorf("AuthService - SignIn: %w", err)
	}

	tokens, err := authService.generateTokens(ctx, uuid)
	if err != nil {
		return tokens, fmt.Errorf("AuthService - SignIn: %w", err)
	}

	return tokens, nil
}

func (authService *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
	accessToken string,
) (tokenDto.TokenResponse, error) {
	tokenInfo, err := authService.tokenManager.ValidateJWT(accessToken)
	if err != nil && !(errors.Is(err, api.ErrAccessTokenExpired)) {
		return tokenDto.TokenResponse{}, fmt.Errorf("AuthService - Refresh: %w", err)
	}

	if ok, err := authService.tokenManager.CheckTokenMatch(tokenInfo, refreshToken); !ok {
		return tokenDto.TokenResponse{}, fmt.Errorf("AuthService - Refresh: %w", err)
	}

	if err := authService.tokenRepo.DeleteToken(ctx, tokenInfo.RefreshId); err != nil {
		return tokenDto.TokenResponse{}, fmt.Errorf("AuthService - Refresh: %w", err)
	}

	uuid := tokenInfo.UUID

	tokens, err := authService.generateTokens(ctx, uuid)

	if err != nil {
		return tokens, fmt.Errorf("AuthService - Refresh: %w", err)
	}

	return tokens, nil
}

func (authService *AuthService) generateTokens(ctx context.Context, uuid string) (tokenDto.TokenResponse, error) {
	refreshToken, err := authService.tokenManager.NewRefreshToken()
	if err != nil {
		return tokenDto.TokenResponse{}, fmt.Errorf("generateTokens: %w", err)
	}

	hashedRefresh, err := hashing.CryptToken(refreshToken)
	if err != nil {
		return tokenDto.TokenResponse{}, fmt.Errorf("generateTokens: %w", err)
	}

	refreshId, err := authService.tokenRepo.InsertToken(ctx, hashedRefresh, uuid)
	if err != nil {
		return tokenDto.TokenResponse{}, fmt.Errorf("generateTokens: %w", err)
	}

	accessToken, err := authService.tokenManager.NewJWTToken(hashedRefresh, uuid, refreshId)
	if err != nil {
		return tokenDto.TokenResponse{}, fmt.Errorf("generateTokens: %w", err)
	}

	return tokenDto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
