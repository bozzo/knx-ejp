#
#    Copyright 2021 Boris Barnier <bozzo@users.noreply.github.com>
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
#
FROM golang:1.18.1-alpine AS builder

# Add Maintainer Info
LABEL maintainer="Boris Barnier <bozzo@users.noreply.github.com>"

RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY *.go ./

# Build the Go app
RUN CGO_ENABLED=0 go build -o knx-ejp .

FROM scratch

COPY --from=builder /app/knx-ejp .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose port 3671 to the outside world
EXPOSE 3671
EXPOSE 13671

ENV MULTICAST_ADDRESS "224.0.23.12"
ENV MULTICAST_PORT "3671"
ENV PROMETHEUS_PORT "13671"

ENV LOG_LEVEL "info"
# ENV LOG_LEVEL "debug"

# ENV LOG_FORMAT "json"
# ENV LOG_FORMAT "text"

# Command to run the executable
ENTRYPOINT ["./knx-ejp"]
