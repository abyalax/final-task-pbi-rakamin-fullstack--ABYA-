package models

type User struct {
	Id          int64   `gorm:"primaryKey" json:"id"`
	Username    string  `gorm:"type:varchar(300);not null" json:"username"`
	Email       string  `gorm:"type:varchar(300);unique;not null" json:"email"`
	Password    string  `gorm:"type:varchar(300);not null" json:"password"`
	NamaLengkap string  `gorm:"type:varchar(300)" json:"nama_lengkap"`
	Photos      []Photo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"photos"`
}

// CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
// UpdatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`

/*
Mohon maaf saya kesulitan pada bagian timestamp di Golang, saya mendapati error seperti ini
"
	sql: Scan error on column index 5, name "created_at": unsupported Scan, storing driver.Value type []uint8 into type *time.Time
"
jadi saya mangatus field created at dan update at secara manual didatabase
*/
