package sk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mxmCherry/openrtb/v15/openrtb2"
	"github.com/prebid/prebid-server/openrtb_ext"
)

func generateStoredRequestsFromJsonFile(jsonFilePath string) {
	jsonFile, err := os.Open(jsonFilePath)

	if err != nil {
		fmt.Println(fmt.Sprintf("Could not open json file '%s'", jsonFilePath))
		return
	}

	defer jsonFile.Close()

	var adUnits []SKAdUnit

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &adUnits)

	if err != nil {
		fmt.Println(fmt.Sprintf("Could not parse json file '%s'", jsonFilePath))
		return
	}

	for _, adUnit := range adUnits {
		fmt.Println(fmt.Sprintf("Processing '%s' Ad-Unit metadata", adUnit.Label))

		id := fmt.Sprintf("id-----%s", adUnit.Label)

		extRequest := openrtb_ext.RequestExt{}

		extRequest.SetPrebid(&openrtb_ext.ExtRequestPrebid{
			CurrencyConversions: &openrtb_ext.ExtRequestCurrency{
				ConversionRates: map[string]map[string]float64{
					"USD": {
						"INR": 75,
					},
				},
			},
			Targeting: &openrtb_ext.ExtRequestTargeting{
				PriceGranularity: openrtb_ext.PriceGranularity{
					Precision: 2,
					Ranges: []openrtb_ext.GranularityRange{
						{Min: 1, Max: 900, Increment: 1},
						{Min: 905, Max: 3150, Increment: 5},
					},
				},
			},
		})

		impId := fmt.Sprintf("request--%s", adUnit.Label)

		impBannerSizes := make([]openrtb2.Format, 0)
		for _, size := range adUnit.Sizes {
			impBannerSizes = append(impBannerSizes, openrtb2.Format{
				W: size[0],
				H: size[1],
			})
		}

		impExt, _ := json.Marshal(adUnit.Bidders)

		impRequest := openrtb2.Imp{
			ID: impId,
			Banner: &openrtb2.Banner{
				Format: impBannerSizes,
			},
			Ext: impExt,
		}

		storedRequest := StoredRequest{
			ID: id,
			Site: StoredRequestSite{
				Page: "sportskeeda.com",
			},
			Currencies: []string{
				"USD",
			},
			Ext: extRequest,
			Imp: []openrtb2.Imp{
				impRequest,
			},
		}

		jsonImpRequest, _ := json.MarshalIndent(impRequest, "", "\t")
		jsonStoredRequest, _ := json.MarshalIndent(storedRequest, "", "\t")

		fileName := fmt.Sprintf("%s.json", adUnit.Label)

		storedImpFile := filepath.Join("stored_requests", "data", "by_id", "stored_imps", fileName)
		storedReqFile := filepath.Join("stored_requests", "data", "by_id", "stored_requests", fileName)

		storedImpFilePath, _ := filepath.Abs(storedImpFile)
		storedReqFilePath, _ := filepath.Abs(storedReqFile)

		fmt.Println(fmt.Sprintf("Going to create '%s'", storedImpFile))

		err = ioutil.WriteFile(storedImpFilePath, jsonImpRequest, 0644)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not write json file '%s'", storedImpFilePath))
		}

		fmt.Println(fmt.Sprintf("Going to create '%s'", storedReqFile))

		err = ioutil.WriteFile(storedReqFilePath, jsonStoredRequest, 0644)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not write json file '%s'", storedReqFilePath))
		}
	}
}

func generateStoredRequestsForAmpAdUnits() {
	jsonFilePath, err := filepath.Abs(filepath.Join("sk", "amp_ad_units.json"))

	if err != nil {
		fmt.Println(fmt.Sprintf("Could not find amp_ad_units.json file in sk directory"))
		return
	}

	generateStoredRequestsFromJsonFile(jsonFilePath)
}

func GenerateStoredRequests() {
	generateStoredRequestsForAmpAdUnits()
}
