package im

import (
	"encoding/json"
	"fmt"
	"social/internal/entity"
	"strconv"
	"sync"
	"time"
)

var clientManager = NewClientManager()

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // appId+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	LoginInfo   chan *login        // 用户登录处理
	Unregister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播 向全部成员发送数据
}

// NewClientManager 创建
func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		LoginInfo:  make(chan *login, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}

	return
}

// 管道事件监听
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)

		case loginInfo := <-manager.LoginInfo:
			// 用户登录
			manager.EventLogin(loginInfo)

		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)

			// case message := <-manager.Broadcast:
			// 	// 广播事件
			// 	clients := manager.GetUserClients()
			// 	for conn := range clients {
			// 		select {
			// 		case conn.Send <- message:
			// 		default:
			// 			close(conn.Send)
			// 		}
			// 	}
		}
	}
}

// EventRegister 用户建立连接事件
func (manager *ClientManager) EventRegister(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
	fmt.Println("EventRegister 用户建立连接", client.Addr)

	// client.Send <- []byte("连接成功")
}

// EventLogin 用户登录
func (manager *ClientManager) EventLogin(loginInfo *login) {

	client := loginInfo.Client
	// 连接存在，在添加
	if manager.InClient(client) {
		userKey := GetUserKey(loginInfo.AppID, loginInfo.UserID)
		manager.AddUsers(userKey, client)
	}

	fmt.Println("EventLogin 用户登录", client.Addr, loginInfo.AppID, loginInfo.UserID)
	// manager.sendAll("有新用户:"+strconv.FormatUint(loginInfo.UserID, 10)+"加入")
}

// EventUnregister 用户断开连接
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}

	// 删除用户连接
	deleteResult := manager.DelUsers(client)

	if deleteResult == false {
		// 不是当前连接的客户端
		return
	}

	// // 清除redis登录数据
	// userOnline, err := cache.GetUserOnlineInfo(client.GetUserKey())
	// if err == nil {
	// 	userOnline.LogOut()
	// 	cache.SetUserOnlineInfo(client.GetUserKey(), userOnline)
	// }

	// close(client.Send)

	fmt.Println("EventUnregister 用户断开连接", client.Addr, client.AppID, client.UserID)

	if client.UserID != 0 {
		manager.sendAll("用户:" + strconv.FormatUint(client.UserID, 10) + "离开")
	}
}

// AddUsers 添加用户
func (manager *ClientManager) AddUsers(key string, client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()
	fmt.Println("添加用户: ", key)
	manager.Users[key] = client
}

// DelUsers 删除用户
func (manager *ClientManager) DelUsers(client *Client) (result bool) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	key := GetUserKey(client.AppID, client.UserID)

	if value, ok := manager.Users[key]; ok {
		// 判断是否为相同的用户
		if value.Addr != client.Addr {
			return
		}
		delete(manager.Users, key)
		result = true
	}

	return
}

// InClient 判断是否存在连接
func (manager *ClientManager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接存在，在添加
	_, ok = manager.Clients[client]

	return
}

// GetUserClients 获取用户的key
func (manager *ClientManager) GetUserClients() (clients []*Client) {

	clients = make([]*Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for _, v := range manager.Users {
		clients = append(clients, v)
	}

	return
}

func (manager *ClientManager) sendAll(data string) {
	clients := manager.GetUserClients()
	for _, client := range clients {
		client.SendMsg([]byte(data))
	}
}

func (manager *ClientManager) sendToUser(AppID uint32, UserID uint64, Nickname string, Touid uint64, data string) {
	// userInfo := models.GetUserinfoByID(UserID)
	item := map[string]interface{}{
		"from_uid":      UserID,
		"from_nickname": Nickname,
		"from_avatar":   "",
		"content":       data,
		"time":          time.Now().Format("2006-01-02 15:04:05"),
	}
	result := map[string]interface{}{
		"type": "chat",
		"data": item,
	}
	resp := entity.Response{
		Code:    200,
		Message: "ok",
		Data:    result,
	}
	resByte, err := json.Marshal(resp)
	if err != nil {
		return
	}
	key := GetUserKey(AppID, Touid)
	client := manager.Users[key]
	client.SendMsg(resByte)
}
