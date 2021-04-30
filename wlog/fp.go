package wlog

import "strings"

type FingerPrints []string

func (fp FingerPrints) String() string {
	n := len(fp)
	if n == 0 {
		return "/"
	}

	if n == 1 {
		return "/" + fp[0]
	}

	for i, c := 0, n; i < c; i++ {
		n += len(fp[i])
	}
	sb := strings.Builder{}
	sb.Grow(n)
	for _, s := range fp {
		sb.WriteRune('/')
		sb.WriteString(s)
	}
	return sb.String()
}

func mustCombineFingerPrint(fp interface{}, appends FingerPrints) FingerPrints {
	if nil == fp {
		return appends // might returns nil
	}

	fpArr, ok := fp.(FingerPrints)
	if !ok {
		return appends
	}

	if appends == nil {
		return fpArr
	}

	nA := len(appends)
	if nA == 0 {
		return fpArr
	}

	return append(append(make(FingerPrints, 0, len(fpArr)+nA), fpArr...), appends...)
}
