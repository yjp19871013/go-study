package main

import (
	"fmt"
	"strings"
)

type Person struct {
	Name string
}

func (p *Person) PrintMe() {
	fmt.Println(p.Name)
}

func (p Person) NameToUpper1() {
	p.Name = strings.ToUpper(p.Name)
}

func (p Person) NameToUpper2() string {
	return strings.ToUpper(p.Name)
}

func (p *Person) NameToUpper3() {
	p.Name = strings.ToUpper(p.Name)
}

type Values []int

func (values Values) PrintValues() {
	fmt.Println(values)
}

func (values Values) AddOne() {
	for i, v := range values {
		values[i] = v + 1
	}
}

type Dict map[string]int

func (d Dict) PrintDict() {
	fmt.Println(d)
}

func (d Dict) AddOneByKey(key string) {
	d[key] = d[key] + 1
}

func main() {
	person := Person{Name: "yjp"}
	person.PrintMe()

	person.NameToUpper1()
	person.PrintMe()

	person.NameToUpper2()
	person.PrintMe()

	person.Name = person.NameToUpper2()
	person.PrintMe()

	person.Name = "yjp"
	person.NameToUpper3()
	person.PrintMe()

	values := Values{1, 2, 3}
	values.AddOne()
	values.PrintValues()

	dict := Dict{"yjp": 1}
	dict.AddOneByKey("yjp")
	dict.PrintDict()
}
