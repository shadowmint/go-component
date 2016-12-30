package component_test

import (
	"reflect"
	"ntoolkit/component"
	"fmt"
	"strings"
	"ntoolkit/errors"
	"strconv"
)

type FakeComponent struct {
	Id    string
	Count int
}

func (fake *FakeComponent) Type() reflect.Type {
	return reflect.TypeOf(fake)
}

func (fake *FakeComponent) Update(_ *component.Context) {
	fake.Count += 1
}

func (fake *FakeComponent) New() component.Component {
	return &FakeComponent{}
}

func (fake *FakeComponent) Serialize() (string, error) {
	return fmt.Sprintf("%s,%d", fake.Id, fake.Count), nil
}

func (fake *FakeComponent) Deserialize(data string) error {
	if len(data) > 0 {
		parts := strings.Split(data, ",")
		if len(parts) != 2 {
			return errors.Fail(component.ErrBadValue{}, nil, "Bad component data")
		}
		fake.Id = parts[0]
		count, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		fake.Count = count
	}
	return nil

}