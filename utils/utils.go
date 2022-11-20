package utils

type Interval struct {
	Start, End float64
}

func (i *Interval) Sat(val float64) float64 {
	min, max := i.GetMinMax()

	if val > max {
		return max
	}
	if val < min {
		return min
	}

	return val
}

func (i *Interval) GetMinMax() (min, max float64) {
	min = i.Start
	max = i.End
	if min > max {
		return max, min
	}
	return
}

func (i *Interval) Length() float64 {
	return i.End - i.Start
}

type IntervalTransformer struct {
	Interval1 Interval
	Interval2 Interval
}

func (t *IntervalTransformer) TransformForward(val float64) float64 {
	return transform(val, &t.Interval1, &t.Interval2)
}

func (t *IntervalTransformer) TransformBackward(val float64) float64 {
	return transform(val, &t.Interval2, &t.Interval1)
}

func transform(val float64, i1, i2 *Interval) float64 {
	return i2.Sat((val-i1.Start)/i1.Length()*i2.Length() + i2.Start)
}
