FROM ubuntu:18.04 
RUN apt-get update
RUN apt-get install golang-go nano -y
ADD . /webApp
WORKDIR /webApp
CMD go run main.go
EXPOSE 3000
