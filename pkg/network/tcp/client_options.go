package tcp

type ClientOptionFunc func(o *clientOptions)

type clientOptions struct {
	addr         string
	maxMsgLength int
}

func defaultClientOptions() *clientOptions {
	return &clientOptions{
		addr:         "",
		maxMsgLength: 1024 * 1024,
	}
}

func WithClientMaxMsgLength(maxMsgLength int) ClientOptionFunc {
	return func(o *clientOptions) { o.maxMsgLength = maxMsgLength }
}
