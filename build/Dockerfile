FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . ./

RUN go mod download
RUN CGO_ENABLED=0 go build -o /usr/bin/imagedownloader ./cmd/imagedownloader

############################### SECOND STAGE ###############################

FROM alpine

COPY --from=builder /usr/bin/imagedownloader /usr/bin/imagedownloader
COPY --from=builder /app/fixtures /fixtures
RUN mkdir /downloads

CMD ["imagedownloader", "start"]
