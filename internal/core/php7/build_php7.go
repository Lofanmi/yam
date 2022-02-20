//go:build php7
// +build php7

package php7

/*
#include "extension.h"
#cgo CFLAGS: -std=c99 -O2 -D_GNU_SOURCE -I../../php-src/7.2.34/main -I../../php-src/7.2.34/Zend -I../../php-src/7.2.34/TSRM -I../../php-src/7.2.34
#cgo LDFLAGS: -Wl,-undefined -Wl,dynamic_lookup
*/
import "C"

var ImportMe struct{}
