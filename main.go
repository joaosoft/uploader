package main

import (
	"github.com/joaosoft/uploader/models/cmd"
)

func main() {
	uploader, err := cmd.NewUploader()
	if err != nil {
		panic(err)
	}

	if err := uploader.Start(); err != nil {
		panic(err)
	}
}
