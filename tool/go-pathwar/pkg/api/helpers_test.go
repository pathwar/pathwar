package api

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetUsers(t *testing.T) {
	Convey("Testing GetUsers", t, func() {
		Convey("without where clause", func() {
			client := NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
			users, err := client.GetUsers(nil)
			So(err, ShouldBeNil)
			So(len(users.Items) > 0, ShouldBeTrue)
			So(users.Items[0].Id, ShouldNotBeNil)
		})
		Convey("with where clause", func() {
			client := NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
			users, err := client.GetUsers(map[string]string{"login": "moul"})
			So(err, ShouldBeNil)
			So(len(users.Items), ShouldEqual, 1)
			So(users.Items[0].Login, ShouldEqual, "moul")
		})
	})
}

func TestGetRawOrganizationUsers(t *testing.T) {
	Convey("Testing GetRawOrganizationUsers", t, func() {
		Convey("without where clause", func() {
			client := NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
			rawOrganizationUsers, err := client.GetRawOrganizationUsers(nil)
			So(err, ShouldBeNil)
			So(len(rawOrganizationUsers.Items), ShouldNotEqual, 0)
			So(rawOrganizationUsers.Items[0].Id, ShouldNotBeNil)
		})
		Convey("with where clause", func() {
			client := NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
			rawOrganizationUsers, err := client.GetRawOrganizationUsers(map[string]string{"role": "pwner"})
			So(err, ShouldBeNil)
			So(len(rawOrganizationUsers.Items), ShouldNotEqual, 0)
			So(rawOrganizationUsers.Items[0].Role, ShouldEqual, "pwner")
		})
	})
}

func TestGetRawLevelInstanceUsers(t *testing.T) {
	Convey("Testing GetRawLevelInstanceUsers", t, func() {
		Convey("without where clause", func() {
			client := NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
			rawLevelInstanceUsers, err := client.GetRawLevelInstanceUsers(nil)
			So(err, ShouldBeNil)
			So(len(rawLevelInstanceUsers.Items), ShouldNotEqual, 0)
			So(rawLevelInstanceUsers.Items[0].Id, ShouldNotBeNil)
		})
		/*Convey("with where clause", func() {
			client := NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
			rawLevelInstanceUsers, err := client.GetRawLevelInstanceUsers(map[string]string{"role": "pwner"})
			So(err, ShouldBeNil)
			So(len(rawLevelInstanceUsers.Items), ShouldNotEqual, 0)
			So(rawLevelInstanceUsers.Items[0].Role, ShouldEqual, "pwner")
		})*/
	})
}

func TestGetRawLevelInstances(t *testing.T) {
	Convey("Testing GetRawLevelInstances", t, func() {
		Convey("without where clause", func() {
			client := NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
			rawLevelInstances, err := client.GetRawLevelInstances(nil)
			So(err, ShouldBeNil)
			So(len(rawLevelInstances.Items), ShouldNotEqual, 0)
			So(rawLevelInstances.Items[0].Id, ShouldNotBeNil)
		})
		Convey("with where clause", func() {
			client := NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
			rawLevelInstances, err := client.GetRawLevelInstances(map[string]bool{"active": true})
			So(err, ShouldBeNil)
			So(len(rawLevelInstances.Items), ShouldNotEqual, 0)
			So(rawLevelInstances.Items[0].Active, ShouldEqual, true)
		})
	})
}
