package model

type User struct {
	ID     int64    `gorm:"primary_key"`
	Source []Source `gorm:"many2many:subscribes;"`
	State  int `gorm:"DEFAULT:0;"`
	EditTime
}

func FindOrInitUser(userID int64) *User {
	db := getConnect()
	defer db.Close()
	var user User
	db.Where(User{ID: userID}).FirstOrCreate(&user)
	return &user
}
