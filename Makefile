
default:	build

build:
	./make.sh build

run:
	./make.sh run

deps:
	./make.sh deps

test:
	./make.sh test

format:
	./make.sh format

install:
	./make.sh deps ; ./make.sh build  ; cp ./godard /usr/local/bin/

convey:
	./bin/goconvey --depth=2


