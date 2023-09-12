package entity

type Friend struct {
	Id          int64  `json:"id"`
	Uid         int64  `json:"uid"`
	FriendId    int64  `json:"friend_id"`
	Remark      string `json:"remark"`
	Shield      int8   `json:"shield"`
	CreatedTime int64  `json:"created_time"`
}

type FriendResp struct {
	Uid      int64  `json:"uid"`
	FriendId int64  `json:"friend_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Remark   string `json:"remark"`
	Shield   int8   `json:"shield"`
}

type FriendApply struct {
	Id          int64  `json:"id"`
	FromId      int64  `json:"from_id"`
	ToId        int64  `json:"to_id"`
	Remark      string `json:"remark"`
	Status      uint8  `json:"status"`
	CreatedTime int64  `json:"created_time"`
}

type FriendApplyReq struct {
	Uid       int64  `json:"uid"`
	FriendId  int64  `json:"friend_id"`
	ApplyInfo string `json:"apply_info"`
}
