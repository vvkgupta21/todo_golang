package model

import "time"

type RegisterModel struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type LoginModel struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type TodoModel struct {
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	DueDate     string `json:"dueDate" db:"due_date"`
}

type TaskModel struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	DueDate     time.Time `json:"dueDate" db:"due_date"`
}
