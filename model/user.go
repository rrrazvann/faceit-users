package model

type User struct {
	ID        uint   `gorm:"primaryKey"                             json:"id"`
	
	Email     string `gorm:"unique;type:VARCHAR(255)"               json:"email"`
	FirstName string `gorm:"type:VARCHAR(255)"                      json:"first_name"`
	LastName  string `gorm:"type:VARCHAR(255)"                      json:"last_name"`
	Nickname  string `gorm:"type:VARCHAR(255)"                      json:"nickname"`
	Password  string `gorm:"type:VARCHAR(255)"                      json:"password"`
	Country   string `gorm:"type:VARCHAR(255)"                      json:"country"`
	
	Timestampable
}
