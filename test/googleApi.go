package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

//Google.Analytics client API
type GoogleApi struct {
	ServiceUrl string
}

func (api *GoogleApi) parseStruct(prefix string, data map[string]interface{}) (values []string) {
	for k, v := range data {
		if k != "" && v != nil {
			t := reflect.ValueOf(v)
			switch t.Kind() {
			case reflect.Map:
				var iface map[string]interface{}
				r := new(bytes.Buffer)
				encoder := json.NewEncoder(r)
				encoder.Encode(t.Interface())
				decoder := json.NewDecoder(r)
				decoder.UseNumber()
				decoder.Decode(&iface)
				values = append(values, api.parseStruct(prefix+k, iface)...)
			case reflect.Slice:
				for i := 0; i < t.Len(); i++ {
					f := t.Index(i).Interface()
					ft := reflect.ValueOf(f)
					switch ft.Kind() {
					case reflect.Map:
						var iface map[string]interface{}
						r := new(bytes.Buffer)
						encoder := json.NewEncoder(r)
						encoder.Encode(ft.Interface())
						decoder := json.NewDecoder(r)
						decoder.UseNumber()
						decoder.Decode(&iface)
						values = append(values, api.parseStruct(prefix+k+fmt.Sprintf("%v", i+1), iface)...)
					default:
						d := fmt.Sprintf("%v", t)
						if d != "" && d != "0" {
							values = append(values, prefix+k+fmt.Sprintf("%v", i+1)+"="+d)
						}
					}
				}
			default:
				d := fmt.Sprintf("%v", t)
				if d != "" && d != "0" {
					values = append(values, prefix+k+"="+d)
				}
			}
		}
	}
	return values
}

func (api *GoogleApi) parseQuery(str string) string {
	u, _ := url.Parse(str)
	q := u.Query()
	u.RawQuery = q.Encode()
	ret := u.String()
	return strings.TrimLeft(ret, "?")
}

func (api *GoogleApi) Send(hit interface{}) error {
	var apidata []string
	var iface map[string]interface{}
	data := new(bytes.Buffer)
	encoder := json.NewEncoder(data)
	encoder.Encode(hit)
	decoder := json.NewDecoder(data)
	decoder.UseNumber()
	decoder.Decode(&iface)

	apidata = api.parseStruct("", iface)
	postdata := api.parseQuery("?" + strings.Join(apidata, "&"))
	cli := new(http.Client)
	req, err := http.NewRequest("POST", api.ServiceUrl, strings.NewReader(postdata))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "SomeOne")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := cli.Do(req)
	if err != nil {
		return err
	} else {
		//fmt.Println(config.ApiUrl)
		//fmt.Println(postdata)
		//fmt.Println(res.Status)
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

//New Google.Analytics client API
func NewApi(config string) *GoogleApi {
	p := GoogleApi{ServiceUrl: config}
	return &p
}
