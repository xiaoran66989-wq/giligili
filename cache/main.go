package cache

import (
	"sort"
	"sync"
)

var (
	rankMu sync.RWMutex
	dailyViews = make(map[uint]uint64)
)

// RecordDailyView records a view for the in-process daily ranking.
func RecordDailyView(id uint) {
	rankMu.Lock()
	dailyViews[id]++
	rankMu.Unlock()
}

// DailyRankIDs returns the top video IDs for the current process day.
func DailyRankIDs(limit int) []uint {
	rankMu.RLock()
	pairs := make([]struct {
		id uint
		views uint64
	}, 0, len(dailyViews))
	for id, views := range dailyViews {
		pairs = append(pairs, struct {
			id uint
			views uint64
		}{id, views})
	}
	rankMu.RUnlock()

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].views == pairs[j].views {
			return pairs[i].id < pairs[j].id
		}
		return pairs[i].views > pairs[j].views
	})
	if limit > len(pairs) {
		limit = len(pairs)
	}
	ids := make([]uint, 0, limit)
	for _, pair := range pairs[:limit] {
		ids = append(ids, pair.id)
	}
	return ids
}

// ResetDailyRank resets the in-process daily ranking.
func ResetDailyRank() error {
	rankMu.Lock()
	dailyViews = make(map[uint]uint64)
	rankMu.Unlock()
	return nil
}
