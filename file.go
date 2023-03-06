package main

import (
	"io"
	"os"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func makeFile(srcPath, dstPath string, info os.FileInfo) file {
	state := info.Sys().(*syscall.Stat_t)
	return file{
		srcPath: srcPath,
		dstPath: dstPath,
		aTime:   time.Unix(state.Atim.Sec, state.Atim.Nsec),
		mTime:   info.ModTime(),
		perm:    info.Mode(),
	}
}

type file struct {
	srcPath string
	dstPath string
	perm    os.FileMode
	aTime   time.Time
	mTime   time.Time
}

func (f *file) paths() (string, string) {
	return f.srcPath, f.dstPath
}

func (f *file) moveFile(remove bool) (int64, error) {

	const flag = os.O_RDWR | os.O_APPEND | os.O_CREATE
	in, err := os.Open(f.srcPath)
	if err != nil {
		return 0, errors.Wrapf(err, "openning file %s", f.srcPath)
	}

	defer func(in *os.File) {
		e := in.Close()
		if err == nil {
			err = e
		}
	}(in)
	out, err := os.OpenFile(f.dstPath, flag, f.perm)
	if err != nil {
		return 0, errors.Wrapf(err, "openning file %s", f.dstPath)
	}
	defer func(out *os.File) {
		e := out.Close()
		if err == nil {
			err = e
		}
	}(out)

	written, err := io.Copy(out, in)
	if err != nil {
		return 0, errors.Wrapf(err, "failed written to file %s", f.dstPath)
	}

	if err = os.Chtimes(f.dstPath, f.aTime, f.mTime); err != nil {
		return 0, err
	}

	if remove {
		if err = os.Remove(f.srcPath); err != nil {
			return 0, err
		}
	}
	return written, out.Sync()
}
