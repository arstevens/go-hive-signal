package configuration

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arstevens/go-hive-signal/internal/analyzer"
	"github.com/arstevens/go-hive-signal/internal/cache"
	"github.com/arstevens/go-hive-signal/internal/comparator"
	"github.com/arstevens/go-hive-signal/internal/connector"
	"github.com/arstevens/go-hive-signal/internal/debriefer"
	"github.com/arstevens/go-hive-signal/internal/gateway"
	"github.com/arstevens/go-hive-signal/internal/localizer"
	"github.com/arstevens/go-hive-signal/internal/manager"
	"github.com/arstevens/go-hive-signal/internal/mapper"
	"github.com/arstevens/go-hive-signal/internal/negotiator"
	"github.com/arstevens/go-hive-signal/internal/register"
	"github.com/arstevens/go-hive-signal/internal/registrator"
	"github.com/arstevens/go-hive-signal/internal/tracker"
	"github.com/arstevens/go-hive-signal/internal/transmuter"
	"github.com/arstevens/go-hive-signal/internal/verifier"
	"github.com/arstevens/go-request/handle"
	"github.com/arstevens/go-request/route"
)

var TimeForTeardown = time.Second * 5

func LinkProgram() (func(), error) {
	endpointRegister, err := register.New()
	if err != nil {
		return nil, fmt.Errorf("Failed to link program in configuration.LinkProgram(): %v", err)
	}
	connectionCache := cache.New()
	identityVerifier := verifier.New(endpointRegister, connectionCache)

	lpseGenerator := debriefer.NewLPSEGenerator(debrieferLoadHistorySize)
	gatewayGenerator := gateway.NewGenerator(gatewayActiveQueueSize, gatewayInactiveQueueSize)
	infoTracker := tracker.New(lpseGenerator, trackerLoadHistorySize)
	managerGenerator := manager.NewGenerator(gatewayGenerator, negotiator.RoundtripLimitedNegotiate,
		infoTracker)
	loadToSizeComparator := comparator.New(infoTracker)
	dataAnalyzer := analyzer.New(infoTracker, loadToSizeComparator)
	swarmMap := mapper.New(managerGenerator)
	swarmTransmuter := transmuter.New(swarmMap, dataAnalyzer)

	requestLocalizer := localizer.New(localizerQueueSize, swarmMap, infoTracker)
	registrationHandler := registrator.New(registratorQueueSize, swarmMap, endpointRegister)
	connectionHandler := connector.New(connectorQueueSize, identityVerifier, swarmTransmuter)

	routeMap := map[int32]handle.RequestHandler{
		localizerRoutingCode:   requestLocalizer,
		registratorRoutingCode: registrationHandler,
		connectorRoutingCode:   connectionHandler,
	}
	unpackersMap := map[int32]handle.UnpackRequest{
		localizerRoutingCode:   localizerUnpacker,
		registratorRoutingCode: registratorUnpacker,
		connectorRoutingCode:   connectorUnpacker,
	}

	listenAddr := fmt.Sprintf(":%d", listenerPort)
	netListener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("Failed to link program in configuration.LinkProgram(): %v", err)
	}
	routeListener := route.NewNetListener(netListener)

	return func() {
		done := make(chan struct{})
		go route.UnpackAndRoute(routeListener, done, routeMap, routeUnpacker, unpackersMap, route.ReadRequestFromNetConn)
		log.Println("Hive Signal started successfully")

		sigterm := make(chan os.Signal)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

		<-sigterm
		close(done)
		time.Sleep(TimeForTeardown)
	}, nil
}
