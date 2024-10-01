package migrator

import (
	"errors"
	"fmt"

	"mig/pkg/migrator/matcher_v1"
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
			matcher_v1.MatchErrorfWithNamedParams,
			matcher_v1.MatchWrapfWithNamedParams,
			matcher_v1.MatchWrapfStderr,
			matcher_v1.MatchSimpleWraps,
			matcher_v1.MatchSimpleErrorsNew,
			matcher_v1.MatchImport,
		}, nil
	case V2:
		return MigrationHandlers{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("migrator: unknown version %v", version))
	}
}
