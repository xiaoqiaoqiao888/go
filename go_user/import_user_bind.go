// UsualPassenger InitData Program
//
// This program used to import all user_bind to Codis
//
// Copyright (C)2017 EmbraceSource Co.,Ltd. All Rights Reserved
//
// Author: yufu.zhao <tony.zhao@embracesource.com> at 2018/7/11 15:07:29

package main

import (
    "runtime"
    "./common"
    "sync"
    "fmt"
    "os"
    "bufio"
    "strings"
    "time"
    "github.com/garyburd/redigo/redis"
    "github.com/axgle/mahonia"
    l4g "github.com/ms2008/log4go"
)

const (
    CONF = "config.conf"
)

var log = make(l4g.Logger)
var buildstamp = ""
var githash = ""


func init() {
    log.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())
    log.AddFilter("file", l4g.WARNING, l4g.NewFileLogWriter(fmt.Sprintf("log/user_bind_%d_%s.log", os.Getpid(), time.Now().Format("20060102150405")), false))
}

func main() {
    defer log.Close()

    args := os.Args
    if len(args)==2 && (args[1]=="--version" || args[1] =="-v") {
        fmt.Printf("Git Commit Hash: %s\n", githash)
        fmt.Printf("UTC Build Time : %s\n", buildstamp)
        return
    }

    log.Warn("导入数据开始...")
    log.Warn("os args: %v", os.Args)

    if args == nil || len(args) < 2 {
        panic("缺少数据文件")
    }

    config := common.InitConfig(CONF)
    log.Warn("config is: %+v", config)
    runtime.GOMAXPROCS(config.GO_MAX_PROCS)//runtime.NumCPU()

    redisPool := common.NewRedisPool(config)
    defer redisPool.Close()

    t1 := time.Now().UnixNano()
    jobs := make(chan string, 1024)
    lineCount := 0

    // 限制 worker 并发量
    work_num := config.THREAD_POOL
    // 生成 worker
    wg := new(sync.WaitGroup)
    wg.Add(work_num)
    for w := 1; w <= work_num; w++ {
        go worker(w, jobs, wg, redisPool, config.PIPELINE_JOB)
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
    log.Warn("导入数据 %d lines 总耗时 %d ms", lineCount, (t2-t1)/1000000) //总耗时 226278 ms
}

func worker(id int, jobs <- chan string, wg *sync.WaitGroup, redisPool *redis.Pool, pipelineJob int) {
    // defer wg.Done()
    defer func(id int) {
        if p := recover(); p != nil {
            err := p.(error)
            log.Critical("worker %3d internal error: %v", id, err)
        }
        wg.Done()
    }(id)

    // 每个 worker 保持一个连接
    c := redisPool.Get()
    defer c.Close()

    var fieldsName = []string{
        "a", // register_name
        "b", // bind_id
        "c", // bind_type
        "e", // if_myself
        "f", // nickname
        "g", // uniq_code
        "l", // pic_path
        "i", // id_no
        "j", // id_type
        "k", // real_name
        "l", // flag1
        "m", // flag2
        "n", // flag3
        "d", // bind_datetime
    }

    quit := make(chan bool)
    sema, msg := make(chan int), make(chan []string)
    chunkSize, chunkData := pipelineJob, make([]string, pipelineJob, pipelineJob)
    jobCount := 0

    rt := new(sync.WaitGroup)
    rt.Add(1)

    // start recv thread
    go recv(id, rt, c, quit, sema, msg)

    // main thread send
    for j := range jobs {
        //log.Info("worker %d processing job: %s", id, j)
        //fmt.Println(j)
        utf8J := mahonia.NewDecoder("gb18030").ConvertString(string(j))

        colList := strings.Split(utf8J, "\x00")  // \0
        // 14 + 1(linenumber)
        if len(colList) != 15 {
            log.Error("worker %3d processing job: number of fields [%2d] does not equal 15", id, len(colList))
            //TODO: 记录 j 到数据文件以便重新导入
            fmt.Fprintln(os.Stderr, j)
            continue
        } else {
            // trim the linenumber
            colList = colList[1:]
        }

        var isLogged = false
        var errMsg error

        // colList[0] = register_name colList[2] = bind_type
        var args = []interface{} { "Bind:" + strings.TrimSpace(colList[0]) + "+" + strings.TrimSpace(colList[2]) }

        for i := 0; i < len(colList); i++ {
            colList[i] = strings.TrimSpace(colList[i])

            // reformat the bind_datetime
            switch i {
            case 13:
                placeHolder := []byte(colList[i])
                placeHolder[len(placeHolder)-6] = 46 // replace to "." due to go's layout string
                placeHolderStr, err := time.Parse("Jan 2 2006 3:04:05.000PM", string(placeHolder))
                if err != nil {
                    isLogged, errMsg = true, err
                    goto NEXT
                }
                colList[i] = placeHolderStr.Format("20060102150405")
            }

            args = append(args, fieldsName[i], colList[i])
        }

        // round-1: init Bind:
        if err := c.Send("hmset", args...); err != nil && ! isLogged {
            isLogged, errMsg = true, err
            goto NEXT
        } else {
            _ = c.Send("ping")
        }

        // round-2: init BindReg:
        if colList[0] != "" {
            if err := c.Send("sadd", "BindReg:"+colList[0], colList[2]); err != nil && ! isLogged {
                isLogged, errMsg = true, err
                goto NEXT
            }
        } else {
            _ = c.Send("ping")
        }

        // round-3: init BindBId:
        if colList[1] != "" {
            if err := c.Send("sadd", "BindBId:"+colList[1], colList[0]); err != nil && ! isLogged {
                isLogged, errMsg = true, err
                goto NEXT
            }
        } else {
            _ = c.Send("ping")
        }

        // round-4: init BindUCd:
        if colList[5] != "" {
            if err := c.Send("sadd", "BindUCd:"+colList[5], colList[0]); err != nil && ! isLogged {
                isLogged, errMsg = true, err
                goto NEXT
            }
        } else {
            _ = c.Send("ping")
        }

NEXT:
        if isLogged {
            log.Error("worker %3d processing job: %s ==> %q", id, colList[0], errMsg)
            //TODO: 记录 j 到数据文件以便重新导入
            fmt.Fprintln(os.Stderr, j)
            continue
        }

        jobCount += 1
        // flush the buffer
        if jobCount % chunkSize == 0 {
            chunkData[chunkSize-1] = j

            err := c.Flush()
            if err != nil {
                log.Critical(err)
                continue
            } else {
                sema <- chunkSize * 4
                // sync send/recv message
                msg <- append([]string{}, chunkData...)
            }
        } else {
            chunkData[(jobCount%chunkSize)-1] = j
        }
        //log.Info("worker %2d processing job: %q", id, colList)
        //runtime.Gosched()
    }

    // handle the tail
    err := c.Flush()
    if err != nil {
        log.Critical(err)
        quit <- true
    } else {
        sema <- (jobCount % chunkSize) * 4
        msg <- append([]string{}, chunkData...)
        quit <- true
    }

    rt.Wait()
    //log.Warn("thread send quit")
}

// thread recv
func recv(id int, rt *sync.WaitGroup, c redis.Conn, quit chan bool, sema chan int, msg chan []string) {
    defer rt.Done()

    for {
        select {
            case <- quit:
                //log.Warn("thread recv abort")
                return

            case n := <- sema:
                msgs := <- msg

                for i:=0; i<n; i++ {
                    _, err := c.Receive()
                    //log.Info(err, msgs[i-1])
                    if err != nil {
                        log.Error("%s", err)
                        fmt.Fprintln(os.Stderr, msgs[i/4])
                    } else {
                        log.Debug("worker %3d processing job: success", id)
                    }
                }
        }
    }
}
