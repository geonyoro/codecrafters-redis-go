package main

type State struct {
	Settings    *Settings
	VariableMap map[string]Variable
	ListMap     map[string]*ListVariable
	StreamMap   map[string]*Stream
}

func NewState() *State {
	vMap := make(map[string]Variable)
	lMap := make(map[string]*ListVariable)
	sMap := make(map[string]*Stream)
	return &State{
		Settings:    &Settings{},
		VariableMap: vMap,
		ListMap:     lMap,
		StreamMap:   sMap,
	}
}

func (s *State) WithReplicaOf(replicaOf string) {
	s.Settings.ReplicaOf = replicaOf
}

func (s *State) updateWithCliArgs(args *CliArgs) {
	if len(args.ReplicaOf) > 0 {
		s.WithReplicaOf(args.ReplicaOf)
	}
}
