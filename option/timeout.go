package option

import "time"

type ReadTimeoutOption time.Duration

func (o ReadTimeoutOption) isOption() {}

func ReadTimeout(duration time.Duration) Option {
	return ReadTimeoutOption(duration)
}

type ReadHeaderTimeoutOption time.Duration

func (o ReadHeaderTimeoutOption) isOption() {}

func ReadHeaderTimeout(duration time.Duration) Option {
	return ReadHeaderTimeoutOption(duration)
}

type WriteTimeoutOption time.Duration

func (o WriteTimeoutOption) isOption() {}

func WriteTimeout(duration time.Duration) Option {
	return WriteTimeoutOption(duration)
}

type IdleTimeoutOption time.Duration

func (o IdleTimeoutOption) isOption() {}

func IdleTimeout(duration time.Duration) Option {
	return IdleTimeoutOption(duration)
}
