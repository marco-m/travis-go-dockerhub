package hello

import "os"

func Hello() string {
	if val, ok := os.LookupEnv("THE_SECRET"); ok {
		return val
	}
	return "hello"
}
