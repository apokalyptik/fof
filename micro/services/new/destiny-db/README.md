destinyAccounts
  ID integer primary key autoincriment
  Account ID string not null default ""
  LastPlayed datetime not null default "0000-00-00 00:00:00"

// Things like... Race, Class, Fastest SRL Time on Trac X
destinyAccountValues
  ID integer not null default 0
  Key string
  Value blob

// Things like VoG Completions, Kills, Etc
destinyAccountStats
  ID integer not null default 0
  Stat string not null default ""
  Value integer not null default 0

// Things like Kills, KF Completions, Rift Matches
destinyAccountDailyStats
  ID integer not null default 0
  Stat string not null default ""
  When date not null default "0000-00-00"
  Value integer not null default 0

// Blobs of compressed json
destinyRawAccountData
  ID integer not null default 0
  Kind string not null default ""
  Raw BLOB not null default ""

destinyCharacters
  ID integer primary key autoincriment
  AccountID string not null default ""
  CharacterID string not null default ""
  LastPlayed datetime not null default "0000-00-00 00:00:00"

// See Account QV
destinyCharacterValues
  ID integer not null default 0
  Key string
  Value blob

// See Account QV
destinyCharacterStats
  ID integer not null default 0
  Stat string not null default ""
  Value integer not null default 0

// See Account QV
destinyCharacterDailyStats
  ID integer not null default 0
  Stat string not null default ""
  When date not null default "0000-00-00"
  Value integer not null default 0

// See Account QV
destinyRawChatacterData
  ID integer not null default 0
  Kind string not null default ""
  Raw BLOB not null default ""


[ slack => userdb <= destinyid ] => { new_user | modified_user | deleted_user }

{ new_user } => []
{ modified_user } => []
{ deleted_user } => []

[ poll_acct_status* ] => { character_update }

{ character_update } => [ pgcr_intake* & attach_pgcr ] => { new_pgcr }
{ character_update } => [ stats_intake* ] => { stat_updated }

{ new_pgcr } => [ record_stats* ]
{ stat_updated } => [ record_stats* ]