package user

type Role string 

const (
	ROLEWORKER Role = "WORKER"
	ROLEMANAGER Role = "MANAGER"
	ROLECUSTOMER Role = "CUSTOMER"
	ROLEDEVELOPER Role = "DEVELOPER"
)

func (r Role) isValid() bool {
	switch r {
	case ROLECUSTOMER, 
		ROLEMANAGER,
		ROLEDEVELOPER,
		ROLEWORKER:
		return  true
	default:
		return  false
	} 
}