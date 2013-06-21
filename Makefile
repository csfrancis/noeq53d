PREFIX=/usr/local

build:
	mkdir -p bin
	go build -o bin/noeq53d
clean:
	if [ -d bin ]; then  rm -f bin/* ;fi
