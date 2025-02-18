package routes

import (
	"log"
	"net/http"
	"strconv"
	"encoding/json"
	"bytes"
	"io"

	"tickethub.com/event/proto"
	"tickethub.com/event/models"
	"github.com/gin-gonic/gin"
)

func CreateEvent (context *gin.Context, grpcClient proto.EventPermClient) {
	var event models.Event

	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = event.CreateEvent()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save event."})
		return
	}

	log.Println(context.Request.Header)
	userIdStr := context.GetHeader("X-User-ID")
  userId, err := strconv.ParseInt(userIdStr, 10, 64)
  if err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID: " + userIdStr})
    return
  }
	
	// res, err := grpcClient.AddEventPerm(context.Request.Context(), &proto.AddEventPermRequest{UserId: userId, EventId: event.Id})
  // log.Println("Bug in sending grpc request??" + err.Error())
	// if err != nil {
  //   context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update event permission" + err.Error()})
  //   return
  // } else {
	// 	log.Println(res.Message)
	// }


	// perm := map[string]string{
	// 	"EventId": strconv.FormatInt(event.Id, 10),
	// 	"UserId":  strconv.FormatInt(userId, 10),
	// }
	perm := map[string]int64{
		"EventId": event.Id,
		"UserId":  userId,
	}
	log.Println("EventId ", event.Id, " UserId ", userId)
	log.Println("EventId " + strconv.FormatInt(event.Id, 10) + " UserId " + strconv.FormatInt(userId, 10))
	permJSON, err := json.Marshal(perm)
	log.Println("PermJSON: ", permJSON)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to marshal JSON"})
		return
	}

	// resp, err := http.Post("http://auth-service:8000/perm", "application/json", bytes.NewBuffer(permJSON))
	// if err != nil {
	// 	context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send request to auth service" + err.Error()})
	// 	return
	// }

	req, err := http.NewRequest("POST", "http://auth-service:8000/perm", bytes.NewBuffer(permJSON))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create request to auth service" + err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send request to auth service" + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		context.JSON(resp.StatusCode, gin.H{"message": "Failed to add permission" + string(bodyBytes)})
		return
	}

	ticket := map[string]int64{
		"EventId": event.Id,
		"NumTickets": event.Ticket,
	}
	log.Println("EventId", event.Id)
	log.Println("NumTickets", event.Ticket)
	ticketJSON, err := json.Marshal(ticket)
	log.Println("TicketJSON: ", ticketJSON)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to marshal JSON"})
		return
	}
	req, err = http.NewRequest("POST", "http://ticket-service:8002/api/ticket", bytes.NewBuffer(ticketJSON))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create request to ticket service" + err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send request to ticket service" + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		context.JSON(resp.StatusCode, gin.H{"message": "Failed to add ticket" + string(bodyBytes)})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
    "message": "Event created successfully",
  })
}

func GetEvents (context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events."})
		return
	}

	context.JSON(200, gin.H{
		"message": "Get events successfully",
		"events": events,
	})
}

func GetEvent (context *gin.Context) {
	eventIdStr := context.Param("id")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID" + eventIdStr})
		return
	}
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Get event successfully",
		"event": event,	
	})
}

func UpdateEvent (context *gin.Context) {
	var event models.Event

	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	eventIdStr := context.Param("id")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID" + eventIdStr})
		return
	}
	event.Id = eventId

	err = event.UpdateEvent()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event."})
		return
	}

	context.JSON(200, gin.H{
		"message": "Update event successfully",
	})
}

func DeleteEvent (context *gin.Context, grpcClient proto.EventPermClient) {
	var event models.Event

	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	eventIdStr := context.Param("id")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID" + eventIdStr})
		return
	}
	event.Id = eventId

	err = event.DeleteEvent()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete event."})
		return
	}

	log.Println(context.Request.Header)
	userIdStr := context.GetHeader("X-User-ID")
  userId, err := strconv.ParseInt(userIdStr, 10, 64)
  if err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID: " + userIdStr})
    return
  }

	perm := map[string]int64{
		"EventId": event.Id,
		"UserId":  userId,
	}
	permJSON, err := json.Marshal(perm)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to marshal JSON"})
		return
	}

	req, err := http.NewRequest("DELETE", "http://auth-service:8000/perm", bytes.NewBuffer(permJSON))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create request to auth service" + err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send request to auth service" + err.Error()})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		context.JSON(resp.StatusCode, gin.H{"message": "Failed to add permission" + string(bodyBytes)})
		return
	}

	context.JSON(200, gin.H{
		"message": "Delete event successfully",
	})
}