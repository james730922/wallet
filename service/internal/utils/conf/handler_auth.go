package conf

type AuthHandler struct {
	LoginMember LoginMemberHandler `yaml:"login_member"`
}
