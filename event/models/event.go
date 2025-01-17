package models

import (
	"tickethub.com/event/config"
)

type Event struct {
	Id							int64
	Name						string	`binding:"required"`
	Location				string	`binding:"required"`
	Date						string	`binding:"required"`
	Ticket 					int64		`binding:"required"`
}

func (event Event) CreateEvent() error {
	query := "INSERT INTO events(name, location, date, ticket) VALUES ($1, $2, $3, $4)"
	args := []interface{}{event.Name, event.Location, event.Date, event.Ticket}

	err := pg.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (event Event) UpdateEvent() error {
	query := "UPDATE events SET name=$1, location=$2, date=$3, ticket=$4 WHERE id=$5"
	args := []interface{}{event.Name, event.Location, event.Date, event.Ticket, event.Id}

	err := pg.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (event Event) DeleteEvent() error {
	query := "DELETE FROM events WHERE id=$1"
	args := []interface{}{event.Id}

	err := pg.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func GetEventById(id int64) (Event, error) {
	query := "SELECT * FROM events WHERE id=$1"
	args := []interface{}{id}

	rows, err := pg.DB.Query(query, args...)
	if err != nil {
		return Event{}, err
	}

	var event Event
	for rows.Next() {
		err = rows.Scan(&event.Id, &event.Name, &event.Location, &event.Date, &event.Ticket)
		if err != nil {
			return Event{}, err
		}
	}

	return event, nil
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"

	rows, err := pg.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var events []Event
	for rows.Next() {
		var event Event
		err = rows.Scan(&event.Id, &event.Name, &event.Location, &event.Date, &event.Ticket)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}