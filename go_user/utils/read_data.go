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
    "time"
)

const (
    FILE = "../data/data.bcp"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    t1 := time.Now().UnixNano()

    // 生成 jobs
    f, err := os.Open(FILE)
    if err != nil {
        //panic(err)
        fmt.Println("Have no data file found！")
        fmt.Println("Press「CTRL + C」to stop this program.")
        fmt.Scanln()
    }
    defer f.Close()

    scanner := bufio.NewScanner(bufio.NewReaderSize(f, 8 * 1024)) // Buffer: 8K
    for scanner.Scan() {
        line := scanner.Text()
        line = "."
        fmt.Print(line)
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }

    t2 := time.Now().UnixNano()
    fmt.Fprintf(os.Stderr, "总耗时 %d ms\n", (t2-t1)/1000000) //总耗时 2118376 ms
}

