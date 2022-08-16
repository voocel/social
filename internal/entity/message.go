package entity

type Message struct {
	ID          int64  `json:"id"`
	Content     string `json:"content"`
	MsgType     int16  `json:"msg_type"`     // 1.单聊 2.群聊
	ContentType int16  `json:"content_type"` // 1.文字 2.普通文件 3.图片 4.音频 5.视频 6.语音聊天 7.视频聊天
	Pic         string `json:"pic"`
}

type Response struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"result"`
}
