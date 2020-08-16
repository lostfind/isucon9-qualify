package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

var allUsersByID map[int64]User
var allUsersByAccountName map[string]User
var simpleUsersByID map[int64]UserSimple

func loadUsers() {
	var allUsers []User
	err := dbx.Select(&allUsers, "SELECT * FROM `users`")
	if err != nil {
		log.Fatal(err)
	}

	simpleUsersByID = make(map[int64]UserSimple)
	allUsersByID = make(map[int64]User)
	allUsersByAccountName = make(map[string]User)

	for _, u := range allUsers {
		simpleUsersByID[u.ID] = UserSimple{
			ID:           u.ID,
			AccountName:  u.AccountName,
			NumSellItems: u.NumSellItems,
		}

		allUsersByID[u.ID] = u
		allUsersByAccountName[u.AccountName] = u
	}
}

func getUserByAccountName(accountName string) (user User, err error) {
	user, ok := allUsersByAccountName[accountName]
	if !ok {
		return user, sql.ErrNoRows
	}
	return
}

func getUserByID(id int64) (user User, err error) {
	user, ok := allUsersByID[id]
	if !ok {
		return user, sql.ErrNoRows
	}
	return
}

func getUserSimpleByID(q sqlx.Queryer, userID int64) (userSimple UserSimple, err error) {
	userSimple, ok := simpleUsersByID[userID]
	if !ok {
		return UserSimple{}, sql.ErrNoRows
	}
	return
}

func updateUser(tx *sqlx.Tx, u User, now time.Time) error {
	numSellItems := u.NumSellItems + 1
	_, err := tx.Exec("UPDATE `users` SET `num_sell_items`=?, `last_bump`=? WHERE `id`=?",
		numSellItems,
		now,
		u.ID,
	)

	if err == nil {
		simpleUser := simpleUsersByID[u.ID]
		simpleUser.NumSellItems = numSellItems
		simpleUsersByID[u.ID] = simpleUser

		user := allUsersByID[u.ID]
		user.NumSellItems = numSellItems
		allUsersByID[u.ID] = user
		allUsersByAccountName[u.AccountName] = user
	}

	return err
}

func insertUser(accountName string, hashedPassword []byte, address string) (id int64, err error) {
	result, err := dbx.Exec("INSERT INTO `users` (`account_name`, `hashed_password`, `address`) VALUES (?, ?, ?)",
		accountName,
		hashedPassword,
		address,
	)

	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()

	if err == nil {
		user := User{
			ID:             id,
			AccountName:    accountName,
			HashedPassword: hashedPassword,
			Address:        address,
		}

		userSimple := UserSimple{
			ID:          id,
			AccountName: accountName,
		}

		allUsersByID[id] = user
		allUsersByAccountName[accountName] = user
		simpleUsersByID[id] = userSimple
	}

	return
}
