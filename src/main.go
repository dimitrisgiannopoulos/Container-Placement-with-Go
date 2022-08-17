package main

import "fmt"

func main() {

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
	var host_used_resources = [][]int{{359, 2685}, {437, 3414}, {305, 2451}}

	// average resource utilization among all hosts (avg CPU, avg MEM)
	var host_avg_resources = []float32{367, 2850}

	// resources capacity per host (CPU, MEM)
	var host_resource_capacities = [][]int{{8000, 7812}, {8000, 7812}, {8000, 7812}}

	// The usage ratio between two resources. This example means we want about the same CPU usage as MEM usage
	var target_ratio [][]float32 = [][]float32{{1, 1}, {1, 1}}
	// To have a target ratio of CPU/MEM = 0.5 we would use this {{1, 0.5}, {2, 1}}

	// decision variable (on what host is each container placed)
	// var x [C]int

	var arr [H]int // the hosts the containers can be placed on
	var n int = H  // items to be combined
	var r int = C  // sample size

	var placements [][]int // all possible placement combinations with repetitions
	var placement_costs []float32
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
	// fmt.Print("\n ", placement_costs) // Since we have the placement costs, it's now only a matter of keeping the best placements (lowest cost, and applying constraints)

	var indices []int
	for i := 0; i < len(placement_costs); i++ {
		indices = append(indices, i)
	}

	// fmt.Print("\n ", indices)

	placement_costs, indices = custom_quickSort(placement_costs, indices, 0, len(placement_costs)-1)
	fmt.Print("\n ", placement_costs)
	// fmt.Print("\n ", indices)

	// fmt.Print("\n ", placements[indices[0]])
	// fmt.Print("\n ", placements[indices[1]])
	// fmt.Print("\n ", placements[indices[2]])

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

/* square of a float32 */
func square(n float32) float32 {
	return n * n
}

/* Resource Utilization cost: The cost from using resources in an unbalanced way.
Unbalanced resource usage across hosts creates bottlenecks.*/
func Ucost_per_resource(per_host_used_resources int, per_container_resources int, per_host_avg_resources float32, H int) float32 {
	return (square(float32(per_host_used_resources+per_container_resources) - per_host_avg_resources)) / float32(H)
}

/* Residual resource balance cost: The cost from depleting one resource, while having another resource left.
Having 100% use of one resource makes the other resources unusable.*/
func Βcost_per_resource(host_resource_r1_capacity int, host_used_resource_r1 int, container_resource_r1 int, host_resource_r2_capacity int, host_used_resource_r2 int, target_ratio float32) float32 {
	var cost float32 = float32(host_resource_r1_capacity-host_used_resource_r1-container_resource_r1) - float32(host_resource_r2_capacity-host_used_resource_r2)*target_ratio
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
func calculate_total_costs(placements [][]int, w1 float32, w2 float32, w3 float32, C int, H int, R int,
	host_used_resources [][]int, container_resources [][]int, host_avg_resources []float32, host_resource_capacities [][]int, target_ratio [][]float32) []float32 {

	var total_costs []float32
	var total_Ucost float32 = 0
	var total_Bcost float32 = 0
	var total_Ccost float32 = 0

	for i := 0; i < len(placements); i++ {
		for k := 0; k < R; k++ {
			for j := 0; j < C; j++ {
				total_Ucost += Ucost_per_resource(host_used_resources[placements[i][j]][k], container_resources[j][k], host_avg_resources[k], H) // Ucost
				// total_Ccost += float32(Ccost_per_resource(placements[i], C, H)) // Ccost
				for l := k; l < R; l++ {
					if k != l {
						// total_Bcost += Βcost_per_resource(host_resource_capacities[placements[i][j]][k], host_used_resources[placements[i][j]][k], container_resources[j][k], host_resource_capacities[placements[i][j]][l], host_used_resources[placements[i][j]][l], target_ratio[k][l]) // Bcost
					}
				}
			}
		}

		weighted_total_cost := w1*total_Ucost + w2*total_Bcost + w3*total_Ccost
		total_costs = append(total_costs, weighted_total_cost)

		total_Ucost = 0
		total_Bcost = 0
		total_Ccost = 0
	}

	return total_costs
}

/* helper function for quicksort */
func partition(arr []float32, ind []int, low, high int) ([]float32, int, []int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			ind[i], ind[j] = ind[j], ind[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	ind[i], ind[high] = ind[high], ind[i]
	return arr, i, ind
}

/* Sorts low to high using the quicksort algorithm and returns the indices of the original array */
func custom_quickSort(arr []float32, ind []int, low, high int) ([]float32, []int) {
	if low < high {
		var p int
		arr, p, ind = partition(arr, ind, low, high)
		arr, ind = custom_quickSort(arr, ind, low, p-1)
		arr, ind = custom_quickSort(arr, ind, p+1, high)
	}
	return arr, ind
}
