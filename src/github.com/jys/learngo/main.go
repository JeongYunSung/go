package main

import (
	"fmt"
	"github.com/jys/learngo/dict"
)

func main() {
	dictionary := dict.Dictionary{"hello": "안녕하세요", "bye": "들어가세요"}
	dictionary["bye"] = "안녕히 가세요"
	dictionary.Add("hi", "안녕")
	fmt.Println(dictionary)
	fmt.Println(dictionary.Search("hi"))
	dictionary1 := dict.Dictionary{"hello": "안녕하세요", "bye": "들어가세요"}
	fmt.Println(dictionary1)
}
