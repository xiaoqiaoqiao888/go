// UsualPassenger InitData Program
//
// This program used to import all user's contacts(UsualPassenger) to Codis
//
// Copyright (C)2017 EmbraceSource Co.,Ltd. All Rights Reserved
//
// Author: yufu.zhao <tony.zhao@embracesource.com> at 2017-04-20 11:34:45

package main

import (
    "runtime"
    "fmt"
    "os"
    "bufio"
    "time"
    "strconv"
)

const (
    FILE = "../data/data.bcp"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    t1 := time.Now().UnixNano()

    lines := fileTolines(FILE)
    for i := 1; i <= 29790; i++ {
        for _, line := range lines {
            line = strconv.Itoa(i) + line
            fmt.Println(line)
        }
    }

    t2 := time.Now().UnixNano()
    fmt.Fprintf(os.Stderr, "总耗时 %d ms\n", (t2-t1)/1000000) //总耗时 2118376 ms
}

func fileTolines(filePath string) []string {
    f, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    var lines []string
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }

    return lines
}