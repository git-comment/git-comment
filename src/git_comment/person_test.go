package git_comment

import (
	"github.com/stvp/assert"
	"testing"
)

func TestCreatePersonFromNameEmail(t *testing.T) {
	person := CreatePerson("Carrie <carrie@example.com>")
	assert.NotNil(t, person)
	assert.Equal(t, person.Email, "carrie@example.com")
	assert.Equal(t, person.Name, "Carrie")
}

func TestCreatePersonFromSpacedNameEmail(t *testing.T) {
	person := CreatePerson("Carrie Ann McBean <carrie.mcbean@example.com>")
	assert.NotNil(t, person)
	assert.Equal(t, person.Email, "carrie.mcbean@example.com")
	assert.Equal(t, person.Name, "Carrie Ann McBean")
}

func TestCreatePersonFromInvalidEmail(t *testing.T) {
	person := CreatePerson("Carrie Ann McBean <carrie.mcbean>")
	assert.NotNil(t, person)
	assert.Equal(t, person.Name, "Carrie Ann McBean <carrie.mcbean>")
	assert.Equal(t, person.Email, "")
}

func TestCreatePersonWithoutName(t *testing.T) {
	person := CreatePerson("<carrie.mcbean@example.com>")
	assert.NotNil(t, person)
	assert.Equal(t, person.Email, "carrie.mcbean@example.com")
	assert.Equal(t, person.Name, "")
}

func TestSerializePersonFull(t *testing.T) {
	data := "Katie Em <katie@example.com>"
	person := CreatePerson(data)
	assert.Equal(t, person.Serialize(), data)
}

func TestSerializePersonNameOnly(t *testing.T) {
	data := "Katie Em"
	person := CreatePerson(data)
	assert.Equal(t, person.Serialize(), data)
}

func TestSerializePersonEmailOnly(t *testing.T) {
	data := "<katie@example.com>"
	person := CreatePerson(data)
	assert.Equal(t, person.Serialize(), data)
}
