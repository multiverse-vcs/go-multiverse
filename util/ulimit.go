package util

import (
	"golang.org/x/sys/unix"
)

// SetFileLimit sets the soft and max file limits for the process.
func SetFileLimit(soft, max uint64) error {
	rlimit := unix.Rlimit{
		Cur: soft,
		Max: max,
	}
	
	return unix.Setrlimit(unix.RLIMIT_NOFILE, &rlimit)
}

// GetFileLimit returns the soft and max file limits for the process.
func GetFileLimit() (uint64, uint64, error) {
	var rlimit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &rlimit); err != nil {
		return 0, 0, err
	}

	return rlimit.Cur, rlimit.Max, nil
}