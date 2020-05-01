package main

import (
	"fmt"
	"hash/fnv"
)

type testS string

type hasher struct{}

func (t *testS) Data() string {
	return string(*t)
}

func (t hasher) GetHash(b []byte) uint64 {
	h := fnv.New64a()
	h.Write(b)
	return h.Sum64()

}
func main() {

	cfg := ExternalConfig{
		HashFn:           hasher{},
		RepetitionFactor: 2,
		ServerCount:      2,
	}
	server1 := testS("server3")
	server2 := testS("server4")

	c := NewConshash(nil, cfg)
	c.AddNode(&server1)
	c.AddNode(&server2)
	adddr := c.GetNode("server4")
	name := (*adddr).Data()
	fmt.Println("Address of the key is ", name)

}
