package main

import (
	"fmt"
	"strings"
)

type URLEndpoint struct {
	url      string
	endpoint string
}

type Endpoints []URLEndpoint

func (eps *Endpoints) String() string {
	var l []string
	for _, ep := range *eps {
		l = append(l, fmt.Sprintf("%s:%d", ep.url, ep.endpoint))
	}
	return strings.Join(l, ",")
}

func (eps *Endpoints) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	ep := URLEndpoint{url: parts[0], endpoint: parts[1]}
	*eps = append(*eps, ep)
	return nil
}
