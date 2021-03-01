package myers

import "fmt"

func Example() {
	edits := Diff(
		"A\nB\nC\nA\nB\nB\nA",
		"C\nB\nA\nB\nA\nC",
	)
	for _, edit := range edits {
		var sign byte
		if edit.Type == Insert {
			sign = '+'
		} else {
			sign = '-'
		}
		fmt.Printf("%c %02d %s\n", sign, edit.Index, edit.Line)
	}
	// Output:
	// - 01 A
	// - 02 B
	// + 03 B
	// - 06 B
	// + 07 C
}
