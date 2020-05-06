FROM golang:1.14 AS builder

ENV GO111MODULE=on
WORKDIR /go/build/app/src
COPY . .
# RUN go get -d -v ./...
RUN go build -o /bin/aws-spot-instance-helper

FROM ubuntu
RUN apt-get update && apt-get install -y ca-certificates
EXPOSE 9777
COPY --from=builder /bin/aws-spot-instance-helper /bin/aws-spot-instance-helper
ENTRYPOINT [ "/bin/aws-spot-instance-helper"]