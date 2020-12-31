package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	// define the url
	url := fmt.Sprintf("%s", "https://wcc.sc.egov.usda.gov/awdbWebService/services")

	// prep payload
	// payload := []byte(strings.TrimSpace(`
	// <soapenv:Envelope
	//    xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
	//    xmlns:ser="http://www.wcc.nrcs.usda.gov/ns/awdbWebService"
	// >
	//     <soapenv:Body>
	//         <ser:getAllForecastsForStation>
	// 			<ser:stationTriplet xmlns="">[string]</ser:stationTriplet>
	// 			<ser:beginPublicationDate xmlns="">[string]</ser:beginPublicationDate>
	// 			<ser:endPublicationDate xmlns="">[string]</ser:endPublicationDate>
	//         </ser:getAllForecastsForStation>
	//     </soapenv:Body>
	// </soapenv:Envelope>`,
	// ))

	// payload := []byte(strings.TrimSpace(`
	// <soapenv:Envelope
	//    xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
	//    xmlns:ser="http://www.wcc.nrcs.usda.gov/ns/awdbWebService"
	// >
	//     <soapenv:Body>
	//         <ser:getData>
	// 			<ser:stationTriplets xmlns="">471:ID:SNTL</ser:stationTriplets>
	// 			<ser:elementCd xmlns="">WTEQ</ser:elementCd>
	// 			<ser:ordinal xmlns="">1</ser:ordinal>
	// 			<ser:duration xmlns="">DAILY</ser:duration>
	// 			<ser:getFlags xmlns="">true</ser:getFlags>
	// 			<ser:beginDate xmlns="">2020-11-29 00:00:00</ser:beginDate>
	// 			<ser:endDate xmlns="">2020-12-02 00:00:00</ser:endDate>
	// 			<ser:alwaysReturnDailyFeb29 xmlns="">true</ser:alwaysReturnDailyFeb29>
	//         </ser:getData>
	//     </soapenv:Body>
	// </soapenv:Envelope>`,
	// ))

	// getDataPayload := []byte(strings.TrimSpace(`
	// <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:q0="http://www.wcc.nrcs.usda.gov/ns/awdbWebService" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	// 	<SOAP-ENV:Body>
	// 		<q0:getData>
	//   			<stationTriplets>471:ID:SNTL</stationTriplets>
	//  		 	<elementCd>WTEQ</elementCd>
	//   			<ordinal>1</ordinal>
	//   			<duration>DAILY</duration>
	//   			<getFlags>true</getFlags>
	//   			<beginDate>2020-11-29</beginDate>
	//   			<endDate>2020-12-01</endDate>
	//   			<alwaysReturnDailyFeb29>true</alwaysReturnDailyFeb29>
	// 		</q0:getData>
	// 	</SOAP-ENV:Body>
	// </SOAP-ENV:Envelope>`,
	// ))

	getStationsPayload := []byte(strings.TrimSpace(`
	<?xml version="1.0" encoding="UTF-8"?>
	<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:q0="http://www.wcc.nrcs.usda.gov/ns/awdbWebService" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<SOAP-ENV:Body>
			<q0:getStations>
	  			<networkCds>SNTL</networkCds>
	  			<elementCds>WTEQ</elementCds>
	  			<ordinals>1</ordinals>
	  			<logicalAnd>true</logicalAnd>
			</q0:getStations>
		</SOAP-ENV:Body>
	</SOAP-ENV:Envelope>
	`))

	// soapAction := "urn:getData" // The format is `urn:<soap_action>`

	httpMethod := "POST"

	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(getStationsPayload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}

	req.Header.Set("Content-type", "text/xml")
	// req.Header.Set("SOAPAction", soapAction)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return
	}

	defer res.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// //son.MarshalIndent()
	// bodyString := string(bodyBytes)
	// fmt.Println("body", bodyString)

	//getDataResult := new(GetData)
	getStationsResult := new(GetStations)

	err = xml.NewDecoder(res.Body).Decode(getStationsResult)
	if err != nil {
		log.Fatal("Error on unmarshaling. ", err.Error())
		return
	}

	r, err := json.MarshalIndent(getStationsResult, "", "  ")
	if err != nil {
		fmt.Println("err marshal intdent: ", err)
		return
	}

	fmt.Println("result: ", string(r))

}

type GetStations struct {
	// Envelope struct {
	Body struct {
		// XMLName         xml.Name
		GetStationsResponse struct {
			// XMLName xml.Name
			Return []string `xml:"return"`
		} `xml:"getStationsResponse"`
	} //`xml:"body"`
	// } `xml:"Envelope"`

}

type GetData struct {
	// Envelope struct {
	Body struct {
		// XMLName         xml.Name
		GetDataResponse struct {
			// XMLName xml.Name
			Return struct {
				BeginDate      string    `xml:"beginDate"`
				Duration       string    `xml:"duration"`
				EndDate        string    `xml:"endDate"`
				Flags          []string  `xml:"flags"`
				StationTriplet string    `xml:"stationTriplet"`
				Values         []float32 `xml:"values"`
			} `xml:"return"`
		} `xml:"getDataResponse"`
	} //`xml:"body"`
	// } `xml:"Envelope"`
}

// <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
//     <Body>
//         <getAllForecastsForStation xmlns="http://www.wcc.nrcs.usda.gov/ns/awdbWebService">
//             <stationTriplet xmlns="">[string]</stationTriplet>
//             <beginPublicationDate xmlns="">[string]</beginPublicationDate>
//             <endPublicationDate xmlns="">[string]</endPublicationDate>
//         </getAllForecastsForStation>
//     </Body>
// </Envelope>
