FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

ENV elara_ml_access_token=ml_access_token
ENV elara_ml_url=vertex_ai_endpoint

RUN go build -o /elara-api


CMD [ "/elara-api" ]