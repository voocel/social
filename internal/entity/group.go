package entity

type Group struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Owner       uint64 `json:"owner"`
	Notice      string `json:"notice"`
	CreatedTime string `json:"created_time"`
}
