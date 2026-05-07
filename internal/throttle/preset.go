package throttle

// Preset returns a named Throttler configuration.
//
// Recognised names:
//
//	"unlimited"  – no effective throttling
//	"low"        – 5 req/s, burst 10
//	"medium"     – 50 req/s, burst 100
//	"high"       – 200 req/s, burst 400
//
// Any unknown name falls back to "unlimited".
func Preset(name string) *Throttler {
	switch name {
	case "low":
		return New(Config{Rate: 5, Burst: 10})
	case "medium":
		return New(Config{Rate: 50, Burst: 100})
	case "high":
		return New(Config{Rate: 200, Burst: 400})
	default:
		return New(Config{})
	}
}
