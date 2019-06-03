package pubsub

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kgrunwald/goweb/ilog"
	mock_log "github.com/kgrunwald/goweb/ilog/mock_ilog"
	"github.com/stretchr/testify/suite"
)

type Msg interface {
	Get() bool
}
type T struct{}

func (*T) Get() bool {
	return true
}

type TestSuite struct {
	suite.Suite
	Ctrl   *gomock.Controller
	Logger ilog.Logger
}

func TestPubSub(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) SetupTest() {
	t.Ctrl = gomock.NewController(t.Suite.T())
	t.Logger = mock_log.NewMockLogger(t.Ctrl)
}

func (t *TestSuite) TestSubscribe() {
	bus := newEventBus(t.Logger)

	bus.Subscribe(func() {})
	t.Equal(1, len(bus.Subscriptions))
}

func (t *TestSuite) TestDispatch() {
	bus := newEventBus(t.Logger)

	c := make(chan bool)
	bus.Subscribe(func(msg *T) {
		t.Equal(reflect.TypeOf(&T{}), reflect.TypeOf(msg))
		c <- true
	})
	bus.Dispatch(&T{})
	<-c
}

func (t *TestSuite) TestSubscribeInterface() {
	bus := newEventBus(t.Logger)

	c := make(chan bool)
	bus.Subscribe(func(msg Msg) {
		t.Equal(reflect.TypeOf(&T{}), reflect.TypeOf(msg))
		c <- true
	})
	bus.Dispatch(&T{})
	<-c
}

func (t *TestSuite) TestNewBus() {
	b := NewBus(t.Logger)
	t.Equal(reflect.TypeOf(&EventBus{}), reflect.TypeOf(b))
}
