package fail

func When(check bool, message interface{}) {
	if check {
		panic(message)
	}
}
