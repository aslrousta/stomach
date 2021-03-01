package myers

import "strings"

// List of edit operation types.
const (
	Insert int = iota
	Delete
)

// Edit is an operation that inserts or deletes consecutive lines.
type Edit struct {
	Type  int    // Operation type: Insert or Delete.
	Index int    // Index of first line to be edited.
	Line  string // Edited line.
}

// Diff computes the required edit operations to transform before to after.
func Diff(before, after string) []*Edit {
	a, b := splitLines(before), splitLines(after)
	trace := shortestEdit(a, b)
	if len(trace) == 0 {
		return nil
	}
	edits := backtrack(a, b, trace)
	for i, n := 0, len(edits); i < n/2; i++ {
		edits[i], edits[n-i-1] = edits[n-i-1], edits[i]
	}
	return edits
}

type step struct {
	v   []int
	max int
}

func (s *step) get(k int) int {
	return s.v[s.max+k]
}

func (s *step) set(k, x int) {
	s.v[s.max+k] = x
}

func (s *step) clone() (w *step) {
	w = &step{v: make([]int, 2*s.max+1), max: s.max}
	copy(w.v, s.v)
	return w
}

func shortestEdit(a, b []string) (trace []*step) {
	n, m := len(a), len(b)
	max := n + m
	v := &step{v: make([]int, 2*max+1), max: max}
	for d := 0; d <= max; d++ {
		trace = append(trace, v.clone())
		for k := -d; k <= d; k += 2 {
			var x, y int
			if k == -d || (k < d && v.get(k-1) < v.get(k+1)) {
				x = v.get(k + 1)
			} else {
				x = v.get(k-1) + 1
			}
			y = x - k
			for x < n && y < m && a[x] == b[y] {
				x, y = x+1, y+1
			}
			v.set(k, x)
			if x == n && y == m {
				return append(trace, v.clone())
			}
		}
	}
	return trace
}

func backtrack(a, b []string, trace []*step) (edits []*Edit) {
	x, y := len(a), len(b)
	for d := len(trace) - 1; d >= 0; d-- {
		v := trace[d]
		k := x - y
		var kp int
		if k == -d || (k < d && v.get(k-1) < v.get(k+1)) {
			kp = k + 1
		} else {
			kp = k - 1
		}
		xp := v.get(kp)
		yp := xp - kp
		for x > xp && y > yp {
			x, y = x-1, y-1
		}
		if d > 0 {
			if x == xp {
				edits = append(edits, &Edit{
					Type:  Insert,
					Index: x,
					Line:  b[yp],
				})
			} else {
				edits = append(edits, &Edit{
					Type:  Delete,
					Index: x,
					Line:  a[xp],
				})
			}
		}
		if x, y = xp, yp; x == 0 && y == 0 {
			break
		}
	}
	return edits
}

func splitLines(text string) (lines []string) {
	lines = strings.SplitAfter(text, "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimRight(lines[i], "\r\n")
	}
	return lines
}
