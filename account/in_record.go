package account

const (
	inAccountId int = iota
	inAccountName
	inFirstName
	inCreatedOn
)

var (
	inHeader = []string{
		"Account ID",
		"Account Name",
		"First Name",
		"Created On",
	}
)

// InRecord is a nice little wrapper from a parsed csv record
type InRecord []string

func (r InRecord) AccountId() string {
	return r[inAccountId]
}

func (r InRecord) AccountName() string {
	return r[inAccountName]
}

func (r InRecord) FirstName() string {
	return r[inFirstName]
}

func (r InRecord) CreatedOn() string {
	return r[inCreatedOn]
}
