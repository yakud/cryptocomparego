package cryptocomparego

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lucazulian/cryptocomparego/context"
	"github.com/pkg/errors"
)

const (
	histominuteBasePath = "data/histominute"
)

// Get the history kline data of any cryptocurrency in any other currency that you need.
type HistominuteService interface {
	Get(context.Context, *HistominuteRequest) (*HistominuteResponse, *Response, error)
}

type HistominuteServiceOp struct {
	client *Client
}

var _ HistominuteService = &HistominuteServiceOp{}

type HistominuteResponse struct {
	Response          string         `json:"Response"`
	Message           string         `json:"Message"` // Error Message
	Type              int            `json:"Type"`
	Aggregated        bool           `json:"Aggregated"`
	Data              []Histominute  `json:"Data"`
	TimeTo            int64          `json:"TimeTo"`
	TimeFrom          int64          `json:"TimeFrom"`
	FirstValueInArray bool           `json:"FirstValueInArray"`
	ConversionType    conversionType `json:"ConversionType"`
}

type Histominute struct {
	Time       int64   `json:"time"`
	Close      float64 `json:"close"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	VolumeFrom float64 `json:"volumefrom"`
	VolumeTo   float64 `json:"volumeto"`
}

type HistominuteRequest struct {
	Fsym          string
	Tsym          string
	E             string
	ExtraParams   string
	Sign          bool
	TryConversion bool
	Aggregate     int // Not Used For Now
	Limit         int
	ToTs          int64
}

func NewHistominuteRequest(fsym string, tsym string, limitMinutes int, fromTime int64) *HistominuteRequest {
	pr := HistominuteRequest{Fsym: fsym, Tsym: tsym}
	pr.E = "CCCAGG"
	pr.Sign = false
	pr.TryConversion = true
	pr.Aggregate = 1
	if limitMinutes < 1 {
		limitMinutes = 1
	}
	if limitMinutes > 2000 {
		limitMinutes = 2000
	}
	pr.Limit = limitMinutes
	pr.ToTs = fromTime
	return &pr
}

func (hr *HistominuteRequest) FormattedQueryString(baseUrl string) string {
	values := url.Values{}

	if len(hr.Fsym) > 0 {
		values.Add("fsym", hr.Fsym)
	}

	if len(hr.Tsym) > 0 {
		values.Add("tsym", hr.Tsym)
	}

	if len(hr.E) > 0 {
		values.Add("e", hr.E)
	}

	if len(hr.ExtraParams) > 0 {
		values.Add("extraParams", hr.ExtraParams)
	}

	values.Add("sign", strconv.FormatBool(hr.Sign))
	values.Add("tryConversion", strconv.FormatBool(hr.TryConversion))
	values.Add("limit", strconv.FormatInt(int64(hr.Limit), 10))
	if hr.ToTs >= 0 {
		values.Add("toTs", strconv.FormatInt(int64(hr.ToTs), 10))
	}

	return fmt.Sprintf("%s?%s", baseUrl, values.Encode())
}

func (s *HistominuteServiceOp) Get(ctx context.Context, histominuteRequest *HistominuteRequest) (*HistominuteResponse, *Response, error) {

	path := histominuteBasePath

	if histominuteRequest != nil {
		path = histominuteRequest.FormattedQueryString(histominuteBasePath)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, *s.client.MinURL, path, nil)
	if err != nil {
		return nil, nil, err
	}

	res := new(HistominuteResponse)
	resp, err := s.client.Do(ctx, req, res)
	if err != nil {
		return nil, resp, err
	}

	if res.Response == "Error" {
		return nil, resp, errors.New(res.Message)
	}

	return res, resp, nil
}
