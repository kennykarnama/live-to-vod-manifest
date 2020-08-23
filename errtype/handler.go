package errtype

func OnErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}
