FROM golang:1.22.3

WORKDIR /go_final_project

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

ENV  CGO_ENABLED=1 GOOS=linux GOARCH=amd64

RUN go build -o /todo_app

COPY scheduler.db ./

CMD ["/todo_app"]