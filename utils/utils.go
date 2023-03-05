// Package utils provides utility functions and types
package utils

// Interval defines a range with start and end values
type Interval struct {
	Start, End float64
}

// Sat returns the nearest value to val that is within the interval range
func (i *Interval) Sat(val float64) float64 {
	// Get the minimum and maximum values of the interval
	min, max := i.GetMinMax()

	// Return the maximum value if val is greater than the maximum
	if val > max {
		return max
	}

	// Return the minimum value if val is less than the minimum
	if val < min {
		return min
	}

	// Otherwise, return val as it is within the interval range
	return val
}

// GetMinMax returns the minimum and maximum values of the interval
func (i *Interval) GetMinMax() (min, max float64) {
	min = i.Start
	max = i.End

	// Swap min and max if min is greater than max
	if min > max {
		return max, min
	}

	return
}

// Length returns the length of the interval
func (i *Interval) Length() float64 {
	return i.End - i.Start
}

// IntervalTransformer maps values from one interval range to another
type IntervalTransformer struct {
	IntervalFrom Interval
	IntervalTo   Interval
}

// TransformForward maps a value from the "from" interval range to the "to" interval range
func (t *IntervalTransformer) TransformForward(val float64) float64 {
	return Transform(val, &t.IntervalFrom, &t.IntervalTo)
}

// TransformBackward maps a value from the "to" interval range to the "from" interval range
func (t *IntervalTransformer) TransformBackward(val float64) float64 {
	return Transform(val, &t.IntervalTo, &t.IntervalFrom)
}

// Transform maps a value from one interval range to another
func Transform(val float64, iFrom, iTo *Interval) float64 {
	// Calculate the normalized value within the "from" interval range
	normVal := (val - iFrom.Start) / iFrom.Length()

	// Scale the normalized value to the "to" interval range and shift it to match the "to" interval start value
	return iTo.Sat(normVal*iTo.Length() + iTo.Start)
}
