/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/spock/spock.go
* @Description: Main. `Emerge WORLD!!`
 */
package main

import (
	"config"
	"ctx"
	"log"
)

var logo = `
  ____                           _    
 / ___|   _ __     ___     ___  | | __
 \___ \  | '_ \   / _ \   / __| | |/ /
  ___) | | |_) | | (_) | | (__  |   < 
 |____/  | .__/   \___/   \___| |_|\_\
         |_|                          
`

// TODO: take this from CLI param;
var defaultConfFileName = "spock.ini"

func main() {
	log.Println(logo + "\n ---------- \n")

	// Load Server Configuration
	log.Println("Loading Server Configuration")

	cf, err := config.New(defaultConfFileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize Application Context
	log.Println("Intialize Application Context")
	err = ctx.Load(cf)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize Engine
	log.Println("3. Initialize Engine")

	// Start Network
	log.Println("4. Starting networking interface")

	// Start Server
	log.Println("Starting Server .... ")

}
