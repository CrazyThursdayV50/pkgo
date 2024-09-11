package resty

import "time"

type Config struct {
	EnableTrace      bool
	EnableLog        bool
	Debug            bool
	Timeout          time.Duration
	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
	CloseConnection  bool
}
