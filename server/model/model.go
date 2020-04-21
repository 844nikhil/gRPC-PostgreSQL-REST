package model

//DBConfig has information required to connect to DB
type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

type User struct {
	UserID       int    `gorm:"column:id"`
	Name         string `gorm:"column:name"`
	Email        string `gorm:"column:email"`
	PhoneNo      string `gorm:"column:phoneno"`
	Password     string `gorm:"column:password"`
	Organisation string `gorm:"column:organization"`
	IsDeleted bool `gorm:"column:IsDeleted"`
}
