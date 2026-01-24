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

//=============================================================================
//===
//=== TradingSystemMap
//===
//=============================================================================

type TradingSystemMap struct {
	TradingSystems map[string]*TradingSystem `json:"tradingSystems"`
}

//=============================================================================

func NewTradingSystemMap() *TradingSystemMap {
	ss := &TradingSystemMap{}
	ss.TradingSystems = map[string]*TradingSystem{}
	return ss
}

//=============================================================================
//===
//=== TradingSystem
//===
//=============================================================================

type TradingSystem struct {
	Name       string `json:"name"`
	DataSymbol string `json:"dataSymbol"`

	TradeLists []*TradeList `json:"tradeLists"`
}

//=============================================================================

func NewTradingSystem() *TradingSystem {
	ts := &TradingSystem{}
	ts.TradeLists = []*TradeList{}
	return ts
}

//=============================================================================
//===
//=== TradeList
//===
//=============================================================================

type TradeList struct {
	FileName     string         `json:"fileName"`
	Trades       []*Trade       `json:"trades"`
	DailyProfits []*DailyProfit `json:"dailyProfits"`
}

//=============================================================================

func NewTradeList() *TradeList {
	tl := TradeList{}
	tl.Trades = []*Trade{}
	tl.DailyProfits = []*DailyProfit{}
	return &tl
}

//=============================================================================
//===
//=== Trade
//===
//=============================================================================

type Trade struct {
	EntryDate   int     `json:"entryDate"`
	EntryTime   int64   `json:"entryTime"`
	EntryPrice  float64 `json:"entryPrice"`
	EntryLabel  string  `json:"entryLabel"`
	ExitDate    int     `json:"exitDate"`
	ExitTime    int64   `json:"exitTime"`
	ExitPrice   float64 `json:"exitPrice"`
	ExitLabel   string  `json:"exitLabel"`
	GrossProfit float64 `json:"grossProfit"`
	Contracts   int64   `json:"contracts"`
	Position    int64   `json:"position"`
}

//=============================================================================

func NewTrade() *Trade {
	return &Trade{}
}

//=============================================================================

type DailyProfit struct {
	Date        int     `json:"date"`
	Time        int64   `json:"time"`
	GrossProfit float64 `json:"grossProfit"`
	Trades      int64   `json:"trades"`
}

//=============================================================================

func NewDailyProfit() *DailyProfit {
	return &DailyProfit{}
}

//=============================================================================
