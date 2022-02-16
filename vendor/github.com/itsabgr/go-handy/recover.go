package handy

//Recover call fn and recover a panic and returns recovered value
//go:noinline
func Recover(fn func()) (recovered interface{}) {
	defer func() {
		recovered = recover()
	}()
	fn()
	return
}
