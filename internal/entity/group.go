package entity

type Group struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Owner        int64  `json:"owner"`
	Notice       string `json:"notice"`
	Introduction string `json:"introduction"`
	CreatedUid   int64  `json:"created_uid"`
}

type JoinGroupReq struct {
	GroupId int64  `json:"group_id"`
	Remark  string `json:"remark"`
}
