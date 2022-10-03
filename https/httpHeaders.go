package https

func SetHeaders(headers map[string]string, name, value string) {
	headers[name] = value
}

func MergeHeaders(headers, headersToMerge map[string]string) {
	for k, v := range headersToMerge {
		if _, ok := headers[k]; !ok {
			headers[k] = v
		}
	}
}
