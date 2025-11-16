CREATE TABLE IF NOT EXISTS todos.todos (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    is_done     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_todos_is_done ON todos.todos(is_done);
