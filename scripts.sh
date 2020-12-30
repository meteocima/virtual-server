#!/bin/bash

RED="\e[31m"
NORMAL="\e[0m"
GREEN="\e[92m"

export all_tests_passed=0

function rebuild_docs() {
  echo > docs/pages/index.md
  echo '{{ useLayout(".layout.njk") }}' >> docs/pages/index.md
  echo '{{ title("CIMA virtual-server") }}' >> docs/pages/index.md
  echo '{{ subtitle("config package") }}' >> docs/pages/index.md
  cat readme.md >> docs/pages/index.md
  for pkg in vpath ctx connection config; do
    godocdown -template docs/layouts/pkg-template.md ./$pkg > docs/pages/$pkg.md
  done
  orgame docs/client docs/pages website
}

function test_all() {
  verbose=$1
  clear
  
  if go test $verbose ./...; then
    printf $GREEN
    figlet 'All tests passed.'
    printf $NORMAL
    
    all_tests_passed=1
  else
    printf $RED
    figlet 'Some test failed.'
    printf $NORMAL
    
    if [[ $all_tests_passed == 1 ]]; then
      spd-say -y Italian+female5 'porco boia'
    fi
    
    all_tests_passed=0
  fi
}

function on_sourcechanges_retest() {
  test_all $1
  rebuild_docs
  while true; do
    event=`inotifywait -qr -e modify -e move -e create -e delete -e delete_self .`
    echo EVENT $event
    previous=$all_tests_passed
    test_all $1
    rebuild_docs
    if [[ $all_tests_passed != $previous ]]; then
      if [[ $all_tests_passed == 1 ]]; then
        spd-say 'All tests passed.'
      else 
        spd-say 'Some test failed.'
      fi
    fi
  done
}