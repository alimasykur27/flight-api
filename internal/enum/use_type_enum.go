package enum

type UseTypeEnum *string

var (
	usePublicStr  = "public"
	usePrivateStr = "private"
)

var (
	USE_PUBLIC  UseTypeEnum = &usePublicStr
	USE_PRIVATE UseTypeEnum = &usePrivateStr
	USE_NIL     UseTypeEnum = nil
)

func ToUseType(s string) UseTypeEnum {
	switch s {
	case "public":
		return USE_PUBLIC
	case "private":
		return USE_PRIVATE
	default:
		return UseTypeEnum(USE_NIL)
	}
}
