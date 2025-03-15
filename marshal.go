package ghb

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Marshaler interface {
	MarshalGHB() (any, error)
}

func marshalBytes(msg any) ([]byte, error) {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("Wrong type %T, expected proto message", protoMsg)
	}
	response, err := marshalMessage(protoMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response body: %v", msg)
	}

	return json.Marshal(response)
}

func marshalMessage(msg proto.Message) (any, error) {
	if marshalable, ok := msg.ProtoReflect().Interface().(Marshaler); ok {
		marshaledMessage, err := marshalable.MarshalGHB()
		if err != nil {
			return nil, err
		}
		return marshaledMessage, nil
	}

	protoKeys, err := jsonToProtoKeys(msg)
	if err != nil {
		return nil, err
	}
	response := make(map[string]interface{})
	reflectedMessage := msg.ProtoReflect()
	fields := reflectedMessage.Descriptor().Fields()

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		name, ok := protoKeys[string(fd.Name())]
		if !ok {
			return nil, fmt.Errorf("key not found %v", fd.Name())
		}

		if fd.IsMap() {
			mapValue := make(map[string]interface{})
			if err := marshalMap(fd, reflectedMessage, mapValue); err != nil {
				return nil, err
			}
			response[name] = mapValue
		} else if fd.IsList() {
			listValue := make([]interface{}, 0)
			if err := marshalList(fd, reflectedMessage, listValue); err != nil {
				return nil, err
			}
			response[name] = listValue
		} else if fd.Kind() == protoreflect.MessageKind {
			nestedMsg := reflectedMessage.Get(fd).Message().Interface()
			if marshalable, ok := nestedMsg.(Marshaler); ok {
				marshaledMessage, err := marshalable.MarshalGHB()
				if err != nil {
					return nil, err
				}
				response[name] = marshaledMessage
			} else {
				nestedValue, err := marshalMessage(nestedMsg)
				if err != nil {
					return nil, err
				}
				response[name] = nestedValue
			}
		} else {
			response[name] = marshalField(fd, reflectedMessage)
		}
	}
	return response, nil
}

func marshalMap(fd protoreflect.FieldDescriptor, reflectedMessage protoreflect.Message, value map[string]interface{}) error {
	mp := reflectedMessage.Get(fd).Map()
	var mapError error
	mp.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		// different type of value
		if fd.MapValue().Kind() == protoreflect.MessageKind {
			val, err := marshalMessage(v.Message().Interface())
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
	return mapError
}

func marshalList(fd protoreflect.FieldDescriptor, reflectedMessage protoreflect.Message, listValue []interface{}) error {
	list := reflectedMessage.Get(fd).List()
	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)
		if fd.Kind() == protoreflect.MessageKind {
			nestedValue, err := marshalMessage(item.Message().Interface())
			if err != nil {
				return err
			}
			listValue = append(listValue, nestedValue)
		} else {
			// For primitive types
			listValue = append(listValue, valueToPrimitive(item, fd.Kind()))
		}
	}
	return nil
}

func marshalField(fd protoreflect.FieldDescriptor, reflectedMessage protoreflect.Message) interface{} {
	value := reflectedMessage.Get(fd)
	return valueToPrimitive(value, fd.Kind())
}

func valueToPrimitive(value protoreflect.Value, kind protoreflect.Kind) interface{} {
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
