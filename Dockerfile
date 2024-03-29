# build stage
FROM golang:alpine as build-env

ARG src_name="yic_email"

WORKDIR ${GOPATH}/src/github.com/fredericalix/${src_name}
ADD . .
RUN CGO_ENABLED=0 go build -o /${src_name}

# run stage
FROM alpine

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=build-env /${src_name} .

CMD [ "./email"]
