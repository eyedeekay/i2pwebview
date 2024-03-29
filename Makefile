
#GO111MODULE=off
#export GO111MODULE=off

i2pd_prerelease_version=c-wrapper-libi2pd-api
i2pd_release_version=2.40.0

export GOPATH=$(HOME)/go

export USE_STATIC=yes
USE_STATIC=yes

export LDFLAGS=-static
LDFLAGS=-static

GXXFLAGS=-static
export GXXFLAGS=-static

CXXFLAGS=-static
export CXXFLAGS=-static

CGO_GXXFLAGS=-static
export CGO_GXXFLAGS=-static

CGO_CFLAGS=-static
export CGO_CFLAGS=-static

CGO_CXXFLAGS=-static
export CGO_CXXFLAGS=-static

CGO_CPPFLAGS=-static
export CGO_CPPFLAGS=-static

#CGO_LDFLAGS=-static
#export CGO_LDFLAGS=-static


#Trying to achieve fully-static builds, this doesn't work yet.
FLAGS= /usr/lib/x86_64-linux-gnu/libboost_system.a /usr/lib/x86_64-linux-gnu/libboost_date_time.a /usr/lib/x86_64-linux-gnu/libboost_filesystem.a /usr/lib/x86_64-linux-gnu/libboost_program_options.a /usr/lib/x86_64-linux-gnu/libssl.a /usr/lib/x86_64-linux-gnu/libcrypto.a /usr/lib/x86_64-linux-gnu/libz.a

example:
	go build -o i2psurf -x -v --tags=netgo \
		-ldflags '-w -linkmode=external -extldflags "-ldl $(FLAGS)"' ./surfi2p