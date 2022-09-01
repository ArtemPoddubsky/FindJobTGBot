FROM golang:1.18 AS builder
COPY . /Bot
WORKDIR /Bot
RUN make all

FROM golang:1.18
COPY --from=builder /Bot ./
CMD ./bot