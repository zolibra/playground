package main

import "fmt"

func main() {
    var ptr *int
    fmt.Println(*ptr) // This will cause a null pointer dereference
}
