package models

import (
	"log"
	"encoding/json"

	"tickethub.com/backup/pg"
)

type Ticket struct {
	EventId					int64
	NumTickets 			int64		
}

func (ticket Ticket) UpdateTicket() error {
	query := "SELECT ticket FROM events WHERE id=$1"
	args := []interface{}{ticket.EventId}
	rows, err := pg.DB.Query(query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var totalTicket int64
		err = rows.Scan(&totalTicket)
		if err != nil {
			return err
		}

		log.Println("Number of tickets: ", totalTicket)

		query = "UPDATE events SET ticket=$1 WHERE id=$2"
		args = []interface{}{totalTicket - ticket.NumTickets, ticket.EventId}
		err = pg.DB.Exec(query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ticket *Ticket) ParseTicketString(jsonData string) error {
	err := json.Unmarshal([]byte(jsonData), &ticket)
	if err != nil {
		log.Println("Error parsing ticket string: ", err)
		return err
	}
	return nil
}