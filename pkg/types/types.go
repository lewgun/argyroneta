package types

import (
	"fmt"
	"strconv"
	"time"
)

//Signal 退出信号
type Signal struct{}

//Blob 表示一大块数据
type Blob []byte

//Domain 域名
type Domain string

//Proxy HTTP代理
type Proxy string

func (p Proxy) String() string {
	return fmt.Sprintf("%s", string(p))
}

//MySQLConf mysql实例配置参数
type MySQLConf struct {
	DBName        string `json:"dbname"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	User          string `json:"user"`
	Password      string `json:"password"`
	ShowSQL       bool   `json:"show_sql"`
	MaxConns      int    `json:"max_conns"`
	WorkerChanLen int    `json:"worker_chan_len"`
}

func (c *MySQLConf) String() string {
	dsn := c.User + ":" + c.Password + "@tcp(" + c.IP + ":" + strconv.Itoa(c.Port) + ")/" +
		c.DBName + "?charset=utf8&parseTime=true&loc=Local"
	return dsn

}

//BoltConf BoltDB 配置参数
type BoltConf struct {
	Path string `json:"path"`
}

func (s *BoltConf) String() string {
	return fmt.Sprintf("[path]: %s", s.Path)
}

//Store 所有的存储系统集合
type Store struct {
	*MySQLConf `json:"mysql"`
	*BoltConf  `json:"bolt"`
}

func (s *Store) String() string {

	return fmt.Sprintf("[mysql]: %v\n\t[bolt]: %v", s.MySQLConf, s.BoltConf)
}

//Log 日志配置
type Log struct {
	Format string `json:"format"`
	Level  string `json:"level"`
}

func (a *Log) String() string {
	return fmt.Sprintf("[log]: %s:%s", a.Format, a.Level)
}

//Rule 定义一条抓取规则
type Rule struct {
	ID int `xorm:" pk autoincr 'id'" json:"id"` //pk

	Domain Domain `xorm:" 'domain'" json:"domain"`

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

//Entry 代表一条保存的实体
type Entry struct {
	ID int `xorm:" pk autoincr 'id'" json:"id"` //pk

	//来源
	Origin string `xorm:" notnull varchar(256) 'origin'" json:"origin"`

	//原始url
	URL string `xorm:" notnull varchar(256) 'url'" json:"url"`

	//分类
	Category string `xorm:" 'category'" json:"category"`

	//标题
	Title string `xorm:" 'title'" json:"title"`

	//子标题
	SubTitle string `xorm:" 'sub_title'" json:"sub_title"`

	//概要
	Summary string `xorm:" 'summary'" json:"summary"`

	//作者
	Author string `xorm:" 'author'" json:"author"`

	//发布时间
	PublishAt time.Time `xorm:" 'publish_at'" json:"publish_at"`

	//抓取时间
	FetchAt time.Time `xorm:" 'fetch_at'" json:"fetch_at"`

	//tags
	Tags string `xorm:" 'tags'" json:"tags"`

	//内容
	ContentID string `xorm:" 'content_id'" json:"content_id"`

	//额外信息
	Extras string `xorm:" varchar(1024) 'extra'" json:"extra"`
}

type TopReq struct {
	N        int    `json:"n"`
	Category string `json:"category"`
}
