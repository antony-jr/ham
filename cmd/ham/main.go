package main

import (
   "os"
  
   "github.com/antony-jr/ham/internal/banner"
   "github.com/antony-jr/ham/internal/cli"

)

func main() {
   banner.Header("v0.0.1-alpha")
   if cli.Run() != nil {
	os.Exit(1)
   }
}
