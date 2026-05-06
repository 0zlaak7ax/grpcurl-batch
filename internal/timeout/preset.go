package timeout

import "time"

// Preset returns a named Limiter configuration.
// Known names: "fast", "default", "slow". Unknown names fall back to Default.
func Preset(name string) *Limiter {
	switch name {
	case "fast":
		return New(Config{
			PerRequest: 5 * time.Second,
			Total:      30 * time.Second,
		})
	case "slow":
		return New(Config{
			PerRequest: 2 * time.Minute,
			Total:      30 * time.Minute,
		})
	default:
		return New(DefaultConfig())
	}
}
