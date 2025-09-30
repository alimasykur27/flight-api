package enum

import "database/sql"

type OwnershipEnum *string

var (
	ownPublicStr  = "public"
	ownPrivateStr = "private"
)

var (
	OWN_PUBLIC  OwnershipEnum = &ownPublicStr
	OWN_PRIVATE OwnershipEnum = &ownPrivateStr
	OWN_NIL     OwnershipEnum = nil
)

func ToOwnership(s string) OwnershipEnum {
	switch s {
	case "public":
		return OWN_PUBLIC
	case "private":
		return OWN_PRIVATE
	default:
		return OwnershipEnum(OWN_NIL)
	}
}

func FromNullStringToOwnership(ns sql.NullString) OwnershipEnum {
	if ns.Valid {
		return ToOwnership(ns.String)
	}
	return OWN_NIL
}
