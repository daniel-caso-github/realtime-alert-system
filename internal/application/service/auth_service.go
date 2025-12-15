package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
)

// Auth service errors.
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
)

// TokenPair represents access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// JWTClaims represents the JWT token claims.
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService handles authentication and authorization logic.
type AuthService struct {
	userRepo  repository.UserRepository
	cacheRepo repository.CacheRepository
	jwtConfig *config.JWTConfig
}

// NewAuthService creates a new authentication service.
func NewAuthService(
	userRepo repository.UserRepository,
	cacheRepo repository.CacheRepository,
	jwtConfig *config.JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
		jwtConfig: jwtConfig,
	}
}

// Login authenticates a user and returns tokens.
func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, *entity.User, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil, ErrUserNotActive
	}

	// Verify password
	passwordHash := valueobject.NewPasswordHashFromHash(user.PasswordHash)
	if !passwordHash.Verify(password) {
		return nil, nil, ErrInvalidCredentials
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, err
	}

	// Update last login
	user.UpdateLastLogin()

	return tokens, user, nil
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, email, password, name string) (*TokenPair, *entity.User, error) {
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, ErrUserAlreadyExists
	}

	// Hash password
	passwordHash, err := valueobject.NewPasswordHash(password)
	if err != nil {
		return nil, nil, err
	}

	// Create user
	user, err := entity.NewUser(email, passwordHash.Value(), name, entity.UserRoleViewer)
	if err != nil {
		return nil, nil, err
	}

	// Save user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, err
	}

	return tokens, user, nil
}

// RefreshToken generates new tokens using a refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Parse and validate refresh token
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if token is blacklisted
	blacklistKey := "blacklist:" + refreshToken
	exists, _ := s.cacheRepo.Exists(ctx, blacklistKey)
	if exists {
		return nil, ErrTokenInvalid
	}

	// Get user
	userID, err := entity.ParseID(claims.UserID)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Blacklist old refresh token
	_ = s.cacheRepo.Set(ctx, blacklistKey, true, s.jwtConfig.RefreshExpiration)

	// Generate new tokens
	return s.generateTokenPair(user)
}

// Logout invalidates the user's tokens.
func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	// Blacklist both tokens
	if accessToken != "" {
		blacklistKey := "blacklist:" + accessToken
		_ = s.cacheRepo.Set(ctx, blacklistKey, true, s.jwtConfig.Expiration)
	}

	if refreshToken != "" {
		blacklistKey := "blacklist:" + refreshToken
		_ = s.cacheRepo.Set(ctx, blacklistKey, true, s.jwtConfig.RefreshExpiration)
	}

	return nil
}

// ValidateToken validates an access token and returns claims.
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check if token is blacklisted
	blacklistKey := "blacklist:" + tokenString
	exists, _ := s.cacheRepo.Exists(ctx, blacklistKey)
	if exists {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// generateTokenPair creates access and refresh tokens.
func (s *AuthService) generateTokenPair(user *entity.User) (*TokenPair, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtConfig.Expiration)

	// Access token claims
	accessClaims := JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    s.jwtConfig.Issuer,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, err
	}

	// Refresh token claims
	refreshClaims := JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtConfig.RefreshExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    s.jwtConfig.Issuer,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    expiresAt,
	}, nil
}

// validateToken parses and validates a JWT token.
func (s *AuthService) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}
