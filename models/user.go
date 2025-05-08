package models

type User struct {
	ID        int64
	Username  string
	FirstName string
	LastName  string

	Email              string
	EmailNotifications bool
	Password           string

	CreatesUnix   int64
	UpdatedUnix   int64
	LastLoginUnix int64

	IsAdmin bool
}
