package entity

type Friend struct {
	Id          int64  `json:"id"`
	Uid         int64  `json:"uid"`
	FriendId    int64  `json:"friend_id"`
	Remark      string `json:"remark"`
	Shield      uint8  `json:"shield"`
	CreatedTime int64  `json:"created_time"`
}
