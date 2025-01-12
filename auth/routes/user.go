package routes

import (
	"net/http"
	"fmt"

	"tickethub.com/auth/models"
	"tickethub.com/auth/utils"
	"github.com/gin-gonic/gin"
)

func Register(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.CreateAccount()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func Login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
	}

	fmt.Println(user.Email, user.Password)
	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials."})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token, "message": "Login successful."})
}

func Verify(context *gin.Context) {
	token := context.GetHeader("Authorization")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized."})
		return
	}

	id, err := utils.VerifyToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": id, "message": "Authorized."})
}