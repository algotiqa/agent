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
	OpenTrade    []*EquityBar   `json:"openTrade"`
	DailyReturns []*DailyReturn `json:"dailyReturns"`
}

//=============================================================================

func NewTradeList() *TradeList {
	tl := TradeList{}
	tl.Trades = []*Trade{}
	return &tl
}

//=============================================================================
//===
//=== Trade
//===
//=============================================================================

type Trade struct {
	EntryDate    int          `json:"entryDate"`
	EntryTime    int64        `json:"entryTime"`
	EntryPrice   float64      `json:"entryPrice"`
	EntryLabel   string       `json:"entryLabel"`
	ExitDate     int          `json:"exitDate"`
	ExitTime     int64        `json:"exitTime"`
	ExitPrice    float64      `json:"exitPrice"`
	ExitLabel    string       `json:"exitLabel"`
	GrossReturn  float64      `json:"grossReturn"`
	MaxContracts int64        `json:"maxContracts"`
	Position     int64        `json:"position"`
	Equity       []*EquityBar `json:"equity"`
}

//=============================================================================

func NewTrade() *Trade {
	return &Trade{}
}

//=============================================================================

type EquityBar struct {
	Date        int     `json:"date"`
	Time        int64   `json:"time"`
	GrossReturn float64 `json:"grossReturn"`
	Contracts   int64   `json:"contracts"`
}

//=============================================================================

func NewEquityBar() *EquityBar {
	return &EquityBar{}
}

//=============================================================================

type DailyReturn struct {
	Date        int     `json:"date"`
	GrossReturn float64 `json:"grossReturn"`
}

//=============================================================================

func NewDailyReturn() *DailyReturn {
	return &DailyReturn{}
}

//=============================================================================
