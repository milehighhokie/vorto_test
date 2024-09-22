package main

import (
    "bufio"
    "fmt"
    "os"
	"math"
)

type load struct {
	loadNumber    int
    pickup        point
    dropoff       point
	distToPickup  float64
	deliveryDist  float64
	distToHome    float64
}

type point struct {
	x float64
	y float64
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

    file, err := os.Open("./data/problem1.txt")
    check(err)
    defer file.Close()

	var loads []load
	var minDist float64
    scanner := bufio.NewScanner(file)
	_ = scanner.Scan()
    for scanner.Scan() {
		var l load
		fmt.Sscanf(scanner.Text(), "%d (%f,%f) (%f,%f)", &l.loadNumber, &l.pickup.x, &l.pickup.y, &l.dropoff.x, &l.dropoff.y)
        l.deliveryDist =  math.Sqrt(math.Pow((l.dropoff.x-l.pickup.x),2) + math.Pow((l.dropoff.y-l.pickup.y),2))
        l.distToPickup =  math.Sqrt(math.Pow((0-l.pickup.x),2) + math.Pow((0-l.pickup.y),2))
        l.distToHome   =  math.Sqrt(math.Pow((l.dropoff.x-0),2) + math.Pow((l.dropoff.y-0),2))
		fmt.Println(l)
		loads = append(loads,l)
		minDist += l.deliveryDist
    }
	fmt.Printf("minimum number of drivers = %f \n", math.Ceil(minDist/720))
}
