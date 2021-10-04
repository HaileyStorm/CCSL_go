package math

import "math"

// DivisorNearLower finds a divisor of x near target and if possible within the range [low,high], and <= isqrt(x).
// If target is a divisor, it is returned.
// If not, then if searchUp is true, the range (target,high) is searched in ascending order and if a divisor
// is not found the range [low,target) is searched in *descending* order; if searchUp is false, the range
// [low,target) is searched in descending order and if a divisor is not found the range (target,high)
// is searched in *ascending* order.
// If a divisor has yet to be found, [high,isqrt(x)] is searched in ascending order if searchUp is true and
// [2,low) is searched in descending order if false; if a divisor is not found in that searched range, the
// other is searched in the order described. If a divisor is still not found, x is prime and 1 is returned.
// If a divisor >= isqrt(x) is desired, use DivisorNearUpper. If the closest divisor to target is desired
// whether or not it is < isqrt(x), use DivisorNear, which runs both functions and returns whichever
// does not panic (it recovers from the panic if 1 < low < target < high < x).
// This is not an efficient primality test. If x is large and there is a good chance x is prime,
// math/big's Int.ProbablyPrime should be used to check primality first.
// This is also not meant to be a highly efficient divisor search, e.g. it does not use a lookup table for
// the smallest n numbers or use a prime factor powers assisted search.
// The following conditions must be true, though they are checked at the relevant point in the logic,
// rather than up-front.
// Low must be in the range [2,target).
// Target must be in the range (low,high).
// High must be in the range (target,isqrt(x)].
func DivisorNearLower(x, target, low, high int, searchUp bool) int {
	if x%target == 0 {
		return target
	}
	if high <= low {
		panic("high must be > low")
	}

	var i int

	HighSqrt := func() int {
		isqrt := int(math.Floor(math.Sqrt(float64(x))))
		if high > isqrt {
			panic("high must be <= the integer square root of x")
		}
		return divisorNearSearchUp(x, high, isqrt+1)
	}
	Low2 := func() int {
		if low < 2 {
			panic("low must be > 1")
		}
		return divisorNearSearchDown(x, low-1, 1)
	}

	if searchUp {
		if i = divisorNearTargetHigh(x, target, low, high); i != 0 {
			return i
		}
		if i = divisorNearTargetLow(x, target, low, high); i != 0 {
			return i
		}
		if i = HighSqrt(); i != 0 {
			return i
		}
		if i = Low2(); i != 0 {
			return i
		}
	} else {
		if i = divisorNearTargetLow(x, target, low, high); i != 0 {
			return i
		}
		if i = divisorNearTargetHigh(x, target, low, high); i != 0 {
			return i
		}
		if i = Low2(); i != 0 {
			return i
		}
		if i = HighSqrt(); i != 0 {
			return i
		}
	}

	// Prime
	return 1
}

// DivisorNearUpper finds a divisor of x near target and if possible within the range [low,high], and >= isqrt(x).
// If target is a divisor, it is returned.
// If not, then if searchUp is true, the range (target,high) is searched in ascending order and if a divisor
// is not found the range [low,target) is searched in *descending* order; if searchUp is false, the range
// [low,target) is searched in descending order and if a divisor is not found the range (target,high)
// is searched in *ascending* order.
// If a divisor has yet to be found, [high,x) is searched in ascending order if searchUp is true and
// [isqrt(x),low) is searched in descending order if false; if a divisor is not found in that searched range, the
// other is searched in the order described. If a divisor is still not found, x is prime and 1 is returned.
// If a divisor <= isqrt(x) is desired, use DivisorNearLower. If the closest divisor to target is desired
// whether or not it is < isqrt(x), use DivisorNear, which runs both functions and returns whichever
// does not panic (it recovers from the panic as long as 1 < low < target < high < x).
// This is not an efficient primality test. If x is large and there is a good chance x is prime,
// math/big's Int.ProbablyPrime should be used to check primality first.
// This is also not meant to be a highly efficient divisor search, e.g. it does not use a lookup table for
// the smallest n numbers or use a prime factor powers assisted search.
// The following conditions must be true, though they are checked at the relevant point in the logic,
// rather than up-front.
// Low must be in the range [isqrt(x),target).
// Target must be in the range (low,high).
// High must be in the range (target,x).
func DivisorNearUpper(x, target, low, high int, searchUp bool) int {
	if x%target == 0 {
		return target
	}
	if high <= low {
		panic("high must be > low")
	}

	var i int

	HighX := func() int {
		if high >= x {
			panic("high must be < x")
		}
		return divisorNearSearchUp(x, high, x)
	}
	LowSqrt := func() int {
		isqrt := int(math.Floor(math.Sqrt(float64(x))))
		if low < isqrt {
			panic("low must be >= the integer square root of x")
		}
		return divisorNearSearchDown(x, low-1, isqrt-1)
	}

	if searchUp {
		if i = divisorNearTargetHigh(x, target, low, high); i != 0 {
			return i
		}
		if i = divisorNearTargetLow(x, target, low, high); i != 0 {
			return i
		}
		if i = HighX(); i != 0 {
			return i
		}
		if i = LowSqrt(); i != 0 {
			return i
		}
	} else {
		if i = divisorNearTargetLow(x, target, low, high); i != 0 {
			return i
		}
		if i = divisorNearTargetHigh(x, target, low, high); i != 0 {
			return i
		}
		if i = LowSqrt(); i != 0 {
			return i
		}
		if i = HighX(); i != 0 {
			return i
		}
	}

	// Prime
	return 1
}

// DivisorNear finds a divisor of x near target  and if possible within the range [low,high].
// If target is a divisor, it is returned.
// If not, then if searchUp is true, the range (target,high) is searched in ascending order and if a divisor
// is not found the range [low,target) is searched in *descending* order; if searchUp is false, the range
// [low,target) is searched in descending order and if a divisor is not found the range (target,high)
// is searched in *ascending* order.
// If a divisor has yet to be found, then if low, target and high are <= isqrt(X):
// [high,isqrt(x)] is searched in ascending order if searchUp is true and [2,low) is searched in
// descending order if false; if a divisor is not found in that searched range, the other is searched in
// the order described. If a divisor is still not found, x is prime and 1 is returned.
// If on the other hand low, target and high are > isqrt(x):
// [high,x) is searched in ascending order if searchUp is true and [isqrt(x),low) is searched in
// descending order if false; if a divisor is not found in that searched range, the other is searched in
// the order described. If a divisor is still not found, x is prime and 1 is returned.
// This function works by calling DivisorNearLower, and returns that value; if that call panics but the
// parameters are well-formed, it calls DivisorNearUpper and returns that value.
// If a divisor <= isqrt(x) is desired, use DivisorNearLower. If a divisor >= isqrt(x) is desired, use
// DivisorNearUpper.
// This is not an efficient primality test. If x is large and there is a good chance x is prime,
// math/big's Int.ProbablyPrime should be used to check primality first.
// This is also not meant to be a highly efficient divisor search, e.g. it does not use a lookup table for
// the smallest n numbers or use a prime factor powers assisted search.
// The following conditions must be true, though they are checked at the relevant point in the logic,
// rather than up-front.
// Low must be in the range [2,target).
// Target must be in the range (low,high).
// High must be in the range (target,x).
// Which is to say: 1 < low < target < high < x.
func DivisorNear(x, target, low, high int, searchUp bool) (d int) {
	defer func() {
		// DivisorNearLower panicked
		if r := recover(); r != nil {
			// We can only recover from the panic if the parameters are well-formed for one function or the other
			if 1 < low && low < target && target < high && high < x {
				// If this panics too, we can't properly recover, so we don't try
				d = DivisorNearUpper(x, target, low, high, searchUp)
			} else {
				panic(r)
			}
		}
	}()
	return DivisorNearLower(x, target, low, high, searchUp)

}

func divisorNearSearchUp(x, min, max int) int {
	for i := min; i < max; i++ {
		if x%i == 0 {
			return i
		}
	}
	return 0
}
func divisorNearSearchDown(x, max, min int) int {
	for i := max; i > min; i-- {
		if x%i == 0 {
			return i
		}
	}
	return 0
}
func divisorNearTargetHigh(x, target, low, high int) int {
	if target >= high {
		panic("target must be < high")
	}
	return divisorNearSearchUp(x, target+1, high)
}
func divisorNearTargetLow(x, target, low, high int) int {
	if target <= low {
		panic("target must be > low")
	}
	return divisorNearSearchDown(x, target-1, low-1)
}
