package main

type Request struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Role    string `json:"role"`
}
