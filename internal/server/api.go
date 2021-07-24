package server

import (
	"../core/xauth"
	"../core/xdatabase"
	"../core/xjson"
	structPerModel "../core/xstruct/model_related"
	discoveryManager "../node_auto_discovery/discovery_manager"
	"../server/session"
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ApiAllData(ctx *fasthttp.RequestCtx) {
	fmt.Println(GetRequestText(ctx))
	CheckThenRun(ctx, "text/json", func() {
		content := ""
		structName := string(ctx.Request.Header.Peek("StructName"))

		// get interfaces and keys by the type of json model
		values, keys, err := xdatabase.ObjectGetAll(structName)
		if err != nil {
			content = JsonError(err.Error())
		} else {
			// create json from the result or just an JsonError
			content += "{"
			content += "\"err\":\"\","
			content += "\"rows\":["
			for i, v := range values {
				content += xjson.StructToJson(v)
				if i != len(values)-1 {
					content += ","
				}
			}
			content += "]"
			content += ",\"keys\":["
			for i := range values {
				content += fmt.Sprintf("\"%s\"", keys[i])
				if i != len(values)-1 {
					content += ","
				}
			}
			content += "]"
			content += "}"
		}
		WriteContent(ctx, content)
	})
}

func ApiAddData(ctx *fasthttp.RequestCtx) {
	fmt.Println(GetRequestText(ctx))
	CheckThenRun(ctx, "text/json", func() {

		// create interface for store in database and reflect.Value to it
		var item interface{}
		itemValue := reflect.ValueOf(&item).Elem()
		structName := string(ctx.Request.Header.Peek("StructName"))

		// temporary reflect.Value from structure
		ivTemp := structPerModel.GetValue(structName)
		ivTempType := ivTemp.Type()

		// fill the temporary reflect.Value with json from headers
		for i := 0; i < ivTemp.NumField() && ivTemp.Field(i).CanInterface(); i++ {
			fieldName := ivTempType.Field(i).Name
			fieldKind := ivTempType.Field(i).Type.Kind()

			header := string(ctx.Request.Header.Peek("modalField_" + fieldName))
			if header != "" {
				if fieldKind == reflect.Int ||
					fieldKind == reflect.Int64 {
					n, err := strconv.ParseInt(header, 10, 64)
					if err != nil {
						WriteContent(ctx, JsonError(err.Error()))
						return
					}
					ivTemp.FieldByName(fieldName).SetInt(n)

				} else if fieldKind == reflect.Uint {
					n, err := strconv.ParseInt(header, 10, 64)
					if err != nil {
						WriteContent(ctx, JsonError(err.Error()))
						return
					}
					ivTemp.FieldByName(fieldName).SetUint(uint64(n))

				} else if fieldKind == reflect.String {
					ivTemp.FieldByName(fieldName).SetString(header)
				}
			}
		}

		switch structName {
		case "Discovery":
			ivTemp.FieldByName("Id").SetString(uuid.NewString())
		}

		// put temporary reflect.Value in the reflect.Value of the interface
		itemValue.Set(ivTemp)

		// try store the interface json model
		err := xdatabase.ObjectStore(item)
		if err != nil {
			WriteContent(ctx, JsonError(err.Error()))
			return
		}

		// exit safe
		WriteContent(ctx, JsonError(""))
	})
}

func ApiDeleteData(ctx *fasthttp.RequestCtx) {
	fmt.Println(GetRequestText(ctx))
	CheckThenRun(ctx, "text/json", func() {

		keys := string(ctx.Request.Header.Peek("Keys2Delete"))
		keyList := strings.Split(keys, ",")

		for _, key := range keyList {
			err := xdatabase.ObjectDeleteByKey(key)
			if err != nil {
				WriteContent(ctx, JsonError(err.Error()))
				return
			}
		}

		// exit safe
		WriteContent(ctx, JsonError(""))
	})
}

func ApiEditData(ctx *fasthttp.RequestCtx) {
	CheckThenRun(ctx, "text/json", func() {
		fmt.Println(GetRequestText(ctx))

		// create interface for store in database and reflect.Value to it
		var item interface{}
		itemValue := reflect.ValueOf(&item).Elem()
		structName := string(ctx.Request.Header.Peek("StructName"))
		key := string(ctx.Request.Header.Peek("Key2Update"))

		// create temporary reflect.Value to an actual type of struct object
		ivTemp := structPerModel.GetValue(structName)
		ivTempType := ivTemp.Type()

		// fill the temporary reflect.Value with json from headers
		for i := 0; i < ivTemp.NumField() && ivTemp.Field(i).CanInterface(); i++ {
			fieldName := ivTempType.Field(i).Name
			fieldKind := ivTempType.Field(i).Type.Kind()

			header := string(ctx.Request.Header.Peek("modalField_" + fieldName))
			if header != "" {
				if fieldKind == reflect.Int ||
					fieldKind == reflect.Int64 {
					n, err := strconv.ParseInt(header, 10, 64)
					if err != nil {
						WriteContent(ctx, JsonError(err.Error()))
						return
					}
					ivTemp.FieldByName(fieldName).SetInt(n)

				} else if fieldKind == reflect.Uint {
					n, err := strconv.ParseInt(header, 10, 64)
					if err != nil {
						WriteContent(ctx, JsonError(err.Error()))
						return
					}
					ivTemp.FieldByName(fieldName).SetUint(uint64(n))

				} else if fieldKind == reflect.String {
					ivTemp.FieldByName(fieldName).SetString(header)
				}
			}
		}
		// put temporary reflect.Value in the reflect.Value of the interface
		itemValue.Set(ivTemp)

		err := xdatabase.ObjectUpdateOnKey(item, key)
		if err != nil {
			WriteContent(ctx, JsonError(err.Error()))
			return
		}

		// exit safe
		WriteContent(ctx, JsonError(""))
	})

	//structName := fmt.Sprintf("%s", ctx.UserValue("structName"))
	//var iv reflect.Value
	//headerId2Edit := string(ctx.Request.Header.Peek("id2edit"))
	//key := ""
	//
	//// pick the model to get the reflect.Value of it
	//switch structName {
	//case xjson.GetStructNameOfInterface(mnc.SnmpCredential{}):
	//	item := mnc.SnmpCredential{
	//		Version:          0,
	//		Network:          "",
	//		Timeout:          0,
	//		Retries:          0,
	//		MessageMaxSize:   0,
	//		Community:        "",
	//		UserName:         "",
	//		SecurityLevel:    0,
	//		AuthPassword:     "",
	//		AuthProtocol:     "",
	//		PrivPassword:     "",
	//		PrivProtocol:     "",
	//		SecurityEngineId: "",
	//		ContextEngineId:  "",
	//		ContextName:      "",
	//	}
	//	iv = reflect.ValueOf(&item).Elem()
	//	key = fmt.Sprintf("%s:%s", mnc.SnmpCredentialIndexName, headerId2Edit)
	//default:
	//	WriteContent(ctx, JsonError("json model not found"))
	//	return
	//}
	//
	//// put json from headers into reflect.Value
	//ivType := iv.Type()
	//for i := 0; i < iv.NumField() && iv.Field(i).CanInterface(); i++ {
	//	fieldName := ivType.Field(i).Name
	//	fieldKind := ivType.Field(i).Type.Kind()
	//
	//	header := string(ctx.Request.Header.Peek("modalField_" + fieldName))
	//	if header != "" {
	//		if fieldKind == reflect.Int ||
	//			fieldKind == reflect.Int64 {
	//			n, err := strconv.ParseInt(header, 10, 64)
	//			if err != nil {
	//				WriteContent(ctx, JsonError(err.Error()))
	//				return
	//			}
	//			iv.FieldByName(fieldName).SetInt(n)
	//
	//		} else if fieldKind == reflect.Uint {
	//			n, err := strconv.ParseInt(header, 10, 64)
	//			if err != nil {
	//				WriteContent(ctx, JsonError(err.Error()))
	//				return
	//			}
	//			iv.FieldByName(fieldName).SetUint(uint64(n))
	//
	//		} else if fieldKind == reflect.String {
	//			iv.FieldByName(fieldName).SetString(header)
	//		}
	//	}
	//}
	//
	//// get the model from interface of reflect.Value
	//switch structName {
	//case xjson.GetStructNameOfInterface(mnc.SnmpCredential{}):
	//	if val, ok := iv.Interface().(mnc.SnmpCredential); ok {
	//		err := val.ObjectUpdateOnKey(key)
	//		if err != nil {
	//			WriteContent(ctx, JsonError(err.Error()))
	//			return
	//		}
	//	} else {
	//		WriteContent(ctx, JsonError("struct type is wrong"))
	//		return
	//	}
	//}
	//
	//// exit safe
	//WriteContent(ctx, JsonError(""))
}

func ApiOneData(ctx *fasthttp.RequestCtx) {
	CheckThenRun(ctx, "text/json", func() {
		fmt.Println(GetRequestText(ctx))

		key := string(ctx.Request.Header.Peek("Key2Select"))
		content := ""

		value, err := xdatabase.ObjectGetByKey(key)
		if err != nil {
			WriteContent(ctx, JsonError(err.Error()))
			return
		} else {
			content += "{\"err\":\"\","
			content += fmt.Sprintf("\"row\":%s,", xjson.StructToJson(value))
			content += fmt.Sprintf("\"key\":\"%s\"", key)
			content += "}"
		}
		WriteContent(ctx, content)
	})
}

func ApiGetFieldsAll(ctx *fasthttp.RequestCtx) {
	CheckThenRun(ctx, "text/json", func() {
		fmt.Println(GetRequestText(ctx))

		structName := string(ctx.Request.Header.Peek("StructName"))

		fieldNames, _ := structPerModel.GetStructFields(structName)

		fieldsJson := "{\"fieldNames\":["
		for i, f := range fieldNames {
			if i != 0 {
				fieldsJson += ","
			}
			fieldsJson += fmt.Sprintf("\"%s\"", f)
		}
		fieldsJson += "]}"

		WriteContent(ctx, fieldsJson)
	})
}

func ApiGetFieldsEdit(ctx *fasthttp.RequestCtx) {
	CheckThenRun(ctx, "text/json", func() {
		fmt.Println(GetRequestText(ctx))

		structName := string(ctx.Request.Header.Peek("StructName"))

		editableFields := make([]string, 0)
		defaultValues := make([]string, 0)
		switch structName {
		case "SnmpCredential":
			editableFields = []string{"Version", "Network", "Timeout", "Retries", "MessageMaxSize", "Community", "UserName", "SecurityLevel", "AuthPassword", "AuthProtocol", "PrivPassword", "PrivProtocol", "SecurityEngineId", "ContextEngineId", "ContextName"}
			defaultValues = []string{"0", "udp", "5000000000", "0", "1400", "", "", "0", "", "SHA", "", "DES", "", "", ""}
		case "WindowsCredential":
			editableFields = []string{"Username", "Password", "AuthType"}
			defaultValues = []string{"", "", "0"}
		case "LinuxCredential":
			editableFields = []string{"Username", "Password", "AuthType"}
			defaultValues = []string{"", "", "0"}
		case "DeviceCredential":
			editableFields = []string{"Username", "Password", "AuthType"}
			defaultValues = []string{"", "", "0"}
		case "Node":
			editableFields = []string{"Ip", "AliveType", "AliveInterval", "PingPacketSize", "PingTimeout"}
			defaultValues = []string{"", "0", "10000000000", "56", "100000000"}
		case "Discovery":
			editableFields = []string{"Ranges", "DiscoveryType", "SnmpCredentialKey", "SnmpPort", "PingPacketSize"}
			defaultValues = []string{"", "0", "", "161", "56"}
		}

		//fieldNames, _ := structPerModel.GetStructFields(structName)

		fieldsJson := "{\"fieldNames\":["
		for i, f := range editableFields {
			if i != 0 {
				fieldsJson += ","
			}
			fieldsJson += fmt.Sprintf("\"%s\"", f)
		}
		fieldsJson += "],"

		fieldsJson += "\"fieldDefaults\":["
		for i, d := range defaultValues {
			if i != 0 {
				fieldsJson += ","
			}
			fieldsJson += fmt.Sprintf("\"%s\"", d)
		}
		fieldsJson += "]}"

		WriteContent(ctx, fieldsJson)
	})
}

func ApiTerminate(ctx *fasthttp.RequestCtx) {
	CheckThenRun(ctx, "text/html", func() {
		fmt.Println(GetRequestText(ctx))

		WriteContent(ctx, "Monitoring Terminated Gracefully.")
	})
	go func() {
		time.Sleep(300 * time.Millisecond)
		os.Exit(0)
	}()
}

func ApiDiscoveryStart(ctx *fasthttp.RequestCtx) {
	CheckThenRun(ctx, "text/json", func() {
		fmt.Println(GetRequestText(ctx))

		discoveryKey := string(ctx.Request.Header.Peek("DiscoveryKey"))
		discoveryManager.DiscoveryAdd(discoveryKey)
		discoveryManager.DiscoveryGetOne(discoveryKey).Start()
	})
}

func ApiAuthenticate(ctx *fasthttp.RequestCtx) {
	CheckThenRun(ctx, "text/json", func() {
		fmt.Println(GetRequestText(ctx))

		username := string(ctx.Request.Header.Peek("Username"))
		password := string(ctx.Request.Header.Peek("Password"))

		auth := xauth.Authenticate(username, password)
		if auth {
			store := session.FetchSessionStore(ctx)
			store.Set("is_auth", true)
			r := "{\"is_auth\":true}"
			WriteContent(ctx, r)
		} else {
			r := "{\"is_auth\":false}"
			WriteContent(ctx, r)
		}
	})
}

func ApiParamGetLast(ctx *fasthttp.RequestCtx) {
	CheckThenRun(ctx, "text/json", func() {
		fmt.Println(GetRequestText(ctx))

		sensorName := string(ctx.Request.Header.Peek("SensorName"))
		nodeKey := string(ctx.Request.Header.Peek("NodeKey"))

		value, timestamp, err := xdatabase.ParamGetLast(sensorName, nodeKey)
		if err == nil {
			WriteContent(ctx, fmt.Sprintf("{\"timestamp\":%d,\"value\":\"%s\"}", timestamp, value))
		} else {
			WriteContent(ctx, JsonError(err.Error()))
		}
	})
}

func ApiParamGetAll(ctx *fasthttp.RequestCtx) {
	CheckThenRun(ctx, "text/json", func() {
		fmt.Println(GetRequestText(ctx))

		sensorName := string(ctx.Request.Header.Peek("SensorName"))
		nodeKey := string(ctx.Request.Header.Peek("NodeKey"))

		values, timestamps := xdatabase.ParamGetAll(sensorName, nodeKey)
		if len(values) != 0 {
			timestampList := "["
			valueList := "["
			for i := range values {
				timestampList += fmt.Sprint(timestamps[i])
				valueList += fmt.Sprintf("\"%s\"", values[i])
				if i != len(values) - 1 {
					timestampList += ","
					valueList += ","
				}
			}
			timestampList += "]"
			valueList += "]"
			WriteContent(ctx, fmt.Sprintf("{\"timestamps\":%s,\"values\":%s}", timestampList, valueList))
		} else {
			WriteContent(ctx, JsonError("no data"))
		}
	})
}
