package fs

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func NewFlockFile(path string) (*os.File, error) {
	return createFlockFile(path)
}

func createFlockFile(flockFile string) (*os.File, error) {
	flockF, err := os.Create(flockFile)
	if err != nil {
		return nil, fmt.Errorf("create lock file [%q] failure, nest error: %w", flockFile, err)
	}
	if err := unix.Flock(int(flockF.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		return nil, fmt.Errorf("acquire lock on file [%q] failure, nest error: %w", flockFile, err)
	}
	return flockF, nil
}
