package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameExists  = errors.New("username already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id`, usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, ErrUsernameExists
		}
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT id, password_hash FROM %s WHERE username=$1`, usersTable)
	err := r.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidPassword
	}
	return user, nil
}
// type RefreshToken struct {
// 	Token     uuid.UUID
// 	UserID    int
// 	IssuedAt  time.Time
// 	ExpiresAt time.Time
// 	Revoked   bool
// }

func (r *AuthPostgres) SaveRefreshToken(refreshToken uuid.UUID, userId int, expires_at time.Time) error {
	query := fmt.Sprintf(`INSERT INTO %s (token, user_id, expires_at) VALUES ($1, $2, $3)`, refreshTokensTable)
	_, err := r.db.Exec(query, refreshToken, userId, expires_at)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) CheckRefreshToken(refreshToken uuid.UUID) (int, error) {
	var userID int
	query := fmt.Sprintf(`SELECT user_id FROM %s WHERE token = $1`, refreshTokensTable)
	err := r.db.Get(&userID, query, refreshToken)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *AuthPostgres) DeleteRefreshToken(refreshToken uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET revoked = true WHERE token = $1`, refreshTokensTable)
	_, err := r.db.Exec(query, refreshToken)
	if err != nil {
		return err
	}
	return nil
}