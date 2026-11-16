package paxostesting

import (
	"io"
)

// readAll reads all bytes from r, returning the result and ignoring errors.
func readAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
