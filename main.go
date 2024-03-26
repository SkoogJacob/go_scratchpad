package main

/**
For this project I'm planning to read the two text files, one containing prices and one tax rates,
and then they will produce a json file that combines the two, thinking that the json file will basically
have structure
[
  TAX_RATE1: [price1, price2...],
  TAX_RATE2: [price1, price2...],
  ...
]
*/

import (
	"fmt"
	"os"
	"scratchpad/vat_calc"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		errMsg := fmt.Sprint("Need the path to both the price list and the tax rate list, instead found ", args)
		panic(errMsg)
	}
	priceListPath := args[0]
	taxListPath := args[1]
	vatTables, err := vat_calc.FromFiles(priceListPath, taxListPath)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, table := range vatTables {
			fmt.Println(table)
		}
	}
	err = vat_calc.ToJson(vatTables, "./resources/output.json")
	if err != nil {
		fmt.Println("Encountered error encoding to json:")
		fmt.Println(err)
	}
}
