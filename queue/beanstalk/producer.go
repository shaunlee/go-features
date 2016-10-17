package beanstalk

import (
	"time"

	json "github.com/bitly/go-simplejson"
	"github.com/kr/beanstalk"
)

type Producer struct {
	conn *beanstalk.Conn
	tube *beanstalk.Tube
}

func NewProducer(tube string) (*Producer, error) {
	if conn, err := newQueue(); err != nil {
		return nil, err
	} else {
		return &Producer{
			conn,
			&beanstalk.Tube{conn, tube},
		}, nil
	}
}

func (p *Producer) Close() error {
	return p.conn.Close()
}

func (p *Producer) Produce(body *json.Json, delay time.Duration) (uint64, error) {
	if data, err := body.Encode(); err != nil {
		return 0, err
	} else {
		return p.tube.Put(data, 1, delay, time.Minute)
	}
}

func ProduceImmediately(tube string, body *json.Json, delay time.Duration) (uint64, error) {
	if p, err := NewProducer(tube); err != nil {
		return 0, err
	} else {
		defer p.Close()
		return p.Produce(body, delay)
	}
}
