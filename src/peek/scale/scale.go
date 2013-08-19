package scale

type f func(float64) float64
type domain f
type scale f

func Domain(a, b float64) domain {
	return uninterpolate(a, b)
}

func uninterpolate(a, b float64) domain {
	if (b - a) > 0 {
		b = 1 / (b - a)
	} else {
		b = 0
	}

	return func(x float64) float64 { return (x - a) * b }
}

func (d domain) Range(a, b float64) scale {
	i := interpolate(a, b)
	return func(x float64) float64 { return i(d(x)) }
}

func interpolate(a, b float64) f {
	b -= a
	return func(x float64) float64 { return a + b*x }
}
