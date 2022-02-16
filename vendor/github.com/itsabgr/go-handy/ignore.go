package handy

//Ignore ignores panics
//go:noinline
func Ignore() {
	recover()
}
