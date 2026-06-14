//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package core

import (
	"bufio"
	"errors"
	"log/slog"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/algotiqa/agent/pkg/app"
)

//=============================================================================

const INFO  = "INFO"
const START = "START"
const BAR   = "BAR"
const TRADE = "TRADE"

//=============================================================================

var config *app.Config

//=============================================================================

func Init(cfg *app.Config) {
	config = cfg
}

//=============================================================================

func ListTradingSystems() ([]string, error) {
	dir := config.Scan.Dir

	files, err := os.ReadDir(dir)

	if err != nil {
		slog.Error("Cannot scan the directory", "dir", dir, "error", err)
		return nil, errors.New("Cannot scan the directory '" + dir + "'. Error: " + err.Error())
	}

	names := map[string]bool{}

	for _, entry := range files {
		fileName := entry.Name()
		if !entry.IsDir() {
			if strings.HasSuffix(fileName, config.Scan.Extension) {
				idx  := strings.Index(fileName,".")
				name := fileName[:idx]

				names[name] = true
			}
		}
	}

	list := slices.Collect(maps.Keys(names))
	slog.Info("ListTradingSystems: Got list of trading systems", "count", len(list))

	return list, nil
}

//=============================================================================

func GetTradingSystem(name string) (*TradingSystem, error) {
	dir := config.Scan.Dir

	files, err := os.ReadDir(dir)

	if err != nil {
		slog.Error("Cannot scan the directory", "dir", dir, "error", err)
		return nil, errors.New("Cannot scan the directory '" + dir + "'. Error: " + err.Error())
	}

	var ts *TradingSystem

	for _, entry := range files {
		fileName := entry.Name()
		if !entry.IsDir() && strings.HasPrefix(fileName, name) && strings.HasSuffix(fileName, config.Scan.Extension) {
			ts1, err1 := handleFile(dir, fileName)
			if err1 == nil {
				if ts == nil {
					ts = ts1
				} else {
					ts.TradeLists = append(ts.TradeLists, ts1.TradeLists...)
				}
			} else {
				slog.Error("Cannot process file", "file", fileName, "error", err1)
				return nil, errors.New("Cannot process file '" + fileName + "'. Error: " + err1.Error())
			}
		}
	}

	slog.Info("GetTradingSystem: Trading system loaded", "name", name)

	return ts, nil
}

//=============================================================================

func handleFile(dir string, fileName string) (*TradingSystem, error) {
	path := dir + string(os.PathSeparator) + fileName
	file, err := os.Open(path)

	if err != nil {
		return nil, errors.New("Cannot open file for reading: " + path + " (cause is: " + err.Error() + " )")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	ts := NewTradingSystem()
	tl := NewTradeList()
	tl.FileName = fileName

	for scanner.Scan() {
		if err = handleLine(ts, tl, scanner.Text()); err != nil {
			return nil, errors.New(err.Error())
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, errors.New("Cannot scan file: " + path + " (cause is: " + err.Error() + " )")
	}

	ts.TradeLists = append(ts.TradeLists, tl)
	return ts, nil
}

//=============================================================================

func handleLine(ts *TradingSystem, tl *TradeList, line string) error {
	tokens := strings.Split(line, "|")

	switch tokens[0] {
	case INFO:
		handleInfo(ts, tokens)
	case START:
		tl.OpenTrade = nil
	case BAR:
		if err := handleBar(tl, tokens); err != nil {
			return err
		}
	case TRADE:
		if err := handleTrade(ts, tl, tokens); err != nil {
			return err
		}
	default:
		return errors.New("Unknown token: " + tokens[0])
	}

	return nil
}

//=============================================================================

func handleInfo(ts *TradingSystem, tokens []string) {
	ts.DataSymbol = tokens[1]
	ts.Name       = tokens[2]
}

//=============================================================================

func handleBar(tl *TradeList, tokens []string) error {
	var err error

	ddate       := tokens[1]
	dtime       := tokens[2]
	grossReturn := tokens[3]
	contracts   := tokens[4]

	eb := NewEquityBar()

	//-----------------------------------------

	eb.Date, err = convertDate(ddate)
	if err != nil {
		return err
	}

	eb.Time, err = strconv.ParseInt(dtime, 10, 32)
	if err != nil {
		return errors.New("Cannot parse time: " + dtime)
	}

	eb.GrossReturn, err = strconv.ParseFloat(grossReturn, 64)
	if err != nil {
		return errors.New("Cannot parse gross return: " + grossReturn)
	}

	eb.Contracts, err = strconv.ParseInt(contracts, 10, 32)
	if err != nil {
		return errors.New("Cannot parse contracts: " + contracts)
	}

	//-----------------------------------------

	tl.OpenTrade = append(tl.OpenTrade, eb)
	return nil
}

//=============================================================================

func handleTrade(ts *TradingSystem, tl *TradeList, tokens []string) error {
	var err error

	entryDate    := tokens[1]
	entryTime    := tokens[2]
	entryPrice   := tokens[3]
	entryLabel   := tokens[4]
	exitDate     := tokens[5]
	exitTime     := tokens[6]
	exitPrice    := tokens[7]
	exitLabel    := tokens[8]
	grossReturn  := tokens[9]
	maxContracts := tokens[10]
	position     := tokens[11]

	tr := NewTrade()

	//-----------------------------------------

	tr.EntryDate, err = convertDate(entryDate)
	if err != nil {
		return err
	}

	tr.EntryTime, err = strconv.ParseInt(entryTime, 10, 32)
	if err != nil {
		return errors.New("Cannot parse entry time: " + entryTime)
	}

	tr.EntryPrice, err = strconv.ParseFloat(entryPrice, 64)
	if err != nil {
		return errors.New("Cannot parse entry price: " + entryPrice)
	}

	tr.EntryLabel = entryLabel

	//-----------------------------------------

	tr.ExitDate, err = convertDate(exitDate)
	if err != nil {
		return err
	}

	tr.ExitTime, err = strconv.ParseInt(exitTime, 10, 32)
	if err != nil {
		return errors.New("Cannot parse exit time: " + exitTime)
	}

	tr.ExitPrice, err = strconv.ParseFloat(exitPrice, 64)
	if err != nil {
		return errors.New("Cannot parse exit price: " + exitPrice)
	}

	tr.ExitLabel = exitLabel

	//-----------------------------------------

	tr.GrossReturn, err = strconv.ParseFloat(grossReturn, 64)
	if err != nil {
		return errors.New("Cannot parse gross return: " + grossReturn)
	}

	tr.MaxContracts, err = strconv.ParseInt(maxContracts, 10, 32)
	if err != nil {
		return errors.New("Cannot parse max contracts: " + maxContracts)
	}

	tr.Position, err = strconv.ParseInt(position, 10, 32)
	if err != nil {
		return errors.New("Cannot parse position: " + position)
	}

	//-----------------------------------------

	tr.Equity    = tl.OpenTrade
	tl.Trades    = append(tl.Trades, tr)
	tl.OpenTrade = nil
	return nil
}

//=============================================================================

func convertDate(date string) (int, error) {
	tokens := strings.Split(date, "/")

	if len(tokens) != 3 {
		return 0, errors.New("Bad format for date: " + date)
	}

	value, err := strconv.ParseInt(tokens[2]+tokens[1]+tokens[0], 10, 32)

	if err != nil {
		return 0, errors.New("Cannot convert date to int: " + date)
	}

	if value < 20000000 || value > 30000000 {
		return 0, errors.New("Date out of range: " + date)
	}

	return int(value), nil
}

//=============================================================================
