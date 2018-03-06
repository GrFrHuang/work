package utils

import (
	"log"
	"testing"
	"strings"
	"fmt"
)

func TestUnDuplicatesSlice(t *testing.T) {
	x := []interface{}{1, 1, 2, 3, 5, 5, 2}
	UnDuplicatesSlice(&x)
	log.Print(x)
}

func TestGetNotEmptyFields(t *testing.T) {
	type X struct {
		Name string
		Id   int
	}
	x := &X{
		Name: "123",
	}
	log.Print(GetNotEmptyFields(x, "Name"))
}


func TestFace(t *testing.T) {
    st:="1,2,3"
    fmt.Println(strings.Split(st,","))

}
