package entity

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Salt        string `json:"salt"`
	Avatar      string `json:"avatar"`
	Summary     string `json:"summary"`
	AccessToken string `json:"access_token"`
	ExpireTime  int64  `json:"expire_time"`
	CreatedTime int64  `json:"created_time"`
	UpdatedTime int64  `json:"updated_time"`
}
