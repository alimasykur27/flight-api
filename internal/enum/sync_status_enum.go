package enum

type SyncStatusEnum int

var (
	SYNC_NEW        SyncStatusEnum = 1
	SYNC_ON_PROCESS SyncStatusEnum = 10
	SYNC_SYNCED     SyncStatusEnum = 20
	SYNC_NO_NEED    SyncStatusEnum = 30
	SYNC_NOT_FOUND  SyncStatusEnum = 40
	SYNC_ERROR      SyncStatusEnum = 50
)

func ToSyncStatus(i int) SyncStatusEnum {
	switch i {
	case 1:
		return SYNC_NEW
	case 10:
		return SYNC_ON_PROCESS
	case 20:
		return SYNC_SYNCED
	case 30:
		return SYNC_NO_NEED
	case 40:
		return SYNC_NOT_FOUND
	case 50:
		return SYNC_ERROR
	default:
		return SYNC_NEW
	}
}

func (s *SyncStatusEnum) String() string {
	switch *s {
	case SYNC_NEW:
		return "new"
	case SYNC_ON_PROCESS:
		return "on_process"
	case SYNC_SYNCED:
		return "synced"
	case SYNC_NO_NEED:
		return "no_need"
	case SYNC_NOT_FOUND:
		return "not_found"
	case SYNC_ERROR:
		return "error"
	default:
		return "unknown"
	}
}

func (s *SyncStatusEnum) Int() int {
	return int(*s)
}
