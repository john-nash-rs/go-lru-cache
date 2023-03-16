package main

import (
	"fmt"
	"go-lru-cache/cache"
)

func main() {
	fmt.Println("Hello World")
	c := cache.New[string](5, 2)
	c.Put("foo0", "bar0")
	fmt.Println("value is ", c.Get("foo0"))
	c.Put("foo1", "bar1")
	c.Put("foo2", "bar2")
	c.Put("foo3", "bar3")
	fmt.Println("value is ", c.Get("foo3"))
	c.Put("foo4", "bar4")
	c.Put("foo5", "bar5")
	c.Put("foo6", "bar6")
	c.Put("foo7", "bar7")
	fmt.Println("value is ", c.Get("foo7"))

}
