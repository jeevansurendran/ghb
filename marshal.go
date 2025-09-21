package ghb

import (
	"encoding/json"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Marshaler interface {
	MarshalGHB() (any, error)
}

type marshalOptions struct {
	timeFormat *timeFormat
}

func marshalBytes(msg any, options *marshalOptions) ([]byte, error) {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("wrong type %T, expected proto message", protoMsg)
	}
	response, err := marshalMessage(protoMsg, options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response body: %v", msg)
	}

	return json.Marshal(response)
}

func marshalMessage(msg proto.Message, options *marshalOptions) (any, error) {
	if isNil(msg) {
		return nil, nil
	}
	switch val := msg.ProtoReflect().Interface().(type) {
	case Marshaler:
		message, err := val.MarshalGHB()
		if err != nil {
			return nil, err
		}
		return message, nil
	case *timestamppb.Timestamp:
		message, err := options.timeFormat.marshal(val)
		if err != nil {
			return nil, err
		}
		return message, nil
	default:
		// Continue below.
	}

	protoKeys, err := jsonToProtoKeys(msg)
	if err != nil {
		return nil, err
	}
	response := make(map[string]any)
	reflectedMessage := msg.ProtoReflect()
	fields := reflectedMessage.Descriptor().Fields()

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		name, ok := protoKeys[string(fd.Name())]
		if !ok {
			return nil, fmt.Errorf("key not found %v", fd.Name())
		}
		if fd.IsMap() {
			mapValue, err := marshalMap(fd, reflectedMessage, options)
			if err != nil {
				return nil, err
			}
			response[name] = mapValue
		} else if fd.IsList() {
			listValue, err := marshalList(fd, reflectedMessage, options)
			if err != nil {
				return nil, err
			}
			response[name] = listValue
		} else if fd.Kind() == protoreflect.MessageKind {
			nestedMsg := reflectedMessage.Get(fd).Message().Interface()
			nestedValue, err := marshalMessage(nestedMsg, options)
			if err != nil {
				return nil, err
			}
			response[name] = nestedValue
		} else {
			response[name] = marshalField(fd, reflectedMessage)
		}
		// TODO: handle required vs optional fields,
		// for now remove all the fields that are nil
		if response[name] == nil {
			delete(response, name)
		}
	}
	return response, nil
}

func marshalMap(fd protoreflect.FieldDescriptor, reflectedMessage protoreflect.Message, options *marshalOptions) (map[string]any, error) {
	mp := reflectedMessage.Get(fd).Map()
	value := make(map[string]any, mp.Len())
	var mapError error
	mp.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		// different type of value
		if fd.MapValue().Kind() == protoreflect.MessageKind {
			val, err := marshalMessage(v.Message().Interface(), options)
			if err != nil {
				mapError = err
				return false
			}
			value[k.String()] = val
		} else {
			value[k.String()] = valueToPrimitive(v, fd.MapValue().Kind())
		}
		return true
	})
	return value, mapError
}

func marshalList(fd protoreflect.FieldDescriptor, reflectedMessage protoreflect.Message, options *marshalOptions) ([]any, error) {
	list := reflectedMessage.Get(fd).List()
	listValue := make([]any, list.Len())
	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)
		if fd.Kind() == protoreflect.MessageKind {
			nestedValue, err := marshalMessage(item.Message().Interface(), options)
			if err != nil {
				return nil, err
			}
			listValue[i] = nestedValue
		} else {
			// For primitive types
			listValue[i] = valueToPrimitive(item, fd.Kind())
		}
	}
	return listValue, nil
}

func isNil(msg proto.Message) bool {
	return msg == nil || reflect.ValueOf(msg).IsNil()
}

func marshalField(fd protoreflect.FieldDescriptor, reflectedMessage protoreflect.Message) any {
	value := reflectedMessage.Get(fd)
	return valueToPrimitive(value, fd.Kind())
}

func valueToPrimitive(value protoreflect.Value, kind protoreflect.Kind) any {
	switch kind {
	case protoreflect.BoolKind:
		return value.Bool()
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return int32(value.Int())
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return value.Int()
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return uint32(value.Uint())
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return value.Uint()
	case protoreflect.FloatKind:
		return float32(value.Float())
	case protoreflect.DoubleKind:
		return value.Float()
	case protoreflect.StringKind:
		return value.String()
	case protoreflect.BytesKind:
		return value.Bytes()
	case protoreflect.EnumKind:
		return int32(value.Enum())
	default:
		return nil
	}
}
