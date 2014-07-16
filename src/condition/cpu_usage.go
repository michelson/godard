package condition 

import (
  system "system"
  "strconv"
  "log"
)

type CpuUsage struct {
  Below float64
}

func NewCpuUsage(options map[string]interface{}) *CpuUsage{
  var below float64
  below = float64(options["below"].(float64))
  c := &CpuUsage{Below: below}
  log.Println("CREATING PROCESS CONDITION BELOW", c.Below)

  return c
}

func (c *CpuUsage) Run(pid int , include_children bool) (float64 , error) {
  val, err := system.CpuUsage(pid) //, include_children)
  log.Println("CPU USAGE:", val )
  return val, err
}

func (c *CpuUsage) Check(value float64 , include_children bool) (bool , error) {
  assert := value < c.Below
  return assert , nil
}

func (c *CpuUsage) FormatValue(value float64) string{
   var int_val int
   int_val = int(value)
   var str string
   str = strconv.Itoa(int_val) 
   return str
}

/*
      def initialize(options = {})
        @Below = options[:below]
      end

      def run(pid, include_children)
        # third col in the ps axu output
        System.cpu_usage(pid, include_children).to_f
      end

      def check(value)
        value < @below
      end
    end
*/