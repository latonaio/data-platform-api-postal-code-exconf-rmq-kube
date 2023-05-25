package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-postal-code-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-postal-code-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"encoding/json"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	database "github.com/latonaio/golang-mysql-network-connector"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
	"golang.org/x/xerrors"
)

type ExistenceConf struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistenceConf(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistenceConf {
	return &ExistenceConf{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistenceConf) Conf(msg rabbitmq.RabbitmqMessage) interface{} {
	var ret interface{}
	ret = map[string]interface{}{
		"ExistenceConf": false,
	}
	input := make(map[string]interface{})
	err := json.Unmarshal(msg.Raw(), &input)
	if err != nil {
		return ret
	}

	_, ok := input["PostalCode"]
	if ok {
		input := &dpfm_api_input_reader.SDC{}
		err = json.Unmarshal(msg.Raw(), input)
		ret = e.confPostalCode(input)
		goto endProcess
	}

	err = xerrors.Errorf("can not get exconf check target")
endProcess:
	if err != nil {
		e.l.Error(err)
	}
	return ret
}

func (e *ExistenceConf) confPostalCode(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.PostalCode {
	exconf := dpfm_api_output_formatter.PostalCode{
		ExistenceConf: false,
	}
	if input.PostalCode.PostalCode == nil {
		return &exconf
	}
	if input.PostalCode.LocalRegion == nil {
		return &exconf
	}
	if input.PostalCode.Country == nil {
		return &exconf
	}
	exconf = dpfm_api_output_formatter.PostalCode{
		PostalCode:    *input.PostalCode.PostalCode,
		LocalRegion:   *input.PostalCode.LocalRegion,
		Country:       *input.PostalCode.Country,
		ExistenceConf: false,
	}

	rows, err := e.db.Query(
		`SELECT PostalCode 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_postal_code_postal_code_data 
		WHERE (PostalCode, LocalRegion, Country) = (?, ?, ?);`, exconf.PostalCode, exconf.LocalRegion, exconf.Country,
	)
	if err != nil {
		e.l.Error(err)
		return &exconf
	}
	defer rows.Close()

	exconf.ExistenceConf = rows.Next()
	return &exconf
}
