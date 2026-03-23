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

package service

import (
	"net/http"

	"github.com/algotiqa/agent/pkg/core"
	"github.com/algotiqa/core/auth"
)

//=============================================================================

func getTradingSystems(c *auth.Context) {
	c.Gin.JSON(http.StatusOK, core.GetTradingSystems())
}

//=============================================================================

func reloadTradingSystem(c *auth.Context) {
	var rr ReloadRequest
	err := c.BindParamsFromBody(&rr)

	if err == nil {
		var ts *core.TradingSystem
		ts, err = core.ReloadTradingSystem(rr.Name)
		if err == nil {
			_ = c.ReturnObject(ts)
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func listTradingSystems(c *auth.Context) {
	names, err := core.ListTradingSystems()
	if err == nil {
		_ = c.ReturnObject(names)
	}

	c.ReturnError(err)
}

//=============================================================================

func handlePing(c *auth.Context) {
	c.Gin.JSON(http.StatusOK, "ok")
}

//=============================================================================
