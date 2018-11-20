FROM golang:stretch
LABEL maintainer="Alexander Mazuruk <a.mazuruk@samsung.com>"

# SLAV tests
RUN go get -u "github.com/onsi/gomega" &&\
	go get -u "github.com/onsi/ginkgo"
# Boruta tests use testify additionally
RUN	go get -u "github.com/stretchr/testify"
# Mocks are used in tests
RUN	go get -u "github.com/golang/mock/gomock"
# Linters (each repo has specific setttings)
RUN curl -L https://raw.githubusercontent.com/alecthomas/gometalinter/master/scripts/install.sh | sh
# Weles dependency management
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Download SLAV projects
# MuxPi and SLAV repositories do not contain go code in repository root thus
# should be downloaded via git.

RUN git clone https://github.com/SamsungSLAV/muxpi -b master /go/src/github.com/SamsungSLAV/muxpi &&\
    git clone https://github.com/SamsungSLAV/perun -b master /go/src/github.com/SamsungSLAV/perun &&\
    git clone https://github.com/SamsungSLAV/slav -b master /go/src/github.com/SamsungSLAV/slav
# As of writing this Dockerfile, perun has no merged go code in repository thus
# must be downlaoded via git. This should be changed soon.
RUN go get -u github.com/SamsungSLAV/boruta &&\
    go get -u github.com/SamsungSLAV/weles

RUN chmod 777 -R /go/ &&\
	mkdir -p /.cache/go-build &&\
	chmod 777 -R /.cache
