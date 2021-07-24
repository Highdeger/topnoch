package xmonitor

import (
	model "../../model/snmp_result"
	"fmt"
	"github.com/k-sone/snmpgo"
	"math"
	"strconv"
)

var (
	OsiLayers = []string{"Physical", "Data Link", "Network", "Transport", "Session", "Presentation", "Application"}
)

// GetValuesForOids is the general snmp-oids get function.
// Result is a slice of *model.SnmpResult.
// Panic has to be handled.
func GetValuesForOids(oidsSlice []string, args snmpgo.SNMPArguments) []*model.SnmpResult {

	snmp, err := snmpgo.NewSNMP(args)
	if err != nil {
		panic("can't create snmp object")
	}

	oids, err := snmpgo.NewOids(oidsSlice)
	if err != nil {
		panic("can't create oid objects")
	}

	if err = snmp.Open(); err != nil {
		panic("can't connect to snmp object")
	}
	defer snmp.Close()

	resultRaw, e := snmp.GetRequest(oids)
	if e != nil {
		return []*model.SnmpResult{}
		//panic("can't get the result for oid objects")
	}
	if resultRaw.ErrorStatus() != snmpgo.NoError {
		panic(fmt.Sprintf("error outside of this context (ErrorStatus: %s - ErrorIndex: %d)", resultRaw.ErrorStatus(), resultRaw.ErrorIndex()))
	}

	results := make([]*model.SnmpResult, 0)
	var v *snmpgo.VarBind
	for i := range resultRaw.VarBinds() {
		v = resultRaw.VarBinds()[i]
		results = append(results, &model.SnmpResult{
			Oid: v.Oid.String(),
			Typ: v.Variable.Type(),
			Val: v.Variable.String(),
		})
	}
	return results
}

// GetTableForOid is the general snmp-table get function.
// Result is a slice of []*model.SnmpResult, every element is a row with columns.
// Panic has to be handled.
func GetTableForOid(tableOid string, tableEntry int, tableColumns []int, args snmpgo.SNMPArguments) [][]*model.SnmpResult {
	result := make([][]*model.SnmpResult, 0)
	entryOid := tableOid + fmt.Sprintf(".%d", tableEntry)
rowIterator:
	for row := 1; true; row++ {
		values := make([]*model.SnmpResult, 0)
		for _, col := range tableColumns {
			r := GetValuesForOids([]string{fmt.Sprintf("%s.%d.%d", entryOid, col, row)}, args)
			if (r[0].Typ == "NoSucheInstance") || (r[0].Typ == "NoSucheObject") {
				break rowIterator
			} else {
				values = append(values, r[0])
			}
		}
		result = append(result, values)
	}
	return result
}

// WalkValuesForOids is the general snmp-oids walk function.
// Result is a slice of *model.SnmpResult.
// Panic has to be handled.
func WalkValuesForOids(oidsSlice []string, args snmpgo.SNMPArguments) []*model.SnmpResult {

	snmp, err := snmpgo.NewSNMP(args)
	if err != nil {
		panic("can't create snmp object")
	}

	oids, err := snmpgo.NewOids(oidsSlice)
	if err != nil {
		panic("can't create oid objects")
	}

	if err = snmp.Open(); err != nil {
		panic("can't connect to snmp object")
	}
	defer snmp.Close()

	resultRaw, e := snmp.GetBulkWalk(oids, 0, 1)
	if e != nil {
		panic("can't get the bulk result for oid objects")
	}
	if resultRaw.ErrorStatus() != snmpgo.NoError {
		panic(fmt.Sprintf("error outside of this context (ErrorStatus: %s - ErrorIndex: %d)", resultRaw.ErrorStatus(), resultRaw.ErrorIndex()))
	}

	results := make([]*model.SnmpResult, 0)
	var v *snmpgo.VarBind
	for i := range resultRaw.VarBinds() {
		v = resultRaw.VarBinds()[i]
		results = append(results, &model.SnmpResult{
			Oid: v.Oid.String(),
			Typ: v.Variable.Type(),
			Val: v.Variable.String(),
		})
	}
	return results
}

// FetchSysDescr fetches the description of a snmp client.
// Result is an octet-string (0-255 characters).
// Panic has to be handled.
func FetchSysDescr(args snmpgo.SNMPArguments) string {
	arr := []string{"1.3.6.1.2.1.1.1.0"}
	r := GetValuesForOids(arr, args)
	if len(r) == 0 {
		return ""
	} else {
		return r[0].Val
	}
}

// FetchSysObjectId fetches the system's object id of a snmp client.
// Result is a string representing an oid.
// Panic has to be handled.
func FetchSysObjectId(args snmpgo.SNMPArguments) string {
	arr := []string{"1.3.6.1.2.1.1.2.0"}
	r := GetValuesForOids(arr, args)
	return r[0].Val
}

// FetchSysUptime fetches the uptime of a snmp client.
// Result is an int in units of centiseconds.
// Panic has to be handled.
func FetchSysUptime(args snmpgo.SNMPArguments) int {
	arr := []string{"1.3.6.1.2.1.1.3.0"}
	raw := GetValuesForOids(arr, args)
	r, e := strconv.Atoi(raw[0].Val)
	if e != nil {
		panic("can't extract int from SysUptime: " + e.Error())
	}
	return r
}

// FetchSysContactNameLocation fetches the contact, name and location of a snmp client.
// Results are strings containing contact info, system's name and location, respectively.
// Panic has to be handled.
func FetchSysContactNameLocation(args snmpgo.SNMPArguments) (contact string, name string, location string) {
	arr := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.5.0", "1.3.6.1.2.1.1.6.0"}
	r := GetValuesForOids(arr, args)
	return r[0].Val, r[1].Val, r[2].Val
}

// FetchSysServices fetches the number of services a snmp client.
// Result is a slice of int containing indices of the found layers.
// Panic has to be handled.
func FetchSysServices(args snmpgo.SNMPArguments) (osiIndices []int) {
	arr := []string{"1.3.6.1.2.1.1.7.0"}
	raw := GetValuesForOids(arr, args)
	r, e := strconv.Atoi(raw[0].Val)
	if e != nil {
		panic("can't extract int from SysUptime: " + e.Error())
	}
	result := make([]int, 0)
	for i := 7; i > 0; i-- {
		layerImpact := math.Pow(2, float64(i-1))
		if layerImpact > float64(r) {
			continue
		} else {
			result = append(result, i-1)
			r -= int(layerImpact)
		}
	}
	return result
}

// FetchAverageCpuLoad fetches the average cpu loads for the last 1, 5 and 15 minutes ago.
// Result is a slice of float32.
// Panic has to be handled.
func FetchAverageCpuLoads(args snmpgo.SNMPArguments) []float32 {
	arr := []string{"1.3.6.1.4.1.2021.10.1.3.1", "1.3.6.1.4.1.2021.10.1.3.2", "1.3.6.1.4.1.2021.10.1.3.3"}
	r := GetValuesForOids(arr, args)
	result := make([]float32, 3)
	for i, v := range r {
		num, err := strconv.ParseFloat(v.Val, 32)
		if err != nil {
			panic(fmt.Sprintf("cpu load number %d can't be converted to float32", i+1))
		}
		result[i] = float32(num)
	}
	return result
}

// FetchSnmpCapabilities fetches the list of snmp capabilities which local snmp application may offer.
// Result is a slice of []*model.SnmpResult, each row has several values related to a capability.
// Panic has to be handled.
func FetchAverageCpuLoadTable(args snmpgo.SNMPArguments) [][]*model.SnmpResult {
	return GetTableForOid("1.3.6.1.4.1.2021.10", 1, []int{1, 2, 3, 4, 5, 6, 100, 101}, args)
}

// FetchCpuAverage fetches the list of snmp capabilities which local snmp application may offer.
// Result is a slice of []*model.SnmpResult, each row has several values related to a capability.
// Panic has to be handled.
func FetchSnmpCapabilities(args snmpgo.SNMPArguments) [][]*model.SnmpResult {
	return GetTableForOid("1.3.6.1.2.1.1.9", 1, []int{2, 3, 4}, args)
}
