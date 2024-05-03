package auth

import (
	_ "embed"
)

const (
	RuleAuthenticate   = "auth"
	RuleAny            = "ruleAny"
	RuleAdminOnly      = "ruleAdminOnly"
	RuleUserOnly       = "ruleUserOnly"
	RuleAdminOrSubject = "ruleAdminOrSubject"
)

const (
	opaPackage string = "diegom7s.rego"
)

var (
	opaAuthentication string
	opaAuthorization  string
)
