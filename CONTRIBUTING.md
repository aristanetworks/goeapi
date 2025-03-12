# Contributing to the GoEAPI Project

## Testing changes locally

To test the changes locally we can make use of `replace directive` option in go.mod file. More details [here](https://go.dev/ref/mod#go-mod-file-replace).

For example, here's how go.mod file looks like after the changes,

```
module goeapi-test

go 1.24

require github.com/aristanetworks/goeapi v1.0.0

replace github.com/aristanetworks/goeapi => /Users/roopesh/Desktop/projects/arista/goeapi

require (
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec // indirect
)
```

