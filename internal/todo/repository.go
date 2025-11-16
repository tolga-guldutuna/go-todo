package todo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Todo, error)
	GetByID(ctx context.Context, id int) (*Todo, error)
	Create(ctx context.Context, title, description string) (*Todo, error)
	Update(ctx context.Context, t Todo) (*Todo, error)
	Delete(ctx context.Context, id int) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) GetAll(ctx context.Context) ([]Todo, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, title, description, is_done, created_at, updated_at 
			FROM todos ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.IsDone, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *repo) GetByID(ctx context.Context, id int) (*Todo, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, title, description,  is_done, created_at, updated_at 
         FROM todos WHERE id = $1`,
		id,
	)

	var t Todo
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.IsDone, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *repo) Create(ctx context.Context, title string, description string) (*Todo, error) {
	row := r.db.QueryRow(ctx,
		`INSERT INTO todos (title, description, is_done, created_at, updated_at)
         VALUES ($1, $2, false, NOW(), NOW())
         RETURNING id, title, description, is_done, created_at, updated_at`,
		title, description,
	)

	var t Todo
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.IsDone, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *repo) Update(ctx context.Context, t Todo) (*Todo, error) {
	t.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		`UPDATE todos
         SET title = $1, description = $2, is_done = $3, updated_at = $4
         WHERE id = $5`,
		t.Title, t.Description, t.IsDone, t.UpdatedAt, t.ID,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *repo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM todos WHERE id=$1`,
		id,
	)
	return err
}
