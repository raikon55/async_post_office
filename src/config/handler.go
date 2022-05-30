package config

import "log"

func HandlerError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}
