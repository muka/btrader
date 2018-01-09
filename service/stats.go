package service

//StatsFilter filter applied on stats retrieval
type StatsFilter struct {
	BaseCoin string
	Asset    []string
}

//StatsView wrap market stats details
type StatsView struct {
}

//Stats elaborate market stats
func Stats(filter StatsFilter) (*StatsView, error) {
	stats := &StatsView{}

	return stats, nil
}
