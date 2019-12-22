// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: errcode.proto

package errcode

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ErrCode int32

const (
	Undefined                               ErrCode = 0
	TODO                                    ErrCode = 666
	ErrNotImplemented                       ErrCode = 777
	ErrDeprecated                           ErrCode = 888
	ErrInternal                             ErrCode = 999
	ErrInvalidInput                         ErrCode = 101
	ErrMissingInput                         ErrCode = 102
	ErrUnauthenticated                      ErrCode = 103
	ErrSSOGetOIDC                           ErrCode = 1001
	ErrSSOInvalidPublicKey                  ErrCode = 1002
	ErrSSOFailedKeycloakRequest             ErrCode = 1003
	ErrSSOInvalidKeycloakResponse           ErrCode = 1004
	ErrSSOLogout                            ErrCode = 1005
	ErrSSOInitKeycloakClient                ErrCode = 1006
	ErrSSOInvalidBearer                     ErrCode = 1007
	ErrSSOKeycloakError                     ErrCode = 1008
	ErrDBNotFound                           ErrCode = 2001
	ErrDBInternal                           ErrCode = 2002
	ErrDBRunMigrations                      ErrCode = 2003
	ErrDBInit                               ErrCode = 2004
	ErrDBConnect                            ErrCode = 2005
	ErrDBAutoMigrate                        ErrCode = 2006
	ErrDBAddForeignKey                      ErrCode = 2007
	ErrComposeInvalidPath                   ErrCode = 3001
	ErrComposeDirectoryNotFound             ErrCode = 3002
	ErrComposeReadConfig                    ErrCode = 3003
	ErrComposeInvalidConfig                 ErrCode = 3004
	ErrComposeMarshalConfig                 ErrCode = 3005
	ErrComposeCreateTempFile                ErrCode = 3006
	ErrComposeWriteTempFile                 ErrCode = 3007
	ErrComposeBuild                         ErrCode = 3008
	ErrComposeBundle                        ErrCode = 3009
	ErrComposeReadDab                       ErrCode = 3010
	ErrComposeParseDab                      ErrCode = 3011
	ErrComposeParseConfig                   ErrCode = 3012
	ErrComposeCreateTempDir                 ErrCode = 3013
	ErrComposeForceRecreateDown             ErrCode = 3014
	ErrComposeRunCreate                     ErrCode = 3015
	ErrComposeRunUp                         ErrCode = 3016
	ErrGetPWInitBinary                      ErrCode = 3017
	ErrWritePWInitFileHeader                ErrCode = 3018
	ErrWritePWInitFile                      ErrCode = 3019
	ErrMarshalPWInitConfigFile              ErrCode = 3020
	ErrWritePWInitConfigFileHeader          ErrCode = 3021
	ErrWritePWInitConfigFile                ErrCode = 3022
	ErrWritePWInitCloseTarWriter            ErrCode = 3023
	ErrCopyPWInitToContainer                ErrCode = 3024
	ErrComposeGetContainersInfo             ErrCode = 3025
	ErrMissingPwinitConfig                  ErrCode = 3026
	ErrGetUserIDFromContext                 ErrCode = 4001
	ErrMissingChallengeValidation           ErrCode = 4002
	ErrInvalidSeason                        ErrCode = 4003
	ErrTeamNotInSeason                      ErrCode = 4004
	ErrGetChallengeSubscription             ErrCode = 4005
	ErrInitSnowflake                        ErrCode = 4006
	ErrUpdateChallengeSubscription          ErrCode = 4007
	ErrCreateChallengeValidation            ErrCode = 4008
	ErrGetChallengeValidation               ErrCode = 4009
	ErrInvalidTeam                          ErrCode = 4010
	ErrChallengeAlreadySubscribed           ErrCode = 4011
	ErrCreateChallengeSubscription          ErrCode = 4012
	ErrFindOrganizations                    ErrCode = 4013
	ErrGetSeasonFromSeasonChallenge         ErrCode = 4014
	ErrGetUserTeamFromSeason                ErrCode = 4015
	ErrInvalidSeasonID                      ErrCode = 4016
	ErrUserHasNoTeamForSeason               ErrCode = 4017
	ErrGetSeasonChallenges                  ErrCode = 4018
	ErrGetSeason                            ErrCode = 4019
	ErrSeasonDenied                         ErrCode = 4020
	ErrAlreadyHasTeamForSeason              ErrCode = 4021
	ErrReservedName                         ErrCode = 4022
	ErrCheckOrganizationUniqueName          ErrCode = 4023
	ErrCreateOrganization                   ErrCode = 4024
	ErrOrganizationAlreadyHasTeamForSeason  ErrCode = 4025
	ErrGetOrganization                      ErrCode = 4026
	ErrGetSeasonChallenge                   ErrCode = 4027
	ErrCannotCreateTeamForSoloOrganization  ErrCode = 4028
	ErrUserNotInOrganization                ErrCode = 4029
	ErrCreateTeam                           ErrCode = 4030
	ErrGetTeam                              ErrCode = 4031
	ErrGetTeams                             ErrCode = 4032
	ErrGetUser                              ErrCode = 4033
	ErrUpdateUser                           ErrCode = 4034
	ErrUpdateTeam                           ErrCode = 4035
	ErrUpdateOrganization                   ErrCode = 4036
	ErrNewUserFromClaims                    ErrCode = 4037
	ErrGetOAuthUser                         ErrCode = 4038
	ErrDifferentUserBetweenTokenAndDatabase ErrCode = 4039
	ErrLoadUserSeasons                      ErrCode = 4040
	ErrGetUserOrganizations                 ErrCode = 4041
	ErrGetSeasons                           ErrCode = 4042
	ErrGetUserBySubject                     ErrCode = 4043
	ErrEmailAddressNotVerified              ErrCode = 4044
	ErrGetDefaultSeason                     ErrCode = 4045
	ErrCommitUserTransaction                ErrCode = 4046
	ErrUpdateActiveSeason                   ErrCode = 4047
	ErrMissingContextMetadata               ErrCode = 4048
	ErrNoTokenProvided                      ErrCode = 4049
	ErrGetTokenWithClaims                   ErrCode = 4050
	ErrNoTokenInContext                     ErrCode = 4051
	ErrGetSubjectFromToken                  ErrCode = 4052
	ErrGetSubjectFromContext                ErrCode = 4053
	ErrGetActiveSeasonMembership            ErrCode = 4054
	ErrGetTokenFromContext                  ErrCode = 4055
	ErrChallengeAlreadyClosed               ErrCode = 4056
	ErrGetAgent                             ErrCode = 4057
	ErrSaveAgent                            ErrCode = 4058
	ErrListChallengeInstances               ErrCode = 4059
	ErrServerListen                         ErrCode = 5001
	ErrServerRegisterGateway                ErrCode = 5002
	ErrInitLogger                           ErrCode = 6001
	ErrStartService                         ErrCode = 6002
	ErrInitServer                           ErrCode = 6003
	ErrGroupTerminated                      ErrCode = 6004
	ErrGetSSOClientFromFlags                ErrCode = 6005
	ErrDumpDatabase                         ErrCode = 6006
	ErrGetDBInfo                            ErrCode = 6007
	ErrGetSSOWhoami                         ErrCode = 6008
	ErrGetSSOLogout                         ErrCode = 6009
	ErrGetSSOClaims                         ErrCode = 6010
	ErrInitDockerClient                     ErrCode = 6011
	ErrInitDB                               ErrCode = 6012
	ErrConfigureDB                          ErrCode = 6013
	ErrInitSSOClient                        ErrCode = 6014
	ErrInitService                          ErrCode = 6015
	ErrAgentGetContainersInfo               ErrCode = 7001
	ErrCheckNginxContainer                  ErrCode = 7002
	ErrRemoveNginxContainer                 ErrCode = 7003
	ErrBuildNginxContainer                  ErrCode = 7004
	ErrStartNginxContainer                  ErrCode = 7005
	ErrParsingTemplate                      ErrCode = 7006
	ErrWriteConfigFileHeader                ErrCode = 7007
	ErrWriteConfigFile                      ErrCode = 7008
	ErrCloseTarWriter                       ErrCode = 7009
	ErrCopyNginxConfigToContainer           ErrCode = 7010
	ErrNginxNewConfigCheckFailed            ErrCode = 7011
	ErrNginxSendCommandNewConfigCheck       ErrCode = 7012
	ErrNginxSendCommandNewConfigRemove      ErrCode = 7013
	ErrNginxSendCommandConfigReplace        ErrCode = 7014
	ErrNginxSendCommandReloadConfig         ErrCode = 7015
	ErrNginxConnectNetwork                  ErrCode = 7016
	ErrContainerConnectNetwork              ErrCode = 7017
	ErrNatPortOpening                       ErrCode = 7018
	ErrBuildNginxConfig                     ErrCode = 7019
	ErrExecuteTemplate                      ErrCode = 7020
	ErrWriteBytesToHashBuilder              ErrCode = 7021
	ErrReadBytesFromHashBuilder             ErrCode = 7022
	ErrGeneratePrefixHash                   ErrCode = 7023
	ErrCleanPathwarInstances                ErrCode = 7024
	ErrParseInitConfig                      ErrCode = 7025
	ErrUpPathwarInstance                    ErrCode = 7026
	ErrUpdateNginx                          ErrCode = 7027
	ErrDockerAPIContainerList               ErrCode = 8001
	ErrDockerAPIContainerRemove             ErrCode = 8002
	ErrDockerAPIImageRemove                 ErrCode = 8003
	ErrDockerAPIContainerCreate             ErrCode = 8004
	ErrDockerAPIContainerExecCreate         ErrCode = 8005
	ErrDockerAPIContainerExecAttach         ErrCode = 8006
	ErrDockerAPIContainerExecStart          ErrCode = 8007
	ErrDockerAPIContainerExecInspect        ErrCode = 8008
	ErrDockerAPIImagePull                   ErrCode = 8009
	ErrDockerAPINetworkList                 ErrCode = 8010
	ErrDockerAPINetworkCreate               ErrCode = 8011
	ErrDockerAPINetworkRemove               ErrCode = 8012
	ErrDockerAPIExitCode                    ErrCode = 8013
	ErrExecuteOnInitHook                    ErrCode = 9001
	ErrRemoveInitConfig                     ErrCode = 9002
)

var ErrCode_name = map[int32]string{
	0:    "Undefined",
	666:  "TODO",
	777:  "ErrNotImplemented",
	888:  "ErrDeprecated",
	999:  "ErrInternal",
	101:  "ErrInvalidInput",
	102:  "ErrMissingInput",
	103:  "ErrUnauthenticated",
	1001: "ErrSSOGetOIDC",
	1002: "ErrSSOInvalidPublicKey",
	1003: "ErrSSOFailedKeycloakRequest",
	1004: "ErrSSOInvalidKeycloakResponse",
	1005: "ErrSSOLogout",
	1006: "ErrSSOInitKeycloakClient",
	1007: "ErrSSOInvalidBearer",
	1008: "ErrSSOKeycloakError",
	2001: "ErrDBNotFound",
	2002: "ErrDBInternal",
	2003: "ErrDBRunMigrations",
	2004: "ErrDBInit",
	2005: "ErrDBConnect",
	2006: "ErrDBAutoMigrate",
	2007: "ErrDBAddForeignKey",
	3001: "ErrComposeInvalidPath",
	3002: "ErrComposeDirectoryNotFound",
	3003: "ErrComposeReadConfig",
	3004: "ErrComposeInvalidConfig",
	3005: "ErrComposeMarshalConfig",
	3006: "ErrComposeCreateTempFile",
	3007: "ErrComposeWriteTempFile",
	3008: "ErrComposeBuild",
	3009: "ErrComposeBundle",
	3010: "ErrComposeReadDab",
	3011: "ErrComposeParseDab",
	3012: "ErrComposeParseConfig",
	3013: "ErrComposeCreateTempDir",
	3014: "ErrComposeForceRecreateDown",
	3015: "ErrComposeRunCreate",
	3016: "ErrComposeRunUp",
	3017: "ErrGetPWInitBinary",
	3018: "ErrWritePWInitFileHeader",
	3019: "ErrWritePWInitFile",
	3020: "ErrMarshalPWInitConfigFile",
	3021: "ErrWritePWInitConfigFileHeader",
	3022: "ErrWritePWInitConfigFile",
	3023: "ErrWritePWInitCloseTarWriter",
	3024: "ErrCopyPWInitToContainer",
	3025: "ErrComposeGetContainersInfo",
	3026: "ErrMissingPwinitConfig",
	4001: "ErrGetUserIDFromContext",
	4002: "ErrMissingChallengeValidation",
	4003: "ErrInvalidSeason",
	4004: "ErrTeamNotInSeason",
	4005: "ErrGetChallengeSubscription",
	4006: "ErrInitSnowflake",
	4007: "ErrUpdateChallengeSubscription",
	4008: "ErrCreateChallengeValidation",
	4009: "ErrGetChallengeValidation",
	4010: "ErrInvalidTeam",
	4011: "ErrChallengeAlreadySubscribed",
	4012: "ErrCreateChallengeSubscription",
	4013: "ErrFindOrganizations",
	4014: "ErrGetSeasonFromSeasonChallenge",
	4015: "ErrGetUserTeamFromSeason",
	4016: "ErrInvalidSeasonID",
	4017: "ErrUserHasNoTeamForSeason",
	4018: "ErrGetSeasonChallenges",
	4019: "ErrGetSeason",
	4020: "ErrSeasonDenied",
	4021: "ErrAlreadyHasTeamForSeason",
	4022: "ErrReservedName",
	4023: "ErrCheckOrganizationUniqueName",
	4024: "ErrCreateOrganization",
	4025: "ErrOrganizationAlreadyHasTeamForSeason",
	4026: "ErrGetOrganization",
	4027: "ErrGetSeasonChallenge",
	4028: "ErrCannotCreateTeamForSoloOrganization",
	4029: "ErrUserNotInOrganization",
	4030: "ErrCreateTeam",
	4031: "ErrGetTeam",
	4032: "ErrGetTeams",
	4033: "ErrGetUser",
	4034: "ErrUpdateUser",
	4035: "ErrUpdateTeam",
	4036: "ErrUpdateOrganization",
	4037: "ErrNewUserFromClaims",
	4038: "ErrGetOAuthUser",
	4039: "ErrDifferentUserBetweenTokenAndDatabase",
	4040: "ErrLoadUserSeasons",
	4041: "ErrGetUserOrganizations",
	4042: "ErrGetSeasons",
	4043: "ErrGetUserBySubject",
	4044: "ErrEmailAddressNotVerified",
	4045: "ErrGetDefaultSeason",
	4046: "ErrCommitUserTransaction",
	4047: "ErrUpdateActiveSeason",
	4048: "ErrMissingContextMetadata",
	4049: "ErrNoTokenProvided",
	4050: "ErrGetTokenWithClaims",
	4051: "ErrNoTokenInContext",
	4052: "ErrGetSubjectFromToken",
	4053: "ErrGetSubjectFromContext",
	4054: "ErrGetActiveSeasonMembership",
	4055: "ErrGetTokenFromContext",
	4056: "ErrChallengeAlreadyClosed",
	4057: "ErrGetAgent",
	4058: "ErrSaveAgent",
	4059: "ErrListChallengeInstances",
	5001: "ErrServerListen",
	5002: "ErrServerRegisterGateway",
	6001: "ErrInitLogger",
	6002: "ErrStartService",
	6003: "ErrInitServer",
	6004: "ErrGroupTerminated",
	6005: "ErrGetSSOClientFromFlags",
	6006: "ErrDumpDatabase",
	6007: "ErrGetDBInfo",
	6008: "ErrGetSSOWhoami",
	6009: "ErrGetSSOLogout",
	6010: "ErrGetSSOClaims",
	6011: "ErrInitDockerClient",
	6012: "ErrInitDB",
	6013: "ErrConfigureDB",
	6014: "ErrInitSSOClient",
	6015: "ErrInitService",
	7001: "ErrAgentGetContainersInfo",
	7002: "ErrCheckNginxContainer",
	7003: "ErrRemoveNginxContainer",
	7004: "ErrBuildNginxContainer",
	7005: "ErrStartNginxContainer",
	7006: "ErrParsingTemplate",
	7007: "ErrWriteConfigFileHeader",
	7008: "ErrWriteConfigFile",
	7009: "ErrCloseTarWriter",
	7010: "ErrCopyNginxConfigToContainer",
	7011: "ErrNginxNewConfigCheckFailed",
	7012: "ErrNginxSendCommandNewConfigCheck",
	7013: "ErrNginxSendCommandNewConfigRemove",
	7014: "ErrNginxSendCommandConfigReplace",
	7015: "ErrNginxSendCommandReloadConfig",
	7016: "ErrNginxConnectNetwork",
	7017: "ErrContainerConnectNetwork",
	7018: "ErrNatPortOpening",
	7019: "ErrBuildNginxConfig",
	7020: "ErrExecuteTemplate",
	7021: "ErrWriteBytesToHashBuilder",
	7022: "ErrReadBytesFromHashBuilder",
	7023: "ErrGeneratePrefixHash",
	7024: "ErrCleanPathwarInstances",
	7025: "ErrParseInitConfig",
	7026: "ErrUpPathwarInstance",
	7027: "ErrUpdateNginx",
	8001: "ErrDockerAPIContainerList",
	8002: "ErrDockerAPIContainerRemove",
	8003: "ErrDockerAPIImageRemove",
	8004: "ErrDockerAPIContainerCreate",
	8005: "ErrDockerAPIContainerExecCreate",
	8006: "ErrDockerAPIContainerExecAttach",
	8007: "ErrDockerAPIContainerExecStart",
	8008: "ErrDockerAPIContainerExecInspect",
	8009: "ErrDockerAPIImagePull",
	8010: "ErrDockerAPINetworkList",
	8011: "ErrDockerAPINetworkCreate",
	8012: "ErrDockerAPINetworkRemove",
	8013: "ErrDockerAPIExitCode",
	9001: "ErrExecuteOnInitHook",
	9002: "ErrRemoveInitConfig",
}

var ErrCode_value = map[string]int32{
	"Undefined":                               0,
	"TODO":                                    666,
	"ErrNotImplemented":                       777,
	"ErrDeprecated":                           888,
	"ErrInternal":                             999,
	"ErrInvalidInput":                         101,
	"ErrMissingInput":                         102,
	"ErrUnauthenticated":                      103,
	"ErrSSOGetOIDC":                           1001,
	"ErrSSOInvalidPublicKey":                  1002,
	"ErrSSOFailedKeycloakRequest":             1003,
	"ErrSSOInvalidKeycloakResponse":           1004,
	"ErrSSOLogout":                            1005,
	"ErrSSOInitKeycloakClient":                1006,
	"ErrSSOInvalidBearer":                     1007,
	"ErrSSOKeycloakError":                     1008,
	"ErrDBNotFound":                           2001,
	"ErrDBInternal":                           2002,
	"ErrDBRunMigrations":                      2003,
	"ErrDBInit":                               2004,
	"ErrDBConnect":                            2005,
	"ErrDBAutoMigrate":                        2006,
	"ErrDBAddForeignKey":                      2007,
	"ErrComposeInvalidPath":                   3001,
	"ErrComposeDirectoryNotFound":             3002,
	"ErrComposeReadConfig":                    3003,
	"ErrComposeInvalidConfig":                 3004,
	"ErrComposeMarshalConfig":                 3005,
	"ErrComposeCreateTempFile":                3006,
	"ErrComposeWriteTempFile":                 3007,
	"ErrComposeBuild":                         3008,
	"ErrComposeBundle":                        3009,
	"ErrComposeReadDab":                       3010,
	"ErrComposeParseDab":                      3011,
	"ErrComposeParseConfig":                   3012,
	"ErrComposeCreateTempDir":                 3013,
	"ErrComposeForceRecreateDown":             3014,
	"ErrComposeRunCreate":                     3015,
	"ErrComposeRunUp":                         3016,
	"ErrGetPWInitBinary":                      3017,
	"ErrWritePWInitFileHeader":                3018,
	"ErrWritePWInitFile":                      3019,
	"ErrMarshalPWInitConfigFile":              3020,
	"ErrWritePWInitConfigFileHeader":          3021,
	"ErrWritePWInitConfigFile":                3022,
	"ErrWritePWInitCloseTarWriter":            3023,
	"ErrCopyPWInitToContainer":                3024,
	"ErrComposeGetContainersInfo":             3025,
	"ErrMissingPwinitConfig":                  3026,
	"ErrGetUserIDFromContext":                 4001,
	"ErrMissingChallengeValidation":           4002,
	"ErrInvalidSeason":                        4003,
	"ErrTeamNotInSeason":                      4004,
	"ErrGetChallengeSubscription":             4005,
	"ErrInitSnowflake":                        4006,
	"ErrUpdateChallengeSubscription":          4007,
	"ErrCreateChallengeValidation":            4008,
	"ErrGetChallengeValidation":               4009,
	"ErrInvalidTeam":                          4010,
	"ErrChallengeAlreadySubscribed":           4011,
	"ErrCreateChallengeSubscription":          4012,
	"ErrFindOrganizations":                    4013,
	"ErrGetSeasonFromSeasonChallenge":         4014,
	"ErrGetUserTeamFromSeason":                4015,
	"ErrInvalidSeasonID":                      4016,
	"ErrUserHasNoTeamForSeason":               4017,
	"ErrGetSeasonChallenges":                  4018,
	"ErrGetSeason":                            4019,
	"ErrSeasonDenied":                         4020,
	"ErrAlreadyHasTeamForSeason":              4021,
	"ErrReservedName":                         4022,
	"ErrCheckOrganizationUniqueName":          4023,
	"ErrCreateOrganization":                   4024,
	"ErrOrganizationAlreadyHasTeamForSeason":  4025,
	"ErrGetOrganization":                      4026,
	"ErrGetSeasonChallenge":                   4027,
	"ErrCannotCreateTeamForSoloOrganization":  4028,
	"ErrUserNotInOrganization":                4029,
	"ErrCreateTeam":                           4030,
	"ErrGetTeam":                              4031,
	"ErrGetTeams":                             4032,
	"ErrGetUser":                              4033,
	"ErrUpdateUser":                           4034,
	"ErrUpdateTeam":                           4035,
	"ErrUpdateOrganization":                   4036,
	"ErrNewUserFromClaims":                    4037,
	"ErrGetOAuthUser":                         4038,
	"ErrDifferentUserBetweenTokenAndDatabase": 4039,
	"ErrLoadUserSeasons":                      4040,
	"ErrGetUserOrganizations":                 4041,
	"ErrGetSeasons":                           4042,
	"ErrGetUserBySubject":                     4043,
	"ErrEmailAddressNotVerified":              4044,
	"ErrGetDefaultSeason":                     4045,
	"ErrCommitUserTransaction":                4046,
	"ErrUpdateActiveSeason":                   4047,
	"ErrMissingContextMetadata":               4048,
	"ErrNoTokenProvided":                      4049,
	"ErrGetTokenWithClaims":                   4050,
	"ErrNoTokenInContext":                     4051,
	"ErrGetSubjectFromToken":                  4052,
	"ErrGetSubjectFromContext":                4053,
	"ErrGetActiveSeasonMembership":            4054,
	"ErrGetTokenFromContext":                  4055,
	"ErrChallengeAlreadyClosed":               4056,
	"ErrGetAgent":                             4057,
	"ErrSaveAgent":                            4058,
	"ErrListChallengeInstances":               4059,
	"ErrServerListen":                         5001,
	"ErrServerRegisterGateway":                5002,
	"ErrInitLogger":                           6001,
	"ErrStartService":                         6002,
	"ErrInitServer":                           6003,
	"ErrGroupTerminated":                      6004,
	"ErrGetSSOClientFromFlags":                6005,
	"ErrDumpDatabase":                         6006,
	"ErrGetDBInfo":                            6007,
	"ErrGetSSOWhoami":                         6008,
	"ErrGetSSOLogout":                         6009,
	"ErrGetSSOClaims":                         6010,
	"ErrInitDockerClient":                     6011,
	"ErrInitDB":                               6012,
	"ErrConfigureDB":                          6013,
	"ErrInitSSOClient":                        6014,
	"ErrInitService":                          6015,
	"ErrAgentGetContainersInfo":               7001,
	"ErrCheckNginxContainer":                  7002,
	"ErrRemoveNginxContainer":                 7003,
	"ErrBuildNginxContainer":                  7004,
	"ErrStartNginxContainer":                  7005,
	"ErrParsingTemplate":                      7006,
	"ErrWriteConfigFileHeader":                7007,
	"ErrWriteConfigFile":                      7008,
	"ErrCloseTarWriter":                       7009,
	"ErrCopyNginxConfigToContainer":           7010,
	"ErrNginxNewConfigCheckFailed":            7011,
	"ErrNginxSendCommandNewConfigCheck":       7012,
	"ErrNginxSendCommandNewConfigRemove":      7013,
	"ErrNginxSendCommandConfigReplace":        7014,
	"ErrNginxSendCommandReloadConfig":         7015,
	"ErrNginxConnectNetwork":                  7016,
	"ErrContainerConnectNetwork":              7017,
	"ErrNatPortOpening":                       7018,
	"ErrBuildNginxConfig":                     7019,
	"ErrExecuteTemplate":                      7020,
	"ErrWriteBytesToHashBuilder":              7021,
	"ErrReadBytesFromHashBuilder":             7022,
	"ErrGeneratePrefixHash":                   7023,
	"ErrCleanPathwarInstances":                7024,
	"ErrParseInitConfig":                      7025,
	"ErrUpPathwarInstance":                    7026,
	"ErrUpdateNginx":                          7027,
	"ErrDockerAPIContainerList":               8001,
	"ErrDockerAPIContainerRemove":             8002,
	"ErrDockerAPIImageRemove":                 8003,
	"ErrDockerAPIContainerCreate":             8004,
	"ErrDockerAPIContainerExecCreate":         8005,
	"ErrDockerAPIContainerExecAttach":         8006,
	"ErrDockerAPIContainerExecStart":          8007,
	"ErrDockerAPIContainerExecInspect":        8008,
	"ErrDockerAPIImagePull":                   8009,
	"ErrDockerAPINetworkList":                 8010,
	"ErrDockerAPINetworkCreate":               8011,
	"ErrDockerAPINetworkRemove":               8012,
	"ErrDockerAPIExitCode":                    8013,
	"ErrExecuteOnInitHook":                    9001,
	"ErrRemoveInitConfig":                     9002,
}

func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}

func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4240057316120df7, []int{0}
}

func init() {
	proto.RegisterEnum("pathwar.errcode.ErrCode", ErrCode_name, ErrCode_value)
}

func init() { proto.RegisterFile("errcode.proto", fileDescriptor_4240057316120df7) }

var fileDescriptor_4240057316120df7 = []byte{
	// 2101 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x58, 0xd9, 0x6f, 0x1c, 0xc7,
	0xd1, 0x97, 0x81, 0xef, 0x33, 0xa1, 0x49, 0x64, 0x96, 0x47, 0xb6, 0xd6, 0x96, 0x6c, 0x8e, 0xed,
	0x38, 0x16, 0x90, 0x63, 0xf5, 0x10, 0x60, 0x81, 0xbc, 0x10, 0xe0, 0x72, 0x49, 0x6a, 0x61, 0x69,
	0x49, 0xf0, 0xb0, 0x80, 0xbc, 0x35, 0x77, 0x6a, 0x67, 0x3b, 0x9c, 0xed, 0x5e, 0xf7, 0xf4, 0x90,
	0x62, 0xfe, 0x03, 0xe7, 0x29, 0xcf, 0x79, 0xcb, 0x1d, 0xdb, 0xb9, 0x6f, 0xdf, 0xd6, 0x7d, 0xdb,
	0x92, 0x7c, 0xdb, 0x39, 0x6c, 0x29, 0x89, 0xe3, 0xfb, 0xca, 0xe1, 0xdc, 0x41, 0xf7, 0x54, 0xcf,
	0xce, 0x2c, 0x29, 0xbf, 0x91, 0xf5, 0xab, 0xaa, 0xae, 0xfe, 0xd5, 0xd1, 0x35, 0xeb, 0x6d, 0x43,
	0xa5, 0xda, 0x32, 0xc4, 0x6a, 0x5f, 0x49, 0x2d, 0xfd, 0xd1, 0x3e, 0xd3, 0xdd, 0x35, 0xa6, 0xaa,
	0x24, 0xde, 0xf9, 0xd9, 0x88, 0xeb, 0x6e, 0xba, 0x5c, 0x6d, 0xcb, 0xde, 0x9e, 0x48, 0x46, 0x72,
	0x8f, 0xd5, 0x5b, 0x4e, 0x3b, 0xf6, 0x3f, 0xfb, 0x8f, 0xfd, 0x2b, 0xb3, 0xff, 0xd4, 0xfd, 0xb7,
	0x7b, 0x23, 0x53, 0x4a, 0x4d, 0xca, 0x10, 0xfd, 0x6d, 0xde, 0xd6, 0x25, 0x11, 0x62, 0x87, 0x0b,
	0x0c, 0x61, 0x8b, 0xbf, 0xd5, 0xfb, 0xbf, 0xc5, 0xd9, 0xc6, 0x2c, 0x7c, 0xf5, 0xff, 0xfd, 0x1d,
	0xde, 0xb5, 0x53, 0x4a, 0xb5, 0xa4, 0x6e, 0xf6, 0xfa, 0x31, 0xf6, 0x50, 0x68, 0x0c, 0xe1, 0x9e,
	0xab, 0x7d, 0xdf, 0xdb, 0x36, 0xa5, 0x54, 0x03, 0xfb, 0x0a, 0xdb, 0xcc, 0xc8, 0x3e, 0xbc, 0xda,
	0x07, 0xef, 0x63, 0x53, 0x4a, 0x35, 0x85, 0x46, 0x25, 0x58, 0x0c, 0xaf, 0x8e, 0xf8, 0xdb, 0xbd,
	0x51, 0x2b, 0x59, 0x65, 0x31, 0x0f, 0x9b, 0xa2, 0x9f, 0x6a, 0x40, 0x12, 0xee, 0xe7, 0x49, 0xc2,
	0x45, 0x94, 0x09, 0x3b, 0xfe, 0x0e, 0xcf, 0x9f, 0x52, 0x6a, 0x49, 0xb0, 0x54, 0x77, 0x51, 0x68,
	0x9e, 0x39, 0x8d, 0xe8, 0x9c, 0x85, 0x85, 0xd9, 0x19, 0xd4, 0xb3, 0xcd, 0xc6, 0x24, 0xbc, 0x36,
	0xe2, 0xef, 0xf2, 0x76, 0x64, 0x32, 0x72, 0x3c, 0x97, 0x2e, 0xc7, 0xbc, 0x7d, 0x27, 0xae, 0xc3,
	0xeb, 0x23, 0xfe, 0x2d, 0xde, 0xae, 0x0c, 0x9c, 0x66, 0x3c, 0xc6, 0xf0, 0x4e, 0x5c, 0x6f, 0xc7,
	0x92, 0xad, 0xcc, 0xe3, 0xdd, 0x29, 0x26, 0x1a, 0xde, 0x18, 0xf1, 0x6f, 0xf3, 0x6e, 0x2e, 0x99,
	0x0f, 0x54, 0x92, 0xbe, 0x14, 0x09, 0xc2, 0x9b, 0x23, 0xfe, 0xb5, 0xde, 0xc7, 0x33, 0x9d, 0x7d,
	0x32, 0x92, 0xa9, 0x86, 0xb7, 0x46, 0xfc, 0x9b, 0xbd, 0x1b, 0x9c, 0x19, 0xd7, 0xce, 0x66, 0x32,
	0xe6, 0x28, 0x34, 0xbc, 0x3d, 0xe2, 0xdf, 0xe0, 0x6d, 0x2f, 0x79, 0xad, 0x23, 0x53, 0xa8, 0xe0,
	0x9d, 0x02, 0xe2, 0x8c, 0xa6, 0x94, 0x92, 0x0a, 0xde, 0x1d, 0x71, 0x24, 0xd6, 0x5b, 0x52, 0x4f,
	0xcb, 0x54, 0x84, 0x70, 0x61, 0x34, 0x97, 0xe5, 0x34, 0x5e, 0x1c, 0xf5, 0x2b, 0x96, 0x9c, 0x46,
	0x7d, 0x3e, 0x15, 0xfb, 0x79, 0xa4, 0x98, 0xe6, 0x52, 0x24, 0xf0, 0xd4, 0xa8, 0x7f, 0x8d, 0xb7,
	0x95, 0x94, 0xb9, 0x86, 0xa7, 0x47, 0x29, 0xec, 0x46, 0x7d, 0x52, 0x0a, 0x81, 0x6d, 0x0d, 0xcf,
	0x8c, 0xfa, 0xd7, 0x7b, 0x60, 0x45, 0x13, 0xa9, 0x96, 0x99, 0x31, 0xc2, 0xb3, 0x03, 0x97, 0x13,
	0x61, 0x38, 0x2d, 0x15, 0xf2, 0x48, 0x18, 0xfe, 0x9e, 0x1b, 0xf5, 0x77, 0x7a, 0xd7, 0xdb, 0xaa,
	0xe8, 0xf5, 0x65, 0x82, 0x8e, 0x60, 0xa6, 0xbb, 0xf0, 0x40, 0x85, 0xb8, 0x25, 0xac, 0xc1, 0x15,
	0xb6, 0xb5, 0x54, 0xeb, 0x79, 0xf4, 0x0f, 0x56, 0xfc, 0x1b, 0xbd, 0xeb, 0x06, 0x1a, 0xf3, 0xc8,
	0xc2, 0x49, 0x29, 0x3a, 0x3c, 0x82, 0x87, 0x2a, 0xfe, 0x4d, 0x5e, 0x65, 0x83, 0x63, 0x42, 0x1f,
	0x1e, 0x42, 0xf7, 0x33, 0x95, 0x74, 0x59, 0x4c, 0xe8, 0x23, 0x15, 0xe2, 0x9e, 0xd0, 0x49, 0x85,
	0x4c, 0xe3, 0x22, 0xf6, 0xfa, 0xd3, 0x3c, 0x46, 0x78, 0x74, 0xc8, 0xf8, 0x80, 0xe2, 0x05, 0xf4,
	0xb1, 0x8a, 0x7f, 0x9d, 0xad, 0x37, 0x42, 0xeb, 0x29, 0x8f, 0x43, 0x78, 0xbc, 0x42, 0xbc, 0xe4,
	0x52, 0x11, 0xc6, 0x08, 0x87, 0x2a, 0x54, 0xef, 0x85, 0x0b, 0x34, 0xd8, 0x32, 0x1c, 0xae, 0x10,
	0x5f, 0x24, 0x9f, 0x63, 0x2a, 0x41, 0x03, 0x1c, 0xa9, 0x94, 0xf9, 0xb2, 0x00, 0x85, 0x7d, 0x74,
	0x28, 0xae, 0x41, 0xd8, 0x0d, 0xae, 0xe0, 0xd8, 0x10, 0x9b, 0xd3, 0x52, 0xb5, 0x71, 0x1e, 0xdb,
	0x56, 0xa9, 0x21, 0xd7, 0x04, 0x1c, 0xaf, 0x50, 0xe5, 0xb8, 0x60, 0x52, 0x91, 0xb9, 0x80, 0x13,
	0x43, 0x77, 0x9a, 0x4f, 0xc5, 0x52, 0x1f, 0x4e, 0xba, 0x20, 0x67, 0x50, 0xcf, 0x1d, 0x30, 0x15,
	0x51, 0xe7, 0x82, 0xa9, 0x75, 0x38, 0xe5, 0xf8, 0xb3, 0xcc, 0x64, 0x90, 0xe1, 0x66, 0x2f, 0xb2,
	0x10, 0x15, 0x9c, 0x76, 0x76, 0x43, 0x30, 0x9c, 0xa9, 0xf8, 0x81, 0xb7, 0xd3, 0xb4, 0x6a, 0x96,
	0x8e, 0x0c, 0xca, 0x6e, 0x67, 0x15, 0xce, 0x56, 0xfc, 0x4f, 0x78, 0x63, 0x65, 0xcb, 0x01, 0x4c,
	0xee, 0xcf, 0x6d, 0x72, 0x7a, 0xc1, 0xc7, 0x13, 0x15, 0xff, 0x56, 0xef, 0xa6, 0x21, 0x38, 0x96,
	0x09, 0x2e, 0xb2, 0x4c, 0xa4, 0xe0, 0xc9, 0x41, 0xfe, 0xfb, 0xeb, 0x99, 0xc6, 0xa2, 0x9c, 0x94,
	0x42, 0x33, 0x2e, 0x50, 0xc1, 0xf9, 0x21, 0x26, 0x67, 0x50, 0xe7, 0x60, 0xd2, 0x14, 0x1d, 0x09,
	0x17, 0x2a, 0x34, 0x32, 0x68, 0xe6, 0xcc, 0xad, 0xf1, 0x3c, 0x08, 0xb8, 0xe8, 0xd2, 0x34, 0x83,
	0x7a, 0x29, 0x41, 0xd5, 0x6c, 0x4c, 0x2b, 0xd9, 0x33, 0x1e, 0xf0, 0xa0, 0x86, 0xaf, 0x05, 0x34,
	0x2e, 0xc8, 0x74, 0xb2, 0xcb, 0xe2, 0x18, 0x45, 0x84, 0x77, 0x99, 0xf2, 0xb5, 0x8d, 0x08, 0x5f,
	0x0f, 0xa8, 0x98, 0xa8, 0xa8, 0x17, 0x90, 0x25, 0x52, 0xc0, 0x37, 0x02, 0xe2, 0x75, 0x11, 0x59,
	0xcf, 0x0c, 0x50, 0x41, 0xc0, 0x37, 0x03, 0x0a, 0xd8, 0x44, 0xea, 0xfc, 0x2d, 0xa4, 0xcb, 0x49,
	0x5b, 0xf1, 0xbe, 0xf5, 0xf8, 0xad, 0x81, 0x47, 0xae, 0x17, 0x84, 0x5c, 0xeb, 0xc4, 0x6c, 0x05,
	0xe1, 0xdb, 0x01, 0xf1, 0xbd, 0xd4, 0x0f, 0x99, 0xc6, 0xcd, 0x6d, 0xbf, 0x13, 0x10, 0xa1, 0x59,
	0xb1, 0x6c, 0x16, 0xf0, 0x77, 0x03, 0x7f, 0xcc, 0xbb, 0x71, 0x28, 0x80, 0x02, 0x7e, 0x6f, 0xe0,
	0x6f, 0xf7, 0xae, 0x19, 0x5c, 0xc8, 0x5c, 0x00, 0xee, 0x73, 0x4c, 0xe4, 0x16, 0x13, 0xb1, 0x42,
	0x16, 0xae, 0xd3, 0xe9, 0xcb, 0x18, 0xc2, 0xfd, 0x2e, 0xc0, 0xa1, 0xb3, 0x4b, 0x01, 0x7e, 0x2f,
	0xa0, 0x29, 0x31, 0xcd, 0x45, 0x38, 0xab, 0x22, 0x26, 0xf8, 0x97, 0x68, 0xa2, 0x7d, 0x3f, 0xf0,
	0x6f, 0xf7, 0x82, 0x2c, 0xb0, 0x8c, 0x2c, 0x93, 0x8b, 0xec, 0xaf, 0xdc, 0x19, 0xfc, 0x20, 0xa0,
	0x7a, 0xa0, 0x8c, 0x99, 0xf0, 0x06, 0x7a, 0xf0, 0x43, 0xc7, 0x7b, 0x29, 0x1d, 0xcd, 0x06, 0xfc,
	0xc8, 0x5d, 0xdb, 0x18, 0xed, 0x65, 0x49, 0x4b, 0x5a, 0x4b, 0xa9, 0xc8, 0xf0, 0xc7, 0x01, 0x95,
	0x49, 0x7e, 0x7a, 0x7e, 0x66, 0x02, 0x3f, 0x09, 0x68, 0xb8, 0xe6, 0x20, 0xfc, 0x34, 0xa0, 0x36,
	0xcc, 0xfe, 0x6f, 0xa0, 0xe0, 0x18, 0xc2, 0xcf, 0x02, 0xea, 0x1a, 0xa2, 0x67, 0x2f, 0x4b, 0xca,
	0xc7, 0xfc, 0xdc, 0x99, 0xcd, 0x63, 0x82, 0x6a, 0x15, 0xc3, 0x16, 0xeb, 0x21, 0xfc, 0x22, 0xa7,
	0xae, 0x8b, 0xed, 0x95, 0x22, 0x2d, 0x4b, 0x82, 0xdf, 0x9d, 0xa2, 0x55, 0xfa, 0x65, 0xe0, 0xc6,
	0x8d, 0xe5, 0xb7, 0xa8, 0x05, 0xbf, 0x0a, 0xfc, 0x4f, 0x7b, 0x77, 0x4c, 0x29, 0x55, 0x94, 0x5e,
	0x29, 0x86, 0x07, 0x82, 0xc1, 0xac, 0x28, 0x79, 0x79, 0xd0, 0x9d, 0xb0, 0x91, 0x03, 0x78, 0xc8,
	0x9d, 0x30, 0xc9, 0x84, 0x90, 0xda, 0xcd, 0xb3, 0xcc, 0xaf, 0x8c, 0x65, 0xc9, 0xd1, 0xc3, 0x2e,
	0x49, 0x86, 0x6c, 0x5b, 0xfd, 0x25, 0xf8, 0x91, 0x80, 0x1e, 0xba, 0x81, 0x17, 0x78, 0x34, 0xf0,
	0x47, 0x3d, 0x2f, 0x3b, 0xdb, 0x0a, 0x1e, 0x0b, 0x68, 0xa5, 0x20, 0x41, 0x02, 0x8f, 0x17, 0x54,
	0x8c, 0x63, 0x38, 0xe4, 0xfc, 0x64, 0x2d, 0x61, 0x65, 0x87, 0xcb, 0x32, 0xeb, 0xea, 0x88, 0xbb,
	0x57, 0x26, 0x2b, 0xc5, 0x72, 0xd4, 0x15, 0x64, 0x0b, 0xd7, 0x8c, 0x03, 0xdb, 0xff, 0x31, 0xe3,
	0xbd, 0x04, 0x8e, 0xb9, 0x5c, 0x19, 0x9e, 0x26, 0x52, 0xdd, 0xb5, 0x07, 0x1c, 0x0f, 0xfc, 0xcf,
	0x78, 0xbb, 0xcd, 0xf3, 0xc9, 0x3b, 0x1d, 0x54, 0x28, 0x6c, 0x2c, 0x75, 0xd4, 0x6b, 0x88, 0x62,
	0x51, 0xae, 0xa0, 0x98, 0x10, 0x61, 0x83, 0x69, 0xb6, 0xcc, 0x12, 0x84, 0x13, 0x8e, 0xeb, 0x7d,
	0x92, 0x85, 0x46, 0x31, 0xe3, 0x35, 0x81, 0x93, 0x41, 0x79, 0xf2, 0x94, 0x7b, 0xe1, 0x94, 0xbb,
	0x45, 0x9e, 0x89, 0x04, 0x4e, 0x07, 0xf4, 0x24, 0x90, 0x45, 0xdd, 0x34, 0xdf, 0x17, 0xcd, 0x43,
	0x7f, 0xc6, 0x55, 0xdd, 0x54, 0x8f, 0xf1, 0x78, 0x22, 0x0c, 0x15, 0x26, 0x49, 0x4b, 0xea, 0xbb,
	0x50, 0xf1, 0x8e, 0x29, 0xcb, 0xb3, 0x05, 0xd3, 0x06, 0x76, 0x58, 0x1a, 0xbb, 0x32, 0x3e, 0x17,
	0x0c, 0x9e, 0xd7, 0x1e, 0xcf, 0x3a, 0x4a, 0x31, 0x91, 0xb0, 0xb6, 0x65, 0xe7, 0x89, 0x32, 0x73,
	0x13, 0x6d, 0xcd, 0x57, 0x91, 0x4c, 0x9f, 0x74, 0x1d, 0xe5, 0xa6, 0x63, 0x36, 0x35, 0xf7, 0xa3,
	0x66, 0x21, 0xd3, 0x0c, 0xce, 0xbb, 0xab, 0xb7, 0xa4, 0xa5, 0x65, 0x4e, 0xc9, 0x55, 0x1e, 0x62,
	0x08, 0x17, 0x0a, 0x65, 0x66, 0x91, 0x03, 0x5c, 0x77, 0x89, 0xf3, 0x8b, 0x2e, 0x52, 0x32, 0x6a,
	0x0a, 0x37, 0x8c, 0x9f, 0x2a, 0x36, 0x68, 0x76, 0x71, 0x93, 0x2b, 0xab, 0x05, 0x4f, 0x17, 0xa6,
	0x42, 0x01, 0x74, 0xb6, 0xcf, 0xb8, 0xb1, 0x38, 0x83, 0xba, 0x78, 0x87, 0xfd, 0xd8, 0x5b, 0x46,
	0x95, 0x74, 0x79, 0x1f, 0x9e, 0x2d, 0xb8, 0xb7, 0x3e, 0x8b, 0xf6, 0xcf, 0xb9, 0xab, 0x0e, 0x8f,
	0x3f, 0xfb, 0x58, 0x85, 0xf0, 0x7c, 0xa1, 0x56, 0x27, 0x22, 0xb3, 0x13, 0xbe, 0xe0, 0x26, 0xc6,
	0x02, 0x5b, 0xc5, 0x4c, 0xf4, 0xa2, 0x73, 0xb2, 0x8f, 0x27, 0x83, 0xc9, 0xdb, 0x14, 0x89, 0x66,
	0xa2, 0x8d, 0x09, 0xbc, 0x34, 0x98, 0x28, 0x6a, 0x15, 0xad, 0x16, 0x0a, 0xb8, 0x67, 0xb7, 0xdb,
	0x3d, 0xad, 0x74, 0x1e, 0x23, 0x23, 0x57, 0x33, 0x4c, 0xe3, 0x1a, 0x5b, 0x87, 0x2f, 0xef, 0xa6,
	0x42, 0x31, 0x8f, 0xc5, 0x3e, 0x19, 0x45, 0xa8, 0xe0, 0xbd, 0xaa, 0x73, 0xa4, 0x99, 0xd2, 0xc6,
	0x8e, 0xb7, 0x11, 0xde, 0xaf, 0x16, 0x34, 0x33, 0x67, 0xf0, 0x41, 0xd5, 0x4d, 0x02, 0x25, 0xd3,
	0xfe, 0x22, 0xaa, 0x1e, 0x17, 0x76, 0xf5, 0xfe, 0x73, 0xb5, 0xc0, 0xe7, 0xc2, 0x6c, 0xb6, 0xe8,
	0x1a, 0x46, 0xa6, 0x63, 0x16, 0x25, 0xf0, 0x17, 0x77, 0x42, 0x23, 0xed, 0xf5, 0xf3, 0x5a, 0xff,
	0x6b, 0x75, 0x30, 0x25, 0xcd, 0x56, 0xda, 0x91, 0xf0, 0xb7, 0xea, 0xa0, 0x85, 0x16, 0x16, 0x66,
	0x0f, 0x74, 0x25, 0xeb, 0x71, 0xf8, 0xb0, 0x2c, 0xa5, 0x2d, 0xfb, 0xef, 0x65, 0x29, 0x15, 0xc4,
	0x3f, 0xaa, 0x54, 0x10, 0x26, 0xec, 0x86, 0x6c, 0xaf, 0xa0, 0xa2, 0xb5, 0xfb, 0x9f, 0x55, 0xda,
	0x80, 0x2d, 0x52, 0x87, 0x7f, 0x55, 0xe9, 0xe1, 0xca, 0xde, 0xf6, 0x54, 0x61, 0xa3, 0x0e, 0xff,
	0xae, 0x16, 0x1f, 0x53, 0x77, 0x13, 0xf8, 0x4f, 0x35, 0x7f, 0xe4, 0x78, 0xce, 0xd0, 0x7f, 0xab,
	0x94, 0x20, 0x9b, 0xaf, 0x8d, 0x9b, 0xc4, 0x0b, 0x35, 0x2a, 0x11, 0x3b, 0xa5, 0x5b, 0x11, 0x17,
	0x07, 0x07, 0x8b, 0xc8, 0x8b, 0x35, 0xea, 0xe7, 0x79, 0xec, 0xc9, 0x55, 0x1c, 0x42, 0x5f, 0x72,
	0xa6, 0x76, 0x03, 0x1d, 0x02, 0x7f, 0xed, 0x40, 0x9b, 0xaf, 0x21, 0xf0, 0x37, 0x35, 0x4a, 0x91,
	0xd9, 0x2e, 0xb9, 0x88, 0xcc, 0x0a, 0x19, 0x9b, 0x3d, 0xf0, 0xb7, 0xb5, 0xe2, 0x6a, 0xb5, 0x61,
	0xf3, 0xfa, 0x5d, 0xad, 0xb8, 0xd8, 0x15, 0x76, 0xae, 0x97, 0x6b, 0x6e, 0xcd, 0x2d, 0x2f, 0x5a,
	0xaf, 0xd4, 0xdc, 0x13, 0x2f, 0xfb, 0xeb, 0x2e, 0x88, 0x0e, 0x8f, 0x8a, 0xdb, 0xd6, 0xa5, 0x1a,
	0xf5, 0x91, 0xc5, 0x5b, 0xb8, 0x96, 0xa9, 0x58, 0x3e, 0xb2, 0x2f, 0x2e, 0xb8, 0x5c, 0xf3, 0xef,
	0xf0, 0x6e, 0x75, 0x2a, 0x0b, 0x28, 0x42, 0x33, 0x59, 0x98, 0x08, 0xcb, 0xda, 0xf0, 0xfb, 0x9a,
	0xbf, 0xdb, 0xbb, 0xed, 0xa3, 0xf4, 0x32, 0x22, 0xe1, 0x0f, 0x35, 0xff, 0x93, 0xde, 0x2d, 0x9b,
	0x28, 0x3a, 0xad, 0x7e, 0xcc, 0xda, 0x08, 0x7f, 0xac, 0xd1, 0xf6, 0x30, 0xac, 0x36, 0x8f, 0xb1,
	0xcc, 0xbf, 0x44, 0x5e, 0x75, 0x54, 0xbb, 0x0b, 0x9a, 0x0f, 0xa5, 0x16, 0xea, 0x35, 0xa9, 0x56,
	0xe0, 0x4f, 0x35, 0x1a, 0xa3, 0xf9, 0x85, 0x87, 0x14, 0x5e, 0x73, 0xd4, 0xb5, 0x98, 0x9e, 0x93,
	0x4a, 0xcf, 0xf6, 0x51, 0x70, 0x11, 0xc1, 0xeb, 0x35, 0xaa, 0xd1, 0x52, 0x76, 0xcd, 0x79, 0x6f,
	0xb8, 0x2c, 0x4c, 0x1d, 0xc4, 0x76, 0x9a, 0x7d, 0x00, 0xd8, 0xec, 0xbd, 0xe9, 0xce, 0xb2, 0xec,
	0xd7, 0xd7, 0x35, 0x26, 0x8b, 0x72, 0x2f, 0x4b, 0xba, 0xd6, 0x05, 0x2a, 0x78, 0xab, 0x46, 0x7b,
	0xa2, 0xf9, 0x0c, 0xb1, 0xb8, 0x69, 0xbf, 0xa2, 0xc6, 0xdb, 0xb5, 0x7c, 0x8c, 0x0a, 0x34, 0x5f,
	0x76, 0x73, 0x0a, 0x3b, 0xfc, 0xa0, 0x51, 0x81, 0x77, 0x5c, 0x71, 0x4c, 0xc6, 0xc8, 0xc4, 0x5c,
	0xf6, 0x5b, 0xc1, 0x60, 0xd4, 0xbc, 0x5b, 0x2c, 0x2a, 0x1c, 0x2c, 0xe5, 0xf0, 0x5e, 0x8d, 0x5e,
	0xc3, 0xa5, 0xfe, 0x90, 0x11, 0xbc, 0x5f, 0xa3, 0x96, 0xc9, 0x9e, 0x02, 0x7b, 0x4b, 0xf8, 0xa0,
	0x46, 0x2d, 0x93, 0x75, 0xe6, 0xc4, 0x5c, 0x33, 0xe7, 0xce, 0xcc, 0x2f, 0x38, 0x34, 0x4e, 0xb7,
	0xd8, 0x88, 0x53, 0x7a, 0x0f, 0x8f, 0x53, 0xdf, 0xe4, 0x1a, 0xcd, 0x1e, 0x8b, 0x90, 0xd0, 0x23,
	0x57, 0xb6, 0xa7, 0xcf, 0xa1, 0xa3, 0xe3, 0x94, 0xf7, 0x8d, 0x1a, 0x86, 0x73, 0xd2, 0x3a, 0xf6,
	0xd1, 0x5a, 0x13, 0x5a, 0xb3, 0x76, 0x17, 0x8e, 0x8f, 0xd3, 0x1a, 0xb6, 0xb9, 0x96, 0x6d, 0x4f,
	0x38, 0x31, 0x4e, 0xf5, 0xb8, 0xb9, 0x52, 0x53, 0x24, 0x7d, 0xf3, 0x26, 0x9f, 0x1c, 0xa7, 0xec,
	0x94, 0xef, 0x35, 0x97, 0xc6, 0x31, 0x9c, 0xda, 0x70, 0x67, 0xaa, 0x30, 0xcb, 0xd9, 0xe9, 0xf1,
	0x61, 0x4e, 0x09, 0xa5, 0xbb, 0x9c, 0xb9, 0x12, 0x4e, 0x9c, 0x9d, 0x1d, 0xa7, 0x1c, 0xe6, 0xf8,
	0xd4, 0x41, 0x93, 0xe0, 0x10, 0xe1, 0x9c, 0x83, 0xa8, 0x1c, 0x67, 0x85, 0xc9, 0xfd, 0x5e, 0x29,
	0x57, 0xe0, 0xde, 0x69, 0xaa, 0xe1, 0xcc, 0x4b, 0xa1, 0x26, 0xee, 0x9b, 0xae, 0x7f, 0xfe, 0xfc,
	0x2b, 0x63, 0x5b, 0x4e, 0x5e, 0x1a, 0xbb, 0xea, 0xfc, 0xa5, 0xb1, 0xab, 0x5e, 0xbe, 0x34, 0x76,
	0xd5, 0x57, 0x2e, 0x8f, 0x6d, 0x39, 0x7f, 0x79, 0x6c, 0xcb, 0xf3, 0x97, 0xc7, 0xb6, 0x7c, 0x61,
	0x97, 0xfb, 0x2d, 0x2a, 0x66, 0x22, 0xdc, 0x13, 0xc9, 0x3d, 0xfd, 0x95, 0x68, 0x0f, 0xfd, 0x2e,
	0xb5, 0x7c, 0xb5, 0xfd, 0xbd, 0xe9, 0x73, 0xff, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xaa, 0x28, 0x8b,
	0x67, 0xc0, 0x12, 0x00, 0x00,
}
