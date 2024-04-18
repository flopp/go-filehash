package filehash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Compute determines the hash/checksum (first 8 hex-bytes of SHA256) of the file identified by "fileName".
func Compute(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%.8x", h.Sum(nil)), nil
}

// Copy determins the hash/checksum of the file identified by "sourceFileName", replaces the last occurence of "placeholder" in "targetFileName" by this hash/checksum (if it occurs in "targetFileName"), and finally copies the source file to modified target location.
// The function returns the modified target filename and possibly an error.
func Copy(sourceFileName string, targetFileName string, placeholder string) (string, error) {
	targetFileName2 := targetFileName
	if placeholder != "" {
		pos := strings.LastIndex(targetFileName, placeholder)
		if pos != -1 {
			hash, err := Compute(sourceFileName)
			if err != nil {
				return "", err
			}
			targetFileName2 = targetFileName[:pos] + hash + targetFileName[pos+len(placeholder):]
		}
	}

	if err := os.MkdirAll(filepath.Dir(targetFileName2), 0770); err != nil {
		return "", err
	}

	source, err := os.Open(sourceFileName)
	if err != nil {
		return "", err
	}
	defer source.Close()

	target, err := os.Create(targetFileName2)
	if err != nil {
		return "", err
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return "", err
	}

	return targetFileName2, nil
}
