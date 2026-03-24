package main

import (
	"log/slog"

	"github.com/MrEraba/todos-open-api/models"
)

type DB struct {
	Users map[string]*models.User
}

func (db *DB) Insert(u *models.User) {
	db.Users[u.ID] = u
}

func main() {
	slog.Info("application starting")

	ivan := models.NewUser("Ivan", "i.almanza@mail.com")
	luis := models.NewUser("Luis", "l.suarez@mail.com")

	database := DB{Users: make(map[string]*models.User)}
	database.Insert(ivan)
	database.Insert(luis)

	slog.Info("users created", "count", len(database.Users))

	for _, u := range database.Users {
		slog.Info("user", "id", u.ID, "name", u.Name, "email", u.Email, "active", u.Active)
	}

	slog.Info("application finished")
}
