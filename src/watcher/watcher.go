package watcher

import (
	condition "condition"
	"log"
	"strings"
	"time"
	"util"
)

type HistoryValue struct {
	Value    string
	Critical bool
}

type ConditionWatch struct {
	Logger           *log.Logger
	Name             string
	Fires            []string
	Every            time.Duration
	Times            []float64
	empty_array      []interface{}
	LastRanAt        float64
	include_children bool
	ProcessCondition []condition.ProcessCondition
	History          []*HistoryValue
}

func NewConditionWatch(name string, options interface{}) *ConditionWatch {

	
	v := options.(map[string]interface{})
	//c.Logger.Println("CREATING", name ,"CONDITION EVERY", v["every"])
	c := &ConditionWatch{}
	
	c.Logger = v["logger"].(*log.Logger)

	if _, ok := v["fires"]; ok {
		c.Fires = append(c.Fires, v["fires"].(string))
	} else {
		c.Fires = append(c.Fires, "restart")
	}
	if _, ok := v["every"]; ok {
		c.Every, _ = util.TimeParse(v["every"].(string))
	}
	if _, ok := v["times"]; ok {
		arr := make([]float64, 2)
		arr[0] = v["times"].(float64)
		arr[1] = v["times"].(float64)
		c.Times = arr
	}
	
	c.Name = name

	/*@include_children = options.delete(:include_children) || false
	  self.clear_history!
	*/

	//c.Logger.Println("WATCH", c.Name)

	conditions := make([]condition.ProcessCondition, 0)

	switch c.Name {
	case "mem_usage":
		conditions = []condition.ProcessCondition{condition.NewMemoryUsage(v)}
	case "cpu_usage":
		conditions = []condition.ProcessCondition{condition.NewCpuUsage(v)}
	case "file_time":
		//conditions = []condition.ProcessCondition{ condition.NewFileTimeUsage( v ) }
	case "running_time":
		//conditions = []condition.ProcessCondition{ condition.NewCpuUsage( v ) }
	}

	c.ProcessCondition = conditions

	c.ClearHistory()
	c.LastRanAt = 0
	return c
}

func (c *ConditionWatch) Run(pid int, tick_number float64) []string {

	fires := make([]string, 0)

	if c.LastRanAt == 0 || (c.LastRanAt+c.Every.Seconds()) <= tick_number {
		c.Logger.Println("TIME DURATION", (c.LastRanAt + c.Every.Seconds()), "VS", tick_number)
		c.LastRanAt = tick_number

		var value float64
		var formatted string
		var checked bool

		value, _ = c.ProcessCondition[0].Run(pid, false)
		formatted = c.ProcessCondition[0].FormatValue(value)
		checked, _ = c.ProcessCondition[0].Check(value, false)

		//c.Logger.Println("VAL", formatted , "CRITIC", checked)
		c.PushHistory(&HistoryValue{Value: formatted, Critical: checked})

		if c.isFired() {
			fires = c.Fires
		}
	}

	return fires

}

func (c *ConditionWatch) ClearHistory() {
	var capacity = int(c.Times[1])
	arr := make([]*HistoryValue, capacity)
	c.History = arr
}

//extracted from utils rotational arr
func (c *ConditionWatch) PushHistory(value *HistoryValue) {
	c.History = append(c.History, value)
	var capacity = int(c.Times[1])
	if len(c.History)+1 > capacity {
		c.History = c.History[1 : capacity+1]
	}
}

func (c *ConditionWatch) isFired() bool {

	var count float64
	for _, h := range c.History {
		if h != nil {
			if !h.Critical {
				count += 1
			}
			//c.Logger.Println("val:", h.Value , "critical", h.Critical, "times:", c.Times[1])
		}
	}

	//c.Logger.Println("HISTORY:", count)
	assert := count >= c.Times[1]

	return assert
}

func (c *ConditionWatch) ToS() string {
	data_arr := make([]string, 0)
	for i, h := range c.History {
		if !c.History[i].Critical {
			data_arr = append(data_arr, h.Value+"*")
		}
	}
	str := strings.Join(data_arr, ", ")
	return str
}
