package entity

import (
	"encoding/json"
	"io"
	"social/pkg/log"
	"sort"
	"time"
)

type Response struct {
	Code uint32  `json:"code"`
	Msg  string  `json:"msg"`
	Data Message `json:"data"`
}

type Message struct {
	ID          int64  `json:"id"`
	Content     string `json:"content"`
	MsgType     int16  `json:"msg_type"`     // 1.单聊 2.群聊 3.系统
	ContentType int16  `json:"content_type"` // 1.文字 2.普通文件 3.图片 4.音频 5.视频 6.语音聊天 7.视频聊天
}

func (r *Response) Resp(data Message) []byte {
	r.Code = 0
	r.Msg = "ok"
	r.Data = data
	b, err := json.Marshal(r)
	if err != nil {
		log.Errorf("Resp marshal err: %v", err)
	}
	return b
}

func (r *Response) ErrResp(msg string) []byte {
	r.Code = 1
	r.Msg = msg
	b, err := json.Marshal(r)
	if err != nil {
		log.Errorf("ErrResp marshal err: %v", err)
	}
	return b
}

//feed := &Feed{
//		Version: "https://google.com",
//		Title:   "feed title",
//		Hubs: []*Hub{
//			&Hub{
//				Type: "WebSub",
//				Url:  "https://websub-hub.example",
//			},
//		},
//	}

type Feed struct {
	Version     string  `json:"version"`
	Title       string  `json:"title"`
	HomePageUrl string  `json:"home_page_url,omitempty"`
	FeedUrl     string  `json:"feed_url,omitempty"`
	Description string  `json:"description,omitempty"`
	UserComment string  `json:"user_comment,omitempty"`
	NextUrl     string  `json:"next_url,omitempty"`
	Icon        string  `json:"icon,omitempty"`
	Favicon     string  `json:"favicon,omitempty"`
	Author      *Author `json:"author,omitempty"`
	Expired     *bool   `json:"expired,omitempty"`
	Hubs        []*Hub  `json:"hubs,omitempty"`
	Items       []*Item `json:"items,omitempty"`
}

// Hub 订阅发布者关于这条feed的实时通知
type Hub struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Item struct {
	Id            string       `json:"id"`
	Url           string       `json:"url,omitempty"`
	ExternalUrl   string       `json:"external_url,omitempty"`
	Title         string       `json:"title,omitempty"`
	ContentHTML   string       `json:"content_html,omitempty"`
	ContentText   string       `json:"content_text,omitempty"`
	Summary       string       `json:"summary,omitempty"`
	Image         string       `json:"image,omitempty"`
	BannerImage   string       `json:"banner_,omitempty"`
	PublishedDate *time.Time   `json:"date_published,omitempty"`
	ModifiedDate  *time.Time   `json:"date_modified,omitempty"`
	Author        *Author      `json:"author,omitempty"`
	Tags          []string     `json:"tags,omitempty"`
	Attachments   []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Url      string        `json:"url,omitempty"`
	MIMEType string        `json:"mime_type,omitempty"`
	Title    string        `json:"title,omitempty"`
	Size     string        `json:"size,omitempty"`
	Duration time.Duration `json:"duration_in_seconds,omitempty"`
}

type Author struct {
	Name   string `json:"name,omitempty"`
	Url    string `json:"url,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

func (f *Feed) Add(item *Item) {
	f.Items = append(f.Items, item)
}

// ToJSON creates a JSON Feed representation of this feed
func (f *Feed) ToJSON() (string, error) {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// WriteJSON writes an JSON representation of this feed to the writer.
func (f *Feed) WriteJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(f)
}

// Sort sorts the Items in the feed with the given less function.
func (f *Feed) Sort(less func(a, b *Item) bool) {
	lessFunc := func(i, j int) bool {
		return less(f.Items[i], f.Items[j])
	}
	sort.SliceStable(f.Items, lessFunc)
}
