package util

import (
  "strings"
  "strconv"
  "errors"
  //"log"
)


const Kilobyte float64 =  1024
const Megabyte float64 =  Kilobyte * 1024
const Gigabyte float64 =  Megabyte * 1024
const Terabyte float64 =  Gigabyte * 1024
const Petabyte float64 =  Terabyte * 1024
const Exabyte  float64 =  Petabyte * 1024

func ParseNumber(input string) (float64, error) {

  count := strings.Count(input, ".")
  //fmt.Println(count)
  var parts []string 
  
  if count == 1 {
    parts = strings.Split(input, ".") 
  }else if count == 2{
    index := strings.LastIndex(input , ".")
    parts = append(parts, input[:index])
    parts = append(parts, input[index+1:])
  }

  var tt float64

  if len(parts) != 2 {
     return tt, errors.New("can't split values")    
  }

  //fmt.Println(parts[0])
  m, err := strconv.ParseFloat(parts[0], 64)

  if err != nil {
     return tt, errors.New("can't convert multipler to integer")    
  }

  var multiplier float64
  multiplier = float64(m)
 
  switch parts[1] {
    case "kilobyte", "kb", "Kb":
      tt = multiplier * Kilobyte
    case "megabyte", "megabytes", "mb":
      tt = multiplier * Megabyte
    case "gigabyte", "gigabytes", "gb":
      tt = multiplier * Gigabyte
    case "terabyte", "terabytes", "tb":
      tt = multiplier * Terabyte
    case "petabyte", "petabytes", "pb":
      tt = multiplier * Petabyte
    case "exabyte", "exabytes","ex":
      tt = multiplier * Exabyte
    default:
      return tt, errors.New("can't find byte type")    
  } 

  return tt, nil
}