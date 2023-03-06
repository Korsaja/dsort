package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

const formatByDate = "02-jan-2006"

var skipDirs = map[string]bool{
	".Trash-1000": true, // folder paperbin (linux)
}

type SortAction struct {
	DirPath  string
	Removed  bool
	SkipDirs []string
}

func DoSort(actions SortAction, logger *slog.Logger) error {
	for _, skipDir := range actions.SkipDirs {
		skipDirs[skipDir] = true
	}

	files, err := walkDirAndPrepare(actions.DirPath, logger)
	if err != nil {
		return err
	}

	var written int64
	for _, instance := range files {
		src, dst := instance.paths()
		logger.Info("move file", slog.String("src", src), slog.String("dst", dst))
		n, err := instance.moveFile(actions.Removed)
		if err != nil {
			logger.Error("move file", err)
		}
		written += n
	}

	logger.Info("done", slog.Int64("total_moved_bytes", written))
	return nil
}

func walkDirAndPrepare(root string, logger *slog.Logger) ([]file, error) {
	var files = make([]file, 0)
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			if _, ok := err.(*fs.PathError); ok {
				logger.Error("path processing", err)
				return filepath.SkipDir
			}
			return err
		}
		if info.IsDir() && skipDirs[info.Name()] {
			return filepath.SkipDir
		}

		if !info.IsDir() && info.Mode().IsRegular() {
			absPath, err := filepath.Abs(root)
			if err != nil {
				return err
			}

			newPath := filepath.Join(absPath, info.ModTime().Format(formatByDate))
			dstPath := filepath.Join(newPath, info.Name())

			// skip dupl file
			if path == dstPath {
				return nil
			}

			files = append(files, makeFile(path, dstPath, info))
			return makeDirIfNotExists(newPath)
		}

		return nil
	})

	logger.Info("folder traversal complete", slog.Int("count_files", len(files)))
	return files, err
}

func makeDirIfNotExists(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return os.Mkdir(dirPath, 0700)
	}
	return err
}
