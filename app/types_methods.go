package main

import (
	"maps"
	"slices"
	"strconv"
)

func (m *MillisVal) GetOrCreateSequence(id string) (SequenceKV, error) {
	sVal, ok := m.Map[id]
	if ok {
		// already exists
		return sVal, nil
	}
	// create it and then ensure the keys are sorted
	sVal = SequenceKV{}
	m.Map[id] = sVal

	// cannot happen, has been validated before
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return SequenceKV{}, err
	}

	// add he id to the list of keys that should be stored in sorted order
	m.Keys = append(m.Keys, id)
	// this is greater, no need to rebalance
	if idInt > m.Last {
		m.Last = idInt
	} else {
		slices.Sort(m.Keys)
	}
	return sVal, nil
}

func (s *Stream) GetOrCreateMillis(id string) *MillisVal {
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

	// add the id to the list of keys that should be stored in sorted order
	s.Keys = append(s.Keys, id)
	// this is greater, no need to rebalance
	if idInt > s.Last {
		s.Last = idInt
	} else {
		slices.Sort(s.Keys)
	}
	return millisVal
}

func (s *Stream) GenerateMillis() string {
	return strconv.Itoa(s.Last + 1)
}

func (s *Stream) GenerateSequence(millis string) string {
	millisVal := s.GetOrCreateMillis(millis)
	return strconv.Itoa(millisVal.Last + 1)
}

func (s *Stream) AddIdWithKV(millis, sequence string, kvArg map[string]string) error {
	millisVal := s.GetOrCreateMillis(millis)
	kv, err := millisVal.GetOrCreateSequence(sequence)
	if err != nil {
		return err
	}
	maps.Copy(kv, kvArg)
	return nil
}

func (s *Stream) IsNewStreamIdValid(millis, sequence string) (bool, error) {
	// ValidationRules: ID must be greater/equal than last entry's ID:
	// millis Must be >= lastEntry milis
	// if millis is equal, sequence must be > lastEntry sequence

	// must be greater than or equal, anything less is invalid
	millisInt, err := strconv.Atoi(millis)
	if err != nil {
		return false, err
	}

	if millisInt < s.Last {
		return false, nil
	}

	if sequence == "*" {
		return true, nil
	}

	// check if we have an entry for this millis, if not, it is valid, since sequence cannot clash
	millisVal, ok := s.Map[millis]
	if !ok {
		return true, nil
	}

	sequenceInt, err := strconv.Atoi(sequence)
	if err != nil {
		return false, err
	}

	// ensure sequence is greater
	if sequenceInt <= millisVal.Last {
		return false, nil
	}

	return true, nil
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
