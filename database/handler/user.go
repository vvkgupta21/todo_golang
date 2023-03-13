package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
	"todo/database"
	"todo/database/dbHelper"
	"todo/middleware"
	"todo/model"
	"todo/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body model.RegisterModel
	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to parse Request body")
	}
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password must be 6 chars long")
		return
	}

	exist, existError := dbHelper.IsUserExist(body.Email)

	if existError != nil {
		utils.RespondError(w, http.StatusInternalServerError, existError, "failed to check user existence")
		return
	}

	if exist {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exist")
		return
	}

	hashPassword, hasErr := utils.HashPassword(body.Password)

	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
		return
	}

	err := dbHelper.CreateUser(database.Todo, body.Name, body.Email, hashPassword)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to create user in table")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"User created successfully"})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body model.LoginModel
	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to Parse Request Body")
		return
	}
	userId, err := dbHelper.GetUserIDByPassword(body.Email, body.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondError(w, http.StatusBadRequest, errors.New("user does not exist"), "user does not exist")
			return
		}
		utils.RespondError(w, http.StatusBadRequest, err, "failed to find user")
		return
	}

	token, err := middleware.GenerateJWT(userId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error in generating jwt token")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var body model.TodoModel

	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to parse request body")
		return
	}

	//const layout = "2006-Jan-02"
	userId := getUserId(r)
	//date, err := time.Parse("2006-01-02T15:04:05", body.DueDate)
	//date, _ := time.Parse(time.RFC3339, "Feb 4, 2014 at 6:05pm (PST)")
	date, err := time.Parse(time.RFC3339, body.DueDate)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dbHelper.AddTodo(database.Todo, userId, body.Title, body.Description, date)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to create todo in table")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"Todo  created successfully"})
}

func GetAllTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	list, err := dbHelper.GetAllTask(userId)
	fmt.Println(list)
	if err != nil {
		return
	}
	err = utils.EncodeJSONBody(w, list)
	if err != nil {
		return
	}
}

func getUserId(r *http.Request) int {
	claims, _ := r.Context().Value("userInfo").(jwt.MapClaims)
	fmt.Println("claims:", claims)
	userIdStr := fmt.Sprintf("%v", claims["userid"])
	userId, _ := strconv.Atoi(userIdStr)
	return userId
}

//func GetTaskById(w http.ResponseWriter, r *http.Request) {
//	idQuery := r.URL.Query().Get("id")
//	if idQuery == "" {
//		return
//	}
//
//	id, err := strconv.Atoi(idQuery)
//
//	if err != nil || id < 0 {
//		return
//	}
//
//	user, err := dbHelper.GetTaskById(id)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			w.WriteHeader(400)
//			errorMessage := fmt.Sprintf("invalid id")
//			err := utils.EncodeJSONBody(w, errorMessage)
//			if err != nil {
//				return
//			}
//			return
//		}
//		w.WriteHeader(500)
//		errorMessage := fmt.Sprintf("SERVER ERROR %v", err)
//		err := utils.EncodeJSONBody(w, errorMessage)
//		if err != nil {
//			return
//		}
//		return
//	}
//	err = utils.EncodeJSONBody(w, user)
//	if err != nil {
//		return
//	}
//}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var body model.TodoModel

	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to parse request body")
		return
	}
	userId := getUserId(r)
	//date, _ := time.Parse(time.RFC3339Nano, body.DueDate.String())
	date, _ := time.Parse("1/2/2006 15:04:05", "12/8/2015 12:00:00")

	id := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(id)
	if err != nil || taskId < 0 {
		return
	}

	err = dbHelper.UpdateTask(body.Title, body.Description, date, userId, taskId)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to update todo in table")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}{
		Title:       body.Title,
		Description: body.Description,
		DueDate:     body.DueDate,
	})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	//date, _ := time.Parse(time.RFC3339Nano, body.DueDate.String())
	date, _ := time.Parse("1/2/2006 15:04:05", "11/02/2023 14:50:00")

	id := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(id)
	if err != nil || taskId < 0 {
		return
	}

	err = dbHelper.DeleteTask(date, userId, taskId)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to update todo in table")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string
	}{
		"Task Deleted Successfully",
	})
}

func Logout(r *http.Request) {
	// you cannot expire jwt token on demand
}
