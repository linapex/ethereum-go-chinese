
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450099258331136>

package librato

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const Operations = "operations"
const OperationsShort = "ops"

type LibratoClient struct {
	Email, Token string
}

//属性字符串
const (
//显示属性
	Color             = "color"
	DisplayMax        = "display_max"
	DisplayMin        = "display_min"
	DisplayUnitsLong  = "display_units_long"
	DisplayUnitsShort = "display_units_short"
	DisplayStacked    = "display_stacked"
	DisplayTransform  = "display_transform"
//特殊仪表显示属性
	SummarizeFunction = "summarize_function"
	Aggregate         = "aggregate"

//公制键
	Name        = "name"
	Period      = "period"
	Description = "description"
	DisplayName = "display_name"
	Attributes  = "attributes"

//测量键
	MeasureTime = "measure_time"
	Source      = "source"
	Value       = "value"

//专用仪表键
	Count      = "count"
	Sum        = "sum"
	Max        = "max"
	Min        = "min"
	SumSquares = "sum_squares"

//批密钥
	Counters = "counters"
	Gauges   = "gauges"

MetricsPostUrl = "https://度量API.libato.com/v1/metrics“
)

type Measurement map[string]interface{}
type Metric map[string]interface{}

type Batch struct {
	Gauges      []Measurement `json:"gauges,omitempty"`
	Counters    []Measurement `json:"counters,omitempty"`
	MeasureTime int64         `json:"measure_time"`
	Source      string        `json:"source"`
}

func (c *LibratoClient) PostMetrics(batch Batch) (err error) {
	var (
		js   []byte
		req  *http.Request
		resp *http.Response
	)

	if len(batch.Counters) == 0 && len(batch.Gauges) == 0 {
		return nil
	}

	if js, err = json.Marshal(batch); err != nil {
		return
	}

	if req, err = http.NewRequest("POST", MetricsPostUrl, bytes.NewBuffer(js)); err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Email, c.Token)

	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			body = []byte(fmt.Sprintf("(could not fetch response body for error: %s)", err))
		}
		err = fmt.Errorf("Unable to post to Librato: %d %s %s", resp.StatusCode, resp.Status, string(body))
	}
	return
}

