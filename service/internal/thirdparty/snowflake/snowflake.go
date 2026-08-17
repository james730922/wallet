package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
)

func NewIDGenerator() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		logger.ApLog().Error(err)
		panic(err)
	}
	return node
}
