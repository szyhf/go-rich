package richTypes

import (
	"errors"
	"strings"
)

var (
	ErrorRedisNil = errors.New("redis: nil")

	ErrorKeyNotExist    = errors.New("key not exist")
	ErrorRebuildNil     = errors.New("rebuild nil.")
	ErrorMemberNotExist = errors.New("member not exist")
	ErrorProtection     = errors.New("key is protected")
	ErrorDeadLoop       = errors.New("rebuilding break for dead loop.")
	ErrorWaitLock       = errors.New("wait for lock.")
)

// 把go-redis的error处理一下
func FromGoRedisErr(err error) error {
	switch true {
	case strings.Compare(err.Error(), "redis: nil") == 0:
		return ErrorRedisNil
	}
	return err
}
