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
  Condition
}

func NewMemoryUsage(options string) *MemoryUsage {
   //below = options //options[:below]
   //Condition.Below = options //options[:below]
   c := &MemoryUsage
   c.below := options
   return c
}

func (c *MemoryUsage) Run(pid int) int { // , include_children bool) {
  return system.MemoryUsage(pid)
}

func (c *MemoryUsage) Check(value string) {
  // rss is on the 5th col
  //  value.kilobytes < @below
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
