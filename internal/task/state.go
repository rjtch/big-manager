package task

type State struct {
	slug string
}

func (s State) string() string {
	return s.slug
}

var (
	Unknown   = State{""}
	Pending   = State{"Pending"}
	Scheduled = State{"Scheduled"}
	Running   = State{"Running"}
	Failed    = State{"Failed"}
	Completed = State{"Completed"}
)

var StateTransitionMap = map[State][]State{
	Pending:   []State{Scheduled},
	Scheduled: []State{Scheduled, Scheduled, Scheduled},
	Running:   []State{Running, Completed, Failed},
	Completed: []State{},
	Failed:    []State{},
}
