package entity

type Group struct {
	Id           int64  `json:"id,omitempty"`
	Name         string `json:"name"`
	Owner        int64  `json:"owner"`
	Avatar       string `json:"avatar"`
	Notice       string `json:"notice"`
	MaxMembers   int    `json:"max_members"`
	Introduction string `json:"introduction"`
	CreatedUid   int64  `json:"created_uid"`
}

type JoinGroupReq struct {
	GroupId int64  `json:"group_id"`
	Remark  string `json:"remark"`
}
