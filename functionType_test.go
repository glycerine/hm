package hm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionTypeBasics(t *testing.T) {
	fnType := NewFnType(TypeVariable('a'), TypeVariable('a'), TypeVariable('a'))
	if fnType.Name() != "→" {
		t.Errorf("FunctionType should have \"→\" as a name. Got %q instead", fnType.Name())
	}

	if fnType.String() != "a → a → a" {
		t.Errorf("Expected \"a → a → a\". Got %q instead", fnType.String())
	}

	ftv := fnType.FreeTypeVar()
	if len(ftv) != 1 {
		t.Errorf("Expected only one free type var")
	}

	for _, fas := range fnApplyTests {
		fn := fas.fn.Apply(fas.sub).(*FunctionType)
		if !fn.Eq(fas.expected) {
			t.Errorf("Expected %v. Got %v instead", fas.expected, fn)
		}
	}

	// bad shit
	f := func() {
		NewFnType(TypeVariable('a'))
	}
	assert.Panics(t, f)
}

var fnApplyTests = []struct {
	fn  *FunctionType
	sub Subs

	expected *FunctionType
}{
	{NewFnType(TypeVariable('a'), TypeVariable('a')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, proton)},
	{NewFnType(TypeVariable('a'), TypeVariable('b')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, neutron)},
	{NewFnType(TypeVariable('a'), TypeVariable('b')), mSubs{'c': proton, 'd': neutron}, NewFnType(TypeVariable('a'), TypeVariable('b'))},
	{NewFnType(TypeVariable('a'), TypeVariable('b')), mSubs{'a': proton, 'c': neutron}, NewFnType(proton, TypeVariable('b'))},
	{NewFnType(TypeVariable('a'), TypeVariable('b')), mSubs{'c': proton, 'b': neutron}, NewFnType(TypeVariable('a'), neutron)},
	{NewFnType(electron, proton), mSubs{'a': proton, 'b': neutron}, NewFnType(electron, proton)},

	// a -> (b -> c)
	{NewFnType(TypeVariable('a'), TypeVariable('b'), TypeVariable('a')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, neutron, proton)},
	{NewFnType(TypeVariable('a'), TypeVariable('a'), TypeVariable('b')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, proton, neutron)},
	{NewFnType(TypeVariable('a'), TypeVariable('b'), TypeVariable('c')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, neutron, TypeVariable('c'))},
	{NewFnType(TypeVariable('a'), TypeVariable('c'), TypeVariable('b')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, TypeVariable('c'), neutron)},

	// (a -> b) -> c
	{NewFnType(NewFnType(TypeVariable('a'), TypeVariable('b')), TypeVariable('a')), mSubs{'a': proton, 'b': neutron}, NewFnType(NewFnType(proton, neutron), proton)},
}
