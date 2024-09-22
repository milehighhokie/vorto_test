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
	delivered     bool 
	distToOthers  [200]float64
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
// func calcDistToOthers(ls []load) (bool, []load) {
// 	var allDelivered bool = true
// 	for i,l := range ls {
// 		allDelivered = allDelivered && l.delivered
// 		l.distToOthers = make([]float64,len(ls))
// 		for j, d := range ls {
// 			l.distToOthers[j] =  math.Sqrt(math.Pow((l.dropoff.x-d.pickup.x),2) + math.Pow((l.dropoff.y-d.pickup.y),2))
// 		}
// 		ls[i] = l
// 	}
// 	return allDelivered, ls
// }
func checkFinished(ls []load) (bool) {
	var allDelivered bool = true
	for _,l := range ls {
		allDelivered = allDelivered && l.delivered
	}
	return allDelivered
}

func calcDistToOthers(ls []load) ([]load) {
	for i,l := range ls {
		for j, d := range ls {
			l.distToOthers[j] =  math.Sqrt(math.Pow((l.dropoff.x-d.pickup.x),2) + math.Pow((l.dropoff.y-d.pickup.y),2))
		}
		ls[i] = l
	}
	return ls
}

func main() {
	filename := os.Args[1]
    file, err := os.Open(filename)
    check(err)
    defer file.Close()

	var loads []load
	var loadList []int
	var minDist float64
	var minDistTotal float64

    scanner := bufio.NewScanner(file)
	_ = scanner.Scan()
    for scanner.Scan() {
		var l load
		fmt.Sscanf(scanner.Text(), "%d (%f,%f) (%f,%f)", &l.loadNumber, &l.pickup.x, &l.pickup.y, &l.dropoff.x, &l.dropoff.y)
        l.deliveryDist =  math.Sqrt(math.Pow((l.dropoff.x-l.pickup.x),2) + math.Pow((l.dropoff.y-l.pickup.y),2))
        l.distToPickup =  math.Sqrt(math.Pow((0-l.pickup.x),2) + math.Pow((0-l.pickup.y),2))
        l.distToHome   =  math.Sqrt(math.Pow((l.dropoff.x-0),2) + math.Pow((l.dropoff.y-0),2))
		loads = append(loads,l)
		minDistTotal += l.deliveryDist
		minDist = min(l.deliveryDist, l.distToPickup, l.distToHome, minDist)
		loadList = append(loadList, l.loadNumber)
    }
	maxTrucks := len(loadList)
	loadsToGo := len(loadList)
	trucks :=  make([]truck, maxTrucks)  
	//fmt.Printf("minimum number of trucks = %d \n", int(math.Ceil(minDistTotal/720)))
	//fmt.Printf("maximum number of trucks = %d \n", maxTrucks)

	var currentStop int
	var possibleStop int
	loads = calcDistToOthers(loads)

	deliverLoop:
	for t:=0; t < maxTrucks; t++ {
		if checkFinished(loads) {
			break deliverLoop
		}
		currentStop, loadList = loadList[0], loadList[1:]
		currentStopIndex := currentStop - 1
		trucks[t].stops = append(trucks[t].stops,loads[currentStopIndex].loadNumber)
		trucks[t].mileage = loads[currentStopIndex].distToPickup + loads[currentStopIndex].deliveryDist + loads[currentStopIndex].distToHome
		loads[currentStopIndex].delivered = true
		loadsToGo--
		for i:=0; i < loadsToGo; i++  {
			possibleStop, loadList = loadList[0], loadList[1:]
			possibleStopIndex := possibleStop - 1
			if trucks[t].mileage - loads[currentStopIndex].distToHome + loads[currentStopIndex].distToOthers[possibleStopIndex] + loads[possibleStopIndex].deliveryDist + loads[possibleStopIndex].distToHome <= 720 {
				loads[possibleStopIndex].delivered = true
				trucks[t].stops = append(trucks[t].stops,possibleStop)
				trucks[t].mileage = trucks[t].mileage - loads[currentStopIndex].distToHome + loads[currentStopIndex].distToOthers[possibleStopIndex] + loads[possibleStopIndex].deliveryDist + loads[possibleStopIndex].distToHome
				currentStop = possibleStop
				currentStopIndex = possibleStopIndex
				loadsToGo--
			} else {
				loadList = append(loadList,possibleStop)
			}
		}
		loadsToGo = len(loadList)
	}
	for _,t := range(trucks) {
	    if len(t.stops) > 0 {
			outString := fmt.Sprint(t.stops)
			fmt.Println(strings.ReplaceAll(outString," ",","))
	    }
	}
}
