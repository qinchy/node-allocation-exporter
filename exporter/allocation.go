package exporter

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

type AllocationData struct {
	ID         string
	requestCpu float64
	requestMem float64
}

// getAllocationData get allocation data
func getAllocationData() []AllocationData {
	var ret []AllocationData
	buildUpload := func(node string, requestCpu, requestMem float64) AllocationData {
		return AllocationData{
			ID:         node,
			requestCpu: requestCpu,
			requestMem: requestMem,
		}
	}

	// 调用获取pod方法
	nodeLists, error := getNodes()
	if error != nil {
		log.Panic("请求获取节点列表接口出错")
		panic("请求获取节点列表接口出错")
	}

	if nodeLists == nil {
		log.Println("nodeLists 为空")
		return nil
	}

	for _, nodeInfo := range nodeLists.Items {
		log.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t,\n",
			nodeInfo.Name,
			nodeInfo.Status.Phase,
			nodeInfo.Status.Addresses,
			nodeInfo.Status.NodeInfo.OSImage,
			nodeInfo.Status.NodeInfo.KubeletVersion,
			nodeInfo.Status.NodeInfo.OperatingSystem,
			nodeInfo.Status.NodeInfo.Architecture,
		)
		// cpu总量
		cpuCapacity := nodeInfo.Status.Capacity.Cpu()
		log.Printf("Cpu Capacity: %s \n", cpuCapacity)
		// cpu可分配量
		cpuAllocatable := nodeInfo.Status.Allocatable.Cpu()
		log.Printf("Cpu Allocatable: %s \n", cpuAllocatable.String())
		// cpu已分配量
		cpuAllocated := cpuCapacity.Value() - cpuAllocatable.Value()
		log.Printf("Cpu cpuAllocated: %d \n", cpuAllocated)

		// mem总量
		memCapacity := nodeInfo.Status.Capacity.Memory()
		log.Printf("Mem Capacity: %s \n", memCapacity)
		// mem可分配量
		memAllocatable := nodeInfo.Status.Allocatable.Memory()
		log.Printf("Mem Allocatable: %s \n", memAllocatable.String())
		// mem已分配量
		memAllocated := memCapacity.Value() - memAllocatable.Value()
		log.Printf("Mem cpuAllocated: %d \n", memAllocated)

		// 构造结构体数据
		ret = append(ret, buildUpload(nodeInfo.Name, float64(cpuAllocated), float64(memAllocated)))
	}

	return ret
}

// buildAllocationData build allocation data
func buildAllocationData() AllocationData {
	return AllocationData{
		ID:         "node2",
		requestCpu: float64(30),
		requestMem: float64(40),
	}
}

// getNodes 获取所有node信息
func getNodes() (*v1.NodeList, error){
	// 获取所有的Node信息
	return clientset.CoreV1().Nodes().List(metav1.ListOptions{})
}
