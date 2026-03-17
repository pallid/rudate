package rudate

// levenshtein computes the Levenshtein distance between two strings.
func levenshtein(a, b string) int {
	la := []rune(a)
	lb := []rune(b)
	m := len(la)
	n := len(lb)

	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// Use single row optimization
	prev := make([]int, n+1)
	curr := make([]int, n+1)

	for j := 0; j <= n; j++ {
		prev[j] = j
	}

	for i := 1; i <= m; i++ {
		curr[0] = i
		for j := 1; j <= n; j++ {
			cost := 1
			if la[i-1] == lb[j-1] {
				cost = 0
			}
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost
			curr[j] = min3(del, ins, sub)
		}
		prev, curr = curr, prev
	}

	return prev[n]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// fuzzyLookup finds the closest dictionary match for a word.
// Returns the matching token and true if found within maxDistance.
// Only considers words of similar length (±maxDistance chars).
func fuzzyLookup(word string, maxDistance int) (Token, bool) {
	wordRunes := []rune(word)
	wordLen := len(wordRunes)

	bestDist := maxDistance + 1
	var bestToken Token
	found := false

	for dictWord, tok := range dictionary {
		dictRunes := []rune(dictWord)
		dictLen := len(dictRunes)

		// Quick length filter — no point checking if lengths differ too much
		if abs(wordLen-dictLen) > maxDistance {
			continue
		}

		// Skip very short words (1-2 chars) — too many false positives
		if dictLen <= 2 {
			continue
		}

		dist := levenshtein(word, dictWord)
		if dist < bestDist {
			bestDist = dist
			bestToken = tok
			bestToken.Raw = word
			found = true
		}
	}

	if found && bestDist <= maxDistance {
		return bestToken, true
	}
	return Token{}, false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
