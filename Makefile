
default:	build

build:
	./make.sh build

run:
	./make.sh run

deps:
	./make.sh deps

test:
	./make.sh test

convey:
	./bin/goconvey --depth=2

	
