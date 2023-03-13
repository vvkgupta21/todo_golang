package dbHelper

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
	"todo/database"
	"todo/model"
	"todo/utils"
)

func CreateUser(db sqlx.Ext, name, email, password string) error {
	SQL := `INSERT INTO users(name, email, password) VALUES ($1, TRIM(LOWER($2)), $3)`
	_, err := db.Exec(SQL, name, email, password)
	return err
}

func IsUserExist(email string) (bool, error) {
	SQL := `SELECT id FROM users WHERE email = TRIM(LOWER($1)) AND archived_at IS NULL`
	var id int
	err := database.Todo.Get(&id, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func GetUserIDByPassword(email, password string) (int, error) {
	SQL := `SELECT
				u.id,
       			u.password
       		FROM
				users u
			WHERE
				u.archived_at IS NULL
				AND u.email = TRIM(LOWER($1))`
	var userID int
	var passwordHash string
	err := database.Todo.QueryRowx(SQL, email).Scan(&userID, &passwordHash)
	if err != nil {
		return 0, err
	}
	// compare password
	if passwordErr := utils.CheckPassword(password, passwordHash); passwordErr != nil {
		return 0, passwordErr
	}
	return userID, nil
}

func AddTodo(db sqlx.Ext, userId int, title, description string, dueDate time.Time) error {
	SQL := `INSERT INTO todo(userid, title, description, due_date) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(SQL, userId, title, description, dueDate)
	return err
}

func GetAllTask(userId int) ([]model.TaskModel, error) {
	SQL := `SELECT id, title, description, due_date FROM todo WHERE userid = $1 AND archived_at is null `
	list := make([]model.TaskModel, 0)
	err := database.Todo.Select(&list, SQL, userId)
	return list, err
}

func GetTaskById(taskId int) (model.TaskModel, error) {
	SQL := `SELECT id ,title, description, due_date FROM todo WHERE id = $1`
	var user model.TaskModel
	err := database.Todo.Get(&user, SQL, taskId)
	return user, err
}

func UpdateTask(title, description string, dueDate time.Time, userId int, taskId int) error {
	SQL := `UPDATE todo SET title = $1, description = $2, due_date = $3 WHERE userid = $4 AND id = $5`
	_, err := database.Todo.Exec(SQL, title, description, dueDate, userId, taskId)
	return err
}

func DeleteTask(archivedAt time.Time, userId int, taskId int) error {
	SQL := `UPDATE todo SET archived_at = $1 WHERE userid = $2 AND id = $3`
	_, err := database.Todo.Exec(SQL, archivedAt, userId, taskId)
	return err
}
