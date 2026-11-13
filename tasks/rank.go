package tasks

import "giligili/cache"

// RestartDailyRank resets the in-process daily ranking.
func RestartDailyRank() error {
	return cache.ResetDailyRank()
}
