package condition 

 /* MB := 1024 ** 2
  FORMAT_STR := "%d%s"
  MB_LABEL := "MB"
  KB_LABEL := "KB"

*/

import(
  system "system"
)

type RunningTime struct {
  Condition
}


func (c *RunningTime) Run(pid int) (int, error) { // , include_children bool) {
  return system.RunningTime(pid)
}

func (c *RunningTime) Check(value string) {
  // rss is on the 5th col
  // value.kilobytes < @below
}
