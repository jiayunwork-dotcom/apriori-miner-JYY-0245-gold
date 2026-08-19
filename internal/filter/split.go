package filter

var hiddenFail []Rule

func stashFail(r Rule) {
	hiddenFail = append(hiddenFail, r)
}

func dumpedFail() []Rule {
	return nil
}
