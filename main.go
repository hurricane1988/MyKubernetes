package main

import (
	"log"
	"os"
)

func main() {
	a := os.Getegid()
	log.Fatalf("当前系统版本信息: %s\n", a)
}
