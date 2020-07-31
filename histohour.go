package cryptocomparego

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/lucazulian/cryptocomparego/context"
	"github.com/pkg/errors"
)

const (
	histohourBasePath = "data/histohour"
)

// Get the history kline data of any cryptocurrency in any other currency that you need.
type HistohourService interface {
	Get(context.Context, *HistohourRequest) (*HistohourResponse, *Response, error)
}

type HistohourServiceOp struct {
	client *Client
}

var _ HistohourService = &HistohourServiceOp{}

type HistohourResponse struct {
	Response          string         `json:"Response"`
	Message           string         `json:"Message"` // Error Message
	Type              int            `json:"Type"`
	Aggregated        bool           `json:"Aggregated"`
	Data              []Histohour    `json:"Data"`
	TimeTo            int64          `json:"TimeTo"`
	TimeFrom          int64          `json:"TimeFrom"`
	FirstValueInArray bool           `json:"FirstValueInArray"`
	ConversionType    conversionType `json:"ConversionType"`
}

type Histohour struct {
	Time       int64   `json:"time"`
	Close      float64 `json:"close"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	VolumeFrom float64 `json:"volumefrom"`
	VolumeTo   float64 `json:"volumeto"`
}

type HistohourRequest struct {
	Fsym          string
	Tsym          string
	E             string
	ExtraParams   string
	Sign          bool
	TryConversion bool
	Aggregate     int // Not Used For Now
	Limit         int
	ToTs          time.Time
}

func NewHistohourRequest(fsym string, tsym string, limit int, fromTime time.Time) *HistohourRequest {
	pr := HistohourRequest{Fsym: fsym, Tsym: tsym}
	pr.E = "CCCAGG"
	pr.Sign = false
	pr.TryConversion = true
	pr.Aggregate = 1
	if limit < 1 {
		limit = 1
	}
	if limit > 2000 {
		limit = 2000
	}
	pr.Limit = limit
	pr.ToTs = fromTime
	return &pr
}

func (hr *HistohourRequest) FormattedQueryString(baseUrl string) string {
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
	if hr.ToTs.Unix() >= 0 {
		values.Add("toTs", strconv.FormatInt(int64(hr.ToTs.Unix()), 10))
	}

	return fmt.Sprintf("%s?%s", baseUrl, values.Encode())
}

func (s *HistohourServiceOp) Get(ctx context.Context, histohourRequest *HistohourRequest) (*HistohourResponse, *Response, error) {

	path := histohourBasePath

	if histohourRequest != nil {
		path = histohourRequest.FormattedQueryString(histohourBasePath)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, *s.client.MinURL, path, nil)
	if err != nil {
		return nil, nil, err
	}

	res := new(HistohourResponse)
	resp, err := s.client.Do(ctx, req, res)
	if err != nil {
		return nil, resp, err
	}

	if res.Response == "Error" {
		return nil, resp, errors.New(res.Message)
	}

	return res, resp, nil
}
