package main

import (
	"os"
	"fmt"
	"io"
	"path/filepath"
	"github.com/urfave/cli"
	"github.com/yeka/zip"
)

func main() {

	app := cli.NewApp()
	app.Name="azip"

	app.Action = func(c *cli.Context) error {
		src := c.Args().Get(0)
		dest := c.Args().Get(1)

		err := os.Rename(src, dest)
		if err != nil {
			return err
		}

		err = compress(dest)
		if err != nil {
			return err
		}

		err = os.Remove(dest)
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func compress(filename string) error {

	destWithoutExt := filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])

	fzip, err := os.Create(fmt.Sprintf("./%s.zip", destWithoutExt))
	if err != nil {
		return err
	}

	zipw := zip.NewWriter(fzip)
	defer zipw.Close()

	password := os.Getenv("ZIP_PASS")

	w, err := zipw.Encrypt(filename, password, zip.AES256Encryption)
	if err != nil {
		return err
	}

	zipw.Flush()

	file,err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)

	return err
}