package github

type Item interface {
	PullRequestWithRepository | IssueWithRepository
}

type store[T Item] map[int64]T

// add adds an Item to the map if the id of the Item isn't present already
func (s store[T]) add(key int64, val T) {
	_, exists := s[key]
	if exists {
		return
	}
	s[key] = val
}
