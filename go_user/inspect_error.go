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
    "fmt"
    "os"
    "bufio"
    "strings"
    "regexp"
    "github.com/axgle/mahonia"
    l4g "github.com/ms2008/log4go"
)

const usage = `
Usage:
    ./inspect-error dcp_file [user|user_bind]
`

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

    var delimiterCnt int
    //log.Warn("os args: %v", os.Args)
    if args == nil || len(args) < 3 {
        fmt.Println(usage)
        panic("缺少数据文件或数据类型")
    } else {
        switch args[2] {

        case "user":
            delimiterCnt = 36 - 1

        case "user_bind":
            delimiterCnt = 14 - 1

        default:
            panic("不支持的数据类型")
        }
    }

    runtime.GOMAXPROCS(runtime.NumCPU())

    f, err := os.Open(args[1])
    if err != nil {
        //panic(err)
        fmt.Println("Have no data file found！")
        fmt.Println("Press「CTRL + C」to stop this program.")
        fmt.Scanln()
    }
    defer f.Close()

    var chunkLine []string
    scanner := bufio.NewScanner(bufio.NewReaderSize(f, 4096))

    for scanner.Scan() {
        j := scanner.Text()
        utf8J := mahonia.NewDecoder("gb18030").ConvertString(string(j))
        colList := strings.Split(utf8J, "\x00")  // \0

        // trim the linenumber
        colList = colList[1:]

        timeEnd := colList[len(colList)-1]
        gb18030j := mahonia.NewEncoder("gb18030").ConvertString(strings.Join(colList, "\x00"))
        chunkLine = append(chunkLine, gb18030j)

        //if strings.HasSuffix(timeEnd, "AM") || strings.HasSuffix(timeEnd, "PM") {
        isTime, _ := regexp.MatchString(":\\d\\d\\d[AP]M$", timeEnd)
        if isTime {
            // merge lines
            row := strings.Join(chunkLine, "")
            if strings.Count(row, "\x00") != delimiterCnt {
                //log.Error("processing job: number of fields [%2d] does not equal 14", strings.Count(row, "\x00"))
                fmt.Fprintln(os.Stdout, row)
            } else {
                fmt.Fprintln(os.Stderr, row)
            }
            // clear chunkLine
            chunkLine = chunkLine[:0]
            continue
        } else {
            continue
        }
    }
    if err := scanner.Err(); err != nil {
        log.Critical(err)
    }
}
