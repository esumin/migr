package migrator

import (
	"errors"
	"fmt"

	common "mig/pkg/migrator/common"
	"mig/pkg/migrator/matcher_v1"
	"mig/pkg/migrator/matcher_v2"
)

type HandleLine func(string) string
type MigrationHandlers []HandleLine

type MigratorVersion string

const V1 MigratorVersion = "v1"
const V2 MigratorVersion = "v2"

func GetMigratorHandlers(version MigratorVersion) (MigrationHandlers, error) {
	switch version {
	case V1:
		return MigrationHandlers{
			common.MatchImport,
			matcher_v1.MatchErrorfWithNamedParams,
			matcher_v1.MatchWrapfWithNamedParams,
			matcher_v1.MatchWrapfStderr,
			matcher_v1.MatchSimpleWraps,
			matcher_v1.MatchSimpleErrorsNew,
		}, nil
	case V2:
		return MigrationHandlers{
			common.MatchImport,
			matcher_v2.HandleLine,
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("migrator: unknown version %v", version))
	}
}
