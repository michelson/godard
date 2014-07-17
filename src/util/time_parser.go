package util

import(
  "time"
  "strings" 
  "strconv"
  //"fmt"
  "errors"
)

func TimeParse(input string) (time.Duration, error) {
  
  parts := strings.Split(input, ".") 
  
  var tt time.Duration

  if len(parts) != 2 {
     return tt, errors.New("can't split values")    
  }

  multiplier, err := strconv.Atoi(parts[0])
  
  if err != nil {
     return tt, errors.New("can't convert multipler to integer")    
  }

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

