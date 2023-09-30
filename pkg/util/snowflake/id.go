package snowflake

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
)

// SnowFlakeID snowflake id
type SnowFlakeID struct {
	*snowflake.Node
}

var snowFlakeIDGenerator *SnowFlakeID

func init() {
	//todo
	rand.Seed(time.Now().UnixNano())
	node, err := snowflake.NewNode(int64(rand.Intn(1000)) + 1)
	if err != nil {
		panic(err.Error())
	}
	snowFlakeIDGenerator = &SnowFlakeID{node}
}

func ID() snowflake.ID {
	id := snowFlakeIDGenerator.Generate()
	return id
}

func IDBase58() string {
	id := snowFlakeIDGenerator.Generate()
	return id.Base58()
}

func IDBase32() string {
	id := snowFlakeIDGenerator.Generate()
	return id.Base32()
}

func IDBase64() string {
	id := snowFlakeIDGenerator.Generate()
	return id.Base64()
}
