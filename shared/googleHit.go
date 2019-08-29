package shared

import (
	"encoding/json"
	"github.com/monoculum/formam"
	"net/http"
)

type GoogleHit struct {
	ProtocolVersion string `json:"v,omitempty"`
	TrackingID      string `json:"tid,omitempty"`
	ClientID        string `json:"cid,omitempty"`
	UserID          string `json:"uid,omitempty"`
	HitType         string `json:"t,omitempty"`
	DocumentPath    string `json:"dp,omitempty"`
}

//ToJSON convert struct to Json byte array
func (hit *GoogleHit) ToJSON() []byte {
	output, err := json.Marshal(hit)
	if err != nil {
		panic(err)
	}
	return output
}

//FromJSON Fill in with data from byte[]
func (hit *GoogleHit) FromJSON(bytes []byte) {
	json.Unmarshal(bytes, hit)
}

//FromHTMLForm Fill in with data from Stream
func (hit *GoogleHit) FromHTMLForm(r *http.Request) {
	r.ParseForm()
	dec := formam.NewDecoder(&formam.DecoderOptions{TagName: "json"})
	err := dec.Decode(r.Form, hit)
	if err != nil {
		panic(err)
	}
}
