package semverparse

// Mode constants define the available semver parsing strategies
const (
	// Strict requires exact semver format (e.g. "2.1.0")
	Strict = "strict"

	// Tolerant allows common variations (e.g. "2.1", "v2.1.0")
	Tolerant = "tolerant"
)
