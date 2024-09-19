package main

import (
	"fmt"
	"strings"
)

type Endpoints map[string]string

func (eps *Endpoints) String() string {
	var l []string
	for url, ep := range *eps {
		l = append(l, fmt.Sprintf("%s:%d", url, ep))
	}
	return strings.Join(l, ",")
}

func (eps *Endpoints) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	(*eps)[parts[0]] = parts[1]
	return nil
}
