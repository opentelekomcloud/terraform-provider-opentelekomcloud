package apigw

type (
	PolicyType string
	PeriodUnit string
)

const (
	PeriodUnitSecond PeriodUnit = "SECOND"
	PeriodUnitMinute PeriodUnit = "MINUTE"
	PeriodUnitHour   PeriodUnit = "HOUR"
	PeriodUnitDay    PeriodUnit = "DAY"

	PolicyTypeExclusive   PolicyType = "API-based"
	PolicyTypeShared      PolicyType = "API-shared"
	PolicyTypeUser        PolicyType = "USER"
	PolicyTypeApplication PolicyType = "APP"

	includeSpecialThrottle int = 1
)

var (
	policyType = map[string]int{
		string(PolicyTypeExclusive): 1,
		string(PolicyTypeShared):    2,
	}
)
