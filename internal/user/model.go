package user

const (
	USER_PUBLIC_PREFIX = "usr_"
	USER_ID_PREFIX = "USER"
	USER_TABLE_NAME = "user"

	USER_ID_CONTEXT_KEY = "userId"
)

type User struct{
	Id 				string	`gorm:"string(255);not null" json:"id" tag_name:"id"`
	Username		string	`gorm:"string(255);not null" json:"username" tag_name:"username"`
	Email			string	`gorm:"string(14);not null" sql:"DEFAULT:null" json:"email" tag_name:"email"`
	IsActive		int		`gorm:"integer(1)" sql:"DEFAULT:1" json:"is_active" tag_name:"is_active"`
	SecretKey		string	`gorm:"string(14);not null" json:"secret_key" tag_name:"secret_key"`
	CreatedAt		int64	`gorm:"integer(11)" sql:"DEFAULT:null" json:"created_at" tag_name:"created_at"`
	DeletedAt		int64	`gorm:"integer(11)" sql:"DEFAULT:null" json:"deleted_at" tag_name:"deleted_at"`
	SecretKeyUpdatedAt		int64	`gorm:"integer(11)" sql:"DEFAULT:null" json:"secret_key_updated_at" tag_name:"secret_key_updated_at"`
}
