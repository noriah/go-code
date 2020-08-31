package embed

import "fmt"

// overFooBar common interface
type overFooBar interface {
	funk() string
	hello() string
}

// overFoo embeds overFooBar
// So it has the method signatures from overFooBar
// If we call any of them now, it will be a panic with nil pointer
type overFoo struct {
	overFooBar
}

// overBar embeds overFooBar as well
// So we must declare the methods or risk a panic
type overBar struct {
	overFooBar
}

// overBaz embeds overBar
type overBaz struct {
	overBar
}

// overQuux embeds overBaz
type overQuux struct {
	overBaz
}

// overFooBarBaz also embeds overBaz
type overFooBarBaz struct {
	overBaz
}

func (*overFoo) funk() string {
	return "regular #funk on *overFoo"
}

func (*overFoo) hello() string {
	return "hello world from *overFoo"
}

func (*overBar) funk() string {
	return "regular #funk on *overBar"
}

func (*overBaz) hello() string {
	return "Hello!!! World!!! From *overBaz"
}

func (b *overBaz) getMe() overFooBar {
	return b
}

func (*overQuux) funk() string {
	return "overloaded *overBaz.funk with *overQuux.funk"
}

func (*overFooBarBaz) funk() string {
	return "overloaded *overBaz.funk with *overFooBarBaz.funk"
}

// TestOverloading tests the above structures
func TestOverloading() {
	var faces = []string{
		"ヾ(⌐■_■)ノ♪", "(╯°□°）╯︵ ┻━┻",
		"┬─┬ ノ( ゜-゜ノ)", "(つ ◕_◕ )つ",
		"¯\\_(ツ)_/¯", "(╯°□°）╯︵( .O.)",
		"(☞ﾟヮﾟ)☞",
	}

	var tester = func(idx int, fb overFooBar) {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Printf(" ### PANIC!!!\n\n(%v)\n%s || I gotchu fam\n", rec, faces[6])
			}
		}()
		fmt.Printf("\n%T (%#v)\n%s || %s", fb, fb, faces[idx], fb.funk())
		fmt.Printf(" || %s\n", fb.hello())
	}

	for i, fb := range []overFooBar{
		&overFoo{}, &overBar{}, &overBaz{},
		&overQuux{}, (&overFooBarBaz{}).getMe(), &overFooBarBaz{},
	} {
		tester(i, fb)
	}
}
