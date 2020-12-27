package Golang

import (
	"fmt"
	"math/rand"
	"time"
)

type LoadBalance struct {
	Client []*Client
	size   int32
}

type Client struct {
	Name string
}

func (c *Client) Do() {
	fmt.Println("do")
}


func NewLoadBalance(size int32) *LoadBalance {
	lb := &LoadBalance{Client: make([]*Client, size), size: size}
	lb.Client = append(lb.Client, &Client{})
	return lb
}

func (lb *LoadBalance) getClient() *Client {
	rand.Seed(time.Now().Unix())
	x := rand.Int31n(100)
	return lb.Client[x%lb.size]
}


func main(){
	lb := NewLoadBalance(4)
	lb.getClient().Do()
}