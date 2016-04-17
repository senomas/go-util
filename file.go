package util

import "os"

// RemoveAll directory
func RemoveAll(path ...string) {
	var err error
	for p := range path {
		err = os.RemoveAll(p)
		Check(err, "Remove error: path '%s'\n%v", p)
	}
}
