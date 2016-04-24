package chuper

import (
	"github.com/PuerkitoBio/fetchbot"
)

type Command interface {
	fetchbot.Command
	Depth() int
	Extras() interface{}
}

type Cmd struct {
	*fetchbot.Cmd
	d      int
	extras interface{}
}

func (c *Cmd) Depth() int {
	return c.d
}

func (c *Cmd) Extras() interface{} {
	return c.extras
}

type CmdBasicAuth struct {
	*fetchbot.Cmd
	d          int
	extras     interface{}
	user, pass string
}

func (c *CmdBasicAuth) Depth() int {
	return c.d
}

func (c *CmdBasicAuth) Extras() interface{} {
	return c.extras
}

func (c *CmdBasicAuth) BasicAuth() (string, string) {
	return c.user, c.pass
}
