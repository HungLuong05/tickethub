package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"tickethub.com/auth/config"
	"tickethub.com/auth/utils"
)

type User struct {
	Id							int64
	Email						string	`binding:"required"`
	Password 				string 	`binding:"required"`
	Salt						string
}

func GenerateRandomSalt(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
			return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (user User) CreateAccount() error {
	salt, err := GenerateRandomSalt(16)
	if err != nil {
		log.Println("Could not generate salt: ", err)
		return errors.New("could not generate salt")
	}
	user.Salt = salt

	hashedPassword, err := utils.HashPassword(user.Password + user.Salt)
	if err != nil {
		log.Println("Could not hash password: ", err)
		return errors.New("could not hash password")
	}

	fmt.Println(user.Email, hashedPassword, user.Salt)
	query := "INSERT INTO users(email, password, salt) VALUES ($1, $2, $3)"
	args := []interface{}{user.Email, hashedPassword, user.Salt}

	err = pg.DB.Exec(query, args...)
	if err != nil {
		log.Println("Could not create user: ", err)
		return errors.New("could not create user")
	}
	fmt.Printf("User created successfully: %v %v\n", query, args)
	return nil
}

func (user User) ValidateCredentials() error {
	all, err := pg.DB.Query("SELECT * FROM users")
	if all == nil || !all.Next() || err != nil {
		return errors.New("no users found")
	} else {
		var id int64
		var email, hashedPassword, salt string
		all.Scan(&id, &email, &hashedPassword, &salt)
		fmt.Println("User found: ", id, email, hashedPassword, salt)
	}
	fmt.Println(all)
	all.Close()

	// query := "SELECT * FROM users WHERE email = $1"
	// row, err := pg.DB.Query(query, user.Email)
	// query := "SELECT * FROM users"
	row, err := pg.DB.Query("SELECT id, password, salt FROM users WHERE email = $1", user.Email)
	if err != nil {
		log.Println("Could not query database: ", err)
		return errors.New("could not query database")
	}
	row.Next()

	var id int64
	var hashedPassword, salt string
	err = row.Scan(&id, &hashedPassword, &salt)
	if err != nil {
		log.Println("Could not scan row: ", err)
		return errors.New("could not scan row")
	}
	row.Close()

	fmt.Println("Check password hash", user.Password, hashedPassword, salt)
	passwordIsValid := utils.CheckPasswordHash(user.Password, hashedPassword, salt)

	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}

func (user User) UpdateEventPerm (eventId int64) error {
	query := "INSERT INTO event_perm(user_id, event_id) VALUES ($1, $2)"
	args := []interface{}{user.Id, eventId}

	err := pg.DB.Exec(query, args...)
	if err != nil {
		log.Println("Could not update event permission: ", err)
		return errors.New("could not update event permission")
	}
	fmt.Printf("Event permission updated successfully: %v %v\n", query, args)
	return nil
}