package util

import "os"

// RemoveAll directory
func RemoveAll(path ...string) {
	var err error
	for _, p := range path {
		err = os.RemoveAll(p)
		Check("Remove error: path '%s'\n%v", p, err)
	}
}
