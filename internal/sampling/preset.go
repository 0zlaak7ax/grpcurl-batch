package sampling

// Preset returns a named Sampler for common use cases.
//
// Known names:
//   - "off"      – never sample (rate 0.0)
//   - "debug"    – sample 10% of requests
//   - "standard" – sample 1% of requests (default)
//   - "full"     – sample every request (rate 1.0)
//
// Unknown names fall back to "standard".
func Preset(name string) Sampler {
	switch name {
	case "off":
		return New(Config{Rate: 0})
	case "debug":
		return New(Config{Rate: 0.10})
	case "full":
		return New(Config{Rate: 1.0})
	default: // "standard" and unknown
		return New(Config{Rate: 0.01})
	}
}
