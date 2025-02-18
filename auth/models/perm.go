package models

import (
	"log"
	"errors"

	"tickethub.com/auth/config"
)

type Perm struct {
	EventId int64 `json:"EventId" binding:"required"`
	UserId  int64 `json:"UserId" binding:"required"`
}

func (perm Perm) AddPermission() error {
	query := "INSERT INTO event_perms(event_id, user_id) VALUES ($1, $2)"	
	args := []interface{}{perm.EventId, perm.UserId}

	err := pg.DB.Exec(query, args...)
	if err != nil {
		log.Println("Could not add permission: ", err)
		return errors.New("could not add permission" + err.Error())
	}
	return nil
}

func (perm Perm) RemovePermission() error {
	query := "DELETE FROM event_perms WHERE event_id = $1 AND user_id = $2"
	args := []interface{}{perm.EventId, perm.UserId}

	err := pg.DB.Exec(query, args...)
	if err != nil {
		log.Println("Could not remove permission: ", err)
		return errors.New("could not remove permission" + err.Error())
	}
	return nil
}

func (perm Perm) VerifyPermission() error {
	query := "SELECT * FROM event_perms WHERE event_id = $1 AND user_id = $2"
	args := []interface{}{perm.EventId, perm.UserId}

	rows, err := pg.DB.Query(query, args...)
	if err != nil {
		log.Println("Could not verify permission: ", err)
		return errors.New("could not verify permission" + err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		return errors.New("permission not found")
	}

	return nil
}