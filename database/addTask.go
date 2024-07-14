package database

import (
    "errors"
    "github.com/jmoiron/sqlx"
    "github.com/Masteker/go_final_project/models"
    "github.com/Masteker/go_final_project/tasks"
    "time"
)

type TaskRepository struct {
    db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
    return &TaskRepository{db: db}
}

func (tr *TaskRepository) AddTask(task models.Task) (int64, error) {
    if task.Title == "" {
        return 0, errors.New("Не указан заголовок задачи")
    }

    // Rest of your AddTask logic here

    result, err := tr.db.Exec(
        `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
        task.Date, task.Title, task.Comment, task.Repeat,
    )
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return id, nil
}