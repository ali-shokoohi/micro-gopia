FROM golang:1.20.5

WORKDIR /app

COPY . ./

# Getting modules
RUN go mod download

# Run the tests in the container
# FROM build-stage AS run-test-stage
RUN CGO_ENABLED=1 go test -v ./...

# Building of project
RUN CGO_ENABLED=0 GOOS=linux go build -o /micro-gopia /app/cmd/app/main.go 

CMD ["/bin/sh","-c","/micro-gopia migrate && /micro-gopia"]