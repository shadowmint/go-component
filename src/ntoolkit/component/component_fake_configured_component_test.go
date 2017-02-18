package component_test

import (
	"reflect"
	"ntoolkit/component"
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
	return component.SerializeState(&fake.Data)
}

func (fake *FakeConfiguredComponent) Deserialize(raw interface{}) error {
	var data FakeConfiguredComponentData
	if err := component.DeserializeState(&data, raw); err != nil {
		return err
	}
	fake.Data = data
	return nil
}