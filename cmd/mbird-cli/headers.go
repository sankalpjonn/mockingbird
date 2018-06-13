package main

import (
	"errors"
	"fmt"
	"regexp"
)

type Headers map[string]string

func (self *Headers) String() string {
	return fmt.Sprintf("%s", *self)
}

func (self *Headers) Set(value string) error {
	r := regexp.MustCompile("(.*)=(.*)")
	match := r.FindStringSubmatch(value)
	if len(match) < 3 {
		return errors.New(`Header flag must be in format "key=value"`)
	}
	if *self == nil {
		*self = map[string]string{}
	}
	(*self)[match[1]] = match[2]
	return nil
}
