package nomad

import "time"

// intToPtr returns the pointer to an int
func intToPtr(i int) *int {
	return &i
}

// timeToPtr returns the pointer to a time stamp
func timeToPtr(t time.Duration) *time.Duration {
	return &t
}

// stringToPtr returns the pointer to a string
func stringToPtr(str string) *string {
	return &str
}

// boolToPtr returns the pointer to a boolean
func boolToPtr(b bool) *bool {
	return &b
}
