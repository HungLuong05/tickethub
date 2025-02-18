package routes

import (
	// "context"
	"strconv"
	"net/http"
	"log"
	"tickethub.com/ticket/models"

	"github.com/gin-gonic/gin"
)

func CreateTicket(context *gin.Context) {
	var ticket models.Ticket

	log.Println("Creating ticket")

	err := context.ShouldBindJSON(&ticket)
	if err != nil {
		context.JSON(400, gin.H{"message": "Invalid request" + err.Error()})
		return
	}
	log.Println("Event ID: ", ticket.EventId)
	log.Println("Number of Tickets: ", ticket.NumTickets)

	err = ticket.CreateTicket()

	if err != nil {
		context.JSON(500, gin.H{"message": "Could not create ticket" + err.Error()})
		return
	}

	context.JSON(201, gin.H{"message": "Ticket created successfully"})
}

func UpdateTicket(context *gin.Context) {
	var ticket models.Ticket

	err := context.ShouldBindJSON(&ticket)

	if err != nil {
		context.JSON(400, gin.H{"message": "Invalid request" + err.Error()})
		return
	}

	eventIdStr := context.Param("id")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
  if err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ticket ID: " + eventIdStr})
    return
  }
  ticket.EventId = eventId

	log.Printf("Event ID: %v\n", ticket.EventId)
	log.Printf("Number of Tickets: %v\n", ticket.NumTickets)

	err = ticket.UpdateTicket()

	if err != nil {
		context.JSON(500, gin.H{"message": "Could not update ticket" + err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "Ticket updated successfully"})
}

func GetTicket(context *gin.Context) {
	var ticket models.Ticket

	// err := context.ShouldBindJSON(&ticket)

	// if err != nil {
	// 	context.JSON(400, gin.H{"message": "Invalid request" + err.Error()})
	// 	return
	// }

	eventIdStr := context.Param("id")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
  if err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ticket ID: " + eventIdStr})
    return
  }
  ticket.EventId = eventId
	log.Printf("Event ID: %v\n", eventId)

	tickets, err := ticket.GetTicket()

	if err != nil {
		context.JSON(500, gin.H{"message": "Could not get ticket" + err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "Ticket retrieved successfully", "tickets": tickets})
}