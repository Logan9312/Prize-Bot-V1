package events

import "gitlab.com/logan9312/discord-auction-bot/database"

func EventSave(event database.Event) error {
	return database.DB.Create(event).Error
}

func Update(event any) error {
	return database.DB.Updates(event).Error
}

func EventCreate(event database.Event) error {
	return database.DB.Create(event).Error
}
