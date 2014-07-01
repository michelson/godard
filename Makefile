
default:	build

build:
	./make.sh build

run:
	./make.sh run

deps:
	./make.sh deps

test:
	ls ./src | grep -v "\." | sed 's/\///g' | xargs go test -cover

convey:
	./bin/goconvey --depth=2

	
