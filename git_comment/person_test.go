package git_comment

import (
	"github.com/stvp/assert"
	"testing"
)

func TestCreatePersonFromNameEmail(t *testing.T) {
	p, err := CreatePerson("Carrie <carrie@example.com>").Dematerialize()
	assert.Nil(t, err)
	person := p.(*Person)
	assert.NotNil(t, person)
	assert.Equal(t, person.Email, "carrie@example.com")
	assert.Equal(t, person.Name, "Carrie")
}

func TestCreatePersonFromSpacedNameEmail(t *testing.T) {
	p, err := CreatePerson("Carrie Ann McBean <carrie.mcbean@example.com>").Dematerialize()
	assert.Nil(t, err)
	person := p.(*Person)
	assert.NotNil(t, person)
	assert.Equal(t, person.Email, "carrie.mcbean@example.com")
	assert.Equal(t, person.Name, "Carrie Ann McBean")
}

func TestCreatePersonFromInvalidEmail(t *testing.T) {
	p, err := CreatePerson("Carrie Ann McBean <carrie.mcbean>").Dematerialize()
	assert.Nil(t, err)
	person := p.(*Person)
	assert.NotNil(t, person)
	assert.Equal(t, person.Name, "Carrie Ann McBean <carrie.mcbean>")
	assert.Equal(t, person.Email, "")
}

func TestCreatePersonWithoutName(t *testing.T) {
	p, err := CreatePerson("<carrie.mcbean@example.com>").Dematerialize()
	assert.Nil(t, err)
	person := p.(*Person)
	assert.NotNil(t, person)
	assert.Equal(t, person.Email, "carrie.mcbean@example.com")
	assert.Equal(t, person.Name, "")
}

func TestCreateEmptyPerson(t *testing.T) {
	_, err := CreatePerson("").Dematerialize()
	assert.NotNil(t, err)
}

func TestSerializePersonFull(t *testing.T) {
	data := "Katie Em <katie@example.com>"
	p, _ := CreatePerson(data).Dematerialize()
	person := p.(*Person)
	assert.Equal(t, person.Serialize(), data)
}

func TestSerializePersonNameOnly(t *testing.T) {
	data := "Katie Em"
	p, _ := CreatePerson(data).Dematerialize()
	person := p.(*Person)
	assert.Equal(t, person.Serialize(), data)
}

func TestSerializePersonEmailOnly(t *testing.T) {
	data := "<katie@example.com>"
	p, _ := CreatePerson(data).Dematerialize()
	person := p.(*Person)
	assert.Equal(t, person.Serialize(), data)
}
