FROM golang:1.22.3-alpine as build
RUN mkdir /src
ADD ./golang-echoer-server/*.go /src
ADD ./golang-echoer-server/go.mod /src
ADD ./golang-echoer-server/go.sum /src
WORKDIR /src
RUN go get -d -v -t
RUN GOOS=linux go build -v -o golang-echoer-server
RUN chmod +x golang-echoer-server

FROM scratch
COPY --from=build /src/golang-echoer-server /usr/local/bin/golang-echoer-server
EXPOSE 8080
CMD ["golang-echoer-server"]