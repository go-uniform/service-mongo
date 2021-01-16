package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("mongo.insert"), mongoInsert)
}

func mongoInsert(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
	p.Info("test", nil)
}