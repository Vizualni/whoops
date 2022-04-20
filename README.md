# Whoops!

[![codecov](https://codecov.io/gh/Vizualni/whoops/branch/master/graph/badge.svg?token=RLFQKJFKI6)](https://codecov.io/gh/Vizualni/whoops)

![whoops logo](./whoops.png)

_Logo inspired by the Awkward Look Monkey Puppet meme_.

-----

Whoops is a small helper Go library to help you with errors in Go. I've decided to create this library because I've noticed that I am using basically the same pattern
all over and I wanted to make a library to help me save some time.



# Error types


## String

A simple string error that you can define as a constant.

```go
import "github.com/vizualni/whoops"
const ErrTheAfwulThingHasHappened = whoops.String("something bad has happened")
```

## Errorf

An error that can be formatted. It can also be defined as a constant.

Usually, you'd do something like:
```go

func foo() error {
	// ...
	var answer = 42
	return fmt.Errorf("the answer is: %d", answer)
}
```

which would make it harder to use error assertion and you'd most likely end up doing a string comparison with errors to figure out which error has happened, and you'd
be unable to use `errors.Is` to help you out.

See example below to get the idea behind the `Errorf` type of error.

```go
import "github.com/vizualni/whoops"
const ErrPayloadSizeTooLarge = whoops.Errorf("payload size too big. got %d bytes")

// ...
return ErrPayloadSizeTooLarge.Format(len(payload))


// later in code you can do
var err error
if whoops.Is(err, ErrPayloadSizeTooLarge) {
	// do your thing
}
```


## Enrich errors with custom fields

Enrich your errors with extra information. e.g. logging the query that failed without query string going into the error message making it unreadable.
It uses Go 1.18 generic's feature to ensure type safety for the fields.

```go
import "yourpackage"
import "github.com/vizualni/whoops"

const ErrFieldUser[yourpackage.User] = "user"
const ErrFieldQuery[yourpackage.Query] = "query"

func process() error {
	var (
		err error
		user yourpackage.User
		query yourpackage.Query
	)
	// ...
	return whoops.Enrich(err, ErrFieldUser.Val(user), ErrFieldQuery.Val(query))
}

// not the best example, but you get the picture
func caller() {
	err := process()
	if err != nil {
		var (
			query yourpackage.Query
			ok bool
		)
		if query, ok = ErrFieldQuery.GetFrom(err); ok {
			// log the query that failed
			log.Error("error %s with query: %s", err, query.Text())	
			// ...
		}
	}
	// ...
}
```


## Grouping errors

If you are processing many things at once and you want to group all those errors together.

### Example 1

```go
import "github.com/vizualni/whoops"

func foo() error {
	var groupErr whoops.Group

	// processN function returns an error only (can return nil as well)
	groupErr.Add(proces1()) 
	groupErr.Add(proces2()) 
	groupErr.Add(proces3()) 

	if groupErr.Err() {
		return groupErr
	}
	// success
	return nil
}
```

### Example 2

```go
import "github.com/vizualni/whoops"
const (
	ErrBad1 = whoops.String("something bad has happened")
	ErrBad2 = whoops.Errorf("format this: %s")
)

// ...
var groupErr whoops.Group
groupErr.Add(ErrBad1) 
groupErr.Add(ErrBad2.Format("foobar")) 

whoops.Is(groupErr, ErrBad1) // returns true
whoops.Is(groupErr, ErrBad2) // returns true
```

## Wrapping errors

Wrapping is very much alike like grouping but it's mostly better for only wrapping two errors together.

```go
import "github.com/vizualni/whoops"
const (
	ErrUnableToIncrementCounter = whoops.String("unable to increment counter")
)

func incrementUserCounter(id int) error {
	// ...
	err := incrementCounterForUser(id)
	if err != nil {
	   return whoops.Wrap(err, ErrUnableToIncrementCounter)
	}
	// ...
}
```

## Errors with stacktrace

All error types have method called `Trace`, which called, taked the current stack trace and creates a new error with the stack trace attached.
This makes things easier for debugging.

```go
import "github.com/vizualni/whoops"
const (
	ErrFoo = whoops.String("foo")
	ErrBar = whoops.Errorf("hello: %s")
)

func bar() error {
	// ...
	if err != nil {
	   return whoops.Trace(err)
	}
	// or
	return ErrFoo.Trace()
	// or
	return ErrBar.Format("Alice").Trace()
	// or
	return Group{ErrFoo, ErrBar.Format("Bob")}.Trace()
	// or
	return Wrap(ErrFoo, ErrBar.Format("Bob")).Trace()
	// ...
}
```


-----

# if err != nil { return err  }

The idea is not to get rid of `err != nil` but rather to eliminate the need to write it all the time.


## Assert

```go
import "github.com/vizualni/whoops"

err := functionCall()
whoops.Assert(err) // panics if err is nil and the value that was panicked is the (wrapped) error
```

## Must

```go
import "github.com/vizualni/whoops"
import "json"

var myMap map[string]any
// ...
var bytes []byte
bytes = whoops.Must(json.Marshal(myMap))
```

## Must1, Must2 && Must3

## Try

```go
import "github.com/vizualni/whoops"
import "json"

// ...
err := whoops.Try(func(){
	bytes := whoops.Must(json.Marshal(myMap))
	bytesWritten := whoops.Must(file.Write(bytes))
})
if err != nil {
	return err
}

```

## TryVal

```go
import "github.com/vizualni/whoops"
import "json"

// ...
n, err := whoops.TryVal(func() int{
	bytes := whoops.Must(json.Marshal(myMap))
	bytesWritten := whoops.Must(file.Write(bytes))
	return bytesWritten
})
if err != nil {
	return err
}

```
