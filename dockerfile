FROM golang:alpine AS build

RUN apk --no-cache add build-base
RUN apk --no-cache add ca-certificates

WORKDIR /build
COPY . .
RUN go mod download
ENV CGO_ENABLED=1
RUN cd /build/cmd/bot && go build -ldflags='-s -w' -trimpath -a -o /build/countula
RUN ldd /build/countula | tr -s [:blank:] '\n' | grep '^/' | xargs -I {} install -D {} /build/{}

ARG version
RUN echo "$version" > /etc/program-version
RUN chattr +i /etc/program-version

# Create the output container from the built image.
FROM scratch
COPY --from=build /build/ /
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs
COPY --from=build /etc/program-version /etc/program-version

ENTRYPOINT [ "/countula" ]