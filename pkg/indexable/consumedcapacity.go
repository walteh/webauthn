package indexable

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ConsumedCapacity struct {
	mutex  *sync.Mutex
	Reads  map[string]float64
	Writes map[string]float64
}

func NewConsumedCapacity() *ConsumedCapacity {
	return &ConsumedCapacity{
		Reads:  map[string]float64{},
		Writes: map[string]float64{},
		mutex:  &sync.Mutex{},
	}
}

func (cc *ConsumedCapacity) Apply(input *types.ConsumedCapacity) {
	if input == nil || input.TableName == nil {
		return
	}

	if input.ReadCapacityUnits == nil && input.Table != nil && input.Table.ReadCapacityUnits != nil {
		input.ReadCapacityUnits = input.Table.ReadCapacityUnits
	}

	if input.WriteCapacityUnits == nil && input.Table != nil && input.Table.WriteCapacityUnits != nil {
		input.WriteCapacityUnits = input.Table.WriteCapacityUnits
	}

	if input.WriteCapacityUnits == nil && input.ReadCapacityUnits == nil && input.CapacityUnits != nil {
		input.ReadCapacityUnits = input.CapacityUnits
	}

	cc.Read(input.TableName, input.ReadCapacityUnits)
	cc.Write(input.TableName, input.WriteCapacityUnits)

}

func (cc *ConsumedCapacity) ApplyArray(input ...types.ConsumedCapacity) {
	for _, v := range input {
		cc.Apply(&v)
	}
}

func (cc *ConsumedCapacity) Add(input *ConsumedCapacity) {
	input.Mutex().Lock()

	for k, v := range input.Reads {
		cc.Read(&k, &v)
	}

	for k, v := range input.Writes {
		cc.Write(&k, &v)
	}

	input.Mutex().Unlock()

}

func (cc *ConsumedCapacity) Read(k *string, f *float64) {
	if cc.Reads == nil {
		cc.Reads = map[string]float64{}
	}
	if k != nil && f != nil {
		cc.Mutex().Lock()

		cc.Reads[*k] += *f

		cc.Mutex().Unlock()
	}
}

func (cc *ConsumedCapacity) Write(k *string, f *float64) {
	if cc.Writes == nil {
		cc.Writes = map[string]float64{}
	}
	if k != nil && f != nil {
		cc.Mutex().Lock()

		cc.Writes[*k] += *f

		cc.Mutex().Unlock()
	}
}

func (cc *ConsumedCapacity) Mutex() *sync.Mutex {
	if cc.mutex == nil {
		cc.mutex = &sync.Mutex{}
	}
	return cc.mutex
}
