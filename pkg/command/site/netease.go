package site

import (
	"github.com/lewgun/argyroneta/pkg/command"
	"github.com/lewgun/argyroneta/pkg/rule"
)

func init() {

	command.Register(command.NetEase, neteaseCMDMaker)
}

func neteaseCMDMaker( r  *rule.Rule) command.Command  {


}