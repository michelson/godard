package process

import (
  "time"
  util "util"
  "strings"
)



type ProcessStatistics struct {
  Events *util.RotationalArray
}

//const Strftime = "%m/%d/%Y %H:%I:%S"
const EventsToPersists = 10


func NewProcessStatistics() *ProcessStatistics {
  c := &ProcessStatistics{}
  c.Events = util.NewRotationalArray(EventsToPersists)
  return c
}

func (c*ProcessStatistics) RecordEvent(event string , reason string){
  //events.push([event, reason, Time.now])
  arr := []string{event , reason, string(time.Now().Unix())}
  c.Events.Push(arr)
}

func (c*ProcessStatistics) ToS() string{

    var str string
    for i := len(c.Events.Array) - 1; i >= 0; i-- {
      str += strings.Join( []string{ c.Events.Array[i][0] , c.Events.Array[i][1], c.Events.Array[i][2]}, "-")
    }
    return str
}
