package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"io.xiu/listAllFiles/domian"
	"io.xiu/listAllFiles/service"
)

func WalkDir(dir string, waitGroup *sync.WaitGroup, fileSizes chan int64, lock *sync.RWMutex) {
	defer waitGroup.Done()
	var files []domian.File
	if cancelled() {
		return
	}
	var subdir string
	for _, entry := range dirents(dir) {
		subdir = filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			waitGroup.Add(1)
			go WalkDir(subdir, waitGroup, fileSizes, lock)
		} else {
			// 获取文件数
			file, _ := entry.Info()
			files = append(files, domian.File{Id: uuid.NewString(), Name: entry.Name(), Path: subdir, Time: file.ModTime().Unix()})
			fileSizes <- file.Size()
		}
	}
	if len(files) != 0 {
		service.Insert(files, lock)
	}
}

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.DirEntry {
	select {
	case sema <- struct{}{}: // 获取令牌
	case <-Done:
		return nil // 取消
	}
	defer func() { <-sema }() // 释放令牌

	// 读取目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du：%v\n", err)
	}
	return entries
}

func cancelled() bool {
	select {
	case <-Done:
		return true
	default:
		return false
	}
}

var Done = make(chan struct{})