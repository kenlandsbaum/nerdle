package errs

import "log"

func LogFatalIfErr(err error) {
	if err != nil {
		log.Fatalf("fatal error %s", err)
	}
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
