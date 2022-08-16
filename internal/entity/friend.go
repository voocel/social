package entity

type Friend struct {
	Id          int    `json:"id"`
	Uid         int    `json:"uid"`
	FriendId    int    `json:"friend_id"`
	Remark      string `json:"remark"`
	Shield      uint8  `json:"shield"`
	CreatedTime int64  `json:"created_time"`
}
