package main

import "time"

var layouts = []string{
	time.RFC1123Z,
    time.RFC1123,
    time.RFC3339,
    time.RFC822,
    time.RFC822Z,
}

func parsePublishedAt(date string) (time.Time, bool) {
	for _, layout := range layouts {
		t, err := time.Parse(layout, date)
		if err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}