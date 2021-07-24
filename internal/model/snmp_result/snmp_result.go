package snmp_result

type SnmpResult struct {
	Oid string
	Typ string
	Val string
}

func (SnmpResult) New(oid, typ, val string) *SnmpResult {
	return &SnmpResult{
		Oid: oid,
		Typ: typ,
		Val: val,
	}
}
