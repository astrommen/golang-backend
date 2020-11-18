FROM golang:buster AS backend
WORKDIR /src/
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/go-back

FROM scratch
COPY --from=build /bin/go-back /bin/go-back
ENTRYPOINT ["/bin/go-back"]