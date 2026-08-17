package resp

type MemberNumResp struct {
	OnlineMemberNum int `json:"online_member_num"`
	TotalMemberNum  int `json:"total_member_num"`
}
