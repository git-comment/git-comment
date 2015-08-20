package result

import (
  "errors"
  "fmt"
  "github.com/stvp/assert"
  "testing"
)

type errorString struct {
  value string
}

func (err *errorString) Error() string {
  return err.value
}


// Test Initialization

func TestFailure(t *testing.T) {
  err := &errorString{"testing error"}
  result := NewFailure(err)

  assert.Nil(t, result.Success)
  assert.Equal(t, result.Failure, err)
}

func TestSuccess(t *testing.T) {
  value := 5
  result := NewSuccess(value)

  assert.Equal(t, result.Success, value)
  assert.Nil(t, result.Failure)
}

func TestNewResultWithSuccess(t *testing.T) {
  example := func() (value interface{}, err error) {
    return 5, nil
  }

  result := NewResult(example())

  assert.Nil(t, result.Failure)
  assert.Equal(t, result.Success, 5)
}

func testNewResultWithFailure(t *testing.T) {
  err := &errorString{"testing error"}

  example := func() (value interface{}, err error) {
    return nil, err
  }

  result := NewResult(example())

  assert.Nil(t, result.Success)
  assert.Equal(t, result.Failure, err)
}

// Test Analysis

func TestFailureAnalysis(t *testing.T) {
  err := &errorString{"testing error"}
  result := NewFailure(err)

  resultantErr := &errorString{"testing new error"}
  resultantResult := result.Analysis(func(value interface{}) Result { return NewSuccess(value) },
                                     func(err error) Result { return NewFailure(resultantErr) })

  assert.Equal(t, resultantResult.Failure, resultantErr)
  assert.Nil(t, resultantResult.Success)
}

func TestSuccessAnalysis(t *testing.T) {
  result := NewSuccess(5)
  resultantResult := result.Analysis(func(value interface{}) Result { return NewSuccess(value.(int) * 2) },
                                     func(err error) Result { return NewFailure(err) })

  assert.Equal(t, resultantResult.Success, 10)
  assert.Nil(t, resultantResult.Failure)
}


// Test FlatMap

func TestFlatMapOnSuccessReturnsNewValue(t *testing.T) {
  result := NewSuccess(5)
  resultantResult := result.FlatMap(func(value interface{}) Result { return NewSuccess(value.(int) * 2) })

  assert.Equal(t, resultantResult.Success, 10)
  assert.Nil(t, resultantResult.Failure)
}

func TestFlatMapOnFailureReturnsFailure(t *testing.T) {
  err := &errorString{"testing error"}
  result := NewFailure(err)
  resultantResult := result.FlatMap(func(value interface{}) Result { return NewSuccess(value.(int) * 2) })

  assert.Equal(t, resultantResult.Failure, err)
  assert.Nil(t, resultantResult.Success)
}

// Test Dematerialize

func TestDematerializeWithSuccess(t *testing.T) {
  result := NewSuccess(5)
  value, err := result.Dematerialize()

  assert.Equal(t, value, 5)
  assert.Nil(t, err)
}

func TestDematerializeWithFailure(t *testing.T) {
  err := &errorString{"testing error"}
  result := NewFailure(err)
  value, resultantErr := result.Dematerialize()

  assert.Nil(t, value)
  assert.Equal(t, resultantErr, err)
}

func TestCombineAllSuccess(t *testing.T) {
  results := []Result{NewSuccess("item1"), NewSuccess("item2"), NewSuccess("item3")}
  transform := func(values ...interface{}) Result {
    return NewSuccess(values)
  }

  result := Combine(transform, results...)

  value := result.Success.([]interface{})
  assert.Equal(t, len(value), 3)
  for index, item := range value {
    assert.Equal(t, item, fmt.Sprintf("item%d", index+1))
  }
}

func TestCombineWithFailures(t *testing.T) {
  err1 := errors.New("ow")
  err2 := errors.New("oww")
  results := []Result{NewSuccess("item1"), NewFailure(err1), NewFailure(err2)}
  transform := func(values ...interface{}) Result {
    return NewSuccess(values)
  }

  result := Combine(transform, results...)

  assert.Nil(t, result.Success)
  assert.Equal(t, result.Failure, err1)
}

// Test Recover

func TestRecoverWithSuccess(t *testing.T) {
  value := NewSuccess(5).Recover(10)

  assert.Equal(t, value, 5)
}

func TestRecoverWithFailure(t *testing.T) {
  err := &errorString{"testing error"}
  value := NewFailure(err).Recover(5)

  assert.Equal(t, value, 5)
}

func TestRecoverWithWithSuccess(t *testing.T) {
  recoveredResult := NewSuccess(10)
  actualResult := NewSuccess(5)
  result := actualResult.RecoverWith(recoveredResult)

  assert.Equal(t, result, actualResult)
}

func TestRecoverWithWithFailure(t *testing.T) {
  err := &errorString{"testing error"}
  recoveredResult := NewSuccess(10)
  result := NewFailure(err).RecoverWith(recoveredResult)

  assert.Equal(t, result, recoveredResult)
}

