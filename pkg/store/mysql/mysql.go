package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type store struct {
	chEntry chan interface{}
	signal  chan types.Signal

	logger *logrus.Logger
	
	*xorm.Engine
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

	m.Engine.ShowSQL = c.ShowSQL
	m.Engine.SetMaxOpenConns(c.MaxConns)
	
	if err = m.Engine.CreateSchemas(&types.Rule{}, &types.Entry{}); err != nil {
		return err 
	}
	
	m.run()
	

	return nil

}

func (m *store) run() {

	go func() {
		for range m.chEntry {
		}

		close(m.signal)
	}()
}

//Close 关闭功能
func (m *store) Close() {

	close(m.chEntry)
	<-m.signal

}

func New(c *types.MySQLConf, logger *logrus.Logger) *store {
	if c == nil || logger == nil {
		return errutil.ErrInvalidParameter
	}
	
	s := &store{}
	if err := s.init(c, logger); err != nil {
		logger.Fatalln("new mysql instance failed.")
	}
	return c 
}