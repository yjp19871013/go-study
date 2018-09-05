package main

import (
	"fmt"
	"os"

	"github.com/tealeg/xlsx"
	"github.com/thedevsaddam/gojsonq"
)

const (
	DEVICE_JSON_FILE = "/home/yjp/go-projects/TestGolang/src/go-study/jsonq/device_list.txt"
)

func main() {
	jq := gojsonq.New().File(DEVICE_JSON_FILE).From("dList").Select("deviceName", "deviceId", "position")

	deviceInfoList, ok := jq.Get().([]interface{})
	if !ok {
		fmt.Println("Convert deviceInfoList error")
	}

	fmt.Println(deviceInfoList)

	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("Sheet 1")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sheet.AddRow().WriteSlice(&[]string{"设备名称", "设备ID", "位置"}, 3)
	for _, deviceInfo := range deviceInfoList {
		deviceInfoMap, ok := deviceInfo.(map[string]interface{})
		if !ok {
			fmt.Println("Convert deviceInfoMap error")
		}

		row := sheet.AddRow()
		row.AddCell().SetValue(deviceInfoMap["deviceName"])
		row.AddCell().SetValue(deviceInfoMap["deviceId"])
		row.AddCell().SetValue(deviceInfoMap["position"])
	}

	xlsxFile.Save("/home/yjp/result.xlsx")
}
