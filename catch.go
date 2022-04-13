package whoops

type panickedError struct {
	error
}

func Catch(fnc func() error) (err error) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}
		perr, ok := recovered.(panickedError)
		if !ok {
			panic(recovered)
		}

		err = perr.error
	}()
	err = fnc()
	return err
}

func Assert(err error) {
	if err != nil {
		panic(panickedError{err})
	}
}

func Must[T any](val T, err error) T {
	return Must2(val, err)
}

func Must1(err error) {
	Assert(err)
}

func Must2[T any](val T, err error) T {
	Assert(err)
	return val
}

func Must3[T1 any, T2 any](val1 T1, val2 T2, err error) (T1, T2) {
	Assert(err)
	return val1, val2
}
