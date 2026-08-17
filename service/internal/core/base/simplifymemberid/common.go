package simplifymemberid

// package simplifymemberid is create for display shorten version of Member ID to user
// so this package can only be used for Member ID's transformation.
// every internal logic should use the original(snowflake ID) member ID.

func newSimplifyID() ISimplifyMemberID {
	return &simplifyMemberID{
		generate:        newCommonGen(),
		mappingOriginal: newCommonMappingOriginal(),
		mappingSimplify: newCommonMappingSimplify(),
	}
}

type ISimplifyMemberID interface {
	Generate() (string, error)
	MappingOriginalID(simplifyID string) (int64, error)
	MappingSimplifyID(id int64) (string, error)
}

type simplifyMemberID struct {
	generate        *commonGen
	mappingOriginal *commonMappingOriginal
	mappingSimplify *commonMappingSimplify
}

const (
	width         = 5
	maxSimplifyID = 99999999
)

// Generate
// warning : do not call this function in any other scenario!
// this function may only be used by member registration
func (s *simplifyMemberID) Generate() (string, error) {
	return s.generate.Handler()
}

func (s *simplifyMemberID) MappingOriginalID(simplifyID string) (int64, error) {
	return s.mappingOriginal.Handler(simplifyID)
}

func (s *simplifyMemberID) MappingSimplifyID(id int64) (string, error) {
	return s.mappingSimplify.Handler(id)
}
