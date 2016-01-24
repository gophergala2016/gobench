# goben.ch
The Golang community's repositories benchmarking and tracking service

## Purpose
- Single point of access to all golang community's packages
- Automatic package benchmarking including: linter, vet, test, bench, etc
- Package changes capture and re-benchmarking
- Package benchmark history
- User notifications

## Inspiration
- Ruby world: https://rubygems.org
- Rust: https://crates.io
- Go world: https://godoc.org

## How to use?
- Most probably you have github account
- If so, you are interested in some packages
- Visit http://goben.ch
- Sign in with your https://github.com account
- We'll take a list of your starred golang repositories
- http://goben.ch pulls package and it's dependencies
- http://goben.ch benchmarks package on different platforms and OS, including
	- Bare metal server with different OS
	- Several Digital Ocean portlets (5$, 10$)
	- GAE
	- Virtual machine
	- Docker in VM
	- Docker on bare metal
- Discover benchmarking reports on a package page http://goben.ch/p/{golang_package_name}
- http://goben.ch starts keeping eye on this package
- http://goben.ch re-benchmarks if changes captured

## Gopher Gala 2016
The Gopher Gala is the world's first Go(lang) distributed hackathon.
We have enetered the competition with idea of building service prototype.
Planned prototype scope:

- Repositories supported
	- [X] https://github.com
- Packages capture
	- [X] User's favorites on github
- Vendoring
	- [X] Standard
- Package benchmarking
	- [X] go test -bench -bechmem
- Package dashboard
	- [X] General information
- Platforms
	- [X] Bare metal, Ubuntu

### What next?
We will continue, definitely! Follow us https://twitter.com/goben_ch
After GopherGala, development will continue here https://github.com/goben-ch/gobench

## TODO List
- Repositories supported
	- [X] https://github.com
	- [ ] https://bitbucket.org
	- [ ] https://labix.org
	- [ ] others
- Packages capture
	- [X] User's favorites on github
	- [ ] Manually
	- [ ] Automatically from package dependencies
	- [ ] other
- Vendoring
	- [X] Standard
	- [ ] GO15VENDOREXPERIMENT
	- [ ] others
- Package benchmarking
	- [X] go test -bench
	- [ ] go test
	- [ ] go vet
	- [ ] go lint
	- [ ] go fmt
	- [ ] others
- Package dashboard
	- [X] General information
	- [ ] Benchmarking results
	- [ ] Comparing performance on platforms
	- [ ] Subcription
	- [ ] List of similar packages, managed by community
- Platforms
	- [X] Bare metal. Ubuntu 14.04, Intel i5, 4 Core
	- [ ] DigitalOcean
	- [ ] GAE
	- [ ] Virtual machine
	- [ ] Docker in VM
	- [ ] Docker on bare metal
	- [ ] others
- Notifications
	- [ ] Dashboard
	- [ ] Slack
	- [ ] others
- Badge generator
	- [ ] design selector and generator


## License
MIT
