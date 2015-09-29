package libgitcomment

import (
	"strings"
	"testing"

	"github.com/stvp/assert"
)

func TestCreatePersonFromNameEmailTime(t *testing.T) {
	p, err := CreatePerson("Carrie <carrie@example.com> 1437612685 -0700").Dematerialize()
	assert.Nil(t, err)
	person := p.(*Person)
	assert.NotNil(t, person)
	assert.Equal(t, person.Email, "carrie@example.com")
	assert.Equal(t, person.Name, "Carrie")
}

func TestCreatePersonFromSpacedNameEmailTime(t *testing.T) {
	p, err := CreatePerson("Carrie Ann McBean <carrie.mcbean@example.com>").Dematerialize()
	assert.Nil(t, err)
	person := p.(*Person)
	assert.NotNil(t, person)
	assert.Equal(t, person.Email, "carrie.mcbean@example.com")
	assert.Equal(t, person.Name, "Carrie Ann McBean")
}

func TestCreatePersonFromInvalidEmail(t *testing.T) {
	_, err := CreatePerson("Carrie Ann McBean <carrie.mcbean>").Dematerialize()
	assert.NotNil(t, err)
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
	data := "Katie Em <katie@example.com> 1437498360 +0400"
	p, err := CreatePerson(data).Dematerialize()
	assert.Nil(t, err)
	assert.NotNil(t, p)
	person := p.(*Person)
	assert.Equal(t, person.Serialize(), data)
}

func TestSerializePersonNameOnly(t *testing.T) {
	data := "Katie Em"
	_, err := CreatePerson(data).Dematerialize()
	assert.NotNil(t, err)
}

func TestSerializePersonEmailOnly(t *testing.T) {
	data := "<katie@example.com>"
	p, err := CreatePerson(data).Dematerialize()
	assert.Nil(t, err)
	assert.NotNil(t, p)
	person := p.(*Person)
	assert.True(t, strings.Contains(person.Serialize(), data))
}
