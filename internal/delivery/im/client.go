package im

import (
	"encoding/json"
	"fmt"
	"social/internal/entity"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// 用户登录
type login struct {
	AppID  uint32
	UserID uint64
	Client *Client
}

// // 用户连接
// type connect struct {
// 	AppID uint32
// 	Uid string
// 	Client *Client
// }

// Client 用户连接
type Client struct {
	ID        string          //ID
	Addr      string          //用户地址
	Socket    *websocket.Conn // 用户连接
	Send      chan []byte     //待发送的数据
	AppID     uint32          // 登录的平台Id app/web/ios
	UserID    uint64          // 用户Id，用户登录以后才有
	FirstTime uint64          // 首次连接事件
	LoginTime uint64          // 登录时间 登录以后才有
}

//NewClient 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	socket.WriteJSON("成功连接到服务器~")
	client = &Client{
		ID:        uuid.New().String(),
		Addr:      addr,
		Socket:    socket,
		Send:      make(chan []byte, 100),
		FirstTime: firstTime,
	}
	return
}

// DisposeFunc func
type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// 读取客户端数据
func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", r)
		}
		fmt.Println("读取客户端数据 关闭send")
		c.Close() //读取完毕后注销该client
		close(c.Send)
	}()
	for {
		c.Socket.SetReadDeadline(time.Now().Add(time.Duration(viper.GetInt("app.heartbeat")) * time.Second))
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			fmt.Println("读取客户端数据 错误", c.Addr, err)
			// c.Close()
			break
		}
		ProcessData(c, message)
	}
}

// 向客户端写数据
func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", r)
		}
		c.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			//如果没有消息
			if !ok {
				// 客户端发送数据错误 关闭连接
				fmt.Println("没有消息了", message)
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//有消息就发送给客户端
			// c.Socket.WriteMessage(websocket.TextMessage, []byte("你好啊!"))
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// Close 关闭连接
func (c *Client) Close() {
	fmt.Println("userid: ", c.UserID)
	clientManager.Unregister <- c
	c.Socket.Close()
}

// Request 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`   // 消息的唯一Id
	Cmd  string      `json:"cmd"`   // 请求命令
	Data interface{} `json:"param"` // 数据 json
}

// ProcessData 处理数据
func ProcessData(client *Client, message []byte) {
	request := &Request{}
	err := json.Unmarshal(message, request)
	if err != nil {
		fmt.Println("处理数据 json Unmarshal", err)
		client.Socket.WriteMessage(websocket.TextMessage, []byte("消息格式错误!"))
		return
	}
	// requestData是[]byte类型
	requestData, err := json.Marshal(request.Data)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
	}
	// seq := request.Seq
	seq := "1234567"
	cmd := request.Cmd

	// 采用 map 注册的方式
	code, msg, data := handleRouter(cmd, client, seq, requestData)
	response := entity.Response{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	responseByte, err := json.Marshal(response)

	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		return
	}
	fmt.Printf("响应数据: %v\n", string(responseByte))
	if code != 0 {
		client.SendMsg(responseByte)
	}
	return
}

// SetUser 用户登录
func (c *Client) SetUser(appID uint32, userID uint64, loginTime uint64) {
	c.AppID = appID
	c.UserID = userID
	c.LoginTime = loginTime
}

// Router 路由注册
func Router(cmd string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[cmd] = value

	return
}

func getHandlers(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]

	return
}

// GetUserKey 获取key
func GetUserKey(AppID uint32, UserID uint64) (key string) {
	key = fmt.Sprintf("%d_%d", AppID, UserID)
	return
}

// 路由处理
func handleRouter(cmd string, client *Client, seq string, requestData []byte) (code uint32, msg string, data interface{}) {

	if value, ok := getHandlers(cmd); ok {
		// 执行action方法
		code, msg, data = value(client, seq, requestData)
	} else {
		code = 400
		msg = "路由不存在"
		fmt.Println("处理数据 路由不存在", client.Addr, "cmd", cmd)
	}
	return
}

// SendMsg 发送给客户端数据
func (c *Client) SendMsg(msg []byte) {
	// fmt.Println(c)
	if c == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r)
		}
	}()

	c.Send <- msg
}
