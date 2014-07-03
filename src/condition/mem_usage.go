package condition 

 /* 
  MB := 1024 ** 2
  FORMAT_STR := "%d%s"
  MB_LABEL := "MB"
  KB_LABEL := "KB"
*/

  import (
    system "system"
  )

type MemoryUsage struct {
  Below int
}

func NewMemoryUsage(value int) *MemoryUsage {
   c := &MemoryUsage{}
   c.Below = value
   return c
}

func (c *MemoryUsage) Run(pid int) (int, error) { // , include_children bool) {
  return system.MemoryUsage(pid)
}

func (c *MemoryUsage) Check(value int) bool {
  //  value.kilobytes < c.Below
  assert := value < c.Below
  return assert
}

func (c *MemoryUsage) FormatValue(value string) {

/*
    if value.kilobytes >= MB
      FORMAT_STR % [(value / 1024).round, MB_LABEL]
    else
      FORMAT_STR % [value, KB_LABEL]
    end
*/
}
