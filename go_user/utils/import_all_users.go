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
    log.AddFilter("file", l4g.WARNING, l4g.NewFileLogWriter(fmt.Sprintf("log/error_%d_%s.log", os.Getpid(), time.Now().Format("20060102150405")), false))
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
    runtime.GOMAXPROCS(config.GO_MAX_PROCS) // runtime.NumCPU()

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
        go worker(w, jobs, wg, redisPool)
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

func worker(id int, jobs <- chan string, wg *sync.WaitGroup, redisPool *redis.Pool) {
    defer wg.Done()
    defer func(id int) {
        if p := recover(); p != nil {
            err := p.(error)
            log.Critical("worker %3d internal error: %v", id, err)
        }
    }(id)

    //每个 worker 保持一个连接
    c := redisPool.Get()
    defer c.Close()

    fieldsName := []string{
        "u_i",    // user_id
        "u_t",    // user_type
        "u_n",    // user_name
        "pwd",    // password
        "i_pwd",  // IVR_passwd
        "pwd_q",  // pwd_question
        "pwd_a",  // pwd_answer
        "r_n",    // real_name
        "sex",    // sex
        "b_d",    // born_date
        "cou",    // country
        "i_t",    // id_type
        "i_n",    // id_no
        "mob",    // mobile_no
        "p_no",   // phone_no
        "eml",    // email
        "addre",  // address
        "pcode",  // postalcode
        "is_a",   // is_active
        "c_code", // check_code
        "l_l_t",  // last_login_time
        "t_t",    // total_times
        "c_cla",  // credit_class
        "if_rec", // if_receive
        "r_time", // regist_time
        "is_va",  // is_valid
        "s_mode", // sale_mode
        "l_chan", // login_channel
        "m_id",   // member_id
        "m_le",   // member_level
        "en_f",   // encourage_flag
        "p_f",    // phone_flag
        "c_i_f",  // check_id_flag
        "u_sta",  // user_status
        "p_lim",  // passenger_limit
        "f_1",    // flag1
    }

    for j := range jobs {
        //log.Info("worker %d processing job: %s", id, j)
        //fmt.Println(j)
        utf8J := mahonia.NewDecoder("gb18030").ConvertString(string(j))

        colList := strings.Split(utf8J, "\x00")  // \0
        // 36 + 1(linenumber)
        if len(colList) != 37 {
            log.Error("worker %3d processing job: number of fileds [%2d] does not equal 37", id, len(colList))
            //TODO: 记录 j 到数据文件以便重新导入
            fmt.Fprintln(os.Stderr, j)
            continue
        } else {
            // trim the linenumber
            colList = colList[1:]
        }

        var args = []interface{}{ colList[2] }
        for i := 0; i < len(colList); i++ {
            colList[i] = strings.TrimSpace(colList[i])
            args = append(args, fieldsName[i], colList[i])
        }

        _, rer := c.Do("hmset", args...)
        if rer != nil {
            log.Error("worker %3d processing job: %s ==> %q", id, colList[2], rer)
            //TODO: 记录 j 到数据文件以便重新导入
            fmt.Fprintln(os.Stderr, j)
            //先推送到chan最后统一处理
        } else {
            //资源消耗非常高(level 需要低一些)
            log.Debug("worker %3d processing job: %s ==> %s", id, colList[2])
        }
        //log.Info("worker %2d processing job: %q", id, colList)
        //runtime.Gosched()
    }
}
