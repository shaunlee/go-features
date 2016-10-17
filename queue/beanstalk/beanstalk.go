package beanstalk

import (
	"github.com/kr/beanstalk"
)

func newQueue() (*beanstalk.Conn, error) {
    return beanstalk.Dial("tcp", "127.0.0.1:11300")
}
