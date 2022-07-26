package main

import (
	"fmt"
	"sort"
	"time"
)

type Slice struct {
	sort.Float64Slice
	idx []int
}

func main() {

	start := time.Now()

	// -------------------------------------------------------------------------------------------------
	// constant: Containers (C) to be placed, physical Hosts (H) to be used and Resources (R) to examine
	const C int = 12 // number of containers in service
	const H int = 3  // number of physical hosts available for placement
	const R int = 2  // number of resources to measure (2 = CPUs, Memory).
	// ATTENTION: The code is not reconfigurable regarding the number of resources. If more resources are taken into account, the calculate_total_costs function needs changes
	// A different Bcost needs to be calculated per two resources and a specific target ratio has to be used
	// -------------------------------------------------------------------------------------------------

	// -------------------------------------------------------------------------------------------------
	// resource values initialization and decision variable declaration

	// resource demands in CPU, MEM per containers
	var container_resources = [][]int{
		{300, 300}, {300, 128}, {200, 128}, {200, 128},
		{200, 128}, {200, 128}, {500, 512}, {200, 128},
		{200, 128}, {200, 450}, {125, 256}, {200, 128},
	}

	// resources currently in use by hosts
	// var host_used_resources = [][]int{{359, 2685}, {437, 3414}, {305, 2451}}	// 3 servers
	// var host_used_resources = [][]int{{359, 2685}, {437, 3414}, {305, 2451}, {319, 2497}, {359, 2797}, {306, 2397}, {317, 2479}, {361, 2811}, {360, 2798}, {387, 3001}} // 10 servers
	var host_used_resources = [][]int{{359, 2685}, {437, 3414}, {305, 2451}, {319, 2497}, {359, 2797}, {306, 2397}, {317, 2479}, {361, 2811}, {360, 2798}, {387, 3001},
		{369, 2865}, {309, 2418}, {413, 3195}, {436, 3366}, {432, 3337}, {378, 2937}, {385, 2989}, {350, 2725}, {409, 3167}, {335, 2616},
		{356, 2768}, {375, 2914}, {331, 2584}, {422, 3259}, {415, 3207}, {417, 3223}, {346, 2694}, {322, 2517}, {370, 2872}, {412, 3187},
		{428, 3309}, {386, 2993}, {360, 2799}, {330, 2582}, {414, 3199}, {332, 2591}, {321, 2509}, {365, 2842}, {363, 2826}, {409, 3163},
		{377, 2929}, {378, 2933}, {331, 2585}, {393, 3049}, {320, 2502}, {347, 2701}, {432, 3337}, {409, 3169}, {417, 3223}, {379, 2943},
		{415, 3212}, {406, 3146}, {434, 3353}, {433, 3347}, {342, 2667}, {309, 2425}, {400, 3096}, {391, 3032}, {320, 2501}, {421, 3253},
		{362, 2818}, {322, 2519}, {420, 3248}, {306, 2399}, {328, 2564}, {423, 3272}, {360, 2801}, {374, 2908}, {309, 2421}, {376, 2922},
		{338, 2638}, {428, 3307}, {434, 3352}, {366, 2847}, {352, 2738}, {323, 2525}, {376, 2920}, {366, 2844}, {315, 2469}, {388, 3013},
		{366, 2845}, {310, 2432}, {412, 3186}, {404, 3126}, {412, 3186}, {406, 3144}, {361, 2810}, {386, 2995}, {361, 2810}, {398, 3082},
		{426, 3291}, {401, 3104}, {364, 2830}, {317, 2478}, {333, 2602}, {329, 2571}, {375, 2916}, {329, 2569}, {362, 2819}, {397, 3074}} // 100 servers

	// average resource utilization among all hosts (avg CPU, avg MEM)
	var host_avg_resources = []float64{367, 2850}

	// resources capacity per host (CPU, MEM)
	// var host_resource_capacities = [][]int{{8000, 7812}, {8000, 7812}, {8000, 7812}}	// 3 servers
	// var host_resource_capacities = [][]int{{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}} // 10 servers
	var host_resource_capacities = [][]int{{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812},
		{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812},
		{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812},
		{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812},
		{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812},
		{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812},
		{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812},
		{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812},
		{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812},
		{8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}, {8000, 7812}} // 100 servers

	// The usage ratio between two resources. This example means we want about the same CPU usage as MEM usage
	var target_ratio [][]float64 = [][]float64{{1, 1}, {1, 1}}
	// To have a target ratio of CPU/MEM = 0.5 we would use this {{1, 0.5}, {2, 1}}

	var arr [H]int // the hosts the containers can be placed on
	var n int = H  // items to be combined
	var r int = C  // sample size

	var placements [][]int // all possible placement combinations with repetitions
	var placement_costs []float64
	// -------------------------------------------------------------------------------------------------

	// -------------------------------------------------------------------------------------------------
	// main code

	// initialize array of hosts, ex. arr[] = {0 1 2} for 3 hosts
	for i := 0; i < H; i++ {
		arr[i] = i
	}

	find_placement_combinations(arr[:], n, r, &placements) // finds all possible combinations with repetitions
	// fmt.Print(" ", placements)                             // combinations[placements][hosts]

	placement_costs = calculate_total_costs(placements, 1.0, 1.0, 1.0, C, H, R, host_used_resources, container_resources, host_avg_resources, host_resource_capacities, target_ratio)
	// fmt.Println(placement_costs) // Since we have the placement costs, it's now only a matter of keeping the best placements (lowest cost, and applying constraints)

	placement_costs_slice := NewSlice(placement_costs)
	sort.Sort(placement_costs_slice)

	// It prints the 3 best costs from low to high and the corresponding indices. The best solutions have the lowest cost.
	for i := 0; i < 3; i++ {
		fmt.Println(i+1, ") cost:", placement_costs_slice.Float64Slice[i], "- index:", placement_costs_slice.idx[i], "- placement:", placements[placement_costs_slice.idx[i]])
	}
	// fmt.Println(placement_costs_slice.Float64Slice, placement_costs_slice.idx)

	elapsed := time.Since(start)
	fmt.Println("program took", elapsed)
	// -------------------------------------------------------------------------------------------------
}

/* The main function that prints all combinations of size r
in arr[] of size n with repetitions. This function mainly
uses CombinationRepetitionUtil() */
func find_placement_combinations(arr []int, n int, r int, placements *[][]int) {
	// Allocate memory
	chosen := make([]int, r+1)

	// Call the recursive function
	CombinationRepetitionUtil(placements, chosen, arr, 0, r, 0, n-1)
}

/* arr[]  ---> Input Array
   chosen[] ---> Temporary array to store indices of current combination
   start & end ---> Starting and Ending indexes in arr[]
   r ---> Size of a combination to be printed */
func CombinationRepetitionUtil(placements *[][]int, chosen []int, arr []int, index int, r int, start int, end int) {
	// Since index has become r, current combination is
	// ready to be printed, print
	if index == r {
		var tmp []int
		for i := 0; i < r; i++ {
			tmp = append(tmp, arr[chosen[i]])
		}
		*placements = append(*placements, tmp)
		return
	}

	// One by one choose all elements (without considering
	// the fact whether element is already chosen or not)
	// and recur
	for i := start; i <= end; i++ {
		chosen[index] = i
		CombinationRepetitionUtil(placements, chosen, arr, index+1, r, i, end)
	}
	return
}

/* square of a float64 */
func square(n float64) float64 {
	return n * n
}

/* Resource Utilization cost: The cost from using resources in an unbalanced way.
Unbalanced resource usage across hosts creates bottlenecks.*/
func Ucost_per_resource(per_host_used_resources int, per_container_resources int, per_host_avg_resources float64, H int) float64 {
	return (square(float64(per_host_used_resources+per_container_resources) - per_host_avg_resources)) / float64(H)
}

/* Residual resource balance cost: The cost from depleting one resource, while having another resource left.
Having 100% use of one resource makes the other resources unusable.*/
func Βcost_per_resource(host_resource_r1_capacity int, host_used_resource_r1 int, container_resource_r1 int, host_resource_r2_capacity int, host_used_resource_r2 int, target_ratio float64) float64 {
	var cost float64 = float64(host_resource_r1_capacity-host_used_resource_r1-container_resource_r1) - float64(host_resource_r2_capacity-host_used_resource_r2)*target_ratio
	if cost > 0 {
		return cost
	} else {
		return -cost
	}
}

/* Comunication cost: The cost of placing containers of the same service on different physical hosts.
Placing containers on near servers minimizes communication cost.
If two containers are placed on the same server they contribute 0 to the cost, otherwise, they contribute 1 */
func Ccost_per_resource(placement []int, C int, H int) int {
	var cost int = 0

	for i := 0; i < C; i++ {
		for j := 0; j < C; j++ {
			if placement[i] != placement[j] && i != j {
				cost++
			}
		}
	}

	cost /= 2 // To account for counting the same relationship twice (bidirectional cost)
	return cost
}

/* The total cost calculated by adding the weighted different costs */
func calculate_total_costs(placements [][]int, w1 float64, w2 float64, w3 float64, C int, H int, R int,
	host_used_resources [][]int, container_resources [][]int, host_avg_resources []float64, host_resource_capacities [][]int, target_ratio [][]float64) []float64 {

	var total_costs []float64
	var total_Ucost []float64
	var total_Bcost []float64
	var total_Ccost []float64

	for i := 0; i < len(placements); i++ {
		total_Ucost = append(total_Ucost, 0)
		total_Ccost = append(total_Ccost, 0)
		total_Bcost = append(total_Bcost, 0)
		for k := 0; k < R; k++ {
			for j := 0; j < C; j++ {
				total_Ucost[i] += Ucost_per_resource(host_used_resources[placements[i][j]][k], container_resources[j][k], host_avg_resources[k], H) // Ucost
				total_Ccost[i] += float64(Ccost_per_resource(placements[i], C, H))                                                                  // Ccost
				for l := k; l < R; l++ {
					if k != l {
						total_Bcost[i] += Βcost_per_resource(host_resource_capacities[placements[i][j]][k], host_used_resources[placements[i][j]][k], container_resources[j][k], host_resource_capacities[placements[i][j]][l], host_used_resources[placements[i][j]][l], target_ratio[k][l]) // Bcost
					}
				}
			}
		}
	}

	min_Ucost, max_Ucost := findMinAndMax(total_Ucost)
	min_Ccost, max_Ccost := findMinAndMax(total_Ccost)
	min_Bcost, max_Bcost := findMinAndMax(total_Bcost)

	for i := 0; i < len(placements); i++ {

		total_Ucost[i] = (total_Ucost[i] - min_Ucost) / (max_Ucost - min_Ucost)
		total_Ccost[i] = (total_Ccost[i] - min_Ccost) / (max_Ccost - min_Ccost)
		total_Bcost[i] = (total_Bcost[i] - min_Bcost) / (max_Bcost - min_Bcost)

		weighted_total_cost := w1*total_Ucost[i] + w2*total_Ccost[i] + w3*total_Bcost[i]
		total_costs = append(total_costs, weighted_total_cost)
	}

	return total_costs
}

func (s Slice) Swap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(n []float64) *Slice {
	s := &Slice{Float64Slice: sort.Float64Slice(n), idx: make([]int, len(n))}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

/* returns the min and max of a float64 array */
func findMinAndMax(a []float64) (min float64, max float64) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}
