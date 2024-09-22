package main

import (
    "bufio"
    "fmt"
    "os"
	"math"
	"strings"
	"sort"
	"cmp"
	"slices"
)

type load struct {
	loadNumber    int
    pickup        point
    dropoff       point
	distToPickup  float64
	deliveryDist  float64
	distToHome    float64
	distToOthers  []loadListEntry
}

type point struct {
	x float64
	y float64
}

type truck struct {
	mileage float64
	stops []int
}

type loadListEntry struct {
	loadNumber int
	loadTotalDistance float64
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func calcDistToOthers(ls []load) ([]load) {
	for i,l := range ls {
		l.distToOthers = make([]loadListEntry,len(ls))
		for j, d := range ls {
			l.distToOthers[j].loadNumber = j + 1
			l.distToOthers[j].loadTotalDistance =  math.Sqrt(math.Pow((l.dropoff.x-d.pickup.x),2) + math.Pow((l.dropoff.y-d.pickup.y),2))
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
	var loadList []loadListEntry

    scanner := bufio.NewScanner(file)
	_ = scanner.Scan()
    for scanner.Scan() {
		var l load
		fmt.Sscanf(scanner.Text(), "%d (%f,%f) (%f,%f)", &l.loadNumber, &l.pickup.x, &l.pickup.y, &l.dropoff.x, &l.dropoff.y)
        l.deliveryDist =  math.Sqrt(math.Pow((l.dropoff.x-l.pickup.x),2) + math.Pow((l.dropoff.y-l.pickup.y),2))
        l.distToPickup =  math.Sqrt(math.Pow((l.pickup.x),2) + math.Pow((l.pickup.y),2))
        l.distToHome   =  math.Sqrt(math.Pow((l.dropoff.x),2) + math.Pow((l.dropoff.y),2))
		//fmt.Println(l)
		loads = append(loads,l)
		loadList = append(loadList, loadListEntry{l.loadNumber,l.deliveryDist+l.distToPickup+l.distToHome})
    }
	maxTrucks := len(loadList)
	loadsToGo := len(loadList)
	trucks :=  make([]truck, maxTrucks)  

	// by sorting to the longest distance first, plan on getting as many as possible on the way back
	//fmt.Printf("before sort of loadList %v \n",loadList)
	sort.Slice(loadList, func(i, j int) bool { return loadList[i].loadTotalDistance < loadList[j].loadTotalDistance})
	//fmt.Printf("after  sort of loadList %v \n",loadList)

	var currentStop loadListEntry
	var possibleStop loadListEntry
	loads = calcDistToOthers(loads)

	deliverLoop:
	for t:=0; t < maxTrucks; t++ {
		if loadsToGo == 0 {
			break deliverLoop
		}
		currentStop, loadList = loadList[0], loadList[1:]
		//currentStopIndex := currentStop.loadNumber - 1
		currentStopIndex := slices.IndexFunc(loads, func(l load) bool {
			return l.loadNumber == currentStop.loadNumber
		})
		trucks[t].stops = append(trucks[t].stops,loads[currentStopIndex].loadNumber)
		trucks[t].mileage = loads[currentStopIndex].distToPickup + loads[currentStopIndex].deliveryDist + loads[currentStopIndex].distToHome
		loadsToGo--
		for i:=0; i < loadsToGo; i++  {
			// pop closest distToOthers
			// fmt.Printf("pre  sort %v \n",loads[currentStopIndex].distToOthers)
			slices.SortFunc(loads[currentStopIndex].distToOthers, func(a, b loadListEntry) int {
				return cmp.Compare(a.loadTotalDistance, b.loadTotalDistance)
			})
			// fmt.Printf("postsort %v \n",loads[currentStopIndex].distToOthers)

			for k:= 0; k < len(loadList); k++ {
				idx := slices.IndexFunc(loadList, func(l loadListEntry) bool {
					return l.loadNumber == loads[currentStopIndex].distToOthers[k].loadNumber
				})
				//fmt.Printf("in loop %d with idx=%d after search for %d \n",k,idx,loads[currentStopIndex].distToOthers[k].loadNumber)
				if idx != -1 {
					possibleStop = loadList[idx]
					loadList = slices.Delete(loadList,idx,idx+1)
					//possibleStopIndex := possibleStop.loadNumber - 1
					possibleStopIndex := slices.IndexFunc(loads, func(l load) bool {
						return l.loadNumber == possibleStop.loadNumber
					})
					possibleNewMileage := trucks[t].mileage - loads[currentStopIndex].distToHome + loads[currentStopIndex].distToOthers[possibleStopIndex].loadTotalDistance + loads[possibleStopIndex].deliveryDist + loads[possibleStopIndex].distToHome
					if possibleNewMileage <= 720 {
						trucks[t].stops = append(trucks[t].stops,possibleStop.loadNumber)
						trucks[t].mileage = possibleNewMileage
						currentStop = possibleStop
						//currentStopIndex = possibleStopIndex
						currentStopIndex = slices.IndexFunc(loads, func(l load) bool {
							return l.loadNumber == currentStop.loadNumber
						})
						loadsToGo--
					} else {
						loadList = append(loadList,possibleStop)
					}		
				}
			}
		}
		loadsToGo = len(loadList)
	}
	for _,t := range(trucks) {
	    if len(t.stops) > 0 {
			outString := fmt.Sprint(t.stops)
			fmt.Print(strings.ReplaceAll(outString," ",","))
			fmt.Printf(" %f \n",t.mileage)
	    }
	}
}
