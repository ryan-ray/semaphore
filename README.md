# Semaphore

Don't use this package in your project. The code is small enough that you should
just copy and paste it, or rewrite it based on what you're doing.

This package is an implementation of the semaphore pattern. I've mostly put it
together because over the years I have had to show engineers new to Go that
structs are not objects, and that you don't need a struct in every situtation if
you just need methods to operate on an underlying type.

That being said, if you want to use this, go for gold. It works somewhat like
the following;

```
aquire, release := semaphore.New(10)

for range 10000 {
    aquire()
    go func() {
        defer release()
        SomeFunctionThatYouNeedToThrottle()
    }()
}
```

There is also a `WithCancel` method which returns an additional function to
cancel the semaphore, closing and draining the underlying channel. Calling the
cancel function will cause all subsequent calls to the aquire and release
functions to return an error.
