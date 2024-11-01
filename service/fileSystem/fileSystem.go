package fs

import (
	"os"
	"path/filepath"
)

var Root = filepath.Join("/tmp", "photos")

func CreatePhotoFile(photo Photo, binaryImage []byte) (err error) {
	err = os.MkdirAll(Root, 0755)
	if err != nil {
		return
	}
	createdFile, err := os.Create(photo.Path())
	if err != nil {
		return
	}
	defer createdFile.Close()
	_, err = createdFile.Write(binaryImage)
	return
}

func DeletePhotoFile(photo Photo) (err error) {
	err = os.Remove(photo.Path())
	return
}
