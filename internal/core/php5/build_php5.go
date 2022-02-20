//go:build php5
// +build php5

package php5

/*
#include "extension.h"
#cgo CFLAGS: -g -std=c99 -O2 -D_GNU_SOURCE -I../../php-src/5.6.40/main -I../../php-src/5.6.40/Zend -I../../php-src/5.6.40/TSRM -I../../php-src/5.6.40
#cgo linux  LDFLAGS: -Wl,--warn-unresolved-symbols -Wl,--unresolved-symbols=ignore-all
#cgo darwin LDFLAGS: -Wl,-undefined -Wl,dynamic_lookup
*/
import "C"

var ImportMe struct{}
