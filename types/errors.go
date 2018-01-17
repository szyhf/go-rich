package richTypes

import (
	"errors"
)

var (
	ErrorKeyNotExist    = errors.New("key not exist")
	ErrorCanNotRebuild  = errors.New("rebuild failed")
	ErrorMemberNotExist = errors.New("member not exist")
	ErrorProtection     = errors.New("key is protected")
	ErrorDeadLoop       = errors.New("rebuilding break for dead loop.")
	ErrorWaitLock       = errors.New("wait for lock.")
)
