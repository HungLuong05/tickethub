package models

import (
	"log"
	"strconv"
	"sync"
	"errors"
	"encoding/json"

	"tickethub.com/ticket/config"
)

var ticketMutex sync.Mutex

type Ticket struct {
	EventId				int64
	NumTickets		int64	
}

func (ticket Ticket) CreateTicket() error {
	eventIdStr := strconv.Itoa(int(ticket.EventId))
	numTicketsStr := strconv.Itoa(int(ticket.NumTickets))

	log.Println("Creating ticket....")
	log.Println("Event ID: ", eventIdStr)
	log.Println("Number of Tickets: ", numTicketsStr)

	err := config.DB.Set("event:" + eventIdStr, numTicketsStr)
	if err != nil {
		return err
	}
	return nil
}

func (ticket Ticket) GetTicket() (int64, error) {
	eventIdStr := strconv.Itoa(int(ticket.EventId))
	// numTicketsStr := strconv.Itoa(int(ticket.numTickets))

	log.Println("Getting ticket....")
	log.Println("Event ID: ", eventIdStr)

	tickets, err := config.DB.Get("event:" + eventIdStr)
	if err != nil {
		return 0, err
	}

	ticketsInt, err := strconv.ParseInt(tickets, 10, 64)
	if err != nil {
		return 0, err
	}
	return ticketsInt, nil
}

func (ticket Ticket) ProtectedUpdateTicket() (bool, error) {
	ticketMutex.Lock()
	defer ticketMutex.Unlock()

	eventIdStr := strconv.Itoa(int(ticket.EventId))

	numTickets, err := config.DB.Get("event:" + eventIdStr)
	if err != nil {
		return false, err
	}

	numTicketsInt, err := strconv.Atoi(numTickets)
	if err != nil {
		return false, err
	}

	if numTicketsInt < int(ticket.NumTickets) {
		return false, nil
	}

	numTicketsStr := strconv.Itoa(numTicketsInt - int(ticket.NumTickets))
	err = config.DB.Set("event:" + eventIdStr, numTicketsStr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ticket Ticket) UpdateTicket() error {
	eventIdStr := strconv.Itoa(int(ticket.EventId))

	log.Println("Updating ticket....")
	log.Println("Event ID: ", eventIdStr)

	accepted, err := ticket.ProtectedUpdateTicket()
	if err != nil {
		return err
	}

	if !accepted {
		return errors.New("Not enough tickets")
	}

	ticketInfo := map[string]interface{}{
		"eventId": ticket.EventId,
		"numTickets": ticket.NumTickets,
	}
	jsonData, err := json.Marshal(ticketInfo)
	if err != nil {
		return err
	}
  log.Println("Ticket Info: ", ticketInfo)
	log.Println("Ticket JSON: ", string(jsonData))

	err = config.Producer.SendMessage("ticket", string(jsonData))
	if err != nil {
		return err
	}
	return nil
}