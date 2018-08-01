package lib

import (
    "testing"
    "strings"
    "fmt"
)

const STR = "abcdefg"
const DELIMITER = "d"


func TestStrings(t *testing.T) {
    return
}

func BenchmarkStrContains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strings.Contains(STR, DELIMITER)
    }
}

func BenchmarkStrContainsAny(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strings.ContainsAny(STR, DELIMITER)
    }
}

func BenchmarkStrSplit(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strings.Split(STR, DELIMITER)
    }
}

func BenchmarkStrReplace(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strings.Replace(STR, DELIMITER, "", -1)
    }
}

func BenchmarkStrNewReplacer(b *testing.B) {
    r := strings.NewReplacer(DELIMITER, "")
    for i := 0; i < b.N; i++ {
        r.Replace(STR)
    }
}

func BenchmarkStrFormat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("%s%s", STR, DELIMITER)
    }
}

func BenchmarkStrPlus(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = STR + DELIMITER
    }
}

func BenchmarkStrJoin(b *testing.B) {
    colList := []string{STR, DELIMITER}
    for i := 0; i < b.N; i++ {
        _ = strings.Join(colList, "")
    }
}
