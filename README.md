# matching-engine
A simple matching engine for crypto exchange in golang

TODO
===
TODO:
* price and amount units
* Cancel order
* Market price
* Logging and Json
* Redis
* Docker

Profilling
===
```
> go test -cpuprofile cpu.prof -memprofile mem.prof -bench .
> go tool pprof cpu.prof
(pprof) pdf
(pprof) list go
```

Reference
===
inspired by https://github.com/bhomnick/bexchange