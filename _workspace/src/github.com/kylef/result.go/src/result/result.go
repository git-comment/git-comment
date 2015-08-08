package result


type Result struct {
  Success interface{}
  Failure error
}

// Create a new failure result
func NewFailure(err error) Result {
  result := Result {
    Success: nil,
    Failure: err,
  }

  return result
}

// Create a new success Result
func NewSuccess(value interface{}) Result {
  result := Result {
    Success: value,
    Failure: nil,
  }

  return result
}

// Creates a new result from the given arguments
func NewResult(value interface{}, err error) Result {
  if err != nil {
    return NewFailure(err)
  }
  return NewSuccess(value)
}

// Transform the success value or error of a result
func (result Result) Analysis(ifSuccess func(interface{}) Result, ifFailure func(error) Result) Result {
  if result.Success != nil {
    return ifSuccess(result.Success)
  }

  return ifFailure(result.Failure)
}

// Transform the success value of a result
func (result Result) FlatMap(transform func(interface{}) Result) Result {
  if result.Success != nil {
    return transform(result.Success)
  }

  return result
}

// Return the underlying success value and error or a result
func (result Result) Dematerialize() (value interface{}, err error) {
  return result.Success, result.Failure
}

// Returns the .Success value, or the given fallback value if the result is a failure
func (result Result) Recover(value interface{}) interface{} {
  if result.Failure != nil {
    return value
  }

  return result.Success
}

// Returns this result if it's a success, or the given result
func (result Result) RecoverWith(recoveredResult Result) Result {
  if result.Failure != nil {
    return recoveredResult
  }

  return result
}

// Compose the success values of results if no failures are present, otherwise
// returns the first failing result
func Combine(transform func(...interface{}) Result, results ...Result) Result {
  values := make([]interface{}, len(results))
  for index, result := range results {
    if result.Failure != nil {
      return result
    }
    values[index] = result.Success
  }
  return transform(values...)
}
