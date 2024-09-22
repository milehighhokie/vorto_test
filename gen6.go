package main

import (
    "bufio"
    "fmt"
    "os"
	"math"
	"strings"
)

type load struct {
	loadNumber    int
    pickup        point
    dropoff       point
	distToPickup  float64
	deliveryDist  float64
	distToHome    float64
	distToOthers  map[int]float64
}

type point struct {
	x float64
	y float64
}

type truck struct {
	mileage float64
	stops []int
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func calcDistToOthers(ls map[int]load) () {
	for _,l := range ls {
		for j, d := range ls {
			l.distToOthers[j] =  math.Sqrt(math.Pow((l.dropoff.x-d.pickup.x),2) + math.Pow((l.dropoff.y-d.pickup.y),2))
		}
	}
	return
}

func main() {
	filename := os.Args[1]
    file, err := os.Open(filename)
    check(err)
    defer file.Close()

	loads := make(map[int]load)
	var loadList []int

    scanner := bufio.NewScanner(file)
	_ = scanner.Scan()
    for scanner.Scan() {
		var l load
		fmt.Sscanf(scanner.Text(), "%d (%f,%f) (%f,%f)", &l.loadNumber, &l.pickup.x, &l.pickup.y, &l.dropoff.x, &l.dropoff.y)
        l.deliveryDist =  math.Sqrt(math.Pow((l.dropoff.x-l.pickup.x),2) + math.Pow((l.dropoff.y-l.pickup.y),2))
        l.distToPickup =  math.Sqrt(math.Pow((l.pickup.x),2) + math.Pow((l.pickup.y),2))
        l.distToHome   =  math.Sqrt(math.Pow((l.dropoff.x),2) + math.Pow((l.dropoff.y),2))
		l.distToOthers = make(map[int]float64)
		loads[l.loadNumber] = l
		loadList = append(loadList, l.loadNumber)
    }
	maxTrucks := len(loadList)
	trucks :=  make([]truck, maxTrucks)  

	var currentStop int
	var possibleStop int
	calcDistToOthers(loads)

	deliverLoop:
	for t:=0; t < maxTrucks; t++ {
		if len(loadList) == 0 {
			break deliverLoop
		}
		currentStop, loadList = loadList[0], loadList[1:]
		trucks[t].stops = append(trucks[t].stops,currentStop)
		trucks[t].mileage = loads[currentStop].distToPickup + loads[currentStop].deliveryDist + loads[currentStop].distToHome
		for i:=0; i < len(loadList); i++  {
			possibleStop, loadList = loadList[0], loadList[1:]
			possibleNewMileage := trucks[t].mileage - loads[currentStop].distToHome + loads[currentStop].distToOthers[possibleStop] + loads[possibleStop].deliveryDist + loads[possibleStop].distToHome
			if possibleNewMileage <= 720 {
				trucks[t].stops = append(trucks[t].stops,possibleStop)
				trucks[t].mileage = possibleNewMileage
				currentStop = possibleStop
			} else {
				loadList = append(loadList,possibleStop)
			}
		}
	}
	for _,t := range(trucks) {
	    if len(t.stops) > 0 {
			outString := fmt.Sprint(t.stops)
			fmt.Print(strings.ReplaceAll(outString," ",","))
			//fmt.Printf(" %f ",t.mileage)	
			fmt.Println()
		}
	}
}
