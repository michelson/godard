package condition 

import (
  system "system"
)

type CpuUsage struct {
  Below float64
}

func NewCpuUsage(options map[string]interface{}) *CpuUsage{
  c := &CpuUsage{Below: 5}
  return c
}

func (c *CpuUsage) Run(pid int , include_children bool) (float64 , error) {
  return system.CpuUsage(pid) //, include_children)
}

func (c *CpuUsage) Check(value float64 , include_children bool) (bool , error) {
  assert := value < c.Below
  return assert , nil
}

func (c *CpuUsage) FormatValue(value float64) string{
   return "oli"
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