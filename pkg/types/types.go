package types

import "fmt"

type Blob []byte

type Domain string

type Store struct {
	Path string `json:"path"`
}

func (s *Store) String() string {
	return fmt.Sprintf("[path]: %s", s.Path)
}

type Auth struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (a *Auth) String() string {
	return fmt.Sprintf("[auth]: %s:%s", a.Account, a.Password)
}

type Log struct {
	Format string `json:"format"`
	Level  string `json:"level"`
}

func (a *Log) String() string {
	return fmt.Sprintf("[log]: %s:%s", a.Format, a.Level)
}

type Proxy string

func (p Proxy) String() string {
	return fmt.Sprintf("%s", string(p))
}

type Spider struct {
	Politeness bool   `json:"politeness"` //是否遵守robots.txt
	UserAgent  string `json:"user_agent"` //代理

}

func (s *Spider) String() string {
	return fmt.Sprintf("[spider]: user agent:%s politeness: %v", s.UserAgent, s.Politeness)
}

type Site struct {
	*Auth   `json:"auth"`
	*Spider `json:"spider"`

	MaxDepth uint32 `json:"max_depth"` //最大抓取深度
	Delay    uint32 `json:"delay"`     //抓取间隔(以s计)

	Seed string `json:"seed"` //种子url

	ProxyOn bool `json:"proxy_on"` //是否启用HTTP 代理

}

func (s *Site) String() string {
	return fmt.Sprintf("%v\t%v\t[max depth]: %d\t[delay]: %d\t[seed]: %s\t[proxy on]: %v",
		s.Auth,
		s.Spider,
		s.MaxDepth,
		s.Delay,
		s.Seed,
		s.ProxyOn)
}
