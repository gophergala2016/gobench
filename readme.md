# goben.ch
The Golang community's repositories benchmarking


## Purpose
- Single point of access to all golang community's packages
- Automatic package benchmarking including: linter, vet, test, bench, etc
- Package benchmark history
- Favorite package notifications


## How to use?
- We believe you have github account
- If so, you do have some packages starred
- Visit http://goben.ch
- Sign in with your https://github.com account
- We'll take a list of your starred golang repositories
- From this moment http://goben.ch tracks these repository changes
- http://goben.ch pulls package and it's dependencies
- Do package benchmarking on different platforms and OS, including
	- Bare metal server with different OS
	- Several Digital Ocean portlets (5$, 10$)
	- GAE
	- Virtual machine
	- Docker in VM
	- Docker on bare metal
- Discover benchmarking reports on a package page http://goben.ch/p/{golang_package_name}

## TODO list
- Repositories supported
	- [X] https://github.com
	- [ ] https://bitbacket.com
	- [ ] https://labix.com
	- [ ] others
- Benchmarking
	- [X] go test -bench
	- [ ] go test
	- [ ] vet
	- [ ] go fmt
	- [ ] others
- Platforms
	- [X] Bare metal. Ubuntu 14.04, Intel i5, 4 Core
	- [ ] DigitalOcean
	- [ ] GAE
	- [ ] Virtual machine
	- [ ] Docker in VM
	- [ ] Docker on bare metal
	- [ ] others

## Inspiration
- Ruby world: https://rubygems.org/gems
- Rust: https://crates.io/

## License
MIT
