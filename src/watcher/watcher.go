package watcher

import(
  "log"
  condition "condition"
  //util "util"
)

type HistoryValue struct {
  Value string
  Critical bool
}

type ConditionWatch struct {
  Logger string
  Name string
  Fires []string
  Every float64
  Times []float64
  empty_array []interface{}
  LastRanAt float64
  include_children bool
  ProcessCondition []condition.ProcessCondition //*condition.Condition
  History []*HistoryValue
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
      log.Println("CREATING", name ,"CONDITION", v["every"])
      c := &ConditionWatch{}
      
      if _,ok := v["fires"]; ok {
        c.Fires = append( c.Fires, v["fires"].(string) ) 
      }
      if _,ok := v["every"]; ok {
        c.Every = v["every"].(float64)
      }
      if _,ok := v["times"]; ok {
        arr := make([]float64, 2)
        arr[0] = v["times"].(float64)
        arr[1] = v["times"].(float64)
        c.Times = arr
      }
      c.Name  = name

      /*@include_children = options.delete(:include_children) || false

      self.clear_history!

      @process_condition = ProcessConditions[@name].new(options)*/
      //process_condition := process_condition[c.Name]condition.ProcessCondition{}
      
      log.Println("WATCH", c.Name)

      conditions := make([]condition.ProcessCondition, 0 )
      
      switch c.Name {
        case "mem_usage":
            conditions = []condition.ProcessCondition{ condition.NewMemoryUsage( v ) }
        case "cpu_usage":
            conditions = []condition.ProcessCondition{ condition.NewCpuUsage( v ) }
        case "file_time":
            //conditions = []condition.ProcessCondition{ condition.NewFileTimeUsage( v ) }
        case "running_time":
            //conditions = []condition.ProcessCondition{ condition.NewCpuUsage( v ) }
      }

      /*for _, cond := range conditions {
        log.Println(cond.Check(100, false))
      }*/

      c.ProcessCondition = conditions

      c.ClearHistory()
      c.LastRanAt = 0
      return c
}

func (c*ConditionWatch) Run(pid int, tick_number float64) []string {
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

    log.Println("RUNNING CONDITION EVERY", c.Every) 

    fires := make([]string, 0)

    if c.LastRanAt == 0 || (c.LastRanAt + c.Every) <= tick_number {
      
      c.LastRanAt = tick_number

      var value float64
      var formatted string
      var checked bool
      
      //for _, cond := range c.ProcessCondition {
        value, _   = c.ProcessCondition[0].Run(pid, false) 
        formatted  = c.ProcessCondition[0].FormatValue(value)   
        checked, _ = c.ProcessCondition[0].Check(value, false)   
      //}

      c.PushHistory( &HistoryValue{Value: formatted, Critical: checked} )
      
      if c.isFired(){
        fires = c.Fires 
      }        
    }

    return fires

}

func (c*ConditionWatch) ClearHistory() {
  // @history = Util::RotationalArray.new(@times.last)
  var capacity = int(c.Times[1])
  arr := make([]*HistoryValue , capacity)
  c.History = arr
}
//extracted from utils rotational arr
func (c*ConditionWatch) PushHistory(value *HistoryValue) {
  c.History = append(c.History, value)
  var capacity = int(c.Times[1])
  if len(c.History)+1 > capacity {
    c.History = c.History[1 : capacity+1]    
  }
}

func (c*ConditionWatch) isFired() bool {
  // @history.count {|v| not v.critical} >= @times.first
  var count float64
  for _, h := range(c.History) {
    if h != nil {
      log.Println("HISTORY:" , "val:", h.Value , "critical", h.Critical)
      if !h.Critical {
        count = count + 1
      }
    }
  }
  assert := count >= c.Times[1]
  return assert
}

func (c*ConditionWatch) ToS() string {
  /* data = @history.collect {|v|  "#{v.value}#{'*' unless v.critical}"}.join(", ")
   "#{@name}: [#{data}]\n"
  */
  var data_arr []string
  //for _, h := range(c.History) {
  for i := 0; i < len(c.History); i++ {
    str := c.History[i].Value 
    if !c.History[i].Critical {
      str += "*"
    }
    data_arr[i] = str
  }
  var str string 
  for _ , r := range(data_arr){
    str = str + ", " + r
  }
  return str
}
