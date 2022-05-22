package main

type Scenarios []Scenario

func (scenarios Scenarios) GetScenario(name string) Scenario {
	var result Scenario = Scenario{}
	for _, value := range scenarios {
		if value.Name == name {
			return value
		}
	}
	return result
}

func (scenarios Scenarios) MakeKeyboard(columns uint) ReplyKeyboardMarkup {
	var result [][]string = make([][]string, 0)
	row := make([]string, 0)
	count := 0
	for _, scenario := range scenarios {
		if count < int(columns) {
			row = append(row, scenario.Name)
			count++
		} else {
			result = append(result, row)
			row = make([]string, 0)
			count = 0
		}
	}
	if len(row) > 0 {
		result = append(result, row)
	}

	return ReplyKeyboardMarkup{result}
}
