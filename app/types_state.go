package main

import "math/rand"

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
		Settings: &Settings{
			MasterReplId: generateMasterReplId(),
		},
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

func generateMasterReplId() string {
	id := ""
	var val int
	for range 40 {
		val = rand.Intn(35)
		var c int
		if val < 10 {
			c = '0' + val
		} else {
			val -= 10
			c = 'a' + val
		}
		id += string(rune(c))
	}
	return id
}
