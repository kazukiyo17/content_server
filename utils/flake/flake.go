package flake

import (
	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake

func Setup() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func Generate() (int64, error) {
	id, err := flake.NextID()
	return int64(id), err
}
