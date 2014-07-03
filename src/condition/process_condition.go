package condition


type ProcessCondition interface {
    Run(pid int, include_children bool) (float64 , error)
    Check(value float64, include_children bool) (bool, error)
    FormatValue(value float64) string
}


  /*
    def run(pid, include_children)
      raise "Implement in subclass!"
    end

    def check(value)
      raise "Implement in subclass!"
    end

    def format_value(value)
      value
    end
  */
