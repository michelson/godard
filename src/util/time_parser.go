package util

import (
	"strconv"
	"strings"
	"time"
	//"fmt"
	"errors"
)

func TimeParse(input string) (time.Duration, error) {

	count := strings.Count(input, ".")
	//fmt.Println(count)
	var parts []string

	if count == 1 {
		parts = strings.Split(input, ".")
	} else if count == 2 {
		index := strings.LastIndex(input, ".")
		parts = append(parts, input[:index])
		parts = append(parts, input[index+1:])
	}

	var tt time.Duration

	if len(parts) != 2 {
		return tt, errors.New("can't split values")
	}

	m, err := strconv.ParseFloat(parts[0], 64)

	if err != nil {
		return tt, errors.New("can't convert multipler to integer")
	}

	var multiplier float64
	multiplier = float64(m)

	switch parts[1] {
	case "seconds", "second", "secs", "s":
		tt = time.Second * time.Duration(multiplier)
	case "minutes", "minute", "min", "mins":
		tt = time.Minute * time.Duration(multiplier)
	case "hours", "hour", "h", "hrs", "hr":
		tt = time.Hour * time.Duration(multiplier)
	default:
		return tt, errors.New("can't find duration type")
	}

	return tt, nil
}

/*func main(){
  fmt.Println(TimeParse("1.seconds"))
  fmt.Println(TimeParse("1.hours"))
  fmt.Println(TimeParse("1.minutes"))
  fmt.Println(TimeParse("ccc.minutes"))
  fmt.Println(TimeParse("12.sparks"))
  fmt.Println(TimeParse("1,minutes"))
}*/

////https://gist.github.com/michelson/708c45269381040c441c
