package xml
// #cgo LDFLAGS: -lxml2
// #cgo CFLAGS: -I/usr/include/libxml2
// #include <libxml/tree.h>
// #include <libxml/xmlmemory.h>
import "C"

import (
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type Encoder struct {}

func xmlChar(str string) *C.uchar {
	return (*C.uchar)(unsafe.Pointer(C.CString(str)))
}

func (e *Encoder) Marshal(o interface{}) ([]byte, error) {
	doc := C.xmlNewDoc(xmlChar("1.0"))

	t := reflect.TypeOf(o)
	val := reflect.ValueOf(o)
	name := nameFromStruct(val, t)
	root := nodeFromStruct(name, val, t, nil)
	
	C.xmlDocSetRootElement(doc, root)
	outBytes := serialize(doc)
	
	C.xmlFreeDoc(doc)
	C.xmlCleanupParser()
	C.xmlMemoryDump()
	return outBytes, nil
}

func serialize(doc C.xmlDocPtr) []byte {
	var xmlBuff *C.uchar
	defer C.free(unsafe.Pointer(xmlBuff));
	
	var bufferSize C.int
	C.xmlDocDumpFormatMemory(doc, &xmlBuff, &bufferSize, 1);
	return C.GoBytes(unsafe.Pointer(xmlBuff), bufferSize)
}

func nameFromStruct(o reflect.Value, t reflect.Type) string {
	if xmlName, ok := t.FieldByName("XMLName"); ok {
		return xmlName.Tag.Get("xml")
	}
	return t.Name()
}

func nameFromField(t reflect.StructField) string {
	tag := t.Tag.Get("xml")
	parts := strings.Split(tag, ",")
	name := t.Name
	if len(parts) > 1 {
		if len(parts[0]) > 0 {
			name = parts[0]
		}
	}
	return name
}

func nsFromStruct(t reflect.Type) string {
	if xmlName, ok := t.FieldByName("XMLNameSpace"); ok {
		return xmlName.Tag.Get("xml")
	}
	return ""
}

func nodeFromStruct(name string, o reflect.Value, t reflect.Type, root C.xmlNodePtr) C.xmlNodePtr {
	node := C.xmlNewNode(nil, xmlChar(name))
	if root == nil {
		root = node
	}

	if ns := nsFromStruct(t); ns != "" {
		parts := strings.Split(ns, "=")
		nsPtr := C.xmlNewNs(root, xmlChar(parts[1]), xmlChar(parts[0]))
		C.xmlSetNs(node, nsPtr)
	}
	
	for i := 0; i < o.NumField(); i++ {
		field := o.Field(i)
		fieldType := t.Field(i)
		
		tag := fieldType.Tag.Get("xml")
		if tag == "-" {
			continue
		}
		
		parts := strings.Split(tag, ",")
		name := nameFromField(fieldType)
		if len(parts) > 1 {
			value := stringValue(field)
			C.xmlNewProp(node, xmlChar(name), xmlChar(value))
			continue
		} else if len(tag) > 0 {
			name = tag
		}

		if fieldType.Name != "XMLName" && fieldType.Name != "XMLNameSpace" {
			child := nodeFromField(name, field, fieldType, root)
			C.xmlAddChild(node, child)
		}
	}

	return node
}

func nodeFromField(name string, field reflect.Value, t reflect.StructField, root C.xmlNodePtr) C.xmlNodePtr {
	var node C.xmlNodePtr
	switch field.Kind() {
	case reflect.Bool:
	case reflect.String:
		node = C.xmlNewNode(nil, xmlChar(name))
		child := C.xmlNewText(xmlChar(stringValue(field)))
		C.xmlAddChild(node, child)
	case reflect.Struct:
		node = nodeFromStruct(name, field, field.Type(), root)
	}
	return node
}

func stringValue(field reflect.Value) string {
	switch field.Kind() {
	case reflect.Bool: 
		return strconv.FormatBool(field.Bool())
	default:
		return field.String()
	}
}