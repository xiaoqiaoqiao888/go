// UsualPassenger InitData Program
//
// This program used to import all user to Codis
//
// Copyright (C)2017 EmbraceSource Co.,Ltd. All Rights Reserved
//
// Author: yufu.zhao <tony.zhao@embracesource.com> at 2018/7/11 15:07:20

package main

import (
    "runtime"
    "../common"
    "sync"
    "fmt"
    "os"
    "bufio"
    "strings"
    "time"
    "github.com/axgle/mahonia"
    l4g "github.com/ms2008/log4go"
)

const (
    CONF = "../config.conf"
)

var log = make(l4g.Logger)
var buildstamp = ""
var githash = ""


func init() {
    log.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())
    log.AddFilter("file", l4g.WARNING, l4g.NewFileLogWriter(fmt.Sprintf("null_%d_%s.log", os.Getpid(), time.Now().Format("20060102150405")), false))
}

func main() {
    defer log.Close()

    args := os.Args
    if len(args)==2 && (args[1]=="--version" || args[1] =="-v") {
        fmt.Printf("Git Commit Hash: %s\n", githash)
        fmt.Printf("UTC Build Time : %s\n", buildstamp)
        return
    }

    log.Warn("扫描数据开始...")
    log.Warn("os args: %v", os.Args)

    if args == nil || len(args) < 2 {
        panic("缺少数据文件")
    }

    config := common.InitConfig(CONF)
    log.Warn("config is: %+v", config)
    runtime.GOMAXPROCS(config.GO_MAX_PROCS)//runtime.NumCPU()

    t1 := time.Now().UnixNano()
    jobs := make(chan string, 1024)
    lineCount := 0

    // 限制 worker 并发量
    work_num := config.THREAD_POOL
    // 生成 worker
    wg := new(sync.WaitGroup)
    wg.Add(work_num)
    for w := 1; w <= work_num; w++ {
        go worker(w, jobs, wg, config.PIPELINE_JOB)
    }

    // 生成 jobs
    f, err := os.Open(args[1])
    if err != nil {
        //panic(err)
        fmt.Println("Have no data file found！")
        fmt.Println("Press「CTRL + C」to stop this program.")
        fmt.Scanln()
    }
    defer f.Close()

    scanner := bufio.NewScanner(bufio.NewReaderSize(f, config.FILE_READBUFFER))
    for scanner.Scan() {
        lineCount += 1
        // add linenumber to the line head
        jobs <- fmt.Sprintf("%d\x00%s", lineCount, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Critical(err)
    }
    close(jobs)

    // 同步 worker
    wg.Wait()

    t2 := time.Now().UnixNano()
    log.Warn("扫描数据 %d lines 总耗时 %d ms", lineCount, (t2-t1)/1000000) //总耗时 226278 ms
}

func worker(id int, jobs <- chan string, wg *sync.WaitGroup, colNum int) {
    defer func(id int) {
        if p := recover(); p != nil {
            err := p.(error)
            log.Critical("worker %3d internal error: %v", id, err)
        }
        wg.Done()
    }(id)

    // main thread send
    for j := range jobs {
        //log.Info("worker %d processing job: %s", id, j)
        //fmt.Println(j)
        utf8J := mahonia.NewDecoder("gb18030").ConvertString(string(j))

        colList := strings.Split(utf8J, "\x00")  // \0
        colCnt := colNum + 1
        // colNum + 1(linenumber)
        if len(colList) != colCnt {
            log.Error("worker %3d processing job: number of fields [%2d] does not equal %2d", id, len(colList), colCnt)
            //TODO: 记录 j 到数据文件以便重新导入
            fmt.Fprintln(os.Stderr, j)
            continue
        } else {
            // trim the linenumber
            colList = colList[1:]
        }
    }

    //log.Warn("thread send quit")
}
