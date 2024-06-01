package models

// Photo model
type Photo struct {
	Id       int64   `gorm:"primaryKey" json:"id"`
	Title    string  `gorm:"type:varchar(255);not null" json:"title"`
	Caption  string  `gorm:"type:varchar(255)" json:"caption"`
	PhotoURL *string `gorm:"type:varchar(255)" json:"photo_url"`
	UserID   int64   `gorm:"not null" json:"user_id"`
	User     User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
}
