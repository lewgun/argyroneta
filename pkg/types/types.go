package types

import (
	"fmt"
	"time"
)

type Blob []byte

type Domain string

type MySQLConf struct {
	DBName string `json:"dbname"`
	IP string `json:"ip"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	ShowSQL bool  `json:"show_sql"`
	MaxConns int `json:"max_conns"`
	WorkerChanLen int `json:"worker_chan_len"`

}

type BoltConf struct {
		Path string `json:"path"`
}	
type Store struct {
	*MySQLConf  `json:"mysql"`
	*BoltConf 	`json:"bolt"`
		 
}

func (s *Store) String() string {
	return ""
	//return fmt.Sprintf("[path]: %s", s.Path)
}

type Log struct {
	Format string `json:"format"`
	Level  string `json:"level"`
}

func (a *Log) String() string {
	return fmt.Sprintf("[log]: %s:%s", a.Format, a.Level)
}

type Rule struct {
	ID int `xorm:" pk autoincr 'id'" json:"id"` //pk

	Domain string `xorm:" 'domain'" json:"domain"`

	//auth
	Account  string `xorm:" 'account'" json:"account"`
	Password string `xorm:" 'password'" json:"password"`

	//spider
	Politeness bool   `xorm:" 'politeness'" json:"politeness"` //是否遵守robots.txt
	UserAgent  string `xorm:" 'user_agent'" json:"user_agent"` //代理

	//filter
	AcceptExpr string `xorm:" 'accept'" json:"accept"`
	RejectExpr string `xorm:" 'reject'" json:"reject"`

	MaxDepth int           `xorm:" 'max_depth'" json:"max_depth"` //最大抓取深度
	Interval time.Duration `xorm:" 'interval'" json:"interval"`   //抓取间隔(以s计)

	Seed string `xorm:" notnull varchar(256) 'seed'" json:"seed"` //种子url

	ProxyOn bool `xorm:" 'proxy_on'" json:"proxy_on"` //是否启用HTTP 代理

}

type Entry struct {
	ID int `xorm:" pk autoincr 'id'" json:"id"` //pk

	//来源
	Origin string `xorm:" notnull varchar(256) 'seed'" json:"seed"`

	//原始url
	URL string `xorm:" notnull varchar(256) 'seed'" json:"seed"`

	//分类
	Category string `xorm:" 'seed'" json:"seed"`

	//标题
	Title string `xorm:" 'seed'" json:"seed"`

	//子标题
	SubTitle string `xorm:" 'seed'" json:"seed"`

	//概要
	Summary string `xorm:" 'seed'" json:"seed"`

	//作者
	Author string `xorm:" 'seed'" json:"seed"`

	//发布时间
	PublishAt time.Time `xorm:" 'seed'" json:"seed"`

	//抓取时间
	FetchAt time.Time `xorm:" 'seed'" json:"seed"`

	//tags
	Tags string `xorm:" 'seed'" json:"seed"`

	//内容
	ContentID string `xorm:" 'content_id'" json:"content_id"`

	//额外信息
	Extras string `xorm:" varchar(1024) 'extra'" json:"extra"`
}
