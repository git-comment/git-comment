package git_comment

import (
  "github.com/stvp/assert"
  "testing"
)

func TestCreatePersonFromNameEmail(t *testing.T) {
  person, err := CreatePerson("Carrie <carrie@example.com>")
  assert.NotNil(t, person)
  assert.Equal(t, person.Email, "carrie@example.com")
  assert.Equal(t, person.Name, "Carrie")
  assert.Nil(t, err)
}

func TestCreatePersonFromSpacedNameEmail(t *testing.T) {
  person, err := CreatePerson("Carrie Ann McBean <carrie.mcbean@example.com>")
  assert.NotNil(t, person)
  assert.Equal(t, person.Email, "carrie.mcbean@example.com")
  assert.Equal(t, person.Name, "Carrie Ann McBean")
  assert.Nil(t, err)
}

func TestCreatePersonFromInvalidEmail(t *testing.T) {
  person, err := CreatePerson("Carrie Ann McBean <carrie.mcbean>")
  assert.NotNil(t, err)
  assert.Nil(t, person)
}

func TestCreatePersonWithoutName(t *testing.T) {
  person, err := CreatePerson("<carrie.mcbean@example.com>")
  assert.NotNil(t, person)
  assert.Equal(t, person.Email, "carrie.mcbean@example.com")
  assert.Equal(t, person.Name, "")
  assert.Nil(t, err)
}
