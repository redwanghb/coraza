package corazawaf

var sseContentTypes []string = []string{
	"application/x-ndjson",
	"text/event-stream",
}

func IsSSEContent(s string) bool {
	for _, sse := range sseContentTypes {
		if sse == s {
			return true
		}
	}
	return false
}
