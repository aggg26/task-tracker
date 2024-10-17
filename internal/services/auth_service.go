package services

import (
	"database/sql"
	"errors"
	"time"
	"trackerApp/internal/models"
	"trackerApp/internal/services/dtos"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

const (
	getByUsername = `SELECT id,username,password FROM users WHERE username=$1`
	createUser    = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
)

const secretKey = "secret key"

type CustomClaims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *AuthService) AddUser(form dtos.UserForm) (int, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	var userId int
	if err := s.db.QueryRow(createUser, form.Username, string(passwordHash)).Scan(&userId); err != nil {
		return 0, err
	}
	return userId, nil
}

func (s *AuthService) GenerateJwt(form dtos.UserForm) (string, error) {
	var user models.User
	if err := s.db.QueryRow(getByUsername, form.Username).Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return "", errors.New("Wrong username or password")
	}

	claims := CustomClaims{
		user.ID,
		user.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(25 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ParseJwt(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, errors.New("token claims are not of type (*CustomClaims)")
	}
	return claims.UserId, nil
}
