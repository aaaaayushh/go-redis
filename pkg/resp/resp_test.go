package resp

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestSerializeDeserialize(t *testing.T) {
	testCases := []struct {
		name  string
		value Value
	}{
		{
			name: "Simple String",
			value: Value{
				DataType: "STRING",
				str:      "Hello, World!",
			},
		},
		{
			name: "Integer",
			value: Value{
				DataType: "INTEGER",
				num:      42,
			},
		},
		{
			name: "Negative Integer",
			value: Value{
				DataType: "INTEGER",
				num:      -15,
			},
		},
		{
			name: "Bulk String",
			value: Value{
				DataType: "BULK",
				bulk:     "This is a bulk string",
			},
		},
		{
			name: "Error",
			value: Value{
				DataType: "ERROR",
				err:      "Error message",
			},
		},
		{
			name: "Null",
			value: Value{
				DataType: "NULL",
				isNull:   true,
			},
		},
		{
			name: "Array",
			value: Value{
				DataType: "ARRAY",
				array: []Value{
					{DataType: "STRING", str: "item1"},
					{DataType: "INTEGER", num: 2},
					{DataType: "BULK", bulk: "item3"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Serialize
			serialized := tc.value.Serialize()
			fmt.Println(string(serialized))
			// Deserialize
			buffer := bytes.NewBuffer(serialized)
			deserializer := NewDeserializer(buffer)
			deserialized, err := deserializer.Read()

			if err != nil {
				t.Fatalf("Error deserializing: %v", err)
			}

			// Compare original and deserialized values
			if !reflect.DeepEqual(tc.value, deserialized) {
				t.Errorf("Deserialized value does not match original.\nOriginal: %+v\nDeserialized: %+v", tc.value, deserialized)
			}
		})
	}
}

func TestDeserializerErrors(t *testing.T) {
	testCases := []struct {
		name        string
		input       []byte
		expectedErr string
	}{
		{
			name:        "Invalid data type",
			input:       []byte("X123\r\n"),
			expectedErr: "invalid data type",
		},
		{
			name:        "Incomplete input",
			input:       []byte("+Hello"),
			expectedErr: "EOF",
		},
		{
			name:        "Invalid integer",
			input:       []byte(":abc\r\n"),
			expectedErr: "strconv.Atoi: parsing \"abc\": invalid syntax",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(tc.input)
			deserializer := NewDeserializer(buffer)
			_, err := deserializer.Read()

			if err == nil {
				t.Fatal("Expected an error, but got nil")
			}

			if err.Error() != tc.expectedErr {
				t.Errorf("Expected error '%s', but got '%s'", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestSerializerMethods(t *testing.T) {
	testCases := []struct {
		name     string
		value    Value
		expected []byte
	}{
		{
			name:     "Serialize Simple String",
			value:    Value{DataType: "STRING", str: "Hello"},
			expected: []byte("+Hello\r\n"),
		},
		{
			name:     "Serialize Integer",
			value:    Value{DataType: "INTEGER", num: 42},
			expected: []byte(":+42\r\n"),
		},
		{
			name:     "Serialize Negative Integer",
			value:    Value{DataType: "INTEGER", num: -15},
			expected: []byte(":-15\r\n"),
		},
		{
			name:     "Serialize Bulk String",
			value:    Value{DataType: "BULK", bulk: "Hello, World!"},
			expected: []byte("$13\r\nHello, World!\r\n"),
		},
		{
			name:     "Serialize Error",
			value:    Value{DataType: "ERROR", err: "Error occurred"},
			expected: []byte("-Error occurred\r\n"),
		},
		{
			name:     "Serialize Null",
			value:    Value{DataType: "NULL"},
			expected: []byte("$-1\r\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.value.Serialize()
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestDeserializerMethods(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected Value
	}{
		{
			name:     "Deserialize Simple String",
			input:    []byte("+Hello\r\n"),
			expected: Value{DataType: "STRING", str: "Hello"},
		},
		{
			name:     "Deserialize Integer",
			input:    []byte(":42\r\n"),
			expected: Value{DataType: "INTEGER", num: 42},
		},
		{
			name:     "Deserialize Negative Integer",
			input:    []byte(":-15\r\n"),
			expected: Value{DataType: "INTEGER", num: -15},
		},
		{
			name:     "Deserialize Bulk String",
			input:    []byte("$13\r\nHello, World!\r\n"),
			expected: Value{DataType: "BULK", bulk: "Hello, World!"},
		},
		{
			name:     "Deserialize Error",
			input:    []byte("-Error occurred\r\n"),
			expected: Value{DataType: "ERROR", err: "Error occurred"},
		},
		{
			name:     "Deserialize Null",
			input:    []byte("$-1\r\n"),
			expected: Value{DataType: "NULL", isNull: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(tc.input)
			deserializer := NewDeserializer(buffer)
			result, err := deserializer.Read()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, but got %+v", tc.expected, result)
			}
		})
	}
}
