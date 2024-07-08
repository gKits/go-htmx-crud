package main

import (
	"crud"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	driver := flag.String("driver", "sqlite3", "select the database driver between postgres, sqlite3")
	conn := flag.String("conn", "crud.db", "set connection string")
	port := flag.Int("port", 4000, "set port")
	flag.Parse()

	dir, _ := os.Getwd()
	fmt.Println(dir)

	repo, err := crud.NewRepository(*driver, *conn)
	if err != nil {
		log.Fatal(err)
	}

	ctrl := crud.NewUserController(repo)

	srv := crud.NewServer(fmt.Sprintf("0.0.0.0:%d", *port))
	srv.Attach("view", &ctrl)

	srv.Run()
}
