package backend

import (
	snow "github.com/bwmarrin/snowflake"
)

// IDGenerator is able to generate an unique ID for each cache item.
type IDGenerator interface {
	Get() int64
}

func newSnowflake(nodeID int64) snowflake {
	node, _ := snow.NewNode(nodeID)
	return snowflake{node}
}

type snowflake struct {
	*snow.Node
}

func (s snowflake) Get() int64 {
	return s.Node.Generate().Int64()
}
