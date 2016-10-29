package teamsnap

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type nameValuePrompt struct {
	Name   string      `json:"name"`
	Value  interface{} `json:"value"`
	Prompt string      `json:"prompt"`
}

type relHrefData struct {
	Rel   string           `json:"rel"`
	HREF  string           `json:"href"`
	Data  nameValuePrompts `json:"data"`
	Links relHrefDatas     `json:"links"`
}

func (link relHrefData) findRelLink(rel string) (string, bool) {
	if link.Rel == rel {
		return link.HREF, true
	}

	return "", false
}

type relHrefDatas []relHrefData

func (links relHrefDatas) findRelLink(rel string) (string, bool) {

	for _, link := range links {
		if href, ok := link.findRelLink(rel); ok {
			return href, ok
		}
	}

	log.WithFields(log.Fields{"package": "teamsnap"}).Warnf("Unable to find Rel link '%s'", rel)
	return "", false
}

type teamSnapResult struct {
	Collection collection
}

type collection struct {
	Version  string
	Links    relHrefDatas
	Template relHrefData
	Queries  relHrefDatas
	Commands relHrefDatas
	Items    relHrefDatas
}

type nameValuePrompts []nameValuePrompt

type nameValueResults map[string]string

func (nvps nameValuePrompts) findValues(names ...string) (nameValueResults, bool) {
	results := make(map[string]string)

	for _, nvp := range nvps {
		for _, name := range names {
			if nvp.Name == name {
				if nvp.Value != nil {
					switch nvp.Value.(type) {
					case bool:
						if nvp.Value.(bool) {
							results[name] = "true"
						} else {
							results[name] = "false"
						}
					case string:
						results[name] = nvp.Value.(string)
					case json.Number:
						if val, err := nvp.Value.(json.Number).Int64(); err == nil {
							results[name] = strconv.FormatInt(val, 10)
						}
					case int:
						results[name] = strconv.FormatInt(int64(nvp.Value.(int)), 10)
					default:
						log.WithFields(log.Fields{"package": "teamsnap"}).Warnf("Unknown value: %v for name: %s", nvp.Value, name)
					}
				} else {
					results[name] = ""
				}
			}
		}
	}

	return results, len(results) == len(names)
}
