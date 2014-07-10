package util

type RotationalArray struct {
  Capacity int
  Array [][]string

}

func NewRotationalArray(size int) *RotationalArray {
  c := &RotationalArray{}
  c.Capacity = size
  c.Array = make([][]string, size)
  return c
}

func (c*RotationalArray) Push(value []string) {
  c.Array = append(c.Array, value)
  if len(c.Array)+1 > c.Capacity {
    c.Array = c.Array[1 : c.Capacity+1]    
  }

  /*

      def push(value)
        super(value)

        self.shift if self.length > @capacity
        self
      end
  */
}

/*
# -*- encoding: utf-8 -*-
module Bluepill
  module Util
    class RotationalArray < Array
      def initialize(size)
        @capacity = size

        super() # no size - intentionally
      end

      def push(value)
        super(value)

        self.shift if self.length > @capacity
        self
      end
      alias_method :<<, :push
    end
  end
end
*/