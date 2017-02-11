FROM resin/rpi-raspbian

RUN apt-get update && apt-get install -y golang \
  && apt-get install -y git \
  && apt-get install -y gcc \
  && apt-get install libsdl1.2-dev \
  && apt-get install build-essential
  RUN export GOPATH=$HOME/go \
   export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
  RUN cd /go/src/goHomeServer \
    go get github.com/jacobsa/go-serial/serial \
    go get github.com/galaktor/gorf24 \
  RUN cd RF24/RPi/RF24 \
      make \
      make install \
