package entity

import "time"

type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Mobile        string    `json:"mobile"`
	Nickname      string    `json:"nickname"`
	Email         string    `json:"email"`
	Avatar        string    `json:"avatar"`
	Summary       string    `json:"summary,omitempty"`
	Sex           int8      `json:"sex"`
	Status        int8      `json:"status"`
	Birthday      time.Time `json:"-"`
	LastLoginTime time.Time `json:"-"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	DeletedAt     time.Time `json:"-"`
}

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//Mobile   string `json:"mobile,omitempty"`
	//Nickname string `json:"nickname,omitempty"`
}
