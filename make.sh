
set -e

cd "`dirname '$0'`"
SCRIPTPATH="`pwd`"
cd - > /dev/null

export GOPATH=$SCRIPTPATH
export GOBIN=

function deps {
  echo "Fetching dependencies to $SCRIPTPATH..."
  printf "###         (01/04)\r"
    go get -u -t github.com/segmentio/go-log
  printf "#####       (02/04)\r"
    go get -u -t github.com/looplab/fsm
  printf "########    (03/04)\r"
    go get -u -t bitbucket.org/kardianos/osext
  printf "############(04/04)\r"
    go get -u -t code.google.com/p/go.tools/cmd/cover
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

$1
