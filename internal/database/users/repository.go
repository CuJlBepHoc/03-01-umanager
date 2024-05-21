package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

func New(userDB *pgx.Conn, timeout time.Duration) *Repository {
	return &Repository{userDB: userDB, timeout: timeout}
}

type Repository struct {
	userDB  *pgx.Conn
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateUserReq) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	u = database.User{
		ID:       req.ID,
		Username: req.Username,
		Password: req.Password,
	}

	query := `
		INSERT INTO users (id, username, password) 
		VALUES ($1, $2, $3)
		RETURNING id, username, password
	`

	err := r.userDB.QueryRow(ctx, query, u.ID, u.Username, u.Password).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return database.User{}, err
	}
	return u, nil
}

func (r *Repository) FindByID(ctx context.Context, userID uuid.UUID) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
		SELECT id, username, password 
		FROM users 
		WHERE id = $1
	`

	err := r.userDB.QueryRow(ctx, query, userID).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return database.User{}, nil // Возвращает пустой объект, если пользователь не найден
		}
		return database.User{}, err
	}

	return u, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
		SELECT id, username, password 
		FROM users 
		WHERE username = $1
	`

	err := r.userDB.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return database.User{}, nil // Возвращает пустой объект, если пользователь не найден
		}
		return database.User{}, err
	}

	return u, nil
}

// ClearTable очищает таблицу пользователей.
func (r *Repository) ClearTable(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.userDB.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	return err
}
