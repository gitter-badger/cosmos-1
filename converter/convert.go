package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/cosmos-io/cosmos/model"
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
		comps := strings.Split(cont.Names[0], "/")
		cName := strings.Replace(comps[len(comps)-1], ".", "_", -1)
		base := MakeContainerSeriesName(token, planet, cName)

		var pathAndValues []*FieldPathAndValue
		pathAndValues = MakeFieldPathAndValue(cont, ".")

		for _, pv := range pathAndValues {
			if pv.Value != nil {
				series := influxdbc.NewSeries(fmt.Sprintf("CONTAINER.%s.%s", base, pv.Path), "txt_value", "num_value")
				t := reflect.TypeOf(pv.Value)
				if t.Kind() == reflect.String {
					series.AddPoint(pv.Value, nil)
				} else {
					series.AddPoint(nil, pv.Value)
				}
				result = append(result, series)
			}
		}
	}

	return result, nil
}

func ConvertFromContainerSeries(token, planet string, series []*influxdbc.Series) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	var regex *regexp.Regexp
	if planet == "" {
		regex = regexp.MustCompile(fmt.Sprintf("^(MIN|HOUR)?\\.?CONTAINER\\.%s\\.\\w+\\.", token))
	} else {
		regex = regexp.MustCompile(fmt.Sprintf("^(MIN|HOUR)?\\.?CONTAINER\\.%s\\.%s\\.", token, planet))
	}

	planetName := planet
	for _, s := range series {
		comps := regex.Split(s.Name, -1)
		containerName := strings.Split(comps[1], ".")[0]

		if planet == "" {
			comps = strings.Split(s.Name, fmt.Sprintf(".%s", containerName))
			comps = strings.Split(comps[0], ".")
			planetName = comps[len(comps)-1]
		}

		key := fmt.Sprintf("%s.%s", planetName, containerName)
		if result[key] == nil {
			result[key] = make(map[string]interface{})
		}
		comps = regexp.MustCompile(fmt.Sprintf(".*%s\\.", containerName)).Split(s.Name, -1)
		result[key][comps[1]] = s.Points[0]
	}

	return result
}

func ConvertFromContainerInfoSeries(token, planetName, containerName string, series []*influxdbc.Series) map[string]interface{} {
	result := make(map[string]interface{})

	var regex *regexp.Regexp
	regex = regexp.MustCompile(fmt.Sprintf("^(MIN|HOUR)?\\.?CONTAINER\\.%s\\.%s\\.%s\\.", token, planetName, containerName))

	for _, s := range series {
		comps := regex.Split(s.Name, -1)
		key := comps[1]

		result[key] = s.Points
	}
	return result
}

func ConvertFromPlanetSeries(token string, series []*influxdbc.Series) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	prefix := fmt.Sprintf("MIN.PLANET.%s.", token)

	for _, s := range series {
		name := strings.TrimPrefix(s.Name, prefix)
		planetName := strings.Split(name, ".")[0]
		key := strings.Replace(name, planetName+".", "", 1)

		if result[planetName] == nil {
			result[planetName] = make(map[string]interface{})
		}
		result[planetName][key] = s.Points
	}
	return result
}

func ConvertFromNewsFeedSeries(series []*influxdbc.Series) [][]interface{} {
	if len(series) > 0 {
		return series[0].Points
	} else {
		return [][]interface{}{}
	}
}
