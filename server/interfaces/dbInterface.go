package interfaces

import (
	"gRPC-PostgreSQL-REST/server/model"
)

//DBEngine is used to call DB Methods from handler
var DBEngine DBInterface

//DBInterface contains all the DB methods
type DBInterface interface {
	DBConnect(model.DBConfig) error
	CreateUser(model.User) error
	CheckUser(string) error
	Authenticate(string, string) (error, int)
	GetUser(int) (model.User, error)
	DeleteUser(int) error
	UpdateUser(string, string, int) error
}
