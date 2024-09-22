package main

import (
    "bufio"
    "fmt"
    "os"
	"math"
	"strings"
	"slices"
)

type load struct {
	loadNumber    int        
    pickup        point      
    dropoff       point      
	distToPickup  float64    
	deliveryDist  float64    
	distToHome    float64    
	distTotal     float64    
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

func dumpTrucks(trucks []truck) () {
	for _,t := range trucks {
	    if len(t.stops) > 0 {
			outString := fmt.Sprint(t.stops)
			fmt.Print(strings.ReplaceAll(outString," ",","))
			//fmt.Printf(" %f ",t.mileage)	
			fmt.Println()
		}
	}
}

func calcScore(trucks []truck) (float64) {
	var score float64
	for _, t := range trucks {
		score = score + t.mileage + 500
	}
	return (score)
}

func firstLoadProcessing(loadList []int, loads map[int]load) ([]truck, float64) {
	maxTrucks := len(loadList)
	trucks :=  make([]truck, maxTrucks) 
	
	var currentStop int
	var possibleStop int
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
	return trucks, calcScore(trucks)
}

func reverseLoadProcessing(loadList []int, loads map[int]load) ([]truck, float64) {
	maxTrucks := len(loadList)
	trucks :=  make([]truck, maxTrucks) 
	slices.Reverse(loadList)

	var currentStop int
	var possibleStop int
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
	return trucks, calcScore(trucks)
}

func closestLoadProcessing(loadList []int, loads map[int]load) ([]truck, float64) {
	maxTrucks := len(loadList)
	trucks :=  make([]truck, maxTrucks) 

	var currentStop int
	var possibleStop int
	deliverLoop:
	for t:=0; t < maxTrucks; t++ {
		if len(loadList) == 0 {
			break deliverLoop
		}
		currentStop, loadList = loadList[0], loadList[1:]
		trucks[t].stops = append(trucks[t].stops,currentStop)
		trucks[t].mileage = loads[currentStop].distToPickup + loads[currentStop].deliveryDist + loads[currentStop].distToHome
		for i:=0; i < len(loadList); i++  {
			// find the closest load for the currentStop
			var minDistPossible=999.999
			var minDistPossibleIdx=-1
			for m:=0; m < len(loadList); m++ {
				if loads[currentStop].distToOthers[loadList[m]] < minDistPossible {
					minDistPossibleIdx = m
					minDistPossible = loads[currentStop].distToOthers[loadList[m]]
				}
			}
			possibleStop = loadList[minDistPossibleIdx]
			loadList = slices.Delete(loadList, minDistPossibleIdx, minDistPossibleIdx+1)
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
	return trucks, calcScore(trucks)
}

func farThenClosestLoadProcessing(loadList []int, loads map[int]load) ([]truck, float64) {
	maxTrucks := len(loadList)
	trucks :=  make([]truck, maxTrucks) 

	var currentStop int
	var possibleStop int
	deliverLoop:
	for t:=0; t < maxTrucks; t++ {
		if len(loadList) == 0 {
			break deliverLoop
		}
		// find the farthest trip to get load for the first stop
		var maxDistPossible=-1.0000000
		var maxDistPossibleIdx=-1
		for m:=0; m < len(loadList); m++ {
			if loads[loadList[m]].distToPickup > maxDistPossible {
				maxDistPossibleIdx = m
				maxDistPossible = loads[loadList[m]].distToPickup
			}
		}
		currentStop = loadList[maxDistPossibleIdx]
		loadList = slices.Delete(loadList, maxDistPossibleIdx, maxDistPossibleIdx+1)
		trucks[t].stops = append(trucks[t].stops,currentStop)
		trucks[t].mileage = loads[currentStop].distToPickup + loads[currentStop].deliveryDist + loads[currentStop].distToHome
		for i:=0; i < len(loadList); i++  {
			// find the closest load for the currentStop
			var minDistPossible=999999.0000000
			var minDistPossibleIdx=-1
			for m:=0; m < len(loadList); m++ {
				if loads[currentStop].distToOthers[loadList[m]] < minDistPossible {
					minDistPossibleIdx = m
					minDistPossible = loads[currentStop].distToOthers[loadList[m]]
				}
			}
			possibleStop = loadList[minDistPossibleIdx]
			loadList = slices.Delete(loadList, minDistPossibleIdx, minDistPossibleIdx+1)
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
	return trucks, calcScore(trucks)
}

func nearThenClosestLoadProcessing(loadList []int, loads map[int]load) ([]truck, float64) {
	maxTrucks := len(loadList)
	trucks :=  make([]truck, maxTrucks) 

	var currentStop int
	var possibleStop int
	deliverLoop:
	for t:=0; t < maxTrucks; t++ {
		if len(loadList) == 0 {
			break deliverLoop
		}
		// find the nearest trip to get load for the first stop
		var minDistPossible=999999.0000000
		var minDistPossibleIdx=-1
		for m:=0; m < len(loadList); m++ {
			if loads[loadList[m]].distToPickup < minDistPossible {
				minDistPossibleIdx = m
				minDistPossible = loads[loadList[m]].distToPickup
			}
		}
		currentStop = loadList[minDistPossibleIdx]
		loadList = slices.Delete(loadList, minDistPossibleIdx, minDistPossibleIdx+1)
		trucks[t].stops = append(trucks[t].stops,currentStop)
		trucks[t].mileage = loads[currentStop].distToPickup + loads[currentStop].deliveryDist + loads[currentStop].distToHome
		for i:=0; i < len(loadList); i++  {
			// find the closest load for the currentStop
			var minDistPossible=999999.0000000
			var minDistPossibleIdx=-1
			for m:=0; m < len(loadList); m++ {
				if loads[currentStop].distToOthers[loadList[m]] < minDistPossible {
					minDistPossibleIdx = m
					minDistPossible = loads[currentStop].distToOthers[loadList[m]]
				}
			}
			possibleStop = loadList[minDistPossibleIdx]
			loadList = slices.Delete(loadList, minDistPossibleIdx, minDistPossibleIdx+1)
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
	return trucks, calcScore(trucks)
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
        l.distTotal    =  l.deliveryDist + l.distToPickup + l.distToHome 
		l.distToOthers = make(map[int]float64)
		loads[l.loadNumber] = l
		loadList = append(loadList, l.loadNumber)
    }
 	calcDistToOthers(loads)

	// []truck{}, 99999999999999.0  //
	firstLoadTrucks, firstLoadScore := []truck{}, 99999999999999.0  //firstLoadProcessing(loadList,loads)
	reverseLoadTrucks, reverseLoadScore := []truck{}, 99999999999999.0  //reverseLoadProcessing(loadList,loads)
	closestLoadTrucks, closestLoadScore := []truck{}, 99999999999999.0  //closestLoadProcessing(loadList,loads)
	farThenClosestLoadTrucks, farThenClosestLoadScore := []truck{}, 99999999999999.0  //farThenClosestLoadProcessing(loadList,loads)
	nearThenClosestLoadTrucks, nearThenClosestLoadScore := nearThenClosestLoadProcessing(loadList,loads)

	lowestScore := slices.Min([]float64{firstLoadScore,reverseLoadScore,closestLoadScore,farThenClosestLoadScore,nearThenClosestLoadScore})
	switch {
		case lowestScore == firstLoadScore:
				dumpTrucks(firstLoadTrucks)
		case lowestScore == reverseLoadScore:
				dumpTrucks(reverseLoadTrucks)	
		case lowestScore == closestLoadScore:
				dumpTrucks(closestLoadTrucks)	
		case lowestScore == farThenClosestLoadScore:
				dumpTrucks(farThenClosestLoadTrucks)
		case lowestScore == nearThenClosestLoadScore:
				dumpTrucks(nearThenClosestLoadTrucks)
	} 
}
