func (a *mockAnalyser) NewExecuter(_ string) (Executer, error) {
		{Name: "Name1", Path: "tool1", Args: "-flag %BASE_BRANCH% ./..."},
		{Name: "Name2", Path: "tool2"},
		ExecuteOut: [][]byte{
			{}, // git clone
			{}, // git fetch
			{}, // install-deps.sh
			[]byte(`/go/src/gopherci`),                   // pwd
			[]byte("main.go:1: error1"),                  // tool 1
			[]byte("/go/src/gopherci/main.go:1: error2"), // tool 2 output abs paths
		},
	expected := []Issue{
		{File: "main.go", HunkPos: 1, Issue: "Name1: error1"},
		{File: "main.go", HunkPos: 1, Issue: "Name2: error2"},
	}
		{"install-deps.sh"},
		{"pwd"},
		{"tool1", "-flag", "FETCH_HEAD", "./..."},
		{"tool2"},