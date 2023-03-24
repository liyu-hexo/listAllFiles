package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	utils "io.xiu/listAllFiles/utils"
)

var file string
var dir string

func main() {
	// 计算耗时
	defer timeCost()()
	flag.StringVar(&file, "f", "", "搜索文件")
	flag.StringVar(&dir, "d", "D:/", "目录")
	flag.Parse()
	lock := &sync.RWMutex{}
	// 确定初始目录
	roots := []string{dir}
	if len(roots) == 0 {
		fmt.Println("命令行后跟上目录路径")
	}
	// 当检测到输入时取消遍历
	go func() {
		// 读一个字节
		os.Stdin.Read(make([]byte, 1))
		close(utils.Done)
	}()
	// 并行遍历文件树的每个根
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		fmt.Printf("当前目录为:%s\n", root)
		n.Add(1)
		go utils.WalkDir(root, file, &n, fileSizes, lock)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	var nfiles, nbytes int64
loop:
	for {
		select {
		case <-utils.Done:
			// 耗尽 fileSizes 以允许已有的  goroutine 结束
			for range fileSizes {
				// 不执行任何操作
			}
			return
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes 关闭
			}
			nfiles++
			nbytes += size
		}
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("耗时%v\n", tc)
	}
}
