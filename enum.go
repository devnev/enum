package enum

import (
	"fmt"
	"strings"
)

// Member is an enum member, a specific value bound to a variable.
type Member[T comparable] struct {
	Value T
}

// iMember is the type constraint for Member used by Enum.
//
// We can't use Member directly in type constraints
// because the users create a new subtype from Member
// instead of using it directly.
//
// We also can't use a normal interface because new types
// don't inherit methods of their base type.
type iMember[T comparable] interface {
	~struct{ Value T }
}

// Enum is a collection of enum members.
//
// Use [New] to construct a new Enum from a list of members.
type Enum[M iMember[V], V comparable] struct {
	members []M
	v2m     map[V]M
}

// New constructs a new [Enum] wrapping the given enum members.
func New[V comparable, M iMember[V]](members ...M) Enum[M, V] {
	e := Enum[M, V]{members, nil}
	e.v2m = make(map[V]M)
	for i, m := range e.members {
		v := e.Value(m)
		e.v2m[v] = e.members[i]
	}
	return e
}

// TypeName is a string representation of the wrapped type.
func (Enum[M, V]) TypeName() string {
	return fmt.Sprintf("%T", *new(V))
}

// Empty returns true if the enum doesn't have any members.
func (e Enum[M, V]) Empty() bool {
	return len(e.members) == 0
}

// Len returns how many members the enum has.
func (e Enum[M, V]) Len() int {
	return len(e.members)
}

// Contains returns true if the enum has the given member.
func (e Enum[M, V]) Contains(member M) bool {
	for _, m := range e.members {
		if m == member {
			return true
		}
	}
	return false
}

// Parse converts a raw value into a member of the enum.
//
// If none of the enum members has the given value, nil is returned.
func (e Enum[M, V]) Parse(value V) (M, bool) {
	m, ok := e.v2m[value]
	return m, ok
}

// Value returns the wrapped value of the given enum member.
func (e Enum[M, V]) Value(member M) V {
	return Member[V](member).Value
}

// Index returns the index of the given member in the enum.
//
// If the given member is not in the enum, it panics.
// Use [Enum.Contains] first if you don't know for sure
// if the member belongs to the enum.
func (e Enum[M, V]) Index(member M) int {
	for i, m := range e.members {
		if e.Value(m) == e.Value(member) {
			return i
		}
	}
	panic("the given Member does not belong to this Enum")
}

// Members returns a slice of the members in the enum.
func (e Enum[M, V]) Members() []M {
	return e.members
}

// Values returns a slice of values of all members of the enum.
func (e Enum[M, V]) Values() []V {
	res := make([]V, 0, len(e.members))
	for _, m := range e.members {
		res = append(res, e.Value(m))
	}
	return res
}

// String implements [fmt.Stringer] interface.
//
// It returns a comma-separated list of values of the enum members.
func (e Enum[M, V]) String() string {
	values := make([]string, 0, len(e.members))
	for _, m := range e.members {
		values = append(values, fmt.Sprintf("%v", e.Value(m)))
	}
	return strings.Join(values, ", ")
}

// GoString implements [fmt.GoStringer] interface.
//
// When you print a member using "%#v" format,
// it will show the enum representation as a valid Go syntax.
func (e Enum[M, V]) GoString() string {
	values := make([]string, 0, len(e.members))
	for _, m := range e.members {
		values = append(values, fmt.Sprintf("%T{%#v}", m, e.Value(m)))
	}
	joined := strings.Join(values, ", ")
	return fmt.Sprintf("enum.New(%s)", joined)
}

// Builder is a constructor for an [Enum].
//
// Use [Builder.Add] to add new members to the future enum
// and then call [Builder.Enum] to create a new [Enum] with all added members.
//
// Builder is useful for when you have lots of enum members, and new ones
// are added over time, as the project grows. In such scenario, it's easy to forget
// to add in the [Enum] a newly created [Member].
// The builder is designed to prevent that.
type Builder[M iMember[V], V comparable] struct {
	members  []M
	finished bool
}

// NewBuilder creates a new [Builder], a constructor for an [Enum].
func NewBuilder[V comparable, M iMember[V]]() Builder[M, V] {
	return Builder[M, V]{make([]M, 0), false}
}

// Add registers a new [Member] in the builder.
func (b *Builder[M, V]) Add(m M) M {
	b.members = append(b.members, m)
	return m
}

// Enum creates a new [Enum] with all members registered using [Builder.Add].
func (b *Builder[M, V]) Enum() Enum[M, V] {
	b.finished = true
	e := New(b.members...)
	return e
}
