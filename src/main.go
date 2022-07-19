package main

func main() {

	// -------------------------------------------------------------------------------------------------
	// constant: Containers (C) to be placed, physical Hosts (H) to be used and Resources (R) to examine
	const C int = 12 // number of containers in service
	const H int = 3  // number of physical hosts available for placement
	const R int = 2  // number of resources to measure (2 = CPUs, Memory)
	// -------------------------------------------------------------------------------------------------

	// -------------------------------------------------------------------------------------------------
	// resource values initialization and decision variable declaration

	// resource demands in CPU, MEM per containers
	var container_resources = [C][R]int{
		{300, 300}, {300, 128}, {200, 128}, {200, 128},
		{200, 128}, {200, 128}, {500, 512}, {200, 128},
		{200, 128}, {200, 450}, {125, 256}, {200, 128},
	}

	// resources currently in use by hosts
	var host_used_resources = [H][R]int{{359, 2685}, {437, 3414}, {305, 2451}}

	// average resource utilization among all hosts (avg CPU, avg MEM)
	var host_avg_resources = [R]float32{367, 2850}

	// resources capacity per host (CPU, MEM)
	var host_resource_capacities = [H][R]int{{8000, 7812}, {8000, 7812}, {8000, 7812}}

	var target_ratio float32 = 0.1

	// decision variable (on what host is each container placed)
	var x [C]int

	var arr = [H]int{0, 1, 2} // the hosts the containers can be placed on
	var n int = H
	var r int = C
	// -------------------------------------------------------------------------------------------------

	// -------------------------------------------------------------------------------------------------
	// main code

	find_placement_combinations(arr, n, r) // finds all possible combinations with repetitions
	// -------------------------------------------------------------------------------------------------

	// minimize
	// (sum(c in containers, r in resources, h in hosts)
	// 	(host_used_resources[h, r] + x[c, h]*container_resources[c, r]*10 - host_avg_resources[r])/H) +
	// (sum(c in containers, h in hosts)
	// 	abs((host_resource_capacities[h, 0] - host_used_resources[h, 0] - container_resources[c, 0]*10 -
	// 	(host_resource_capacities[h, 1] - host_used_resources[h, 1])*target_ratio)) ) ; //(r1+1)%R in case of 3 resources

	// subject to{
	//  forall(c in containers) sum(h in hosts) x[c, h] == 1; //all containers are placed exactly in one host

	//  forall(c in containers)
	//    forall(r in resources)
	// 	 forall(h in hosts)
	// 		 x[c, h]*container_resources[c, r]*10 <= host_resource_capacities[h, r] - host_used_resources[h, r]; //container resources are less than available

	//   forall (c in containers)
	//     forall(h in hosts)
	//       x[c, h] >= x[(c+1)%C, h];
	// }
}

// Resource Utilization cost: The cost from using resources in an unbalanced way.
// Unbalanced resource usage across hosts creates bottlenecks.
func Ucost(host_used_resources int, container_resources int, host_avg_resources float32, H int) float32 {
	return (float32(host_used_resources+container_resources) - host_avg_resources) / float32(H)
}

// Residual resource balance cost: The cost from depleting one resource, while having another resource left.
// Having 100% use of one resource makes the other resources unusable.
func Βcost(host_resource_r1_capacity int, host_used_resource_r1 int, container_resource_r1 int, host_resource_r2_capacity int, host_used_resource_r2 int, target_ratio float32) float32 {
	var cost float32 = float32(host_resource_r1_capacity-host_used_resource_r1-container_resource_r1) - float32(host_resource_r2_capacity-host_used_resource_r2)*target_ratio
	if cost > 0 {
		return cost
	}
	return 0
}

// Comunication cost: The cost of placing containers of the same service on different physical hosts.
// Placing containers on near servers minimizes communication cost.
// If two containers are placed on the same server they contribute 0 to the cost, otherwise, they contribute 1
func Ccost(x []int, C int, H int) int {
	var cost int = 0

	for i := 0; i < H; i++ {
		for j := 0; j < C; j++ {
			if x[j] != i {
				cost++
			}
		}
	}
	return cost
}
