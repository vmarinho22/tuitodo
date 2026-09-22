package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"tuitodo/internal/domain"

	_ "modernc.org/sqlite"
)

var ErrCategoryInUse = errors.New("category still has parent tasks")

const SettingLocale = "locale"

const schema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY,
	parent_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
	category_id INTEGER REFERENCES categories(id),
	title TEXT NOT NULL,
	completed_at TEXT,
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS tasks_parent_id_idx ON tasks(parent_id);
CREATE INDEX IF NOT EXISTS tasks_category_id_idx ON tasks(category_id);
CREATE INDEX IF NOT EXISTS tasks_completed_at_idx ON tasks(completed_at);

CREATE TABLE IF NOT EXISTS settings (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL
);
`

type SQLiteStore struct {
	db *sql.DB
}

func DefaultDatabasePath() (string, error) {
	if runtime.GOOS == "darwin" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("user config dir: %w", err)
		}
		return filepath.Join(configDir, "tuitodo", "tuitodo.db"), nil
	}
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("user home dir: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".local", "share")
	}
	return filepath.Join(dataDir, "tuitodo", "tuitodo.db"), nil
}

func Open(databasePath string) (*SQLiteStore, error) {
	if databasePath != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(databasePath), 0o755); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate sqlite: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	db.SetMaxOpenConns(1)
	return &SQLiteStore{db: db}, nil
}

func (sqliteStore *SQLiteStore) Close() error {
	return sqliteStore.db.Close()
}

func (sqliteStore *SQLiteStore) InsertCategory(category domain.Category) (domain.Category, error) {
	if err := domain.ValidateTitle(category.Name); err != nil {
		return domain.Category{}, err
	}
	result, err := sqliteStore.db.Exec(
		`INSERT INTO categories (name, created_at) VALUES (?, ?)`,
		category.Name,
		category.CreatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return domain.Category{}, fmt.Errorf("insert category: %w", err)
	}
	category.ID, err = result.LastInsertId()
	if err != nil {
		return domain.Category{}, fmt.Errorf("category id: %w", err)
	}
	return category, nil
}

func (sqliteStore *SQLiteStore) RenameCategory(categoryID int64, name string) error {
	if err := domain.ValidateTitle(name); err != nil {
		return err
	}
	_, err := sqliteStore.db.Exec(`UPDATE categories SET name = ? WHERE id = ?`, name, categoryID)
	if err != nil {
		return fmt.Errorf("rename category: %w", err)
	}
	return nil
}

func (sqliteStore *SQLiteStore) DeleteCategory(categoryID int64) error {
	var parentCount int
	if err := sqliteStore.db.QueryRow(
		`SELECT COUNT(*) FROM tasks WHERE parent_id IS NULL AND category_id = ?`,
		categoryID,
	).Scan(&parentCount); err != nil {
		return fmt.Errorf("count category parent tasks: %w", err)
	}
	if parentCount > 0 {
		return ErrCategoryInUse
	}
	_, err := sqliteStore.db.Exec(`DELETE FROM categories WHERE id = ?`, categoryID)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}

func (sqliteStore *SQLiteStore) CategoryByID(categoryID int64) (domain.Category, error) {
	row := sqliteStore.db.QueryRow(
		`SELECT id, name, created_at FROM categories WHERE id = ?`,
		categoryID,
	)
	category, err := scanCategory(row)
	if err != nil {
		return domain.Category{}, fmt.Errorf("category by id: %w", err)
	}
	return category, nil
}

func (sqliteStore *SQLiteStore) ListCategories() ([]domain.Category, error) {
	rows, err := sqliteStore.db.Query(`SELECT id, name, created_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	categories := make([]domain.Category, 0)
	for rows.Next() {
		category, err := scanCategory(rows)
		if err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (sqliteStore *SQLiteStore) InsertParentTask(task domain.Task) (domain.Task, error) {
	if err := domain.ValidateTitle(task.Title); err != nil {
		return domain.Task{}, err
	}
	if task.CategoryID == nil {
		return domain.Task{}, fmt.Errorf("insert parent task: category is required")
	}
	result, err := sqliteStore.db.Exec(
		`INSERT INTO tasks (parent_id, category_id, title, completed_at, created_at) VALUES (NULL, ?, ?, ?, ?)`,
		*task.CategoryID,
		task.Title,
		nullableTime(task.CompletedAt),
		task.CreatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("insert parent task: %w", err)
	}
	task.ID, err = result.LastInsertId()
	if err != nil {
		return domain.Task{}, fmt.Errorf("parent task id: %w", err)
	}
	return task, nil
}

func (sqliteStore *SQLiteStore) InsertSubtask(subtask domain.Task, parent domain.Task) (domain.Task, error) {
	if err := domain.ValidateSubtask(subtask, parent); err != nil {
		return domain.Task{}, err
	}
	subtask.ParentID = &parent.ID
	subtask.CategoryID = nil
	result, err := sqliteStore.db.Exec(
		`INSERT INTO tasks (parent_id, category_id, title, completed_at, created_at) VALUES (?, NULL, ?, ?, ?)`,
		parent.ID,
		subtask.Title,
		nullableTime(subtask.CompletedAt),
		subtask.CreatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("insert subtask: %w", err)
	}
	subtask.ID, err = result.LastInsertId()
	if err != nil {
		return domain.Task{}, fmt.Errorf("subtask id: %w", err)
	}
	return subtask, nil
}

func (sqliteStore *SQLiteStore) UpdateTaskTitle(taskID int64, title string) error {
	if err := domain.ValidateTitle(title); err != nil {
		return err
	}
	_, err := sqliteStore.db.Exec(`UPDATE tasks SET title = ? WHERE id = ?`, title, taskID)
	if err != nil {
		return fmt.Errorf("update task title: %w", err)
	}
	return nil
}

func (sqliteStore *SQLiteStore) SaveTaskCompletions(parent domain.Task, subtasks []domain.Task) error {
	tx, err := sqliteStore.db.Begin()
	if err != nil {
		return fmt.Errorf("begin save task completions: %w", err)
	}
	defer tx.Rollback()

	if err := updateTaskCompletion(tx, parent); err != nil {
		return err
	}
	for _, subtask := range subtasks {
		if err := updateTaskCompletion(tx, subtask); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit save task completions: %w", err)
	}
	return nil
}

func updateTaskCompletion(tx *sql.Tx, task domain.Task) error {
	_, err := tx.Exec(`UPDATE tasks SET completed_at = ? WHERE id = ?`, nullableTime(task.CompletedAt), task.ID)
	if err != nil {
		return fmt.Errorf("save task completion: %w", err)
	}
	return nil
}

func (sqliteStore *SQLiteStore) DeleteTaskByID(taskID int64) error {
	_, err := sqliteStore.db.Exec(`DELETE FROM tasks WHERE id = ?`, taskID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

func (sqliteStore *SQLiteStore) ParentTaskByID(taskID int64) (domain.Task, error) {
	row := sqliteStore.db.QueryRow(
		`SELECT id, parent_id, category_id, title, completed_at, created_at FROM tasks WHERE id = ? AND parent_id IS NULL`,
		taskID,
	)
	task, err := scanTask(row)
	if err != nil {
		return domain.Task{}, fmt.Errorf("parent task by id: %w", err)
	}
	return task, nil
}

func (sqliteStore *SQLiteStore) SubtasksByParentID(parentID int64) ([]domain.Task, error) {
	rows, err := sqliteStore.db.Query(
		`SELECT id, parent_id, category_id, title, completed_at, created_at FROM tasks WHERE parent_id = ? ORDER BY id`,
		parentID,
	)
	if err != nil {
		return nil, fmt.Errorf("subtasks by parent id: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

func (sqliteStore *SQLiteStore) ListPendingParentTasks(categoryID *int64) ([]domain.Task, error) {
	return sqliteStore.listParentTasks(categoryID, false)
}

func (sqliteStore *SQLiteStore) ListCompletedParentTasks(categoryID *int64) ([]domain.Task, error) {
	return sqliteStore.listParentTasks(categoryID, true)
}

func (sqliteStore *SQLiteStore) listParentTasks(categoryID *int64, completed bool) ([]domain.Task, error) {
	query := `SELECT id, parent_id, category_id, title, completed_at, created_at FROM tasks WHERE parent_id IS NULL`
	args := make([]any, 0, 2)
	if completed {
		query += ` AND (
			completed_at IS NOT NULL
			OR EXISTS (
				SELECT 1 FROM tasks AS sub
				WHERE sub.parent_id = tasks.id AND sub.completed_at IS NOT NULL
			)
		)`
	} else {
		query += ` AND completed_at IS NULL`
	}
	if categoryID != nil {
		query += ` AND category_id = ?`
		args = append(args, *categoryID)
	}
	if completed {
		query += ` ORDER BY COALESCE(
			completed_at,
			(SELECT MAX(sub.completed_at) FROM tasks AS sub WHERE sub.parent_id = tasks.id)
		) DESC`
	} else {
		query += ` ORDER BY created_at`
	}

	rows, err := sqliteStore.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list parent tasks: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCategory(row scanner) (domain.Category, error) {
	var category domain.Category
	var createdAt string
	if err := row.Scan(&category.ID, &category.Name, &createdAt); err != nil {
		return domain.Category{}, err
	}
	parsed, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return domain.Category{}, err
	}
	category.CreatedAt = parsed
	return category, nil
}

func scanTask(row scanner) (domain.Task, error) {
	var task domain.Task
	var parentID sql.NullInt64
	var categoryID sql.NullInt64
	var completedAt sql.NullString
	var createdAt string
	if err := row.Scan(&task.ID, &parentID, &categoryID, &task.Title, &completedAt, &createdAt); err != nil {
		return domain.Task{}, err
	}
	if parentID.Valid {
		task.ParentID = &parentID.Int64
	}
	if categoryID.Valid {
		task.CategoryID = &categoryID.Int64
	}
	if completedAt.Valid {
		parsed, err := time.Parse(time.RFC3339Nano, completedAt.String)
		if err != nil {
			return domain.Task{}, err
		}
		task.CompletedAt = &parsed
	}
	parsedCreatedAt, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return domain.Task{}, err
	}
	task.CreatedAt = parsedCreatedAt
	return task, nil
}

func scanTasks(rows *sql.Rows) ([]domain.Task, error) {
	tasks := make([]domain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.Format(time.RFC3339Nano)
}

func (sqliteStore *SQLiteStore) Setting(key string) (string, error) {
	var value string
	err := sqliteStore.db.QueryRow(`SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get setting: %w", err)
	}
	return value, nil
}

func (sqliteStore *SQLiteStore) SetSetting(key, value string) error {
	_, err := sqliteStore.db.Exec(
		`INSERT INTO settings (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key,
		value,
	)
	if err != nil {
		return fmt.Errorf("set setting: %w", err)
	}
	return nil
}
