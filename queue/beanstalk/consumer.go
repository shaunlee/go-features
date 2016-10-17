package beanstalk

import (
	"errors"
	"time"

	json "github.com/bitly/go-simplejson"
	"github.com/kr/beanstalk"
)

var GiveUp = errors.New("Give up")

type Consumer struct {
	conn *beanstalk.Conn
	tube *beanstalk.TubeSet
}

type ConsumeFunc func(*json.Json) error

func NewConsumer(tube string) (*Consumer, error) {
	if conn, err := newQueue(); err != nil {
		return nil, err
	} else {
		return &Consumer{
			conn,
			beanstalk.NewTubeSet(conn, tube),
		}, nil
	}
}

func (p *Consumer) Close() error {
	return p.conn.Close()
}

func (p *Consumer) Consume(timeout time.Duration, fn ConsumeFunc) error {
	if id, body, err := p.tube.Reserve(timeout); err != nil {
		return err
	} else {
		if data, err := json.NewJson(body); err != nil {
			p.conn.Bury(id, 1)
			return err
		} else {
			if err := fn(data); err != nil {
				p.conn.Release(id, 1, time.Minute)
				return err
			} else {
				p.conn.Delete(id)
			}
			return nil
		}
	}
}

func (p *Consumer) ConsumeLoop(timeout time.Duration, fn ConsumeFunc) error {
	for {
		if err := p.Consume(timeout, fn); err != nil {
			if err == GiveUp {
				break
			} else {
				return err
			}
		}
	}
	return nil
}

func ConsumeImmediately(tube string, timeout time.Duration, fn ConsumeFunc) error {
	if p, err := NewConsumer(tube); err != nil {
		return err
	} else {
		defer p.conn.Close()
		return p.Consume(timeout, fn)
	}
}

func ConsumeImmediatelyLoop(tube string, timeout time.Duration, fn ConsumeFunc) error {
	if p, err := NewConsumer(tube); err != nil {
		return err
	} else {
		defer p.conn.Close()
		return p.ConsumeLoop(timeout, fn)
	}
}
