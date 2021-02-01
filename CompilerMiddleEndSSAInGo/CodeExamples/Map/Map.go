package main

func main() {

	m := make(map[string]int)

	m["k1"] = 7
	m["k2"] = 13

	v1 := m["k1"]
	_ = v1

	_, prs := m["k2"]
	_ = prs

	delete(m, "k2")
}

// See how different map operations are carried out in SSA:
//	1. Creating a map
//	2. Assigning a value to a map
//	3. Reading from a map
//	4. Reading from a map into multiple variables
//	5. Deleting a key from a map
