package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"petmatch/internal/config"
	"petmatch/internal/models"
	"petmatch/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailInUse            = errors.New("email already registered")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrShelterNotApproved    = errors.New("shelter account pending approval")
	ErrUnsupportedRole       = errors.New("unsupported role")
	ErrAdminCredentialsUnset = errors.New("admin credentials must not be empty")
)

type AuthService struct {
	users          *repositories.UserRepository
	jwtSecret      []byte
	adminEmail     string
	adminPassword  string
	tokenExpiresIn time.Duration
}

type RegisterInput struct {
	Name        string
	Email       string
	Password    string
	Role        string
	ShelterName *string
	Phone       *string
	City        *string
}

type LoginOutput struct {
	Token string
	User  models.User
}

func NewAuthService(repo *repositories.UserRepository, cfg config.Config) (*AuthService, error) {
	if strings.TrimSpace(cfg.AdminEmail) == "" || strings.TrimSpace(cfg.AdminPassword) == "" {
		return nil, ErrAdminCredentialsUnset
	}

	service := &AuthService{
		users:          repo,
		jwtSecret:      []byte(cfg.JWTSecret),
		adminEmail:     cfg.AdminEmail,
		adminPassword:  cfg.AdminPassword,
		tokenExpiresIn: 24 * time.Hour,
	}

	if err := service.ensureDefaultAdmin(); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *AuthService) Register(input RegisterInput) (*models.User, error) {
	role, err := parseRole(input.Role)
	if err != nil {
		return nil, err
	}

	existing, err := s.users.FindByEmail(strings.ToLower(input.Email))
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailInUse
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         input.Name,
		Email:        strings.ToLower(input.Email),
		PasswordHash: string(hash),
		Role:         role,
		ShelterName:  input.ShelterName,
		Phone:        input.Phone,
		City:         input.City,
		IsApproved:   role != models.RoleShelter,
	}

	if err := s.users.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*LoginOutput, error) {
	user, err := s.users.FindByEmail(strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.Role == models.RoleShelter && !user.IsApproved {
		return nil, ErrShelterNotApproved
	}

	token, err := s.generateToken(*user)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) ParseToken(rawToken string) (*models.User, error) {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidCredentials
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	id, err := strconv.Atoi(sub)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	user, err := s.users.FindByID(uint(id))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *AuthService) ApproveShelter(user *models.User) error {
	user.IsApproved = true
	return s.users.Update(user)
}

func (s *AuthService) generateToken(user models.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   fmt.Sprint(user.ID),
		"role":  user.Role,
		"name":  user.Name,
		"email": user.Email,
		"iat":   now.Unix(),
		"exp":   now.Add(s.tokenExpiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ensureDefaultAdmin() error {
	admin, err := s.users.FindByEmail(strings.ToLower(s.adminEmail))
	if err != nil {
		return err
	}

	if admin != nil {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(s.adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:         "Platform Admin",
		Email:        strings.ToLower(s.adminEmail),
		PasswordHash: string(hash),
		Role:         models.RoleAdmin,
		IsApproved:   true,
	}

	return s.users.Create(user)
}

func parseRole(role string) (models.UserRole, error) {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case string(models.RoleAdopter):
		return models.RoleAdopter, nil
	case string(models.RoleShelter):
		return models.RoleShelter, nil
	default:
		return "", ErrUnsupportedRole
	}
}
