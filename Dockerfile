FROM golang:1.18
COPY main.go main.go
EXPOSE 8080
CMD ["go", "run", "main.go"]
