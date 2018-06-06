FROM golang:1.8

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# COPY ./app/. /go/src/app/

EXPOSE 8080

RUN go get github.com/labstack/echo
RUN go get github.com/jinzhu/gorm
RUN go get github.com/jinzhu/now
RUN go get github.com/pilu/fresh
RUN go get github.com/joho/godotenv
RUN go get github.com/go-ozzo/ozzo-validation

# CMD go run /go/src/app/main.go
# CMD fresh -c runner.conf
