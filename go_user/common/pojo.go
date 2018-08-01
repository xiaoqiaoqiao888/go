package common

import (
    "sync"
    "sort"
)

// Raw 结果集
type Cmap struct {
    Data map[int]uint64
    Lock sync.RWMutex
}

func NewCmap(size int) *Cmap {
    cmap := new(Cmap)
    cmap.Data = make(map[int]uint64, size)
    return cmap
}

func (d *Cmap) Get(k int) uint64 {
    d.Lock.RLock()
    defer d.Lock.RUnlock()
    return d.Data[k]
}

func (d *Cmap) Set(k int, v uint64) uint64 {
    d.Lock.Lock()
    defer d.Lock.Unlock()
    d.Data[k] = v
    return d.Data[k]
}

func (d *Cmap) Incr(k int, v uint64) uint64 {
    d.Lock.Lock()
    defer d.Lock.Unlock()
    d.Data[k] += v
    return d.Data[k]
}

// Raw 结果集
type Csmap struct {
    Data map[string]int
    Lock sync.RWMutex
}

func NewCsmap(size int) *Csmap {
    csmap := new(Csmap)
    csmap.Data = make(map[string]int, size)
    return csmap
}

func (d *Csmap) Get(k string) int {
    d.Lock.RLock()
    defer d.Lock.RUnlock()
    return d.Data[k]
}

func (d *Csmap) Set(k string, v int) int {
    d.Lock.Lock()
    defer d.Lock.Unlock()
    d.Data[k] = v
    return d.Data[k]
}

// sorttd string map
type sortedMap struct {
    m map[string]int
    s []string
}

func (sm *sortedMap) Len() int {
    return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
    return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
    sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func SortedKeys(m map[string]int) []string {
    sm := new(sortedMap)
    sm.m = m
    sm.s = make([]string, len(m))
    i := 0
    for key, _ := range m {
        sm.s[i] = key
        i++
    }
    sort.Sort(sm)
    return sm.s
}

// sorttd int map
type sortedIntMap struct {
    m map[int]uint64
    s []int
}

func (sm *sortedIntMap) Len() int {
    return len(sm.m)
}

func (sm *sortedIntMap) Less(i, j int) bool {
    return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedIntMap) Swap(i, j int) {
    sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func SortedIntKeys(m map[int]uint64) []int {
    sm := new(sortedIntMap)
    sm.m = m
    sm.s = make([]int, len(m))
    i := 0
    for key, _ := range m {
        sm.s[i] = key
        i++
    }
    sort.Sort(sm)
    return sm.s
}