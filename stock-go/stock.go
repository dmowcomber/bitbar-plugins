package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	// "github.com/guptarohit/asciigraph"

	"github.com/logrusorgru/aurora"
	finance "github.com/piquette/finance-go"
	financeChart "github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

var (
	// TODO: make these configurable
	stockTicker = "GOOG"
	graphRange  = day

	upTriangle   = aurora.Green("▲")
	downTriangle = aurora.Red("▼")
)

const (
	day   = "day" // default
	week  = "week"
	month = "month"
)

func main() {
	// print a refresh button at the bottom of the dropdown
	defer func() {
		fmt.Println("---")
		fmt.Println("Refresh Me| terminal=false refresh=true")
	}()

	q, err := quote.Get(stockTicker)
	if err != nil {
		fmt.Printf("%s: err\n", stockTicker)
		fmt.Println("---")
		fmt.Println(err.Error())
		return
	}

	if !(q.MarketState == finance.MarketStateRegular || q.MarketState == finance.MarketStatePre) {
		fmt.Println("💤")
		fmt.Println("---")
	} else {
		fmt.Printf("%s\n", stockTicker)
		fmt.Println("---")
	}

	triangle := upTriangle
	chartColor := chart.ColorGreen
	if q.RegularMarketChangePercent <= 0 {
		triangle = downTriangle
		chartColor = chart.ColorRed
	}

	graph := getChart(stockTicker)

	// TODO: fix this. the graph is broken during POST
	smallGraphText := getGraphText(graph, nil, chartColor, true)

	prefix := fmt.Sprintf("%.2f%s", q.RegularMarketPrice, triangle)

	fmt.Printf("%s$%.2f %s\n", prefix, q.RegularMarketChange, smallGraphText)
	fmt.Printf("%s%.2f%% %s\n", prefix, q.RegularMarketChangePercent, smallGraphText)

	// output in dropdown only
	fmt.Println("---")
	fmt.Printf("pre market price: %6.2f\n", aurora.White(q.PreMarketPrice))
	fmt.Printf("post market price: %6.2f\n", aurora.White(q.PostMarketPrice))
	fmt.Printf("Bid: %6.2f\n", aurora.White(q.Bid))
	fmt.Printf("Ask: %6.2f\n", aurora.White(q.Ask))
	fmt.Printf("Day High: %6.2f\n", aurora.White(q.RegularMarketDayHigh))
	fmt.Printf("Day Low: %6.2f\n", aurora.White(q.RegularMarketDayLow))
	fmt.Printf("52 Week High: %6.2f\n", aurora.White(q.FiftyTwoWeekHigh))
	fmt.Printf("52 Week Low: %6.2f\n", aurora.White(q.FiftyTwoWeekLow))

	fmt.Printf("RegularMarketPreviousClose: %.2f\n", aurora.White(q.RegularMarketPreviousClose))

	fmt.Println("---")
	fmt.Print("MarketState: ")
	textColor := aurora.Red(q.MarketState)
	if q.MarketState == finance.MarketStateRegular {
		textColor = aurora.Green("OPEN")
	} else if q.MarketState == finance.MarketStatePost || q.MarketState == finance.MarketStatePre {
		textColor = aurora.Yellow(q.MarketState)
	}
	fmt.Printf("%s\n", textColor)

	fmt.Println("---")
	fmt.Println(getGraphText(graph, nil, chartColor, false))
}

func getChart(symbol string) *chartItem {
	now := time.Now()

	// note truncate will truncate to UTC midnight
	today := now.Truncate(24 * time.Hour)
	// start := today
	interval := datetime.FiveMins

	// TODO: is this the correct time?
	// TODO: just use the market state (post, pre, etc) to figure out when to show yesterday vs today..
	today9am := today.Add(14 * time.Hour)
	start := today9am

	// if the day hasn't started yet, show the previous day, otherwise the chart will have an error: "infinite x-range delta"
	if start.After(now) {
		start = start.Add(-24 * time.Hour)
	}

	if graphRange == week {
		// last week
		start = now.AddDate(0, 0, -7)
		interval = datetime.OneHour
	} else if graphRange == month {
		// last week
		start = now.AddDate(0, -1, 0)
		interval = datetime.OneDay
	}

	iter := financeChart.Get(&financeChart.Params{
		Symbol:   symbol,
		Start:    datetime.New(&start),
		End:      datetime.New(&now),
		Interval: interval,
	})

	priceValues := make([]float64, 0, iter.Count())
	timeValues := make([]time.Time, 0, iter.Count())
	for iter.Next() {
		chartBar := iter.Bar()
		chartValue, _ := chartBar.Close.Float64()
		// chartTime := time.Unix(int64(chartBar.Timestamp), 0)
		if chartValue == 0.0 {
			continue
		}
		priceValues = append(priceValues, chartValue)
		timeValues = append(timeValues, time.Unix(int64(chartBar.Timestamp), 0))
	}

	return &chartItem{
		name:        symbol,
		priceValues: priceValues,
		timeValues:  timeValues,
	}
}

type chartItem struct {
	name        string
	priceValues []float64
	timeValues  []time.Time
}

func getGraphText(chart1, chart2 *chartItem, chartColor drawing.Color, small bool) string {
	width := 512
	height := 200

	if small {
		width = 30
		height = 26
	}

	// default format to dates
	valueFormatter := chart.TimeValueFormatterWithFormat("01-02")
	if graphRange == day {
		// format to hours for a single day graph
		valueFormatter = chart.TimeValueFormatterWithFormat("3:04PM")
	}

	graph := chart.Chart{
		Width:  width,
		Height: height,
		XAxis: chart.XAxis{
			ValueFormatter: valueFormatter,
		},

		Background: chart.Style{
			FillColor:   chart.ColorWhite,
			FontColor:   chart.ColorWhite,
			StrokeColor: chart.ColorWhite,
			DotColor:    chart.ColorWhite,
		},
		Canvas: chart.Style{
			FillColor:   chart.ColorWhite,
			FontColor:   chart.ColorWhite,
			StrokeColor: chart.ColorWhite,
			DotColor:    chart.ColorWhite,
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Style: chart.Style{
					StrokeColor: chartColor, // will supercede defaults
					FillColor:   chartColor, //.WithAlpha(64), // will supercede defaults
				},
				Name:    chart1.name,
				XValues: chart1.timeValues,
				YValues: chart1.priceValues,
			},
		},
	}

	if chart2 != nil {
		timeSeries2 := chart.TimeSeries{
			Name:    chart2.name,
			XValues: chart2.timeValues,
			YValues: chart2.priceValues,
		}
		graph.Series = append(graph.Series, timeSeries2)
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	if small {
		// hide axises and legend
		graph.XAxis = chart.HideXAxis()
		graph.YAxis = chart.HideYAxis()
		graph.Elements = []chart.Renderable{}

		graph.Background = chart.Style{
			Padding:   chart.NewBox(2, 0, 0, 2),
			FillColor: chart.ColorBlack,
		}
		graph.Canvas = chart.Style{
			FillColor: chart.ColorBlack,
		}
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Println("---")
		fmt.Println(err.Error())
	}

	imageBase64 := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return fmt.Sprintf("| image=%s\n", imageBase64)
}
