package fleet

import (
	"io/ioutil"
	"path"
	"reflect"
	"sync"
	"testing"

	"github.com/coreos/fleet/schema"
	"github.com/brendangibat/deis/deisctl/units"
)

func TestUnits(t *testing.T) {
	t.Parallel()

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name: "deis-router@1.service",
		},
		&schema.Unit{
			Name: "deis-router@2.service",
		},
		&schema.Unit{
			Name: "deis-router@3.service",
		},
	}

	c := &FleetClient{Fleet: &stubFleetClient{testUnits: testUnits, unitsMutex: &sync.Mutex{}}}

	targets, err := c.Units("router")

	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"deis-router@1.service", "deis-router@2.service", "deis-router@3.service"}

	if !reflect.DeepEqual(targets, expected) {
		t.Fatalf("Expected %v, Got %v", expected, targets)
	}
}

func TestNextUnit(t *testing.T) {
	t.Parallel()

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name: "deis-router@1.service",
		},
		&schema.Unit{
			Name: "deis-router@3.service",
		},
	}

	c := &FleetClient{Fleet: &stubFleetClient{testUnits: testUnits, unitsMutex: &sync.Mutex{}}}

	num, err := c.nextUnit("router")

	if err != nil {
		t.Fatal(err)
	}

	expected := 2

	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
}

func TestLastUnit(t *testing.T) {
	t.Parallel()

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name: "deis-router@1.service",
		},
		&schema.Unit{
			Name: "deis-router@3.service",
		},
	}

	c := &FleetClient{Fleet: &stubFleetClient{testUnits: testUnits, unitsMutex: &sync.Mutex{}}}

	num, err := c.lastUnit("router")

	if err != nil {
		t.Fatal(err)
	}

	expected := 3

	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
}

func TestFormatUnitName(t *testing.T) {
	t.Parallel()

	unitName, err := formatUnitName("router", 1)
	if err != nil {
		t.Fatal(err)
	}
	if unitName != "deis-router@1.service" {
		t.Fatalf("invalid unit name: %v", unitName)
	}

}

func TestNewUnit(t *testing.T) {
	t.Parallel()

	name, err := ioutil.TempDir("", "deisctl-fleetctl")
	unitFile := `[Unit]
Description=deis-controller`

	unit := "deis-controller"

	ioutil.WriteFile(path.Join(name, unit+".service"), []byte(unitFile), 777)

	uf, err := NewUnit(unit[5:], []string{name})

	if err != nil {
		t.Fatal(err)
	}

	result := uf.Contents["Unit"]["Description"][0]
	expected := unitFile[19:]
	if result != expected {
		t.Errorf("Expected: %s, Got %s", expected)
	}
}

func TestReadTemplate(t *testing.T) {
	t.Parallel()

	name, err := ioutil.TempDir("", "deisctl-fleetctl")
	expected := []byte("test")

	if err != nil {
		t.Error(err)
	}

	for _, unit := range units.Names {
		ioutil.WriteFile(path.Join(name, unit+".service"), expected, 777)
		output, err := readTemplate(unit[5:], []string{name})

		if err != nil {
			t.Error(err)
		}

		if string(output) != string(expected) {
			t.Errorf("Unit %s: Expected %s, Got %s", unit, expected, output)
		}
	}
}

func TestReadTemplateError(t *testing.T) {
	t.Parallel()

	name, err := ioutil.TempDir("", "deisctl-fleetctl")

	if err != nil {
		t.Error(err)
	}

	_, err = readTemplate("foo", []string{name})
	expected := "Could not find unit template for foo"
	errorf := err.Error()

	if errorf != expected {
		t.Errorf("Expected %s, Got %s", expected, errorf)
	}
}
