// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 13.
//!+

// Modify dup2 to print the names of all files in which each duplicated line occurs.
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type set map[string]bool

func main() {
    counts := make(map[string]int)
    dupFiles := make(map[string]set)
    files := os.Args[1:]
    if len(files) == 0 {
        countLines(os.Stdin, counts, dupFiles)
    } else {
        for _, arg := range files {
            f, err := os.Open(arg)
            if err != nil {
                fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
                continue
            }
            countLines(f, counts, dupFiles)
            f.Close()
        }
    }
    for line, n := range counts {
        if n > 1 {
            filesSet := dupFiles[line]
            fileNames := make([]string, 0, len(filesSet))
            for f := range filesSet {
                //fmt.Println("hoho: ", f)
                fileNames = append(fileNames, f)
            }
            fmt.Printf("%d\t%s\t%s\n", n, line, strings.Join(fileNames, ","))
        }
    }
}

func countLines(f *os.File, counts map[string]int, dup_files map[string]set) {
    input := bufio.NewScanner(f)
    for input.Scan() {
        s := input.Text()
        counts[s]++
        if val, ok := dup_files[s]; ok {
            val[f.Name()] = true
        } else {
            val = make(set)
            val[f.Name()] = true
            dup_files[s] = val
        }
    }
    // NOTE: ignoring potential errors from input.Err()
}

//!-
