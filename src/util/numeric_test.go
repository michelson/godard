package util

import( 
  "testing"
)


func TestParseNumber(t *testing.T)  {
  v, _ := ParseNumber("1.kb")
  if v != 1024 {
      t.Error("Expected 1024, got ", v)
  }

  v2, _ := ParseNumber("3.megabytes")
  if v2 != 3145728 {
      t.Error("Expected 3145728, got ", v2)
  }

  v3, err3 := ParseNumber("3.5.gigabytes")
  if v3 != 3758096384 {
      t.Error("Expected 3758096384, got ", v3, err3)
  }

  v4, _ := ParseNumber("-4.exabytes")
  if v4 != -4611686018427387904 {
      t.Error("Expected -4611686018427387904, got ", v4)
  }

  /*
  2.kilobytes   # => 2048
  3.megabytes   # => 3145728
  3.5.gigabytes # => 3758096384
  -4.exabytes   # => -4611686018427387904
  */
}
