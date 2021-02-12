package common

// StringSearcher is string search structure with internal cache
// This allows single iteration from base slice for multiple string searches
type StringSearcher interface {
	AddToIndex(keys ...string)
	Contains(key string) bool
}

type stringSearcher struct {
	index map[string]interface{}
}

func (s *stringSearcher) AddToIndex(keys ...string) {
	for _, k := range keys {
		s.index[k] = nil
	}
}

func (s *stringSearcher) Contains(key string) bool {
	_, ok := s.index[key]
	return ok
}

// Create new StringSearcher instance
func NewStringSearcher() StringSearcher {
	return &stringSearcher{
		index: make(map[string]interface{}),
	}
}
