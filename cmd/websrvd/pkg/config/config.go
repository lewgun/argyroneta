package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/lewgun/argyroneta/pkg/errutil"
	"github.com/lewgun/argyroneta/pkg/types"
)

type WebServer struct {
	HTTPPort int    `json:"port"`
	Prefix   string `json:"prefix"`
	RunMode  string `json:"run_mode"`
}

type Config struct {
	*types.Log   `json:"log"`
	*WebServer   `json:"webserver"`
	*types.Store `json:"store"`
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
