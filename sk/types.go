package sk

import (
	"github.com/mxmCherry/openrtb/v15/openrtb2"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type StoredRequestSite struct {
	Page string `json:"page"`
}

type StoredRequest struct {
	ID         string                 `json:"id"`
	Site       StoredRequestSite      `json:"site"`
	Currencies []string               `json:"cur"`
	Ext        openrtb_ext.RequestExt `json:"ext"`
	Imp        []openrtb2.Imp         `json:"imp"`
}

type SKAdUnit struct {
	Label   string                 `json:"label"`
	Sizes   [][]int64              `json:"sizes"`
	Bidders map[string]interface{} `json:"bidders"`
}
