package data_structure_and_algorithm

func convert(s string, numRows int) string {
	if numRows == 1 {
		return s
	}

	ans := make([][]byte, numRows)
	var out []byte
	n := len(s)
	roundNum := 2*numRows - 2
	for i := 0; i < n; i++ {
		tempI := i % roundNum
		remain := tempI % numRows
		divider := tempI / numRows
		row := remain
		if divider > 0 {
			row = numRows - 1 - divider - remain
		}
		ans[row] = append(ans[row], s[i])
	}

	for _, v := range ans {
		out = append(out, v...)
	}
	return string(out)
}