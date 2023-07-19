package main

import (
	"fmt"
	"github.com/korvised/go-ecommerce/config"
	"os"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {

	cfg := config.LoadConfig(envPath())

	fmt.Println(cfg.App().Url())
	fmt.Println(cfg.Db().Url())

}
