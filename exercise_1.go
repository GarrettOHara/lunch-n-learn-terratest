package main

import "fmt"

func main() {
    var name string
    fmt.Println("What's your name?")
    fmt.Scanln(&name)
    fmt.Printf("Hi, %s! I'm Go!\n", name)
}
