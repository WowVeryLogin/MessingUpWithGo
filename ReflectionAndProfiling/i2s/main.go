package main

import (
	"errors"
	"reflect"
)

func i2sImpl(in reflect.Value, val reflect.Value) error {
	val = reflect.Indirect(val)
    if val.Kind() == reflect.Slice {
        if in.Kind() != reflect.Slice {
            return errors.New("Error")
        }

        val.Set(reflect.MakeSlice(val.Type(), in.Len(), in.Len()))
        for i := 0; i < in.Len(); i++ {
            err := i2sImpl(in.Index(i).Elem(), val.Index(i).Addr())
            if err != nil {
                return err
            }
        }
    } else {
        if in.Kind() != reflect.Map {
            return errors.New("Error")
        }

        for i := 0; i < val.NumField(); i++ {
            valueField := val.Field(i)
            typeField := val.Type().Field(i)
            dataField := in.MapIndex(reflect.ValueOf(typeField.Name)).Elem()
            switch valueField.Kind() {
                case reflect.Map:
                    for _, key := range valueField.MapKeys() {
                        i2sImpl(dataField, valueField.MapIndex(key).Addr())
                    }
                case reflect.Struct:
                    fallthrough
                case reflect.Slice:
                    err := i2sImpl(dataField, valueField.Addr())
                    if err != nil {
                        return err
                    }
                case reflect.String:
                    fallthrough
                case reflect.Bool:
                    if dataField.Type() != typeField.Type {
                        return errors.New("Error")
                    }
                    valueField.Set(dataField.Convert(typeField.Type))
                default:
                    if (dataField.Type().Name() == "string" || dataField.Type().Name() == "bool") {
                        return errors.New("Error")
                    }
                    valueField.Set(dataField.Convert(typeField.Type))
            }
        }
    }

    return nil
}

func i2s(data interface{}, out interface{}) error {
    val := reflect.ValueOf(out)
    if val.Kind() != reflect.Ptr {
        return errors.New("Error")
    }
    err := i2sImpl(reflect.ValueOf(data), val)
    if err != nil {
        return err
    }
    return nil
}
