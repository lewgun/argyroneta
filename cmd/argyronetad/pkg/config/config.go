package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"bytes"

	"github.com/lewgun/argyroneta/pkg/errutil"
)

var (
	C *Config
)

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

type Config struct {
	*Store `json:"store"`

	HTTPProxies []Proxy          `json:"http_proxies"`
	Sites       map[string]*Site `json:"sites"`
}

func (c *Config) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("store:\n")
	buf.WriteString("\t" + c.Store.String())
	buf.WriteString("\n")
	buf.WriteString("proxys:\n")

	for _, p := range c.HTTPProxies {
		buf.WriteString("\t" + p.String() + "\n")

	}

	buf.WriteString("sites:\n")
	for domain, site := range c.Sites {
		buf.WriteString("\t" + domain + ":\n\t\t " + site.String() + "\n")
	}

	return buf.String()

}

func (c *Config) init(path string) error {
	var err error
	if err = c.parse(path); err != nil {
		return fmt.Errorf("Can't load config from: %s with error: %v ", path, err)
	}

	if err = c.adjust(); err != nil {
		return fmt.Errorf("Adjust config failed.")
	}

	fmt.Println(c)

	return c.check()
}

func (c *Config) adjust() error {

	return nil
}

func (c *Config) parse(path string) error {
	if path == "" {
		return errutil.ErrInvalidParameter
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)

	return err
}

//check检测配置参数是否完备
func (c *Config) check() error {
	return nil
}

//Load 创建一个配置
func SharedInstance(path string) *Config {

	if C != nil {
		return C
	}

	C = &Config{}
	err := C.init(path)
	if err != nil {
		panic(err)
	}
	return C
}
