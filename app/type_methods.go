package main

import (
	"maps"
	"slices"
	"strconv"
)

func (m *MillisVal) AddSequence(id string) SequenceKV {
	sVal, ok := m.Map[id]
	if ok {
		// already exists
		return sVal
	}
	// create it and then ensure the keys are sorted
	sVal = SequenceKV{}
	m.Map[id] = sVal

	// cannot happen, has been validated before
	idInt, _ := strconv.Atoi(id)

	// add he id to the list of keys that should be stored in sorted order
	m.Keys = append(m.Keys, id)
	// this is greater, no need to rebalance
	if idInt > m.Last {
		m.Last = idInt
	} else {
		slices.Sort(m.Keys)
	}
	return sVal
}

func (s *Stream) AddMillis(id string) *MillisVal {
	millisVal, ok := s.Map[id]
	if ok {
		// already exists
		return millisVal
	}
	// create it and then ensure the keys are sorted
	millisVal = NewMillisVal()
	s.Map[id] = millisVal

	// cannot happen, has been validated before
	idInt, _ := strconv.Atoi(id)

	// add he id to the list of keys that should be stored in sorted order
	s.Keys = append(s.Keys, id)
	// this is greater, no need to rebalance
	if idInt > s.Last {
		s.Last = idInt
	} else {
		slices.Sort(s.Keys)
	}
	return millisVal
}

func (s *Stream) AddIdWithKV(millis, sequence string, kvArg map[string]string) {
	millisVal := s.AddMillis(millis)
	kv := millisVal.AddSequence(sequence)
	maps.Copy(kv, kvArg)
}

func (s *Stream) IsNewStreamIdValid(millis, sequence int) bool {
	// ValidationRules: ID must be greater/equal than last entry's ID:
	// millis Must be >= lastEntry milis
	// if millis is equal, sequence must be > lastEntry sequence

	// must be greater than or equal, anything less is invalid
	if millis < s.Last {
		return false
	}

	// check if we have an entry for this millis, if not, it is valid, since sequence cannot clash
	millisVal, ok := s.Map[strconv.Itoa(millis)]
	if !ok {
		return true
	}

	// ensure sequence is greater
	if sequence <= millisVal.Last {
		return false
	}

	return true
}

func (s *State) GetOrCreateStreamForKey(key string) *Stream {
	streamMap := *s.StreamMap
	stream, ok := streamMap[key]
	if ok {
		return stream
	}
	// make the stream
	stream = NewStream()
	streamMap[key] = stream
	return stream
}
