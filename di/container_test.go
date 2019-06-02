package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testSvc struct{}

func (*testSvc) TestMethod() {}

func newTestSvc() *testSvc { return &testSvc{} }

var testSvcName = svc.GetTypeName(reflect.TypeOf(newTestSvc()))

type ContainerTestSuite struct {
	suite.Suite
}

func (s *ContainerTestSuite) SetupTest() {
	svc = newServiceContainer()
}

func TestContainerTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}

func (s *ContainerTestSuite) TestInit() {
	s.NotNil(svc, "Service container not initialized")
	s.NotNil(svc.Constructors, "Constructors map not initialized")
	s.NotNil(svc.Services, "Services map not initialized")
	s.NotNil(svc.Constructors["di.Container"], "Container does not have container constructor")
}

func (s *ContainerTestSuite) TestGetContainer() {
	s.EqualValues(svc, GetContainer(), "Should receive service container instance")
}

func (s *ContainerTestSuite) TestGetTypeName() {
	typeName := reflect.TypeOf(newTestSvc())
	s.Equal(svc.GetTypeName(typeName), testSvcName, "Wrong type name for pointer to struct")

	typeName = reflect.TypeOf(*newTestSvc())
	s.Equal(svc.GetTypeName(typeName), testSvcName, "Wrong type name for struct")
}

func (s *ContainerTestSuite) TestGetRegisteredService() {
	c := GetContainer()
	c.Register(newTestSvc)

	ret := c.Get("di.testSvc")
	s.Equal(reflect.TypeOf(ret), reflect.TypeOf(newTestSvc()))
}

func (s *ContainerTestSuite) TestGetReturnsSingleton() {
	c := GetContainer()
	c.Register(newTestSvc)

	inst1 := c.Get(testSvcName)
	inst2 := c.Get(testSvcName)
	s.EqualValues(inst1, inst2, "Get did not return singleton")
}

func (s *ContainerTestSuite) TestContainerPanicsOnGetNonExistentService() {
	c := GetContainer()
	s.Panics(func() { c.Get(testSvcName) }, "Container should have panicked")
}

func (s *ContainerTestSuite) TestGetMethod() {
	c := GetContainer()
	c.Register(newTestSvc)
	m := c.GetMethod(testSvcName, "TestMethod")
	s.Equal(m.IsNil(), false, "Could not find test method")
}

func (s *ContainerTestSuite) TestGetMethodPanicsIfServiceNotRegistered() {
	c := GetContainer()
	s.Panics(func() { c.GetMethod(testSvcName, "TestMethod") }, "Container should panic when getting method of non-registered service")
}

func (s *ContainerTestSuite) TestCallFunc() {
	c := GetContainer()
	c.Register(newTestSvc)
	values := svc.Call(func(t *testSvc, t2 *testSvc) (bool, error) {
		s.EqualValues(t, t2, "Should have two references to singleton instance")
		return false, nil
	})

	s.Equal(len(values), 2, "Did not get all return values")
	s.Equal(values[0].Bool(), false, "Did not get correct return value 1")
	s.Nil(values[1].Interface(), "Did not get correct return value 2")
}

func (s *ContainerTestSuite) TestInvoke() {
	c := GetContainer()
	c.Register(newTestSvc)
	c.Invoke(func(t *testSvc, t2 *testSvc) (bool, error) {
		s.EqualValues(t, t2, "Should have two references to singleton instance")
		return false, nil
	})
}
