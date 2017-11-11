package main

import (
	"log"
	"os"
)

// InitCache Initialize Filesystem cache system
func InitCache() {

	log.Print("Init cache")
	_, err := os.Stat("_cache")
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("_cache", os.FileMode(0755))
			if err != nil {
				log.Fatal("Can't create _cache folder")
			} else {
				log.Print("_cache folder created")
			}
		}
	}

}
