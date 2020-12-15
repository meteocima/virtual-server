#!/bin/bash

RED="\e[31m"
NORMAL="\e[0m"
GREEN="\e[92m"

export all_tests_passed=0

function test_all() {
  verbose=$1
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
  while inotifywait -r -e modify -e move -e create -e delete -e delete_self .; do
    previous=$all_tests_passed
    test_all
    if [[ $all_tests_passed != $previous ]]; then
      if [[ $all_tests_passed == 1 ]]; then
        spd-say -y Italian+female5 'Tutti i test sono ok'
      else 
        spd-say -y Italian+female5 'Test falliti!'
      fi
    fi
  done
}