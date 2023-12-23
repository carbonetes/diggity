package helper

import "encoding/xml"

func ToXML(data interface{}) ([]byte, error) {
	return xml.Marshal(data)
}
