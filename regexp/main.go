package main

import (
	"fmt"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`^([a-zA-Z0-9])(([\-.]|[_]+)?([a-zA-Z0-9]+))*(@){1}[a-z0-9]+[.]{1}(([a-z]{2,3})|([a-z]{2,3}[.]{1}[a-z]{2,3}))$`)
	fmt.Println(re.MatchString("hello@example.com"))
	fmt.Println(re.MatchString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa!"))
}
