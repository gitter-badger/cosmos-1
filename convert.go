package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/brannpark/cosmos/model"
	"github.com/cosmos-io/influxdbc"
)

type FieldPathAndValue struct {
	Path  string
	Value interface{}
}

func findFieldPathAndValue(obj interface{}, path string, pathDelimeter string, ret []*FieldPathAndValue) []*FieldPathAndValue {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	typeKind := objType.Kind()

	if typeKind == reflect.Ptr {
		if objValue.IsNil() == false {
			return findFieldPathAndValue(objValue.Elem().Interface(), path, pathDelimeter, ret)
		} else {
			return append(ret, &FieldPathAndValue{Path: path, Value: nil})
		}
	}

	if typeKind == reflect.Struct {
		fieldCnt := objType.NumField()
		for i := 0; i < fieldCnt; i++ {
			typeField := objType.Field(i)
			valueField := objValue.Field(i)

			var newPath string
			if path != "" {
				newPath = fmt.Sprintf("%s%s%s", path, pathDelimeter, typeField.Name)
			} else {
				newPath = typeField.Name
			}

			ret = findFieldPathAndValue(valueField.Interface(), newPath, pathDelimeter, ret)
		}
	} else if typeKind == reflect.Array || typeKind == reflect.Slice {
		sliceLen := objValue.Len()
		if sliceLen == 0 {
			return append(ret, &FieldPathAndValue{Path: path, Value: nil})
		}
		for i := 0; i < sliceLen; i++ {
			var newPath string
			if path != "" {
				newPath = fmt.Sprintf("%s%s%d", path, pathDelimeter, i)
			} else {
				newPath = fmt.Sprintf("%d", i)
			}
			ret = findFieldPathAndValue(objValue.Index(i).Interface(), newPath, pathDelimeter, ret)
		}
	} else {
		return append(ret, &FieldPathAndValue{Path: path, Value: objValue.Interface()})
	}

	return ret
}

func MakeFieldPathAndValue(obj interface{}, pathDelimeter string) []*FieldPathAndValue {
	return findFieldPathAndValue(obj, "", ".", make([]*FieldPathAndValue, 0))
}

func MakeContainerSeriesName(token, planet, containerId string) string {
	return fmt.Sprintf("%s.%s.%s", token, planet, containerId)
}

func ConvertToContainerSeries(token, planet string, body []byte) ([]*influxdbc.Series, error) {
	var containers []*model.Container
	err := json.Unmarshal(body, &containers)
	if err != nil {
		return nil, err
	}

	result := make([]*influxdbc.Series, 0)
	for _, cont := range containers {
		base := MakeContainerSeriesName(token, planet, *cont.Id)
		// series := influxdbc.NewSeries(base, "value")
		// series.AddPoint("")
		// result = append(result, series)

		var pathAndValues []*FieldPathAndValue
		pathAndValues = MakeFieldPathAndValue(cont, ".")

		for _, pv := range pathAndValues {
			if pv.Value != nil {
				series := influxdbc.NewSeries(fmt.Sprintf("%s.%s", base, pv.Path), "txt_value", "num_value")
				t := reflect.TypeOf(pv.Value)
				if t.Kind() == reflect.String {
					series.AddPoint(pv.Value, 0)
				} else {
					series.AddPoint("", pv.Value)
				}
				result = append(result, series)
			}
		}
	}

	return result, nil
}

func ConvertToPlanetSeries(token string, body []byte) (*influxdbc.Series, error) {
	var raw map[string]interface{}
	err := json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}

	cols := make([]string, len(raw))
	points := make([][]interface{}, 1)

	series := influxdbc.NewSeries(token)
	series.Columns = cols
	series.Points = points
	points[0] = make([]interface{}, len(raw))

	i := 0
	for k, v := range raw {
		cols[i] = k
		points[0][i] = v
		i = i + 1
	}

	return series, nil
}

func ConvertFromContainerSeries(planet string, series []*influxdbc.Series) map[string]map[string][]interface{} {
	result := make(map[string]map[string][]interface{})

	for _, s := range series {
		comps := regexp.MustCompile(fmt.Sprintf(".*%s\\.", planet)).Split(s.Name, -1)
		containerId := strings.Split(comps[1], ".")[0]
		if result[containerId] == nil {
			result[containerId] = make(map[string][]interface{})
		}
		comps = regexp.MustCompile(fmt.Sprintf(".*%s\\.", containerId)).Split(s.Name, -1)
		key := comps[1]
		result[containerId][key] = s.Points[0]
	}

	return result
}

func ConvertFromContainerInfoSeries(containerId string, series []*influxdbc.Series) map[string][][]interface{} {
	result := make(map[string][][]interface{})

	for _, s := range series {
		comps := regexp.MustCompile(fmt.Sprintf(".*%s\\.", containerId)).Split(s.Name, -1)
		fmt.Println(comps)
		key := comps[1]

		result[key] = s.Points
	}
	return result
}
