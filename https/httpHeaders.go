package https

// SetHeaders sets the header name and value in the provided headers map.
func SetHeaders(headers map[string]string, name, value string) {
	headers[name] = value
}

// MergeHeaders merges the headers from the `headersToMerge` map into the `headers` map.
// If a header with a key already exists in the `headers` map, it is skipped and not overwritten.
func MergeHeaders(headers, headersToMerge map[string]string) {
	for k, v := range headersToMerge {
		if _, ok := headers[k]; !ok {
			headers[k] = v
		}
	}
}
