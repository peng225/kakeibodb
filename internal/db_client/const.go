package db_client

const DBName = "kakeibo"

const EventTableName = "event"
const TagTableName = "tag"
const EventToTagTableName = "event_to_tag"
const PatternTableName = "pattern"
const PatternToTagTableName = "pattern_to_tag"

const EventColID = "id"
const EventColDate = "dt"
const EventColMoney = "money"
const EventColDescription = "description"
const TagColID = "id"
const TagColName = "name"
const PatternColID = "id"
const PatternColKey = "key_string"
const PatternToTagColPID = "pattern_id"
const PatternToTagColTID = "tag_id"

const EventDescLength = 32
