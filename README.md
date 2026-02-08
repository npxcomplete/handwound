# handwound
A clock / time interface for the go time package that enables testing of time based state changes and eventing.

For the performance sensitive, the function call overhead is negligible:

```
~/npxcomplete/w/handwound$ go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/npxcomplete/handwound
cpu: AMD Ryzen 7 7700X 8-Core Processor
BenchmarkRawClockNow-16         34159422                36.14 ns/op
BenchmarkSystemClockNow-16      33887870                36.00 ns/op
```
