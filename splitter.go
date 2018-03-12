package main

import (
	"os"
	"path"
	"time"
)

type splitter struct {
	target *Target
	files map[string]*os.File
	dir string
	day string
}

func newSplitter(target *Target, dir string) *splitter {
	s := new(splitter)
	s.target = target
	s.dir = dir
	s.files = make(map[string]*os.File)
	s.day = day(time.Now())
	return s
}

func (s *splitter) split(obj map[string]interface{}) *os.File {
	if s.target.Split != nil && s.target.Split.By == "date" {
		return s.splitByDay(obj)
	}
	return s.file(s.day)
}

func (s *splitter) splitByDay(obj map[string]interface{}) *os.File {
	t := obj[s.target.Split.Value].(time.Time)
	key := day(t)
	return s.file(key)
}

func (s *splitter) file(key string) *os.File {
	f, ok := s.files[key]
	if ok {
		return f
	}
	p := path.Join(s.dir, key + ".json")
	f, err := os.Create(p)
	check(err)
	s.files[key] = f
	return f
}

func (s *splitter) close() {
	for _, f := range s.files {
		f.Close()
	}
}