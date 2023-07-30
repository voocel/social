package entity

type Group struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Owner       int64  `json:"owner"`
	Notice      string `json:"notice"`
	CreatedTime string `json:"created_time"`
}
