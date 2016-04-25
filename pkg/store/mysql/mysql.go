package mysql

import (
	"strconv"

	"github.com/lewgun/argyroneta/pkg/types"

	"github.com/Sirupsen/logrus"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/lewgun/argyroneta/pkg/errutil"
)

type store struct {
	initialized bool
	chEntry     chan interface{}
	signal      chan types.Signal

	logger *logrus.Logger

	*xorm.Engine
}

var M *store

func init() {
	M = &store{}
}

func (m *store) init(c *types.MySQLConf, logger *logrus.Logger) error {

	var err error

	//"root:123new@tcp(125.64.93.75:3306)/oss?charset=utf8"
	dsn := c.User + ":" + c.Password + "@tcp(" + c.IP + ":" + strconv.Itoa(c.Port) + ")/" +
		c.DBName + "?charset=utf8&parseTime=true&loc=Local"

	m.Engine, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		return err
	}

	m.signal = make(chan types.Signal)
	m.chEntry = make(chan interface{}, c.WorkerChanLen)

	m.Engine.ShowSQL(c.ShowSQL)
	m.Engine.SetMaxOpenConns(c.MaxConns)

	if err = m.Engine.CreateTables(&types.Rule{}, &types.Entry{}); err != nil {
		return err
	}

	m.run()

	logger.Info("mysql is running")

	return nil

}

func (m *store) run() {

	go func() {
		for range m.chEntry {
		}

		close(m.signal)
	}()
}

//PowerOff 关闭功能
func (m *store) PowerOff() {

	close(m.chEntry)
	<-m.signal

}

////SharedInstInit initialize the shared instance, it can be called only once.
func SharedInstInit(c *types.MySQLConf, logger *logrus.Logger) error {

	if M.initialized {
		return fmt.Errorf("the bolt had initialized, please use the global variable 'mysql.M' instead")
	}

	if c == nil || logger == nil {
		return errutil.ErrInvalidParameter
	}

	if err := M.init(c, logger); err != nil {
		return err
	}
	return nil
}
