package main

import (
    "bufio"
    "fmt"
    "os"
	"math"
	"strings"
	"slices"
	"sync"
	"math/rand"
)
var wg sync.WaitGroup

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

func displayResults(trucks []truck) () {
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

func farthestPickCurrent(loadList []int, loads map[int]load) ([]int, int) {
	// find the farthest trip to get load for the first stop
	var maxDistPossible=-1.0000000
	var maxDistPossibleIdx=-1
	for m:=0; m < len(loadList); m++ {
		if loads[loadList[m]].distToPickup > maxDistPossible {
			maxDistPossibleIdx = m
			maxDistPossible = loads[loadList[m]].distToPickup
		}
	}
	currentStop := loadList[maxDistPossibleIdx]
	loadList = slices.Delete(loadList, maxDistPossibleIdx, maxDistPossibleIdx+1)
	return loadList, currentStop
}

func nearestPickCurrent(loadList []int, loads map[int]load) ([]int, int) {
	// find the nearest trip to get load for the first stop
	var minDistPossible=999999.0000000
	var minDistPossibleIdx=-1
	for m:=0; m < len(loadList); m++ {
		if loads[loadList[m]].distToPickup < minDistPossible {
			minDistPossibleIdx = m
			minDistPossible = loads[loadList[m]].distToPickup
		}
	}
	currentStop := loadList[minDistPossibleIdx]
	loadList = slices.Delete(loadList, minDistPossibleIdx, minDistPossibleIdx+1)
	return loadList, currentStop

}

func randomPickCurrent(loadList []int, loads map[int]load) ([]int, int) {
	// pick a random first stop
	randomIdx := rand.Intn(len(loadList))
	currentStop := loadList[randomIdx]
	loadList = slices.Delete(loadList, randomIdx, randomIdx+1)
	return loadList, currentStop
}

func groupBunchPickPossible(loadList []int, loads map[int]load, currentStop int) ([]int, int) {
	// find the biggestest grouping of loads for the next stop
	groupDef := 23.456789
	var maxGroupingCount=-1
	var maxGroupingIdx=-1
	for m:=0; m < len(loadList); m++ {
		var thisGroupingCount int 
		for g:=0; g < len(loadList); g++ {
			if loads[currentStop].distToOthers[loadList[g]] < groupDef {
				thisGroupingCount++
			}
		}
		if thisGroupingCount > maxGroupingCount {
			maxGroupingIdx = m
			maxGroupingCount = thisGroupingCount
		}
	}
	possibleStop := loadList[maxGroupingIdx]
	loadList = slices.Delete(loadList, maxGroupingIdx, maxGroupingIdx+1)
	return loadList, possibleStop
}

func groupBunchPickCurrent(loadList []int, loads map[int]load) ([]int, int) {
	// find the biggestest grouping of loads for the first stop
	groupDef := 23.456789
	var maxGroupingCount=-1
	var maxGroupingIdx=-1
	for m:=0; m < len(loadList); m++ {
		var thisGroupingCount int 
		for g:=0; g < len(loadList); g++ {
			if loads[loadList[m]].distToOthers[loadList[g]] < groupDef {
				thisGroupingCount++
			}
		}
		if thisGroupingCount > maxGroupingCount {
			maxGroupingIdx = m
			maxGroupingCount = thisGroupingCount
		}
	}
	currentStop := loadList[maxGroupingIdx]
	loadList = slices.Delete(loadList, maxGroupingIdx, maxGroupingIdx+1)
	return loadList, currentStop
}

func findClosestPossibleStop(loadList []int, loads map[int]load, currentStop int) ([]int, int) {
	// find the closest load for the currentStop
	var minDistPossible=999999.0000000
	var minDistPossibleIdx=-1
	for m:=0; m < len(loadList); m++ {
		if loads[currentStop].distToOthers[loadList[m]] < minDistPossible {
			minDistPossibleIdx = m
			minDistPossible = loads[currentStop].distToOthers[loadList[m]]
		}
	}
	possibleStop := loadList[minDistPossibleIdx]
	loadList = slices.Delete(loadList, minDistPossibleIdx, minDistPossibleIdx+1)
	return loadList, possibleStop
}

func findRandomPossibleStop(loadList []int, loads map[int]load, currentStop int) ([]int, int) {
	randomIdx := rand.Intn(len(loadList))
	possibleStop := loadList[randomIdx]
	loadList = slices.Delete(loadList, randomIdx, randomIdx+1)
	return loadList, possibleStop
}

func genericLoadProcessing(loadList []int, loads map[int]load, trucksPointer *[]truck, pickCurrentStop func([]int,map[int]load) ([]int,int), pickPossibleStop func([]int,map[int]load,int) ([]int,int) ) () {
	maxTrucks := len(loadList)
	trucks :=  make([]truck, maxTrucks) 

	var currentStop int
	var possibleStop int
	deliverLoop:
	for t:=0; t < maxTrucks; t++ {
		if len(loadList) == 0 {
			break deliverLoop
		}
		loadList, currentStop = pickCurrentStop(loadList, loads)

		trucks[t].stops = append(trucks[t].stops,currentStop)
		trucks[t].mileage = loads[currentStop].distToPickup + loads[currentStop].deliveryDist + loads[currentStop].distToHome
		for i:=0; i < len(loadList); i++  {
			loadList, possibleStop = pickPossibleStop(loadList, loads, currentStop)

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
	*trucksPointer = trucks
	wg.Done()
	return
}

func main() {
	filename := os.Args[1]
    file, err := os.Open(filename)
    check(err)
    defer file.Close()

	loads1 := make(map[int]load)
	loads2 := make(map[int]load)
	loads3 := make(map[int]load)
	loads4 := make(map[int]load)
	var loadList1 []int
	var loadList2 []int
	var loadList3 []int
	var loadList4 []int

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
		loads1[l.loadNumber] = l
		loadList1 = append(loadList1, l.loadNumber)
		loads2[l.loadNumber] = l
		loadList2 = append(loadList2, l.loadNumber)
		loads3[l.loadNumber] = l
		loadList3 = append(loadList3, l.loadNumber)
		loads4[l.loadNumber] = l
		loadList4 = append(loadList4, l.loadNumber)
	}
 	calcDistToOthers(loads1)
 	calcDistToOthers(loads2)
 	calcDistToOthers(loads3)
 	calcDistToOthers(loads4)

	farThenClosestLoadTrucks  := make([]truck, len(loadList1)) 
	nearThenClosestLoadTrucks := make([]truck, len(loadList2))
	nearestThenBunchTrucks   := make([]truck, len(loadList3))
	bunchedFirstLoadTrucks    := make([]truck, len(loadList4)) 

	wg.Add(1)
	go  genericLoadProcessing(loadList1, loads1, &farThenClosestLoadTrucks, farthestPickCurrent, findClosestPossibleStop)
	farThenClosestLoadScore := calcScore(farThenClosestLoadTrucks)

	//wg.Add(1)
	//go  genericLoadProcessing(loadList2, loads2, &nearestThenBunchTrucks, nearestPickCurrent, groupBunchPickPossible)
	nearestThenBunchScore := 999999.0 // calcScore(nearestThenBunchTrucks)

	wg.Add(1)
	go  genericLoadProcessing(loadList3, loads3, &nearThenClosestLoadTrucks, nearestPickCurrent, findClosestPossibleStop)
	nearThenClosestLoadScore := calcScore(nearThenClosestLoadTrucks)

	wg.Add(1)
	go  genericLoadProcessing(loadList4, loads4, &bunchedFirstLoadTrucks, groupBunchPickCurrent, findClosestPossibleStop)
	bunchedFirstLoadScore := calcScore(bunchedFirstLoadTrucks)

	wg.Wait()

	lowestScore := slices.Min([]float64{nearestThenBunchScore,farThenClosestLoadScore,nearThenClosestLoadScore,bunchedFirstLoadScore})
	switch {
		case lowestScore == bunchedFirstLoadScore:
				displayResults(bunchedFirstLoadTrucks)	
		case lowestScore == nearestThenBunchScore:
				displayResults(nearestThenBunchTrucks)	
		case lowestScore == farThenClosestLoadScore:
				displayResults(farThenClosestLoadTrucks)
		case lowestScore == nearThenClosestLoadScore:
				displayResults(nearThenClosestLoadTrucks)
	} 
}
