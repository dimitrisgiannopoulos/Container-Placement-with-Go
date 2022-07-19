package main

import "fmt"

func main() {
	var C int = 12 // number of containers in service
	var H int = 3  // number of physical hosts available for placement
	var R int = 2  // number of resources to measure (2 = CPUs, Memory)

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

	// decision variable (on what host do I place each container)
	var x bool;
	
 
 
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
	}
}
