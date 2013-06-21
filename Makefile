PREFIX=/usr/local

build:
	mkdir -p bin
	go build -o bin/noeqd

clean:
	if [ -d bin ]; then  rm -f bin/* ;fi

apt-builder-fetch:
	# noop
apt-builder-deps:
	add-apt-repository ppa:gophers/go
	apt-get update
	mk-build-deps
apt-builder-build:
	debuild -i -us -uc -b
