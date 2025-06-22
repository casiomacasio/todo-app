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
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrTokenRevoked = errors.New("refresh token revoked")
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
	query := fmt.Sprintf(`SELECT id, username, password_hash FROM %s WHERE username=$1`, usersTable)
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


func (r *AuthPostgres) RevokeRefreshToken(tokenUUID uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`UPDATE %s SET revoked = true WHERE id = $1 AND revoked = false`, refreshTokensTable)

	res, err := r.db.Exec(query, tokenUUID)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *AuthPostgres) RevokeRefreshTokenByUserId(userId int) (bool, error) {
	query := fmt.Sprintf(`UPDATE %s SET revoked = true WHERE user_id = $1 AND revoked = false`, refreshTokensTable)

	res, err := r.db.Exec(query, userId)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *AuthPostgres) SaveRefreshToken(hashed_token string, userId int, expires_at time.Time) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf(`INSERT INTO %s (hashed_token, user_id, expires_at) VALUES ($1, $2, $3) RETURNING id`, refreshTokensTable)
	row := r.db.QueryRow(query, hashed_token, userId, expires_at)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUserIdAndHashByRefreshTokenId(refreshToken uuid.UUID) (int, string, error) {
	type expid struct {
		UserID    int       `db:"user_id"`
		ExpiredAt time.Time `db:"expires_at"`
		HashedToken string `db:"hashed_token"`
		Revoked bool `db:"revoked"`
	}
	var expd expid

	query := fmt.Sprintf(`SELECT user_id, hashed_token, expires_at, revoked FROM %s WHERE id = $1`, refreshTokensTable)
	err := r.db.Get(&expd, query, refreshToken)
	if err != nil {
		return 0, "", err
	}
	if expd.Revoked {
		return 0, "", ErrTokenRevoked
	}
	if time.Now().After(expd.ExpiredAt) {
		return 0, "", ErrRefreshTokenExpired
	}

	return expd.UserID, expd.HashedToken, nil
}


func (r *AuthPostgres) DeleteRefreshToken(refreshToken uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET revoked = true WHERE token = $1`, refreshTokensTable)
	_, err := r.db.Exec(query, refreshToken)
	if err != nil {
		return err
	}
	return nil
}