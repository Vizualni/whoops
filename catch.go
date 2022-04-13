package whoops

type panickedError struct {
	error
}

func _catchPanic(orig *error) {
	recovered := recover()
	if recovered == nil {
		return
	}
	perr, ok := recovered.(panickedError)
	if !ok {
		panic(recovered)
	}
	*orig = perr.error
}

func Catch(fnc func()) (err error) {
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
	fnc()
	return nil
}

func CatchVal[T any](fnc func() T) (t T, err error) {
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
	t = fnc()
	return
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

func DeferWrap(wrapWith error) func() {
	return func() {
		recovered := recover()
		if recovered == nil {
			return
		}
		perr, ok := recovered.(panickedError)
		if !ok {
			panic(recovered)
		}
		panic(panickedError{
			error: Wrap(perr.error, wrapWith),
		})
	}
}
