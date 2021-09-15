#!/bin/bash

stop(){
  ps aux|grep uenglish|grep -v grep |awk '{system("kill -9 "$2)}'
  echo "Stop uenglish ..."
}

start(){
  RE=`ps aux|grep uenglish|grep -v grep`
  if [ "$RE" == "" ]; then
    ./uenglish config.properties >> uenglish.log &
    echo "uenglish start ...."
  else
    echo "Already uenglish is running"
  fi

}

restart(){
  stop
  start
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart)
    restart
    ;;
  *)
    echo $"Usage: $0 {start|stop|restart}"
    RETVAL=1
esac
exit $RETVAL
