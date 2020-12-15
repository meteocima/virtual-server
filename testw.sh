#!/bin/bash

RED="\e[31m"
STOP="\e[0m"
GREEN="\e[92m"

while inotifywait -r -e modify -e move -e create -e delete -e delete_self .; do
  if go test ./...; then
    printf $GREEN
    figlet evviva
    printf $STOP
    spd-say -y Italian+female5 'evviva'
  else
    printf $RED
    figlet è rotto
    printf $STOP
    spd-say -y Italian+female5 'porco boia'
  fi
done