package main

func mainFor() {
	defer func() {
		for {
		}
	}()
}

func mainSelect() {
	defer func() { select {} }()
}

func mainChannel() {
	defer func() { <-make(chan bool) }()
}
