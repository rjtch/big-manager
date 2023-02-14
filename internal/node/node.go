package node

// represent the physical aspect of a worker. It is an object that represents
// any machine in our cluster. For node  example, the manager is one type of in Cube.
// The worker, of which there can be more than node one, is another type of.
// The manager will make extensive use of node objects to representnode workers.
type Node struct {
	Name            string
	Ip              string
	Memory          int
	MemoryAllocated int
	Disk            int
	DiskAllocated   int
	TaskCount       int
}
