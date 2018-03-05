package strRadix

func (e *Edges) replace(atIndex int, edge Edge) {
	// attention : first we're checking for star : causes malfunction because it's not a slice of pointers
	edge.setStars()
	// replace it in the slice
	(*e)[atIndex] = edge
	// we're always sorting in reverse, so stars are last siblings
	e.sort()
}

func (e *Edges) add(edge Edge) {
	// attention : first we're checking for star : causes malfunction because it's not a slice of pointers
	edge.setStars()
	// add it to the slice
	*e = append(*e, edge)
	// we're always sorting in reverse, so stars are last siblings
	e.sort()
}

func (e Edges) less(i, j int) bool {
	switch compare(e[i].label, e[j].label) {
	case -1:
		return false
	case 0, 1:
		return true
	default:
		panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return true
	}
}

func (e Edges) swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Edges) sort() {
	n := len(e)
	quickSort(e, 0, n, maxDepth(n))
}

func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func quickSort(data Edges, a, b, maxDepth int) {
	for b-a > 12 { // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 {
			heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ {
			if data.less(i, i-6) {
				data.swap(i, i-6)
			}
		}
		insertionSort(data, a, b)
	}
}

func insertionSort(data Edges, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.less(j, j-1); j-- {
			data.swap(j, j-1)
		}
	}
}

func siftDown(data Edges, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.less(first+child, first+child+1) {
			child++
		}
		if !data.less(first+root, first+child) {
			return
		}
		data.swap(first+root, first+child)
		root = child
	}
}

func heapSort(data Edges, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data.swap(first, first+i)
		siftDown(data, lo, i, first)
	}
}

func doPivot(data Edges, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(data, lo, lo+s, lo+2*s)
		medianOfThree(data, m, m-s, m+s)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree(data, lo, m, hi-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && data.less(a, pivot); a++ {
	}
	b := a
	for {
		for ; b < c && !data.less(pivot, b); b++ { // data[b] <= pivot
		}
		for ; b < c && data.less(pivot, c-1); c-- { // data[c-1] > pivot
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		data.swap(b, c-1)
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if !data.less(pivot, hi-1) { // data[hi-1] = pivot
			data.swap(c, hi-1)
			c++
			dups++
		}
		if !data.less(b-1, pivot) { // data[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if !data.less(m, pivot) { // data[m] = pivot
			data.swap(m, b-1)
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			for ; a < b && !data.less(b-1, pivot); b-- { // data[b] == pivot
			}
			for ; a < b && data.less(a, pivot); a++ { // data[a] < pivot
			}
			if a >= b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			data.swap(a, b-1)
			a++
			b--
		}
	}
	// Swap pivot into middle
	data.swap(pivot, b-1)
	return b - 1, c
}

func medianOfThree(data Edges, m1, m0, m2 int) {
	// sort 3 elements
	if data.less(m1, m0) {
		data.swap(m1, m0)
	}
	// data[m0] <= data[m1]
	if data.less(m2, m1) {
		data.swap(m2, m1)
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if data.less(m1, m0) {
			data.swap(m1, m0)
		}
	}
	// now data[m0] <= data[m1] <= data[m2]
}
