FROM golang:latest

WORKDIR /go/src

COPY ./src /go/src

#RUN go mod init test_effective_mobile
RUN go get -u gorm.io/gorm
RUN go get -u gorm.io/driver/postgres
RUN go get -u github.com/labstack/echo/v4
RUN go get -u github.com/sirupsen/logrus
RUN go mod tidy

EXPOSE 6050