# Terratest Lunch N' Learn

## Pre-requisites

- Install Terraform: https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli
- Install Golang: https://go.dev/doc/install


## History of Go

Go attempts to combine the ease of programming of an interpreted, dynamically typed 
language with the efficiency and safety of a statically typed, compiled language. 
It also aimed to be better adapted to current hardware, with support for networked 
and multicore computing. Finally, working with Go is intended to be fast: it should 
take at most a few seconds to build a large executable on a single computer.

![Go Gopher](./images/go-gopher.jpeg)

## Golang Introduction

Basic program in Go: [source](https://www.digitalocean.com/community/tutorials/how-to-write-your-first-program-in-go)
```Golang
package main

import "fmt"

func main() {
    var name string
    fmt.Println("What's your name?")
    fmt.Scanln(&name)
    fmt.Printf("Hi, %s! I'm Go!", name)
}
```
Output:
```
~/software/lunch-n-learn-terratest main % go run example.go
What's your name?
Garrett
Hi, Garrett! I'm Go!
```

Let's analyze the components: 

- `package`: a named collection of one or more related `.go` files.
    - Helps you isolate and reuse code.
    - Every .go file that you write should begin with a package
    - Can only have **1** package per directory.
- `&` the *ampersand* Gets/returns the memory address of a variable
    - In this instance, we reference the memory address of `name`
    which is an integer 
    - **Why use a  pointer?** `Scanln()` requires a memory address
    to write the user input to.
        - see more here: https://pkg.go.dev/fmt#Scanln
- The rest is simple.
