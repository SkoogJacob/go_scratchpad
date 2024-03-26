package vat_calc

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type VatTable struct {
	InputPrices []float64 `json:"input_prices"`
	TaxRate     float64   `json:"tax_rate"`
	VatPrices   []float64 `json:"vat_prices"`
}

func roundFloat(number float64, accuracy uint8) float64 {
	p := math.Pow(10, float64(accuracy))
	return math.Round(number*p) / p
}

func New(inputPrices []float64, taxRate float64) VatTable {
	var vats = make([]float64, len(inputPrices))

	for index, val := range inputPrices {
		vats[index] = roundFloat(val*(1+(taxRate/100)), 2)
	}
	return VatTable{
		InputPrices: inputPrices,
		TaxRate:     taxRate,
		VatPrices:   vats,
	}
}

func FromFiles(priceListPath string, taxRatePath string) ([]VatTable, error) {
	var prices, taxes []float64
	prices, err := readFile(priceListPath)
	if err != nil {
		return []VatTable{}, err
	}
	taxes, err = readFile(taxRatePath)
	if err != nil {
		return []VatTable{}, err
	}
	vatTables := make([]VatTable, len(taxes))
	for index, tax := range taxes {
		vatTables[index] = New(prices, tax)
	}
	return vatTables, nil
}

func readFile(filePath string) ([]float64, error) {
	fileReader, err := os.Open(filePath)
	var values = make([]float64, 0, 8)
	if err != nil {
		_ = fileReader.Close()
		return []float64{}, err
	}
	scanner := bufio.NewScanner(fileReader)
	for scanner.Scan() {
		line := scanner.Text()
		f, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
		if err != nil {
			fmt.Printf("Unable to parse %s into a float with error:\n%v\n", line, err)
			continue
		}
		values = append(values, f)
	}
	if len(values) == 0 {
		_ = fileReader.Close()
		return []float64{}, errors.New(fmt.Sprint("unable to parse any viable values from ", filePath))
	}
	_ = fileReader.Close()
	return values, nil
}

func ToJson(vatTables []VatTable, targetPath string) error {
	baseWriter, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		_ = baseWriter.Close()
		panic(fmt.Sprintf("Unable to open file %s, got error:\n%v\n", targetPath, err))
	}
	encoder := json.NewEncoder(baseWriter)
	err = encoder.Encode(vatTables)
	_ = baseWriter.Close()
	return err
}
