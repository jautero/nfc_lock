nfc_lock
========

(more) secure electric lock with DESFire EV1 NFC tags


## Go

Make sure your GOPATH is set

    mkdir $HOME/go
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin

Install dependencies

    go get github.com/fuzxxl/nfc/2.0/nfc
    go get github.com/fuzxxl/freefare/0.3/freefare

### RasPi

See https://xivilization.net/~marek/blog/2014/06/10/go-1-dot-2-for-raspberry-pi/ for Go 1.2 (1.0 will not work)

    sudo apt-get install apt-transport-https

Then edit `/etc/apt/sources.list.d/xivilization-raspbian.list` and switch to https

