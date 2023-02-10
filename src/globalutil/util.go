package globalutil

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return Errorf(err, "Failed to open zip file %q", src)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return Errorf(err, "Failed to open file %q", f.Name)
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return Errorf(err, "Failed to create directory %q", path)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return Errorf(err, "Failed to create directory %q", filepath.Dir(path))
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return Errorf(err, "Failed to create file %q", path)
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return Errorf(err, "Failed to copy file %q", path)
			}
		}
	}
	return nil
}

func FromWindowsPath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func CreateDirAndFile(path string) (*os.File, error) {
	path = filepath.FromSlash(path)

	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, Errorf(err, "Failed to create directory %q", filepath.Dir(path))
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, Errorf(err, "Failed to create file %q", path)
	}

	return file, nil
}
