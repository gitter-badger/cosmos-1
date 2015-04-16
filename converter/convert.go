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

func ConvertFromContainerSeries(planet string, series []*influxdbc.Series) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	var regex *regexp.Regexp

	if planet == "" {
		regex = regexp.MustCompile("^(min|hour)?\\.?[^\\.]+\\.[^\\.]+\\.")
	} else {
		regex = regexp.MustCompile(fmt.Sprintf("^(min|hour)?\\.?[^\\.]+\\.%s\\.", planet))
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

func ConvertFromContainerInfoSeries(containerName string, series []*influxdbc.Series) map[string][][]interface{} {
	result := make(map[string][][]interface{})

	var regex *regexp.Regexp
	regex = regexp.MustCompile(fmt.Sprintf("^(min|hour)?\\.?[^\\.]+\\.[^\\.]+\\.%s\\.", containerName))

	for _, s := range series {
		comps := regex.Split(s.Name, -1)
		key := comps[1]

		result[key] = s.Points
	}
	return result
}

func ConvertFromPlanetSeries(series []*influxdbc.Series) []map[string]interface{} {
	result := make([]map[string]interface{}, len(series))
	for i, s := range series {
		result[i] = make(map[string]interface{})
		result[i]["Name"] = strings.Split(s.Name, ".")[1]
	}
	return result
}
