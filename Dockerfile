FROM golang:1.12-alpine AS build

# git is required, keeping small image with -no-cache
RUN apk add --no-cache git

WORKDIR /anargu/miauth/
COPY . /anargu/miauth/
RUN go mod tidy
RUN go build -o /bin/miauth cmd/main.go

# This results in a single layer image
FROM alpine:3.9.4
# adding ca certificates
RUN apk --update add ca-certificates
COPY --from=build /bin/miauth /bin/miauth
# copying assets dir
COPY --from=build /anargu/miauth/public /public

ENTRYPOINT ["/bin/miauth"]