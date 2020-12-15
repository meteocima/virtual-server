#!/bin/bash

RED="\e[31m"
STOP="\e[0m"
GREEN="\e[92m"
export status=1
while inotifywait -r -e modify -e move -e create -e delete -e delete_self .; do
  if go test ./...; then
    printf $GREEN
    figlet evviva
    printf $STOP
    if [[ $status == 0 ]]; then
      spd-say -y Italian+female5 'evviva'
    fi
    status=1
  else
    printf $RED
    figlet rotto
    printf $STOP
    
    if [[ $status == 1 ]]; then
      spd-say -y Italian+female5 'porco boia'
    fi
    
    status=0
  fi
done