package scheduler

// scheduler methods to manage task to be assigned to nodes
// this scheduler uses round-robin algorithm that kept a list of workers and identified which worker got the most
// recent task. Then, when the next task came in, the scheduler simply picked the next worker in its
// list
type scheduler interface {
	SelecCandidateNodes()
	Score()
	Pick()
}
