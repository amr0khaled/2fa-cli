package res

import (
	"encoding/base32"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type MfaAccount struct {
	Name   string
	Email  string
	key    *Key
	digits uint8
}

func NewMFaAccount(name string, email string, key *Key) *MfaAccount {
	return &MfaAccount{
		Name:  name,
		Email: email,
		key:   key,
	}
}

type MfaAccountInterface interface {
	GetKey() Key
	GetDigits() uint8
}

func (m *MfaAccount) GetKey() *Key {
	return m.key
}

func (m *MfaAccount) GetDigits() uint8 {
	return m.digits
}

type Key struct {
	armour string
	raw    []byte
}

func NewKey(secret []byte, armour ...string) *Key {
	var base string
	if !(len(armour) > 0) && !(len(secret) > 0) {
		check(errors.New("Must pass an argument to generate a new key"))
	}
	if !(len(secret) > 0) {
		base = armour[0]
		_secret, err := base32.StdEncoding.DecodeString(base)
		check(err)
		secret = _secret
	} else {
		base = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)
	}

	return &Key{
		armour: base,
		raw:    secret,
	}
}

type Iter[T any] struct {
	elements []T
	index    int
	length   int
	pre      func(el T) T
}

type IterInterface[T any] interface {
	Peek() T
	PeekN(i int) T
	Consume() T
	Current() T
	IsFin() bool
	IsOutRange() bool
	Items() []T
}

func NewIter[T any](elements []T) *Iter[T] {
	return &Iter[T]{
		elements: elements,
		index:    0,
		length:   len(elements),
		pre:      func(el T) T { return el },
	}
}

func (t *Iter[T]) AddPre(fn func(el T) T) {
	t.pre = fn
}

func (t *Iter[T]) Peek() T {
	var zero T
	if t.IsOutRange() {
		return zero
	}
	return t.pre(t.elements[t.index+1])
}

func (t *Iter[T]) PeekN(i int) T {
	return t.pre(t.elements[t.index+i])
}

func (t *Iter[T]) Consume() T {
	value := t.pre(t.elements[t.index])
	if !t.IsOutRange() {
		t.index++
	}
	return value
}

func (t *Iter[T]) Current() T {
	return t.pre(t.elements[t.index])
}

func (t *Iter[T]) IsFin() bool {
	return !t.IsOutRange()
}

func (t *Iter[T]) IsOutRange() bool {
	return t.index+1 >= t.length
}

func (t *Iter[T]) Items() []T {
	return t.elements
}

type handler func(account *MfaAccount, iter *Iter[string])

var handlers map[string]handler = map[string]handler{
	"email":  handleEmail,
	"digits": handleDigits,
	"armour": handleArmour,
}

func handleEmail(account *MfaAccount, iter *Iter[string]) {
	account.Email = iter.Consume()
}

func handleDigits(account *MfaAccount, iter *Iter[string]) {
	num, err := strconv.Atoi(iter.Consume())
	check(err)
	account.digits = uint8(num)
}

func handleArmour(account *MfaAccount, iter *Iter[string]) {
	key := NewKey([]byte{}, iter.Consume())
	account.key = key
}

func NewAccountFromFile(name string, content string) *MfaAccount {
	fmt.Printf("%s\n", content)
	lines := NewIter(strings.Split(content, "\n"))
	triming := func(el string) string {
		return strings.TrimSpace(el)
	}
	lines.AddPre(triming)
	account := &MfaAccount{
		Name: name,
	}

	for lines.IsFin() {
		line := NewIter(strings.Split(lines.Current(), ":"))
		line.AddPre(triming)
		for line.IsFin() {
			handler, ok := handlers[line.Current()]
			if ok {
				line.Consume()
				handler(account, line)
			}
			line.Consume()
		}
		lines.Consume()
	}
	return account
}

type KeyInterface interface {
	GetArmour() string
	GetRaw() []byte
}

func (k *Key) GetArmour() string {
	return k.armour
}

func (k *Key) GetRaw() []byte {
	return k.raw
}

var accounts []*MfaAccount

func init() {
	if res == nil {
		res = NewRes()
	}
	keys := readKeys()
	for name, key := range keys {
		account := NewAccountFromFile(name, key.Content)
		accounts = append(accounts, account)
	}
	fmt.Printf("%+vss", accounts)

}
