#!/bin/bash

set -eux

export GOPATH="/tmp/.gobuild"
SRCDIR="${GOPATH}/src/github.com/chonla/cotton"

[ -d ${GOPATH} ] && rm -rf ${GOPATH}
mkdir -p ${GOPATH}/{src,pkg,bin}
mkdir -p ${SRCDIR}
cp tf.go ${SRCDIR}
(
    echo ${GOPATH}
    cd ${SRCDIR}
    go get .
    go install .
)