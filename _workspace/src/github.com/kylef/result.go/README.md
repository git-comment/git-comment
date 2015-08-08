# Result

Result is a go structure to represent a Result type. It can represent either a
success or a failure.

Using Result, we can build APIs like:


```go
latestCommitMessage := openRepository(url)
    .flatMap(func (repo interface{})   { return repo.headReference() })
    .flatMap(func (ref interface{})    { return ref.commit() })
    .flatMap(func (commit interface{}) { return commit.message() })

fmt.Printf("%+v\n", latestCommitMessage)
```

Instead of like:

```go
repository, err := openRepository(url)
if err != nil {
    return
}

reference, err := repository.headReference()
if err != nil {
    return
}

commit, err := reference.commit()
if err != nil {
    return
}

message, err := commit.message()
if err != nil {
    return
}

fmt.Printf("%s\n", message)
```

## Usage

#### Creating a success result

```go
result := NewSuccess(5)
```

#### Creating a failure result

```go
result := NewFailure(err)
```

#### Creating a new result from a value or error

Constructing a result from the "result" of the `Open` standard os function.

```go
func Open(name string) (file *File, err error)
```

```go
result := NewResult(Open("HelloWorld.txt"))
```

#### Transform a success result into another result

This example will transform a the original result and return a result with a
value that has doubled, when the original result is a failure, the failure will
be returned.


```go
result := result.FlatMap(func(value interface{}) Result {
    return NewSuccess(value.(int) * 2)
})
```

#### Dematerializing a result into a result tuple

```go
value, err := result.Dematerialize()
```

## Credits

Result in Go was heavily inspired by [Result in Swift](https://github.com/antitypical/Result) and parts of the README we're inspired by [Gift](https://github.com/modocache/Gift), a Swift git binding.

## License

Result is licensed under the BSD license. See [LICENSE](LICENSE) for more
info.

