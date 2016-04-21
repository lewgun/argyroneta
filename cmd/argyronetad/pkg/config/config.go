package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/lewgun/argyroneta/pkg/errutil"
	"github.com/lewgun/argyroneta/pkg/types"
)

type Config struct {
	*types.Store `json:"store"`
	*types.Log   `json:"log"`

	HTTPProxies []types.Proxy                `json:"http_proxies"`
	Sites       map[types.Domain]*types.Site `json:"sites"`
}

func (c *Config) String() string {
	buf := &bytes.Buffer{}

	buf.WriteString("store:\n")
	buf.WriteString("\t" + c.Store.String())
	buf.WriteString("\n")

	buf.WriteString("log:\n")
	buf.WriteString("\t" + c.Log.String())
	buf.WriteString("\n")

	buf.WriteString("proxys:\n")
	for _, p := range c.HTTPProxies {
		buf.WriteString("\t" + p.String() + "\n")

	}

	buf.WriteString("sites:\n")
	for domain, site := range c.Sites {
		buf.WriteString("\t" + string(domain) + ":\n\t\t " + site.String() + "\n")
	}

	return buf.String()

}

func (c *Config) Init(path string) error {
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

//New 创建一个配置
func New() *Config {
	return &Config{}

}
