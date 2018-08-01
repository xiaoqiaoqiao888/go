// UsualPassenger InitData Program
//
// This program used to import all user's contacts(UsualPassenger) to Codis
//
// Copyright (C)2017 EmbraceSource Co.,Ltd. All Rights Reserved
//
// Author: yufu.zhao <tony.zhao@embracesource.com> at 2017-04-14 10:08:35

package main

import (
    "sync"
    "fmt"
    "time"
    "math/rand"
    "github.com/garyburd/redigo/redis"
    l4g "github.com/ms2008/log4go"
)

var log = make(l4g.Logger)
var buildstamp = ""
var githash = ""


func init() {
    log.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())
}

func main() {
    defer log.Close()

    c, err := redis.Dial("tcp", "127.0.0.1:16380")
    if err != nil {
        //return nil, err
        panic(err)
    }
    defer c.Close()

    quit := make(chan bool)
    sema := make(chan int)
    msg := make(chan [10]string)

    rt := new(sync.WaitGroup)
    rt.Add(1)

    // thread recv
    go func(rt *sync.WaitGroup) {
        defer rt.Done()

        for {
            select {
                case <- quit:
                    log.Warn("thread recv abort")
                    return

                case n := <- sema:
                    msgs := <- msg

                    for i:=1; i<= n; i++ {
                        _, err := c.Receive()
                        log.Info(err, msgs[i-1])
                        if err != nil {
                            log.Error("%s", err)
                        }
                    }
            }
        }
    }(rt)

    rand.Seed(time.Now().Unix())
    jobNum := rand.Intn(30)
    chunkSize := 10
    chunkData := [10]string{}

    // main thread send
    for i := 1; i <= jobNum; i++ {
        // assembled to redis pipeline
        err := c.Send("set", fmt.Sprintf("test:%d", i), "test string")
        log.Info(err)
        if err != nil {
            log.Critical(err)
            continue
        }

        // flush the buffer
        if i % chunkSize == 0 {
            chunkData[9] = fmt.Sprintf("test:%d", i)

            err := c.Flush()
            if err != nil {
                log.Critical(err)
                continue
            } else {
                sema <- chunkSize
                // sync send/recv message
                msg <- chunkData
            }
        } else {
            chunkData[(i%chunkSize)-1] = fmt.Sprintf("test:%d", i)
        }
    }

    // handle the tail
    err = c.Flush()
    if err != nil {
        log.Critical(err)
        quit <- true
    } else {
        sema <- jobNum % chunkSize
        msg <- chunkData
        quit <- true
    }

    rt.Wait()
    log.Warn("thread send quit")
}
