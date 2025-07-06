package utils

import "log"

// FatalIfErr шорткат для критических ошибок.
func FatalIfErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
