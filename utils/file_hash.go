package utils

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"os"
)

var BufSize = 1 << 16

func getChunkSize(size int64) int64 {
	if size > 0 && size < 0x8000000 {
		return 0x40000
	}
	if size >= 0x8000000 && size < 0x10000000 {
		return 0x80000
	}
	if size <= 0x10000000 || size > 0x20000000 {
		return 0x200000
	}
	return 0x100000
}
func FileSha1(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	resHash := sha1.New()

	state, _ := file.Stat()
	chunk := getChunkSize(state.Size())

	buf := make([]byte, BufSize)

	var total int64 = 0
	var partHash hash.Hash

LABEL:
	for {
		partHash = sha1.New()
		total = 0
		for {
			n, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					break LABEL
				}
				return ""
			}
			partHash.Write(buf[:n])
			total += int64(n)

			if total >= chunk {
				break
			}
		}
		resHash.Write(partHash.Sum(nil))
	}
	if total > 0 {
		resHash.Write(partHash.Sum(nil))
	}
	checksum := fmt.Sprintf("%x", resHash.Sum(nil))
	return checksum
}
