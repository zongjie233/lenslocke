package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/zongjie233/lenslocked/rand"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	// 保证token的最小字节数
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	//更新用户session, 如果user_id部分存在冲突则进行更新,将token_hash设置为传入的新值,确保用户永远拥有一个最新的session
	row := ss.DB.QueryRow(`
			Insert     
			Into
				sessions
				(user_id,token_hash)     
			VALUES
				($1,$2)     
					ON CONFLICT (user_id) DO UPDATE
						
				SET
					token_hash = $2 RETURNING id`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	// 存储session
	if err != nil {
		return nil, fmt.Errorf("create:%w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User
	row := ss.DB.QueryRow(`
				SELECT
			users.id,
			users.email,
			users.password_hash 
		FROM
			sessions 
		JOIN
			users 
				ON users.id= sessions.user_id 
		WHERE
			sessions.token_hash = $1`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`
	DELETE FROM sessions 
	WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
