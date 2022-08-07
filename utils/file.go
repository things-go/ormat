package utils

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

// FileModTime returns file modified time and possible error.
func FileModTime(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.ModTime().Unix(), nil
}

// FileSize returns file size in bytes and possible error.
func FileSize(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

// IsDir returns true if given path is a dir,
// or returns false when it's a file or does not exist.
func IsDir(filePath string) bool {
	f, err := os.Stat(filePath)
	return err == nil && f.IsDir()
}

// IsFile returns true if given path is a file,
// or returns false when it's a directory or does not exist.
func IsFile(filePath string) bool {
	f, err := os.Stat(filePath)
	return err == nil && !f.IsDir()
}

// FileMode returns file mode and possible error.
func FileMode(name string) (os.FileMode, error) {
	fInfo, err := os.Lstat(name)
	if err != nil {
		return 0, err
	}
	return fInfo.Mode(), nil
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(paths string) bool {
	_, err := os.Stat(paths)
	return err == nil || os.IsExist(err)
}

// HasPermission returns a boolean indicating whether that permission is allowed.
func HasPermission(name string) bool {
	_, err := os.Stat(name)
	return !os.IsPermission(err)
}

// FileCopy copies file from source to target path.
func FileCopy(src, dest string) error {
	// Gather file information to set back later.
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Handle symbolic link.
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		// NOTE: os.Chmod and os.Chtimes don't recoganize symbolic link,
		// which will lead "no such file or directory" error.
		return os.Symlink(target, dest)
	}

	sr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sr.Close()

	dw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dw.Close()

	if _, err = io.Copy(dw, sr); err != nil {
		return err
	}

	// Set back file information.
	if err := os.Chtimes(dest, si.ModTime(), si.ModTime()); err != nil {
		return err
	}
	return os.Chmod(dest, si.Mode())
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it
// and its upper level paths.
func WriteFile(filename string, data []byte) error {
	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0655)
}

// FilePaths returns all root dir (contain sub dir) file full path
func FilePaths(root string) ([]string, error) {
	var result = make([]string, 0)

	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	return result, err
}
