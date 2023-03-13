package handler

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"todo/database"
	"todo/database/dbHelper"
	"todo/middleware"
	"todo/model"
	"todo/utils"
)

func RegisterUser(ctx *gin.Context) {
	var body model.RegisterModel
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Failed to parse Request body"})
		return
	}
	if len(body.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "password must be 6 chars long"})
		return
	}

	exist, existError := dbHelper.IsUserExist(body.Email)

	if existError != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": existError.Error(), "message": "failed to check user existence"})
		return
	}

	if exist {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "user already exist"})
		return
	}

	hashPassword, hasErr := utils.HashPassword(body.Password)

	if hasErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": hasErr.Error(), "message": "failed to secure password"})
		return
	}

	err := dbHelper.CreateUser(database.Todo, body.Name, body.Email, hashPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": existError.Error(), "message": "Failed to create user in table"})
		return
	}

	ctx.JSON(http.StatusOK, struct {
		Message string `json:"message"`
	}{"User created successfully"})
}

func LoginUser(ctx *gin.Context) {
	var body model.LoginModel
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Failed to Parse Request Body"})
		return
	}
	userId, err := dbHelper.GetUserIDByPassword(body.Email, body.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "user does not exist"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "failed to find user"})
		return
	}

	token, err := middleware.GenerateJWT(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "error in generating jwt token"})
		return
	}
	ctx.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func CreateTask(ctx *gin.Context) {
	var body model.TodoModel

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to parse request body"})
		return
	}
	userId := getUserId(ctx)
	date, err := time.Parse(time.RFC3339, body.DueDate)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dbHelper.CreateTodo(database.Todo, userId, body.Title, body.Description, date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Failed to create todo in table"})
		return
	}
	ctx.JSON(http.StatusOK, struct {
		Message string `json:"message"`
	}{"Todo  created successfully"})
}

func GetAllTask(ctx *gin.Context) {
	userId := getUserId(ctx)
	list, err := dbHelper.GetAllTask(userId)
	fmt.Println(list)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusCreated, list)
}

func getUserId(c *gin.Context) int {
	claims, _ := c.Get("userInfo")
	userInfo, _ := claims.(jwt.MapClaims)
	fmt.Println("userInfo:", userInfo)
	userIdStr := fmt.Sprintf("%v", userInfo["userid"])
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

func UpdateUser(ctx *gin.Context) {
	var body model.TodoModel

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Failed to parse request body"})
		return
	}
	userId := getUserId(ctx)
	date, _ := time.Parse("1/2/2006 15:04:05", "12/8/2015 12:00:00")

	id := ctx.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil || taskId < 0 {
		return
	}

	err = dbHelper.UpdateTask(body.Title, body.Description, date, userId, taskId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Failed to update todo in table"})
		return
	}

	ctx.JSON(http.StatusCreated, struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}{
		Title:       body.Title,
		Description: body.Description,
		DueDate:     body.DueDate,
	})
}

func DeleteTask(ctx *gin.Context) {
	userId := getUserId(ctx)
	date, _ := time.Parse("1/2/2006 15:04:05", "11/02/2023 14:50:00")

	id := ctx.Param("id")

	taskId, err := strconv.Atoi(id)
	if err != nil || taskId < 0 {
		return
	}

	err = dbHelper.DeleteTask(date, userId, taskId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Failed to update todo in table}"})
		return
	}

	ctx.JSON(http.StatusCreated, struct {
		Message string
	}{
		"Task Deleted Successfully",
	})
}
