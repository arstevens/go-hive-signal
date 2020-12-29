package registrator

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestRegistrationHandler(t *testing.T) {
	dataspaceRequests := 20
	originRequests := 20
	preloadDataspaces := 20

	dspaces := make([]string, preloadDataspaces)
	swarmMap := SwarmMapTest{smap: make(map[string]bool)}
	for i := 0; i < preloadDataspaces; i++ {
		dspace := "/dataspace/" + strconv.Itoa(i)
		swarmMap.smap[dspace] = true
		dspaces[i] = dspace
	}
	originReg := OriginRegistratorTest{origins: make(map[string]bool)}

	requests := make([]RegistrationRequestTest, dataspaceRequests+originRequests)
	killCount := 0
	for i := 0; i < dataspaceRequests; i++ {
		if i%2 == 0 {
			dspace := "/dataspace/" + strconv.Itoa(1000+i)
			requests[i] = RegistrationRequestTest{isOrigin: false, isAdd: true, datafield: dspace}
		} else {
			dspace := "/dataspace/" + strconv.Itoa(killCount)
			requests[i] = RegistrationRequestTest{isOrigin: false, isAdd: false, datafield: dspace}
			killCount++
		}
	}

	killCount = 0
	for i := dataspaceRequests; i < dataspaceRequests+originRequests; i++ {
		if i < (dataspaceRequests + (originRequests / 2)) {
			origin := "/origin/" + strconv.Itoa(i)
			requests[i] = RegistrationRequestTest{isOrigin: true, isAdd: true, datafield: origin}
		} else {
			origin := "/origin/" + strconv.Itoa(dataspaceRequests+killCount)
			killCount++
			requests[i] = RegistrationRequestTest{isOrigin: true, isAdd: false, datafield: origin}
		}
	}

	fconn := FakeConn{}
	queueSize := 3
	regHandler := New(queueSize, &swarmMap, &originReg)

	for _, request := range requests {
		regHandler.AddJob(&RegistrationRequestTest{
			isAdd:     request.IsAdd(),
			isOrigin:  request.IsOrigin(),
			datafield: request.GetDataField(),
		}, &fconn)
	}
	time.Sleep(time.Second)
}

type FakeConn struct{}

func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { return nil }

type SwarmMapTest struct {
	smap map[string]bool
}

func (sm *SwarmMapTest) AddSwarm(dspace string) error {
	sm.smap[dspace] = true
	return nil
}
func (sm *SwarmMapTest) RemoveSwarm(dspace string) error {
	if _, ok := sm.smap[dspace]; ok {
		delete(sm.smap, dspace)
	}
	return nil
}

type OriginRegistratorTest struct {
	origins map[string]bool
}

func (ot *OriginRegistratorTest) AddOrigin(s string) error {
	fmt.Printf("Adding origin %s\n", s)
	ot.origins[s] = true
	return nil
}

func (ot *OriginRegistratorTest) RemoveOrigin(s string) error {
	if _, ok := ot.origins[s]; !ok {
		return fmt.Errorf("Cannot remove nonexisted registration entity %s", s)
	}
	fmt.Printf("Removing origin %s\n", s)
	delete(ot.origins, s)
	return nil
}

type RegistrationRequestTest struct {
	isAdd     bool
	isOrigin  bool
	datafield string
}

func (rt *RegistrationRequestTest) IsAdd() bool          { return rt.isAdd }
func (rt *RegistrationRequestTest) IsOrigin() bool       { return rt.isOrigin }
func (rt *RegistrationRequestTest) GetDataField() string { return rt.datafield }
