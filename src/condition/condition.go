package condition


type Condition struct {
  Options map[string]string
  Below bool

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
}