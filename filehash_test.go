package filehash

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCompute(t *testing.T) {
	// Expect error with non-existing file
	_, err := Compute("bad-filename.txt")
	if err == nil {
		t.Errorf(`Compute("bad-filename.txt") did not return an error`)
	}

	// Expect certain hash/checksum with known file
	tmpfile, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf(`Could not create temporary file: %v`, err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("My Temporary File")
	tmpfile.Close()
	hash, err := Compute(tmpfile.Name())
	if err != nil {
		t.Errorf(`Could not compute hash for file %s: %v`, tmpfile.Name(), err)
	}
	expectedHash := "95c39c37ef89acb2"
	if hash != expectedHash {
		t.Errorf(`Unexpected hash for  file %s: %s; expected: %s`, tmpfile.Name(), hash, expectedHash)
	}
}

func TestCopy(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Errorf(`Could not create temporary folder: %v`, err)
	}
	defer os.RemoveAll(dir)

	sourcePath := filepath.Join(dir, "source.txt")
	sourceFile, err := os.Create(sourcePath)
	if err != nil {
		t.Errorf(`Could not create temporary file: %v`, err)
	}
	sourceFile.WriteString("My Source File")
	sourceFile.Close()

	hash, err := Compute(sourcePath)
	if err != nil {
		t.Errorf(`Could not compute hash for file %s: %v`, sourcePath, err)
	}

	targetPath0 := filepath.Join(dir, "target.txt")
	targetPath1 := filepath.Join(dir, "target-HASH.txt")
	targetPath2 := filepath.Join(dir, "target-HASH-HASH.txt")

	_, err = Copy("bad-filename.txt", targetPath0, "HASH")
	if err == nil {
		t.Errorf(`No error when trying to copy non-existing file`)
	}

	// no placeholder occurence
	expectedPath0 := targetPath0
	if path, err := Copy(sourcePath, targetPath0, "HASH"); err != nil {
		t.Errorf(`Unexpected error when copying '%s' to '%s': %v`, sourcePath, targetPath0, err)
	} else if path != expectedPath0 {
		t.Errorf(`Unexpected modified path when copying '%s' to '%s': %s; expected: %s`, sourcePath, targetPath0, path, expectedPath0)
	}

	// one placeholder occurence
	expectedPath1 := filepath.Join(dir, fmt.Sprintf("target-%s.txt", hash))
	if path, err := Copy(sourcePath, targetPath1, "HASH"); err != nil {
		t.Errorf(`Unexpected error when copying '%s' to '%s': %v`, sourcePath, targetPath1, err)
	} else if path != expectedPath1 {
		t.Errorf(`Unexpected modified path when copying '%s' to '%s': %s; expected: %s`, sourcePath, targetPath1, path, expectedPath1)
	}

	// two placeholder occurences
	expectedPath2 := filepath.Join(dir, fmt.Sprintf("target-HASH-%s.txt", hash))
	if path, err := Copy(sourcePath, targetPath2, "HASH"); err != nil {
		t.Errorf(`Unexpected error when copying '%s' to '%s': %v`, sourcePath, targetPath2, err)
	} else if path != expectedPath2 {
		t.Errorf(`Unexpected modified path when copying '%s' to '%s': %s; expected: %s`, sourcePath, targetPath2, path, expectedPath2)
	}
}
