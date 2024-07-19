// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	ad_listing "section-02-02/clients/ad-listing"
)

func main() {

	// TODO #5 setup output for logger to write it to a file
	logger := log.Default()

	c := ad_listing.NewClient(ad_listing.BaseUrl, 3, logger)

	ads, err := c.GetAdByCate(context.TODO(), ad_listing.CatePty)
	if err != nil {
		panic("GetAdByCate " + err.Error())
	}

	logger.Printf("Number of Ads: %v", ads.Total)
	err = writeFile(ads.Ads)
	if err != nil {
		fmt.Println("Cannot write to file with err" + err.Error())
	}
}

func writeFile(fileData interface{}) error {
	file, err := os.Create("myLogFile.log")

	if err != nil {
		return err
	}

	defer file.Close()

	b, err := json.MarshalIndent(fileData, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}
	return nil
}
