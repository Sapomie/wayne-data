package cons

const KeyPrefix = "WayneData"

const (
	RedisKeyDayEss     = KeyPrefix + "EssDay"
	RedisKeyTenEss     = KeyPrefix + "EssTen"
	RedisKeyMonthEss   = KeyPrefix + "EssMonth"
	RedisKeyQuarterEss = KeyPrefix + "EssQuarter"
	RedisKeyHalfEss    = KeyPrefix + "EssHalf"
	RedisKeyYearEss    = KeyPrefix + "EssYear"

	RedisKeyProgressPrefix = KeyPrefix + "Progress"

	RedisKeyParent  = KeyPrefix + "Parent"
	RedisKeyTask    = KeyPrefix + "Task"
	RedisKeyStuff   = KeyPrefix + "Stuff"
	RedisKeyProject = KeyPrefix + "Project"
	RedisKeyTag     = KeyPrefix + "Tag"

	RedisKeyBook   = KeyPrefix + "Book"
	RedisKeySeries = KeyPrefix + "Series"
	RedisKeyMovie  = KeyPrefix + "Movie"

	RedisKeyRun                 = KeyPrefix + "Run"
	RedisKeyRunZonePrefix       = KeyPrefix + "RunZone"
	RedisKeyAnaerobic           = KeyPrefix + "Anaerobic"
	RedisKeyAnaerobicZonePrefix = KeyPrefix + "AnaerobicZone"
)
