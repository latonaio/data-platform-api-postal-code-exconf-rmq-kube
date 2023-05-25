package dpfm_api_input_reader

import (
	"data-platform-api-postal-code-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *SDC) ConvertToPostalCode() *requests.PostalCode {
	data := sdc.PostalCode
	return &requests.PostalCode{
		PostalCode:  data.PostalCode,
		LocalRegion: data.LocalRegion,
		Country:     data.Country,
	}
}
