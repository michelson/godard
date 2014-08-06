
set -e

cd "`dirname '$0'`"
SCRIPTPATH="`pwd`"
cd - > /dev/null

export GOPATH=$SCRIPTPATH
export GOBIN=

function deps {
  echo "Fetching dependencies to $SCRIPTPATH..."
  printf "#####       (02/06)\r"
    go get -u -t github.com/looplab/fsm
  printf "#######     (03/06)\r"
    go get -u -t bitbucket.org/kardianos/osext
  printf "##########  (04/06)\r"
    go get -u -t code.google.com/p/go.tools/cmd/cover
  printf "############(06/06)\r"
    go get -u -t github.com/davecheney/profile
  printf "\n"
}

function build {
  go build godard
}

#function run {
#  go build godard ; ./godard
#}

function run {
  go build godard ; ./godard load --config=./godard.cfg
}

function test {
  ls ./src | grep -v "\." | sed 's/\///g' | xargs go test -cover
}

function format {
  ls ./src | grep -v "\." | sed 's/\///g' | xargs go fmt
}

$1
