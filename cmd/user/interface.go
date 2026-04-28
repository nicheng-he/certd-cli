package user

import (
	"database/sql"
	"time"
)

var (
	limit  int
	offset int
)

type User struct {
	ID          int            `json:"id"`
	UserName    string         `json:"user_name"`
	NickName    string         `json:"nick_name"`
	Email       sql.NullString `json:"email"`
	Remark      sql.NullString `json:"remark"`
	CreatedTime time.Time      `json:"created_time"`
}
