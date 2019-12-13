package pwapi

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwsso"
)

func testingSeasons(t *testing.T, svc Service) *pwdb.SeasonList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.SeasonList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	assert.NoError(t, err, "list seasons")
	return &list
}

func testingAgents(t *testing.T, svc Service) *pwdb.AgentList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.AgentList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	if err != nil {
		t.Fatalf("list agents: %v", err)
	}

	return &list
}

func testingSeasonChallenges(t *testing.T, svc Service) *pwdb.SeasonChallengeList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.SeasonChallengeList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	assert.NoError(t, err, "list season challenges")
	return &list
}

func testingSoloSeason(t *testing.T, svc Service) *pwdb.Season {
	t.Helper()

	seasons := testingSeasons(t, svc)
	for _, season := range seasons.Items {
		if season.Name == "Solo Mode" {
			return season
		}
	}

	t.Fatalf("no such solo season")
	return nil
}

func testingTeams(t *testing.T, svc Service) *pwdb.TeamList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.TeamList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	assert.NoError(t, err, "list teams")
	return &list
}

func testingChallenges(t *testing.T, svc Service) *pwdb.ChallengeList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.ChallengeList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	assert.NoError(t, err, "list challenges")
	return &list
}

func testingSvcDB(t *testing.T, svc Service) *gorm.DB {
	t.Helper()

	typed := svc.(*service)
	return typed.db
}

func testingSetContextToken(ctx context.Context, t *testing.T) context.Context {
	t.Helper()

	return context.WithValue(ctx, userTokenCtx, pwsso.TestingToken(t))
}

func checkErr(t *testing.T, name string, err error) {
	t.Helper()

	if !assert.NoError(t, err) {
		t.Fatal(name)
	}
}

func testSameErrcodes(t *testing.T, name string, expected, got error) {
	t.Helper()

	assert.Equal(t, errcode.Code(expected), errcode.Code(got))
}

func testIsTrue(t *testing.T, name string, got bool) {
	t.Helper()

	if !got {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected %v to be true.", prefix, got)
	}
}

func testIsNil(t *testing.T, name string, got interface{}) {
	t.Helper()

	if got != nil {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected %+v to be nil.", prefix, got)
	}
}

func testIsNotNil(t *testing.T, name string, got interface{}) {
	t.Helper()

	if got == nil {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected %+v to be not nil.", prefix, got)
	}
}

func testSameErrs(t *testing.T, name string, expected, got error) {
	t.Helper()

	if !errors.Is(got, expected) {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected %+v, got %+v.", prefix, expected, got)
	}
}

func testSameAnys(t *testing.T, name string, expected, got interface{}) {
	t.Helper()

	if expected != got {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected %#v, got %#v.", prefix, expected, got)
	}
}

func testSameInt64s(t *testing.T, name string, expected, got int64) {
	t.Helper()

	if expected != got {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected %d, got %d.", prefix, expected, got)
	}
}

func testSameInts(t *testing.T, name string, expected, got int) {
	t.Helper()

	if expected != got {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected %d, got %d.", prefix, expected, got)
	}
}

func testDifferentInt64s(t *testing.T, name string, expected, got int64) {
	t.Helper()

	if expected == got {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected different values, got %d two times.", prefix, expected)
	}
}

func testSameStrings(t *testing.T, name string, expected, got string) {
	t.Helper()

	if expected != got {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected %q, got %q.", prefix, expected, got)
	}
}

func testSameDeep(t *testing.T, name string, expected, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(expected, got) {
		prefix := ""
		if name != "" {
			prefix = name + ": "
		}
		t.Errorf("%sExpected %#v, got %#v.", prefix, expected, got)
	}
}
