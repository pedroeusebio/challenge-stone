package main

import "fmt"

func main() {
	scores := []int{}
	c := cap(scores)
	fmt.Println(c)

	for i := 0; i < 25; i++ {
		scores = append(scores, i)

		fmt.Println(scores, cap(scores), len(scores))

		if cap(scores) != c {
			c = cap(scores)
			fmt.Println(c)
		}
	}
}
