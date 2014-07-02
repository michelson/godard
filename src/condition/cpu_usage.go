package condition 

import (
  system "system"
)

type CpuUsage struct {
  Condition
  //Below string
}

func (c *CpuUsage) Run(pid int , include_children bool) (float64 , error) {
  // rss is on the 5th col
  return system.CpuUsage(pid) //, include_children)
}

/*
      def initialize(options = {})
        @below = options[:below]
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