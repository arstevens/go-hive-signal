package register

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestRegistrationHandler(t *testing.T) {
	dataspaceRequests := 20
	originRequests := 20
	preloadDataspaces := 20
	totalSwarms := 5

	swarms := make([]string, totalSwarms)
	for i := 0; i < totalSwarms; i++ {
		swarms[i] = "Swarm#" + strconv.Itoa(i)
	}

	swarmMap := SwarmMapTest{smap: make(map[string]string)}
	for i := 0; i < preloadDataspaces; i++ {
		dspace := "/dataspace/" + strconv.Itoa(i)
		swarmMap.smap[dspace] = swarms[i%totalSwarms]
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
	smap map[string]string
}

func (st *SwarmMapTest) GetSwarmID(d string) (string, error) {
	sm, ok := st.smap[d]
	if !ok {
		return "", fmt.Errorf("No swarm associated with name %s", d)
	}
	return sm, nil
}

func (st *SwarmMapTest) AddDataspace(sid string, dspace string) error {
	fmt.Printf("Adding dataspace %s\n", dspace)
	st.smap[dspace] = sid
	return nil
}

func (st *SwarmMapTest) RemoveDataspace(sid string, dspace string) error {
	if _, ok := st.smap[dspace]; !ok {
		return fmt.Errorf("No dataspace with name %s", dspace)
	}
	fmt.Printf("Removing dataspace %s\n", dspace)
	delete(st.smap, dspace)
	return nil
}

func (st *SwarmMapTest) GetMinDataspaceSwarm() (string, error) {
	if len(st.smap) == 0 {
		return "", fmt.Errorf("No dataspaces available")
	}
	n := rand.Intn(len(st.smap))
	i := 0
	for _, swarm := range st.smap {
		if i == n {
			return swarm, nil
		}
		i++
	}
	return "", nil
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
