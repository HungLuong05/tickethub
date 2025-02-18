package routes

import (
	"net/http"
	"strconv"
	"log"
	"io"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"tickethub.com/auth/utils"
	"tickethub.com/auth/models"
)

func AddPerm(context *gin.Context) {
	var perm models.Perm

	bodyBytes, err := io.ReadAll(context.Request.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read request body"})
		return
	}

	log.Println("Request Body:", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &perm)
  if err != nil {
    log.Println("Could not parse request data:", err)
    context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
    return
  }

	log.Println("EventId", perm.EventId, "UserId", perm.UserId)

	err = perm.AddPermission()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add permission." + err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Permission added successfully"})
}

func VerifyPerms(context *gin.Context) {
	token := context.GetHeader("Authorization")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Token not found"})
		return
	}

	id, err := utils.VerifyToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized" + err.Error()})
		return
	}

	// context.Header("User-ID", strconv.FormatInt(id, 10))
	// context.JSON(http.StatusOK, gin.H{"message": "Authorized."})

	var perm models.Perm

	// err = context.ShouldBindJSON(&perm)
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
	// 	return
	// }

	perm.EventId, err = strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}
	perm.UserId = id
	log.Println(perm.EventId, perm.UserId)

	err = perm.VerifyPermission()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized." + err.Error()})
		return
	}

	context.Header("User-ID", strconv.FormatInt(perm.UserId, 10))
	context.JSON(http.StatusOK, gin.H{"message": "Authorized."})
}

func RemovePerm(context *gin.Context) {
	var perm models.Perm

	bodyBytes, err := io.ReadAll(context.Request.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read request body"})
		return
	}

	log.Println("Request Body:", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &perm)
  if err != nil {
    log.Println("Could not parse request data:", err)
    context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
    return
  }

	log.Println("EventId", perm.EventId, "UserId", perm.UserId)

	err = perm.RemovePermission()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not remove permission."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Permission removed successfully"})
}