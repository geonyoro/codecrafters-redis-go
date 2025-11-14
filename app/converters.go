package main

func StringArraytoAny(vals []string) (anys []any) {
	for _, val := range vals {
		anys = append(anys, val)
	}
	return anys
}
