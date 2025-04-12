package model

type ERole string

const (
	ERoleAdmin ERole = "admin"
	ERoleStaff ERole = "staff"
)

func ERoleFromString(role string) ERole {
	switch role {
	case "admin":
		return ERoleAdmin
	case "staff":
		return ERoleStaff
	default:
		return ERoleStaff
	}
}

func (e ERole) String() string {
	return string(e)
}

type User struct {
	Username string `gorm:"column:username;primaryKey"`
	Password string `gorm:"column:password"`
	Role     ERole  `gorm:"column:role"`
}

func (User) TableName() string {
	return "user"
}
