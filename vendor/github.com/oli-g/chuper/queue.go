package chuper

import (
	"net/url"

	"github.com/PuerkitoBio/fetchbot"
)

type Queue struct {
	*fetchbot.Queue
}

type Enqueuer interface {
	Enqueue(string, string, int, ...interface{}) error

	EnqueueWithBasicAuth(string, string, int, string, string, ...interface{}) error

	EnqueueCommand(Command) error
}

func (q *Queue) Enqueue(method, URL string, depth int, extras ...interface{}) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	cmd := &Cmd{
		Cmd: &fetchbot.Cmd{U: u, M: method},
		d:   depth,
	}

	if len(extras) != 0 {
		cmd.extras = extras
	}

	if err = q.Send(cmd); err != nil {
		return err
	}
	return nil
}

func (q *Queue) EnqueueWithBasicAuth(
	method string,
	URL string,
	depth int,
	user string,
	password string,
	extras ...interface{}) error {
	if user == "" && password == "" {
		return q.Enqueue(method, URL, depth)
	}

	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	cmd := &CmdBasicAuth{
		Cmd:  &fetchbot.Cmd{U: u, M: method},
		d:    depth,
		user: user,
		pass: password,
	}

	if len(extras) != 0 {
		cmd.extras = extras
	}

	if err = q.Send(cmd); err != nil {
		return err
	}

	return nil
}

func (q *Queue) EnqueueCommand(cmd Command) error {
	return q.Send(cmd)
}
