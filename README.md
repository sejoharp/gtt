# introduction
timetracker with golang

# features

# requirements
* mongodb
* go

# installation
* clone this repo
* change to the repo
* local build `go build`
* optional: installation to bin directory `go install`

# status
travis ci: [![Build Status](https://travis-ci.org/zippelmann/gtt.svg?branch=master)](https://travis-ci.org/zippelmann/gtt)

coveralls: [![Coverage Status](https://coveralls.io/repos/zippelmann/gtt/badge.svg)](https://coveralls.io/r/zippelmann/gtt)

# useful commands
* automatic testing: `ginkgo watch -r --randomizeAllSpecs --trace --race --failFast --compilers=2 --cover --notify`

# used projects - thanks for the great tools
* [ginkgo](onsi.github.io/ginkgo/)
* [gomega](onsi.github.io/gomega/)
* [goji](http://goji.io)
* [mgo](https://labix.org/mgo)
* [gocov](https://github.com/axw/gocov)
* [gover](https://github.com/modocache/gover)
* [goveralls](https://github.com/mattn/goveralls)
* [jwt-go](github.com/dgrijalva/jwt-go)

# license
MIT license
