package sdks

import (
	"strings"

	"github.com/google/uuid"
)

func NewIDFunc() string {
	return newIDFunc()
}

func SetIDFunc(f func() string) {
	newIDFunc = f
}

var newIDFunc = func() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
