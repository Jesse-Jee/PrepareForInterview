package Practice

import "fmt"

func f(n int, in <-chan int, out chan<- int) {
	for {
		<-in
		fmt.Println(n)
		out<-1
	}
}


func main() {
	c := [4]chan int{}
	for i := 0;i<4;i++{
		c[i] = make(chan int)
	}

	go f(1,c[3],c[0])
	go f(2,c[0],c[1])
	go f(3,c[1],c[2])
	go f(4,c[2],c[3])
	c[3]<-1
	select{}
}
