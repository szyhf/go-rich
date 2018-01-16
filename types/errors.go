package richTypes

import (
	"errors"
	"fmt"
)

var (
	ErrorKeyNotExist    = fmt.Errorf("key not exist")
	ErrorCanNotRebuild  = errors.New("rebuild failed")
	ErrorMemberNotExist = fmt.Errorf("member not exist")
	ErrorProtection     = fmt.Errorf("key is protected")
)
