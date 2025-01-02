package limiter

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) bool
	AllowN(key string, n int) bool
}

// Config 限流配置
type Config struct {
	Enable    bool
	Rules     Rules
	WhiteList WhiteList
}

type Rules struct {
	Global GlobalRule
	IP     Rule
	User   Rule
	API    APIRules
}

type Rule struct {
	QPS   int
	Burst int
}

type APIRules struct {
	Routes map[string]Rule
}

type WhiteList struct {
	IPs   []string
	Users []int64
}
