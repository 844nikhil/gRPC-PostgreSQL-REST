package database

import (
	"errors"
	"gRPC-PostgreSQL-REST/server/model"
)

func (dc *DBRepo) CreateUser(user model.User) error {
	user.IsDeleted = false
	rows := dc.GormDB.Debug().Create(&user).RowsAffected
	if rows == 0 {
		return errors.New("DB Error")
	}
	return nil
}

//Check if the email already exists
func (dc *DBRepo) CheckUser(email string) error {
	rows := dc.GormDB.Debug().Exec(`SELECT * FROM "users" where "email"=? and "IsDeleted" = ?`, email, false).RowsAffected
	if rows == 0 {
		return nil
	}
	return errors.New("Email already exists")
}

//Authenticate user with email and password
func (dc *DBRepo) Authenticate(email, password string) (error, int) {
	var user model.User
	rows := dc.GormDB.Debug().Where(`"email"=? and "password" =? and "IsDeleted" = ?`, email, password, false).First(&user).RowsAffected
	if rows == 1 {
		return nil, user.UserID
	}
	return errors.New("Invalid user"), 0

}

func (dc *DBRepo) GetUser(userID int) (model.User, error) {
	var user model.User
	rows := dc.GormDB.Debug().Where(`"id"=? and "IsDeleted" =? `, userID, false).First(&user).RowsAffected
	if rows == 1 {
		return user, nil
	}
	return user, errors.New("Invalid user")
}

func (dc *DBRepo) DeleteUser(userID int) error {
	rows := dc.GormDB.Debug().Table("users").Where(`"id" = ? and "IsDeleted" = ?`, userID, false).Updates(map[string]interface{}{"IsDeleted": true}).RowsAffected
	if rows == 0 {
		return errors.New("Couldn't delete user")
	}
	return nil
}

func (dc *DBRepo) UpdateUser(phoneNo, organisation string, userID int) error {
	rows := dc.GormDB.Debug().Table("users").Where(`"id" = ? and "IsDeleted" = ?`, userID, false).Updates(map[string]interface{}{"phoneno": phoneNo, "organization": organisation}).RowsAffected
	if rows == 0 {
		return errors.New("Couldn't update user")
	}
	return nil
}
