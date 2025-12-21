package main

import (
	"fmt"
	"time"
)

type order struct {
	id      string
	amount  int32
	created time.Time
}

func (o order) changeAmt(amt int) {
	o.amount = int32(amt)
}

const (
	a = iota + 10
	b
	c
	d
)

func main() {
	o1 := order{
		id:      "1",
		amount:  10,
		created: time.Now(),
	}
	fmt.Println(o1)
	o1.changeAmt(55)
}
