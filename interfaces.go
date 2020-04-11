package main

import (
	"fmt"
	"sort"
)

type Data interface {
	Data() string
}
type hashFunction interface {
	GetHash([]byte) uint64
}

type ExternalConfig struct {
	HashFn           hashFunction
	RepetitionFactor int
	ServerCount      int
}

type Store struct {
	Config     ExternalConfig
	memberList map[string]*Data
	hashMapper map[uint64]*Data
	//hashRing map[uint64]*Data
	hashList []uint64
}

type Conshash interface {
	AddNode(*Data)
	RemoveNode(*Data)
	GetNode(key string) int
}

func NewConshash(data []Data, config ExternalConfig) *Store {
	s := Store{
		Config:     config,
		memberList: make(map[string]*Data),
		hashMapper: make(map[uint64]*Data),
		//	hashRing:   ,
		//hashList:   make([]),
	}
	for _, d := range data {
		s.AddNode(d)
	}

	return &s
}

func (s *Store) AddNode(d Data) {
	for i := 0; i < s.Config.RepetitionFactor; i++ {
		key := fmt.Sprintf("%s%d", (d).Data(), i)
		hash := s.Config.HashFn.GetHash([]byte(key))
		s.hashList = append(s.hashList, hash)
		s.hashMapper[hash] = &d
	}
	sort.Slice(s.hashList, func(i int, j int) bool {
		return s.hashList[i] < s.hashList[j]
	})
	s.memberList[(d).Data()] = &d

}

func (s *Store) GetNode(key string) *Data {
	searchFn := func(i int) bool {
		hash := s.Config.HashFn.GetHash([]byte(key))
		return s.hashList[i] >= hash
	}

	node := sort.Search(len(s.hashList), searchFn)
	if node >= s.Config.ServerCount {
		node = 0
	}
	partition := s.hashMapper[s.hashList[node]]
	server := s.memberList[(*partition).Data()]
	return server
}


