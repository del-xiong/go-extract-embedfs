package extractfs

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
)

// extract fs to local path
func ExtractToPath(embedfs *embed.FS, localSavePath string) (*ExtractedInfo, error) {
	_, err := os.Stat(localSavePath)
	if err != nil && os.IsNotExist(err) {
		return &ExtractedInfo{}, errors.New(fmt.Sprintf("path(%s) not exists or not accessible", localSavePath))
	} else if entries, err := embedfs.ReadDir("."); err != nil {
		return &ExtractedInfo{}, err
	} else if fileCount, dirCount, err := ExtractDirEntryToPath(embedfs, entries, ".", localSavePath); err != nil {
		return &ExtractedInfo{}, err
	} else {
		return &ExtractedInfo{
			ExtractedFiles:   fileCount,
			ExtractedFolders: dirCount,
		}, nil
	}

	return &ExtractedInfo{}, nil
}

// extract fs.DirEntry to local path
func ExtractDirEntryToPath(embedfs *embed.FS, entries []fs.DirEntry, fsPath string, localSavePath string) (int64, int64, error) {
	var fileCount int64
	var dirCount int64
	for _, entry := range entries {
		if entry.IsDir() {
			dirCount++
			newPath := path.Join(fsPath, entry.Name())
			if entries, err := embedfs.ReadDir(newPath); err != nil {
				return fileCount, dirCount, err
			} else if subdirFileCount, subdirCount, err := ExtractDirEntryToPath(embedfs, entries, newPath, path.Join(localSavePath, entry.Name())); err != nil {
				return fileCount, dirCount, err
			} else {
				fileCount += subdirFileCount
				dirCount += subdirCount
			}
		} else {
			fileCount++
			_, err := os.Stat(localSavePath)
			if err != nil && os.IsNotExist(err) {
				if err = os.MkdirAll(localSavePath, os.ModePerm); err != nil {
					return fileCount, dirCount, err
				}
			}
			fileStream, err := embedfs.Open(path.Join(fsPath, entry.Name()))
			if err != nil {
				return fileCount, dirCount, err
			}
			defer fileStream.Close()
			targetStream, err := os.OpenFile(path.Join(localSavePath, entry.Name()), os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return fileCount, dirCount, err
			}
			_, err = io.Copy(targetStream, fileStream)
			if err != nil {
				return fileCount, dirCount, err
			}
			if err = targetStream.Close(); err != nil {
				return fileCount, dirCount, err
			} else if err = fileStream.Close(); err != nil {
				return fileCount, dirCount, err
			}
		}
	}
	return fileCount, dirCount, nil
}
