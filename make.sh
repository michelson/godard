
set -e

cd "`dirname '$0'`"
SCRIPTPATH="`pwd`"
cd - > /dev/null

export GOPATH=$SCRIPTPATH
export GOBIN=

function deps {
  echo "Fetching dependencies to $SCRIPTPATH..."
  printf "###         (01/03)\r"
    go get -u -t github.com/segmentio/go-log
  printf "#####       (02/03)\r"
    go get -u -t github.com/looplab/fsm
  printf "########    (03/03)\r"
    go get -u -t bitbucket.org/kardianos/osext
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

$1
