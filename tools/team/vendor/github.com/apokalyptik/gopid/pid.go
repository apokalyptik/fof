// Package pid provides an easy, reusable, library for locking and creating a
// pidfile common to *nix style servives.
//
//		import "github.com/apokalyptik/gopid"
//		var pidFile = "/var/run/my.pid"
//		_, err := pid.Do(pidFile)
//		if err != nil {
//			log.Fatalf("error creating pidfile: %s", err.Error())
//		}
package pid

import (
	"fmt"
	"os"
	"syscall"
)

// Do does what you want: it creates (or opens) a pidfile, exclusively locks it
// and writes the current executing programs PID to it.  It returns the file
// descriptor and an error.
func Do(filename string, permissions ...uint32) (*os.File, error) {
	if len(permissions) == 0 {
		permissions = []uint32{0666}
	}
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.FileMode(permissions[0]))
	if err != nil {
		return nil, err
	}
	err = syscall.Flock(int(fp.Fd()), syscall.LOCK_NB|syscall.LOCK_EX)
	if err != nil {
		return nil, err
	}
	syscall.Ftruncate(int(fp.Fd()), 0)
	syscall.Write(int(fp.Fd()), []byte(fmt.Sprintf("%d", os.Getpid())))
	return fp, nil
}
