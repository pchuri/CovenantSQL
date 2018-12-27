# Stage: builder
FROM golang:1.11-stretch as builder

ARG BUILD_ARG

WORKDIR /go/src/github.com/CovenantSQL/CovenantSQL
COPY . .
RUN GOOS=linux GOLDFLAGS="-linkmode external -extldflags -static" ./build.sh ${BUILD_ARG}
RUN rm -f bin/*.test

