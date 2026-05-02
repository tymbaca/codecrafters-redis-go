package command

import "github.com/codecrafters-io/redis-starter-go/pkg/enc"

type Echo struct {
	Val enc.Value
	Context
}
