package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

type User struct {
	TelegramUserID int64
	TelegramChatID int64
	Username       string
	IsAdmin        bool
}

type Subject struct {
	ID          int64
	Name        string
	Description string
}

type QueueEntry struct {
	UserID    int64
	SubjectID int
	JoinedAt  time.Time
}

func Init(dbPool *pgxpool.Pool) *Repository {
	return &Repository{
		DB: dbPool,
	}
}

/* User */
func (r *Repository) UpsertUser(user User) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO users (id, chat_id, username, is_admin)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET chat_id = EXCLUDED.chat_id,
		    username = EXCLUDED.username;
	`, user.TelegramUserID, user.TelegramChatID, user.Username, user.IsAdmin)
	return err
}

func (r *Repository) GetUserByID(userID int64) (*User, error) {
	row := r.DB.QueryRow(context.Background(), `
		SELECT id, chat_id, username, is_admin FROM users WHERE id = $1
	`, userID)

	var u User
	err := row.Scan(&u.TelegramUserID, &u.TelegramChatID, &u.Username, &u.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *Repository) CreateOrUpdateUser(user User) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO users (id, chat_id, username, is_admin)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET chat_id = EXCLUDED.chat_id,
		    username = EXCLUDED.username;
	`, user.TelegramUserID, user.TelegramChatID, user.Username, user.IsAdmin)
	return err
}



/* Subject */
func (r *Repository) CreateSubject(name, description string) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO subjects (name, description) VALUES ($1, $2)
	`, name, description)
	return err
}

func (r *Repository) GetAllSubjects() ([]Subject, error) {
	rows, err := r.DB.Query(context.Background(), `
		SELECT id, name, description FROM subjects ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var s Subject
		err := rows.Scan(&s.ID, &s.Name, &s.Description)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	return subjects, nil
}

/* Queue */

func (r *Repository) AddUserToQueue(subjectID int, userID int64) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO subject_queue (subject_id, user_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING;
	`, subjectID, userID)
	return err
}

func (r *Repository) RemoveUserFromQueue(subjectID int, userID int64) error {
	_, err := r.DB.Exec(context.Background(), `
		DELETE FROM subject_queue WHERE subject_id = $1 AND user_id = $2
	`, subjectID, userID)
	return err
}

func (r *Repository) GetQueueForSubject(subjectID int) ([]User, error) {
	rows, err := r.DB.Query(context.Background(), `
		SELECT u.id, u.full_name, u.username, u.is_admin
		FROM subject_queue q
		JOIN users u ON u.id = q.user_id
		WHERE q.subject_id = $1
		ORDER BY q.joined_at
	`, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.TelegramUserID, &u.Username, &u.IsAdmin); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
