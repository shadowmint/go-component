package component_test

import (
	"reflect"
	"ntoolkit/component"
	"ntoolkit/errors"
)

type FakeConfiguredComponentData struct {
	Items []FakeConfiguredComponentItem
}

type FakeConfiguredComponentItem struct {
	Id    string
	Count int
}

type FakeConfiguredComponent struct {
	Data FakeConfiguredComponentData
}

func (fake *FakeConfiguredComponent) Type() reflect.Type {
	return reflect.TypeOf(fake)
}

func (fake *FakeConfiguredComponent) New() component.Component {
	return &FakeConfiguredComponent{}
}

func (fake *FakeConfiguredComponent) Serialize() (interface{}, error) {
	return fake.Data, nil
}

func (fake *FakeConfiguredComponent) Deserialize(raw interface{}) error {
	var data FakeConfiguredComponentData
	if err := component.AsObject(&data, raw); err != nil {
		if errors.Is(err, component.ErrNullValue{}) {
			return nil
		} else {
			return err
		}
	}
	fake.Data = data
	return nil
}