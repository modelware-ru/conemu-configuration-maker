package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
)

type Data struct {
	Cols     int `json:"cols"`
	Commands []Command`json:"commands"`
}

type Command struct {
	Title   string `json:"title"`
	Command string `json:"cmd"`
	ConEmu
}

type ConEmu struct {
	Order   int // start from 1
	Parent  int // start from 1
	Percent int
	IsVert  bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Errorf("wrong arguments number\n")
		os.Exit(1)
	}

	configFileName := os.Args[1]

	configJson, err := os.Open(configFileName)

	if err != nil {
		fmt.Errorf("no such file\n")
		os.Exit(1)
	}

	defer configJson.Close()

	byteValue, _ := ioutil.ReadAll(configJson)

	var data Data

	json.Unmarshal(byteValue, &data)

	cellCount := len(data.Commands)
	rows := int(math.Ceil(float64(cellCount) / float64(data.Cols)))
	//completedRows := int(math.Floor(float64(cellCount) / float64(data.Cols)))

	order := 1
	for i := 0; i < rows; i++ {
		data.Commands[i*data.Cols].Order = order
		data.Commands[i*data.Cols].IsVert = true

		if i == 0 {
			data.Commands[i*data.Cols].Parent = 0
			data.Commands[i*data.Cols].Percent = 100
		} else {
			data.Commands[i*data.Cols].Parent = order - 1
			data.Commands[i*data.Cols].Percent = int(math.Floor(float64(100*(rows-i)) / float64(rows-i+1)))
		}
		order++

	}

	var parent int
	var percentIndex int
	for i := 0; i < cellCount; i++ {
		if data.Commands[i].Order != 0 {
			parent = data.Commands[i].Order
			percentIndex = 1
			continue
		}

		data.Commands[i].Order = order
		data.Commands[i].IsVert = false
		data.Commands[i].Parent = parent
		parent = order

		data.Commands[i].Percent = int(math.Floor(float64(100*(data.Cols-percentIndex)) / float64(data.Cols-percentIndex+1)))

		//if i == 0 {
		//	data.Commands[i*data.Cols].Parent = -1
		//	data.Commands[i*data.Cols].Percent = 100
		//} else {
		//	data.Commands[i*data.Cols].Parent = (i - 1) * data.Cols
		//	data.Commands[i*data.Cols].Percent = int(math.Floor(float64(100*(rows-i)) / float64(rows-i+1)))
		//}
		percentIndex++
		order++
	}

	sort.Slice(data.Commands, func(i, j int) bool {return data.Commands[i].ConEmu.Order < data.Commands[j].ConEmu.Order })

	for i := 0; i < cellCount; i++ {
		fmt.Printf("%s -cur_console", data.Commands[i].Command)
		if i != 0 {
			direct := 'H'
			if data.Commands[i].IsVert {
				direct = 'V'
			}
			fmt.Printf(":s%dT%d%c", data.Commands[i].Parent, data.Commands[i].Percent, direct)
		}
		fmt.Printf(":t:\"%s\"\n", data.Commands[i].Title)
	}
}
