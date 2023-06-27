# syntax=docker/dockerfile:1

##
## Build the application from source
##
# https://docs.docker.com/language/golang/build-images/

FROM golang:1.19 AS build-stage

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /gin-html-templates

##
## Deploy the application binary into a lean image
##

FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /gin-html-templates /gin-html-templates
COPY ./templates ./templates
COPY ./assets ./assets
EXPOSE 8082

ENTRYPOINT ["/gin-html-templates"]
