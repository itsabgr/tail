package handy

//Catch call fn with recovered value is its not nil
func Catch(fn func(recovered interface{})) {
	recovered := recover()
	if recovered != nil {
		fn(recovered)
	}
}
