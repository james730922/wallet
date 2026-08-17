package scanpay

var (
	dao *storage
)

func newDAO() {
	dao = &storage{
		ScanPayDAO: newScanPayDAO(),
		Mapping:    newScanPayMappingDAO(),
		Record:     newScanPayRecordDAO(),
		//ScanPayMemberDAO: newScanPayMemberDAO(),
	}
}

type storage struct {
	ScanPayDAO IScanPayDAO
	Mapping    IScanPayMappingDAO
	Record     IScanPayRecordDAO
	//ScanPayMemberDAO IScanPayMemberDAO
}
