package utils

// Here we create a simple function that will take care of errors, helping with some code clean up
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
