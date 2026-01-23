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
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/algotiqa/agent/pkg/app"
	"golang.org/x/exp/maps"
)

//=============================================================================

const INFO = "INFO"
const TRADE = "TRADE"
const DAILY = "DAILY"

//=============================================================================

var tradingSystems *TradingSystemMap

//=============================================================================

func GetTradingSystems() []*TradingSystem {
	log.Println("Getting trading systems for client")
	return maps.Values(tradingSystems.TradingSystems)
}

//=============================================================================

func StartPeriodicScan(cfg *app.Config) *time.Ticker {

	ticker := time.NewTicker(cfg.Scan.PeriodHour * time.Hour)

	go func() {
		time.Sleep(2 * time.Second)
		run(cfg)

		for range ticker.C {
			run(cfg)
		}
	}()

	return ticker
}

//=============================================================================

func run(cfg *app.Config) {
	dir := cfg.Scan.Dir
	log.Println("Fetching files from: " + dir)

	files, err := os.ReadDir(dir)

	if err != nil {
		log.Println("Scan error: ", err)
	} else {
		tsMap := NewTradingSystemMap()

		for _, entry := range files {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), cfg.Scan.Extension) {
				ts := handleFile(dir, entry.Name())
				if ts != nil {
					if tsIn, ok := tsMap.TradingSystems[ts.Name]; ok {
						tsIn.TradeLists = append(tsIn.TradeLists, ts.TradeLists...)
					} else {
						tsMap.TradingSystems[ts.Name] = ts
					}
				}
			}
		}
		tradingSystems = tsMap
	}
}

//=============================================================================

func handleFile(dir string, fileName string) *TradingSystem {
	log.Println("Handling: " + fileName)

	path := dir + string(os.PathSeparator) + fileName
	file, err := os.Open(path)

	if err != nil {
		log.Println("Cannot open file for reading: " + path + " (cause is: " + err.Error() + " )")
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	ts := NewTradingSystem()
	tl := NewTradeList()

	for scanner.Scan() {
		if err = handleLine(ts, tl, scanner.Text()); err != nil {
			log.Println(err.Error())
			return nil
		}
	}

	if err = scanner.Err(); err != nil {
		log.Println("Cannot scan file: " + path + " (cause is: " + err.Error() + " )")
		return nil
	}

	ts.TradeLists = append(ts.TradeLists, tl)
	return ts
}

//=============================================================================

func handleLine(ts *TradingSystem, tl *TradeList, line string) error {
	tokens := strings.Split(line, "|")

	switch tokens[0] {
	case INFO:
		handleInfo(ts, tokens)
	case TRADE:
		if err := handleTrade(tl, tokens); err != nil {
			return err
		}
	case DAILY:
		if err := handleDaily(tl, tokens); err != nil {
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
	ts.Name = tokens[2]
}

//=============================================================================

func handleTrade(tl *TradeList, tokens []string) error {
	var err error

	entryDate := tokens[1]
	entryTime := tokens[2]
	entryPrice := tokens[3]
	entryLabel := tokens[4]
	exitDate := tokens[5]
	exitTime := tokens[6]
	exitPrice := tokens[7]
	exitLabel := tokens[8]
	grossProfit := tokens[9]
	contracts := tokens[10]
	position := tokens[11]

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

	tr.GrossProfit, err = strconv.ParseFloat(grossProfit, 64)
	if err != nil {
		return errors.New("Cannot parse gross profit: " + grossProfit)
	}

	tr.Contracts, err = strconv.ParseInt(contracts, 10, 32)
	if err != nil {
		return errors.New("Cannot parse contracts: " + contracts)
	}

	tr.Position, err = strconv.ParseInt(position, 10, 32)
	if err != nil {
		return errors.New("Cannot parse position: " + position)
	}

	//-----------------------------------------

	tl.Trades = append(tl.Trades, tr)
	return nil
}

//=============================================================================

func handleDaily(tl *TradeList, tokens []string) error {
	var err error

	ddate := tokens[1]
	dtime := tokens[2]
	grossProfit := tokens[3]
	trades := tokens[4]

	dp := NewDailyProfit()

	//-----------------------------------------

	dp.Date, err = convertDate(ddate)
	if err != nil {
		return err
	}

	dp.Time, err = strconv.ParseInt(dtime, 10, 32)
	if err != nil {
		return errors.New("Cannot parse time: " + dtime)
	}

	dp.GrossProfit, err = strconv.ParseFloat(grossProfit, 64)
	if err != nil {
		return errors.New("Cannot parse gross profit: " + grossProfit)
	}

	dp.Trades, err = strconv.ParseInt(trades, 10, 32)
	if err != nil {
		return errors.New("Cannot parse trades: " + trades)
	}

	//-----------------------------------------

	tl.DailyProfits = append(tl.DailyProfits, dp)
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
		log.Println("Bad value for day: " + date)
		return 0, errors.New("Date out of range: " + date)
	}

	return int(value), nil
}

//=============================================================================
