FROM golang:alpine AS build

RUN apk --no-cache add build-base
RUN apk --no-cache add ca-certificates

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG version
RUN echo "$version" > ./internal/utils/VERSION

ENV CGO_ENABLED=0
RUN cd /build/cmd/bot && go build -a -o /build/countula

# Create the output container from the built image.
FROM scratch
COPY --from=build /build/countula /countula
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "/countula" ]