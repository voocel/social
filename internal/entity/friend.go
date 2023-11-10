package entity

import (
	"time"
)

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
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Remark   string `json:"remark"`
	Shield   int8   `json:"shield"`
}

type FriendApply struct {
	Id     int64  `json:"id"`
	FromId int64  `json:"from_id"`
	ToId   int64  `json:"to_id"`
	Remark string `json:"remark"`
	Status uint8  `json:"status"`
}

type FriendApplyResp struct {
	Id         int64     `json:"id"`
	FromId     int64     `json:"from_id"`
	FromName   string    `json:"from_name"`
	FromAvatar string    `json:"from_avatar"`
	ToId       int64     `json:"to_id"`
	Remark     string    `json:"remark"`
	Status     int8      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type FriendApplyReq struct {
	Uid       int64  `json:"uid"`
	FriendId  int64  `json:"friend_id"`
	ApplyInfo string `json:"apply_info"`
}
