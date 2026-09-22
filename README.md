# Algotiqa Agent

> [!TIP]
> The main documentation of the platform is located [HERE](https://github.com/algotiqa/docs).

## Introduction

This is an agent that collects metrics from external trading systems and exposes them through a REST API. It is used by the Algotiqa platform when the trading system's runtime is external to the platform itself (like TradeStation, MultiCharts, MetaTrader, NinjaTrader, etc...).


## File format

The agent scans a folder and read all files with a given extension. These files must have a CSV (Comma Separated Value) like format and each line may be an INFO, START, BAR, TRADE or DAILY line. The separator is a "|" character (the pipe).

There is only 1 INFO line per file and it is the first line. Its purpose is to provide general information and has the following structure:
- **INFO** : fixed text to identify the INFO line
- **root** : root name of the instrument (i.e. if the instrument is MNQH25 then the root is MNQ)
- **name** : Name of the trading system that is generating this file

The START line identifies the start of a trade and has the following structure:
- **START** : fixed text to identify the START line
- **date**  : date when the trade started (i.e. when the trading system entered the market)
- **time**  : time when the trade started

The BAR line tracks the strategy returns on each bar of the trade. The structure is:
- **BAR**       : fixed text to identify the BAR line
- **date**      : date of the bar
- **time**      : time of the bar
- **return**    : gross return of the trade at this bar (either a profit or loss) respect to when the trade was opened
- **contracts** : number of contracts held at this bar

The TRADE line indicates a new trade and has the following structure:
- **entryDate**  : Date when the trade started (i.e. entered the market)
- **entryTime**  : Time when the trade started
- **entryPrice** : Market price when entering the market
- **entryLabel** : Generic text used with the buy/sellshort commands. For example, in *buy("LE") 1 contracts* the label is *LE*
- **exitDate**   : Date when the trade ended
- **exitTime**   : Time when the trade ended
- **exitPrice**  : Market price when exiting the market
- **exitLabel**  : Generic text used with the sell/buyToCover commands
- **return**     : Trade's gross return: profit (if positive) or loss (if negative)
- **contracts**  : Max number of contracts bought or sold during the trade's life
- **operation**  : Type of operation (buy=1, sell=-1)

The DAILY line is used to keep track of daily returns, which are used to calculate the correlation between trading systems:
- **DAILY**     : fixed text to identify the DAILY line
- **date**      : date related to the end of the session
- **time**      : time related to the end of the session
- **return**    : gross return of the trade at the end of the session (either a profit or loss)

All values (entry/exit date/time/price, profit) refers to the platform that is running the trading system. The timezone of the date+time is the one of the exchange where the product is traded. The format of the fields is:

- **date** : format is *DD/MM/YYYY*
- **time** : This is an integer and a value of *900* represents the time at *09:00*

The [example_strategy.prod.trl](example_strategy.prod.trl) is an example of a file that conforms to this specification.


## Integration with trading platforms

### MultiCharts

This function, called **writeTrades** exports trades from a strategy (signal) running on MultiCharts:

```
Inputs: tag(string);

var: fileName(""), tt(0), pos(0), suffix(""), startEquity(0), barEquity(0), prevEquity(0), dailyReturn(0), mp(0);

once begin
	if StrLen(tag) <> 0 then begin
		fileName = "\reports\" + getstrategyname + "."+ tag +".trl";
		filedelete(fileName);
		Print(File(fileName),"INFO", "|",symbolroot, "|",getstrategyname);
	end;
end;


tt = totaltrades;
mp = MarketPosition;

if StrLen(tag) <> 0 then begin
	if tt<>tt[1] then begin
		for pos = tt - tt[1] downto 1 begin
			Print(File(fileName ),"BAR", "|", Date2String(exitdate(pos)), "|", exittime(pos):0:0, "|", positionProfit(pos):0:2, "|", 0:0:0);
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

	if (mp[1] = 0 and mp <> 0) or (mp[1] <> 0 and mp[1] <> mp) then begin
		Print(File(fileName ),"START", "|", Date2String(date), "|", time:0:0);
		startEquity = i_OpenEquity;
	end;
	
	if mp<>0 then begin
		barEquity = i_OpenEquity - startEquity;
		Print(File(fileName ),"BAR", "|", Date2String(date), "|", time:0:0, "|", barEquity:0:2, "|", CurrentContracts:0:0);
	end;
	
	if SessionLastBar then begin
		dailyReturn = i_OpenEquity - prevEquity;
		prevEquity  = i_OpenEquity;
		Print(File(fileName ),"DAILY", "|", Date2String(date), "|", time:0:0, "|", dailyReturn:0:2);
	end;
end;	

writeTrades = True;
```

This function can be added to a strategy adding these lines at the end of it:

```
input: tag("");
if tag<>"" then writeTrades(tag);
```
The *tag* can be any string. Typical strings are *dev* and *prod* (please, see the next session).

## Development and production modes

When a strategy is live, it is usually configured to look back up to 6 months ago to keep resource consumption at minimum. This configuration will generate trades only in this period of time (what we call *production* or simply *prod*) but for analysis reasons we would like to have also the trades generated during the development stage, which typically covers a 10+ years period (what we call *development* or simply *dev*). This need is covered adopting the following process:

- During the development of a strategy, the *dev* tag is used to generate all trades. The development period starts in the past (like 2010) and should arrive up to the most recent day (like today). The *writeTrades* function will generate a file like *mystrategy.dev.trl*

- When the strategy is completed and moved into production, the *prod* tag is used to generate a second export file, like *mystrategy.prod.trl*

- The agent is able to load multiple files related to the same strategy and to serves them to the platform

- The platform gets the files from the agent and sorts them depending on their filename: the presence of the *.dev.* string into the filename causes those trades to be imported before all others

- Trades in *prod* files will be imported *only if* their date is after the last imported date of all dev trades. This creates a continuous flow of trades


## Building

When using Linux as a development platform, to build for Windows just issue:
```
GOOS=windows GOARCH=amd64 go build
```
