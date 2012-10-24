package rog

import (
    "time"
    "runtime"
)

type Timer <-chan time.Time

func NewTimer(msec float64) <-chan time.Time {
    return time.Tick(Milliseconds(msec))
}

func CheckTimer(c Timer) bool {
    select {
        case <-c:
            return true
        default: runtime.Gosched()
    }
    return false
}

type Ticker struct {
    callback func()
    rate float64
    timer Timer
}

func NewTicker(rate float64, callback func()) *Ticker {
    ticker := Ticker {
        callback: callback,
        rate: rate,
        timer: NewTimer(rate),
    }

    go func() {
        for {
            <-ticker.timer
            callback()
        }
    } ()

    return &ticker
}

func (self *Ticker) Rate() float64 {
    return self.rate
}

func (self *Ticker) SetRate(rate float64) {
    self.rate = rate
    self.timer = NewTimer(rate)
}

func Milliseconds(msecs float64) time.Duration {
    return time.Millisecond * time.Duration(msecs)
}