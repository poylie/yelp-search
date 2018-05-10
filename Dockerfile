FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/poylie/yelp-search

ENV USER SKYNET\poylie
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET Z3plbbWvm3N1NZtW

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://SKYNET\poylie@localhost:5432/yelp-search?sslmode=disable

WORKDIR /go/src/github.com/poylie/yelp-search

RUN godep go build

EXPOSE 8888
CMD ./yelp-search