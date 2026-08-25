FROM golang:1.23 AS build
WORKDIR /src
ENV GOPROXY=off GOSUMDB=off
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
COPY web ./web
RUN go build -mod=vendor -o /out/venueops ./cmd/venueops

FROM golang:1.23
WORKDIR /app
ENV GOPROXY=off GOSUMDB=off
COPY --from=build /src/go.mod ./go.mod
COPY --from=build /src/go.sum ./go.sum
COPY --from=build /src/vendor ./vendor
COPY --from=build /src/cmd ./cmd
COPY --from=build /src/internal ./internal
COPY --from=build /src/web ./web
COPY --from=build /out/venueops ./venueops
EXPOSE 8080
CMD ["./venueops", "-addr", "0.0.0.0:8080", "-data", "/app/data", "-web", "/app/web", "-seed"]
