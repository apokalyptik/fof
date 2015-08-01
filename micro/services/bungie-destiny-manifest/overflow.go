package main

import "strconv"

const MaxUint = ^uint32(0)
const MinUint = 0
const MaxInt = int64(MaxUint >> 1)
const MinInt = -MaxInt - 1

func overflow(i int64) int64 {
	if i < 0 {
		return i
	}
	if i < MaxInt {
		return i
	}
	return (MinInt - (MaxInt - i)) - 1
}

func overflowString(i string) (int64, error) {
	ni, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		return ni, err
	}
	return overflow(ni), nil
}
