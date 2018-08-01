// UsualPassenger InitData Program
//
// This program used to import all user's contacts(UsualPassenger) to Codis
//
// Copyright (C)2017 EmbraceSource Co.,Ltd. All Rights Reserved
//
// Author: yufu.zhao <tony.zhao@embracesource.com> at 2017-04-14 10:08:35

package main

import (
    "runtime"
    "runtime/pprof"
    "../common"
    "sync"
    "fmt"
    "os"
    "bufio"
    "strings"
    "time"
    _ "github.com/axgle/mahonia"
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
}

func main() {
    defer log.Close()

    args := os.Args
    if len(args)==2 && (args[1]=="--version" || args[1] =="-v") {
        fmt.Printf("Git Commit Hash: %s\n", githash)
        fmt.Printf("UTC Build Time : %s\n", buildstamp)
        return
    }

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
        go worker(w, jobs, wg)
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

    // start profiling
    pf, err := os.OpenFile("cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        log.Critical(err)
    }
    defer pf.Close()
    pprof.StartCPUProfile(pf)

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

    // stop profiling
    pprof.StopCPUProfile()
    pf.Close()

    t2 := time.Now().UnixNano()
    log.Warn("扫描数据 %d lines 总耗时 %d ms", lineCount, (t2-t1)/1000000) //总耗时 226278 ms
}

func worker(id int, jobs <- chan string, wg *sync.WaitGroup) {
    defer func(id int) {
        if p := recover(); p != nil {
            err := p.(error)
            log.Critical("worker %3d internal error: %v", id, err)
        }
    }(id)
    defer wg.Done()

    for j := range jobs {
        //utf8J := mahonia.NewDecoder("gb18030").ConvertString(string(j))
        //fmt.Fprintln(os.Stdout, "\x1F\x1F\n")

        if strings.Contains(j, "\x1F") {
            log.Error("UR char found")
            //TODO: 记录 j 到数据文件以便重新导入
            //fmt.Fprintln(os.Stderr, fmt.Sprintf("%d\x00%s", lineCount, j))
        }
        //runtime.Gosched()
    }
}

