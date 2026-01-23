# Algotiqa Agent

## Introduction

This is an agent that collects metrics from external trading systems and exposes them through a REST API. It is used by 
the Algotiqa platform when the trading system's runtime is external to the platform itself (like TradeStation and MultiCharts).

## File format

The agent scans a folder and read all files with a given extension. These files must have a CSV (Comma Separated Value) format
and each line may be an INFO line or a TRADE line. The separator is a "|" character (the pipe).

There can be only 1 INFO line per file with the following structure:
- **INFO** : fixed text to identify the INFO line
- **root** : root name of the instrument (i.e. if the instrument is MNQH25 then the root is MNQ)
- **name** : Name of the trading system

The TRADE lines have the following structure:
- **entryDate**  : Date when the trade started (i.e. entered the market). The format is *dd/mm/yyyy*
- **entryTime**  : Time when the trade started. This is an integer abd a value of *900* represents the time at *09:00*
- **entryPrice** : Market price when entering the market
- **entryLabel** : Generic text used with the buy/sellshort commands. For example, in *buy("LE") 1 contracts* the label is *LE*
- **exitDate**   : Date when the trade ended
- **exitTime**   : Time when the trade ended
- **exitPrice**  : Market price when exiting the market
- **exitLabel**  : Generic text used with the sell/buyToCover commands
- **profit**     : Trade's profit (if positive) or loss (if negative)
- **contracts**  : Number of contracts bought or sold
- **operation**  : Type of operation (buy=1, sell=-1)

All values (entry/exit date/time/price, profit) refers to the platform that is running the trading system, not the broker.

Here is an example of a file's content that is accepted:

```
INFO|MNQ|_Ac_MNQ_Bias_30m_v1
TRADE|31/01/2025|900|22009.50000000|Short|31/01/2025|900|22069.50000000|Stop Loss|-360.00|3|-1
TRADE|03/02/2025|1330|21699.00000000|Buy#1|03/02/2025|1500|21639.00000000|Stop Loss|-360.00|3|1
TRADE|04/02/2025|700|21683.25000000|Short|04/02/2025|900|21743.25000000|Stop Loss|-360.00|3|-1
TRADE|04/02/2025|1330|21854.00000000|Buy#1|04/02/2025|1730|21794.00000000|Stop Loss|-360.00|3|1
TRADE|05/02/2025|1330|21950.00000000|Buy#1|06/02/2025|700|21977.50000000|Short|165.00|3|1
TRADE|06/02/2025|700|21977.50000000|Short|06/02/2025|900|22037.50000000|Stop Loss|-360.00|3|-1
TRADE|06/02/2025|1330|22050.00000000|Buy#1|06/02/2025|1400|21990.00000000|Stop Loss|-360.00|3|1
TRADE|10/02/2025|1330|22111.00000000|Buy#1|10/02/2025|1730|22051.00000000|Stop Loss|-360.00|3|1
TRADE|11/02/2025|700|21970.25000000|Short|11/02/2025|900|22030.25000000|Stop Loss|-360.00|3|-1
TRADE|13/02/2025|700|22052.75000000|Short|13/02/2025|800|22112.75000000|Stop Loss|-360.00|3|-1
TRADE|13/02/2025|1330|22231.25000000|Buy#1|14/02/2025|700|22319.50000000|Short|529.50|3|1
TRADE|14/02/2025|700|22319.50000000|Short|14/02/2025|900|22379.50000000|Stop Loss|-360.00|3|-1
TRADE|14/02/2025|1330|22425.50000000|Buy#1|17/02/2025|700|22494.25000000|Sell#1|412.50|3|1
TRADE|19/02/2025|1330|22470.25000000|Buy#1|19/02/2025|2000|22410.25000000|Stop Loss|-360.00|3|1
TRADE|03/03/2025|700|21307.50000000|Short|03/03/2025|930|21045.00000000|Profit Target|1575.00|3|-1
TRADE|05/03/2025|700|20721.00000000|Short|05/03/2025|1030|20600.00000000|Cover|726.00|3|-1
TRADE|05/03/2025|1330|20811.75000000|Buy#1|06/03/2025|230|20751.75000000|Stop Loss|-360.00|3|1
TRADE|12/03/2025|700|19792.50000000|Short|12/03/2025|800|19884.25000000|Stop Loss|-550.50|3|-1
TRADE|12/03/2025|1330|19748.25000000|Buy#1|13/03/2025|0|19688.00000000|Stop Loss|-361.50|3|1
```

## Integration with trading platforms

### MultiCharts

This function exports trades from a strategy (signal) running on Multicharts:

```
Inputs: tag(string);

var: fileName(""), tt(0), pos(0), suffix("");

once begin
	if StrLen(tag) <> 0 then begin
		fileName = "\reports\" + getstrategyname + "."+ tag +".trl";
		filedelete(fileName);
		Print(File(fileName),"INFO", "|",symbolroot, "|",getstrategyname);
	end;
end;

tt = totaltrades;

if StrLen(tag) <> 0 and tt<>tt[1] then begin
	for pos = tt - tt[1] downto 1 begin
		Print(File(fileName ),"TRADE", 
			"|", Date2String(entrydate(pos)), 
			"|", entrytime(pos):0:0, 
			"|", entryprice(pos):0:8,
			"|", entryname(pos),
			"|", Date2String(exitdate(pos)),
			"|", exittime(pos):0:0,
			"|", exitprice(pos):0:8,
			"|", exitname(pos),
			"|", positionprofit(pos):0:2,
			"|", maxcontracts(pos):0:0,
			"|", marketposition(pos):0:0
			);
	end;
end;	

writeTrades = True;
```

This function can be added to a strategy adding these lines at the end of it:

```
input: tag("");
if tag<>"" then writeTrades(tag);
```


## Building

When using Linux as a development platform, to build for Windows just issue:
```
GOOS=windows GOARCH=amd64 go build
```
