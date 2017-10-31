package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
)

const (
	appName = "Password Generator"
	appDesc = "Generates random passwords in plaintext/base64"
)

func main() {
	app := cli.App(appName, appDesc)

	length := app.Int(cli.IntOpt{
		Name:   "l",
		Desc:   "Length of generated passwords",
		EnvVar: "LENGTH",
		Value:  32,
	})
	count := app.Int(cli.IntOpt{
		Name:   "n",
		Desc:   "Number of passwords to generate",
		EnvVar: "COUNT",
		Value:  1,
	})
	safe := app.String(cli.StringOpt{
		Name:   "c",
		Desc:   "Characters that are Safe",
		EnvVar: "SAFE_CHARACTERS",
		Value:  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]~`",
	})

	app.Action = func() {
		for i := 0; i < *count; i++ {
			p, err := genPass(*length, []byte(*safe))
			if err != nil {
				log.Fatal(err)
			}
			b := base64.URLEncoding.EncodeToString(p)
			fmt.Println(string(p), ":", b)
		}
	}

	app.Run(os.Args)
}

func genPass(length int, chars []byte) ([]byte, error) {
	password := make([]byte, length)
	random := make([]byte, length+(length/4))
	l := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, random); err != nil {
			if err != nil {
				return nil, err
			}
		}
		for _, c := range random {
			if c >= maxrb {
				continue
			}
			password[i] = chars[c%l]
			i++
			if i == length {
				return password, nil
			}
		}
	}
}
