package main

var cityID2Open = [12000]bool{}

func init() {
	// readConfig()
	for i := 0; i < len(cityID2Open); i++ {
		if true { // city i is opened in configs
			cityID2Open[i] = true
		}
	}
}

func isPassed(cityID int) bool {
	return cityID2Open[cityID]
}
