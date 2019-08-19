package conduction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
)

func ParseSimulationFromFile(fileName string) (sim *Simulation, err error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		err = fmt.Errorf("Error while reading case file: %v \n", err)
		return
	}

	err = json.Unmarshal([]byte(file), &sim)
	if err != nil {
		err = fmt.Errorf("Error while loading case file: %v \n", err)
		return
	}
	return
}

func (s *Simulation) GetDomainLimits() (start float64, end float64) {
	start = math.Inf(1)
	end = math.Inf(-1)

	// TODO : Implement domains inspection
	for _, domain := range s.Domains {
		if domain.Start < start {
			start = domain.Start
		}

		if domain.End > end {
			end = domain.End
		}
	}
	return
}
