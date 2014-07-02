package watcher

import(
  "log"
  condition "condition"
)

type HistoryValue struct {
  value string
  critical string
}

type ConditionWatch struct {
  Logger string
  Name string
  Fires string
  Every string
  Times []float64
  empty_array []interface{}
  include_children bool
  ProcessCondition *condition.Condition
}

func NewConditionWatch(name string, options interface{}) *ConditionWatch{

  /*
      @name = name

      @logger = options.delete(:logger)
      @fires  = options.has_key?(:fires) ? Array(options.delete(:fires)) : [:restart]
      @every  = options.delete(:every)
      @times  = options.delete(:times) || [1,1]
      @times  = [@times, @times] unless @times.is_a?(Array) # handles :times => 5
      @include_children = options.delete(:include_children) || false

      self.clear_history!

      @process_condition = ProcessConditions[@name].new(options)
  */

      v := options.(map[string]interface{})
      log.Println("CREATING CONDITION", v["every"])
      c := &ConditionWatch{}
      
      if _,ok := v["fires"]; ok {
        c.Fires = v["fires"].(string)  
      }
      if _,ok := v["every"]; ok {
        c.Every = v["every"].(string)
      }
      if _,ok := v["times"]; ok {
        arr := make([]float64, 2)
        arr[0] = v["times"].(float64)
        arr[1] = v["times"].(float64)
        c.Times = arr
      }
      c.Name  = name


      return c
}

func (*ConditionWatch) Run(pid string, tick_number int64) {
/*
    def run(pid, tick_number = Time.now.to_i)
      if @last_ran_at.nil? || (@last_ran_at + @every) <= tick_number
        @last_ran_at = tick_number

        value = @process_condition.run(pid, @include_children)
        @history << HistoryValue.new(@process_condition.format_value(value), @process_condition.check(value))
        self.logger.info(self.to_s)

        return @fires if self.fired?
      end
      EMPTY_ARRAY
    end
    */
}

func (*ConditionWatch) ClearHistory() {

  // @history = Util::RotationalArray.new(@times.last)
  
}

func (*ConditionWatch) Fired() {
  // @history.count {|v| not v.critical} >= @times.first
}

func (*ConditionWatch) to_s() {
  // data = @history.collect {|v| "#{v.value}#{'*' unless v.critical}"}.join(", ")
  // "#{@name}: [#{data}]\n"
}
