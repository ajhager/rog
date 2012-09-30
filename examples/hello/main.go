package main

import (
    "github.com/ajhager/rog"
    "github.com/ajhager/rog/wde"
)

func main() {
    rog.Open(20, 11, 2, "rog", wde.Backend())
    for rog.IsOpen() {
        rog.Set(5, 5, nil, nil, "Hello, 世界!")
        if rog.Key() == rog.Escape {
            rog.Close()
        }
        rog.Flush()
    }
}
