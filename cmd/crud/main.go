package main

import (
	"crud"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	driver := flag.String("driver", "memory", "select the database driver between postgres, sqlite3")
	conn := flag.String("conn", "", "set connection string (not needed for memory)")
	flag.Parse()

	dir, _ := os.Getwd()
	fmt.Println(dir)

	repo, err := crud.NewRepository(*driver, *conn)
	if err != nil {
		log.Fatal(err)
	}

	ctrl := crud.NewUserController(repo)

	srv := crud.NewServer("0.0.0.0:5050")
	srv.Attach("view", &ctrl)

	srv.Run()
}
