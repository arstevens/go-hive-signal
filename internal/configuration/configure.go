package configuration

import (
	"log"
	"time"

	"github.com/arstevens/go-hive-signal/internal/analyzer"
	"github.com/arstevens/go-hive-signal/internal/cache"
	"github.com/arstevens/go-hive-signal/internal/comparator"
	"github.com/arstevens/go-hive-signal/internal/debriefer"
	"github.com/arstevens/go-hive-signal/internal/manager"
	"github.com/arstevens/go-hive-signal/internal/negotiator"
	"github.com/arstevens/go-hive-signal/internal/register"
	"github.com/arstevens/go-hive-signal/internal/tracker"
	"github.com/arstevens/go-hive-signal/internal/transmuter"
	"github.com/arstevens/go-hive-signal/pkg/protomsg"
)

const (
	MissingSettingError = "[ERROR] Missing required setting %s"
	InvalidOptionError  = "[ERROR] Unknown value %s specified for %s"
)

var requestQueueSizeKey = "RequestBufferSize"
var (
	connectorQueueSize   = 30
	registratorQueueSize = 30
	localizerQueueSize   = 30

	gatewayActiveQueueSize = 0

	trackerLoadHistorySize   = 0
	debrieferLoadHistorySize = 0
)

var (
	localizerRoutingCode   int32 = 0
	registratorRoutingCode int32 = 1
	connectorRoutingCode   int32 = 2
)

var (
	localizerUnpacker   = protomsg.UnpackLocalizeRequest
	registratorUnpacker = protomsg.UnpackRegistrationRequest
	connectorUnpacker   = protomsg.UnpackConnectionRequest
	routeUnpacker       = protomsg.UnpackRouteWrapper
)

var (
	RouterKey      = "Router"
	ListenerKey    = "Listener"
	MessagingKey   = "Messaging"
	AnalyzerKey    = "Analyzer"
	CacheKey       = "Cache"
	ComparatorKey  = "Comparator"
	ConnectorKey   = "Connector"
	DebrieferKey   = "Debriefer"
	GatewayKey     = "Gateway"
	LocalizerKey   = "Localizer"
	ManagerKey     = "Manager"
	NegotiatorKey  = "Negotiator"
	RegisterKey    = "Register"
	RegistratorKey = "Registrator"
	TrackerKey     = "Tracker"
	TransmuterKey  = "Transmuter"
)

var (
	listenerPort = 10000
)

var UnitOfTime = time.Millisecond

type ComponentConfigurator func(map[string]interface{})

func GenerateConfiguratorMap() map[string]ComponentConfigurator {
	confMap := make(map[string]ComponentConfigurator)
	confMap[RouterKey] = ConfigureRouter
	confMap[ListenerKey] = ConfigureListener
	confMap[MessagingKey] = ConfigureMessaging
	confMap[AnalyzerKey] = ConfigureAnalyzer
	confMap[CacheKey] = ConfigureCache
	confMap[ComparatorKey] = ConfigureComparator
	confMap[ConnectorKey] = ConfigureConnector
	confMap[DebrieferKey] = ConfigureDebriefer
	confMap[GatewayKey] = ConfigureGateway
	confMap[LocalizerKey] = ConfigureLocalizer
	confMap[ManagerKey] = ConfigureManager
	confMap[NegotiatorKey] = ConfigureNegotiator
	confMap[RegisterKey] = ConfigureRegister
	confMap[RegistratorKey] = ConfigureRegistrator
	confMap[TrackerKey] = ConfigureTracker
	confMap[TransmuterKey] = ConfigureTransmuter

	return confMap
}

func ConfigureRouter(config map[string]interface{}) {
	LocalizerRoutingKey := "LocalizerRoutingCode"
	RegistratorRoutingKey := "RegistratorRoutingCode"
	ConnectorRoutingKey := "ConnectorRoutingCode"

	if lr, ok := config[LocalizerRoutingKey]; ok {
		localizerRoutingCode = int32(lr.(float64))
	}
	if rr, ok := config[RegistratorRoutingKey]; ok {
		registratorRoutingCode = int32(rr.(float64))
	}
	if cr, ok := config[ConnectorRoutingKey]; ok {
		connectorRoutingCode = int32(cr.(float64))
	}
}

func ConfigureListener(config map[string]interface{}) {
	ListenerPortKey := "PortNumber"
	if lp, ok := config[ListenerPortKey]; ok {
		listenerPort = int(lp.(float64))
	}
}

func ConfigureMessaging(config map[string]interface{}) {
	negotiator.UnmarshalMessage = protomsg.UnmarshalNegotiateMessage
	MessageFormatKey := "MessageEncodingFormat"

	if mf, ok := config[MessageFormatKey]; ok {
		ProtobufOption := "protobuf"
		if mf.(string) == ProtobufOption {
		} else {
			log.Fatalf(InvalidOptionError, mf.(string), MessageFormatKey)
		}
	}
}

func ConfigureDebriefer(config map[string]interface{}) {
	DebrieferLoadPreferrenceKey := "LoadPreferrenceHistoryLength"
	if dlp, ok := config[DebrieferLoadPreferrenceKey]; ok {
		debrieferLoadHistorySize = int(dlp.(float64))
	}
}

func ConfigureConnector(config map[string]interface{}) {
	if rqs, ok := config[requestQueueSizeKey]; ok {
		connectorQueueSize = int(rqs.(float64))
	}
}

func ConfigureRegistrator(config map[string]interface{}) {
	if rqs, ok := config[requestQueueSizeKey]; ok {
		registratorQueueSize = int(rqs.(float64))
	}
}

func ConfigureLocalizer(config map[string]interface{}) {
	if rqs, ok := config[requestQueueSizeKey]; ok {
		localizerQueueSize = int(rqs.(float64))
	}
}

func ConfigureAnalyzer(config map[string]interface{}) {
	DistancePollTimeKey := "SwarmFitCalculationFrequency"
	if dpt, ok := config[DistancePollTimeKey]; ok {
		distancePollTime := time.Duration(int64(dpt.(float64)) * int64(UnitOfTime))
		analyzer.DistancePollTime = distancePollTime
	}
}

func ConfigureCache(config map[string]interface{}) {
	ConnectionTTLKey := "ConnectionRecordTTL"
	DisconnectionTTLKey := "DisconnectionRecordTTL"
	GarbageCollectionPeriodKey := "RecordCleanupPeriod"

	if cTTL, ok := config[ConnectionTTLKey]; ok {
		connectionTTL := time.Duration(int64(cTTL.(float64)) * int64(UnitOfTime))
		cache.ConnectionTTL = connectionTTL
	}
	if dTTL, ok := config[DisconnectionTTLKey]; ok {
		disconnectionTTL := time.Duration(int64(dTTL.(float64)) * int64(UnitOfTime))
		cache.DisconnectionTTL = disconnectionTTL
	}
	if gcp, ok := config[GarbageCollectionPeriodKey]; ok {
		garbageCollectionPeriod := time.Duration(int64(gcp.(float64)) * int64(UnitOfTime))
		cache.GarbageCollectionPeriod = garbageCollectionPeriod
	}
}

func ConfigureComparator(config map[string]interface{}) {
	DefaultPreferredLoadKey := "DefaultPreferredLoad"

	if dpl, ok := config[DefaultPreferredLoadKey]; ok {
		defaultPreferredLoad := int(dpl.(float64))
		comparator.DefaultPreferredLoad = defaultPreferredLoad
	}
}

func ConfigureGateway(config map[string]interface{}) {
	DefaultQueueCapacityKey := "MaxActiveConnectionsOnStandby"

	if dq, ok := config[DefaultQueueCapacityKey]; ok {
		gatewayActiveQueueSize = int(dq.(float64))
	}
}

func ConfigureManager(config map[string]interface{}) {
	ChangeTriggerLimitKey := "ChangesUntilSizeUpdate"
	DebriefProcedureKey := "DebriefProcedure"

	if dp, ok := config[DebriefProcedureKey]; ok {
		PreferredLoadOption := "PreferredLoad"

		debriefProcedure := dp.(string)
		if debriefProcedure == PreferredLoadOption {
			manager.DebriefProcedure = debriefer.LoadPreferrenceDebrief
		} else {
			log.Fatalf(InvalidOptionError, debriefProcedure, DebriefProcedureKey)
		}
	} else {
		log.Fatalf(MissingSettingError, DebriefProcedureKey)
	}

	if ctl, ok := config[ChangeTriggerLimitKey]; ok {
		manager.ChangeTriggerLimit = int(ctl.(float64))
	}
}

func ConfigureNegotiator(config map[string]interface{}) {
	RoundtripLimitKey := "MaxRoundtripsDuringNegotiation"

	if rl, ok := config[RoundtripLimitKey]; ok {
		negotiator.RoundtripLimit = int(rl.(float64))
	}
}

func ConfigureRegister(config map[string]interface{}) {
	HostKey := "Hostname"
	UserKey := "Username"
	PasswordKey := "Password"
	PortKey := "PortNumber"

	if h, ok := config[HostKey]; ok {
		register.Host = h.(string)
	}
	if u, ok := config[UserKey]; ok {
		register.User = u.(string)
	}
	if p, ok := config[PasswordKey]; ok {
		register.Password = p.(string)
	}
	if pt, ok := config[PortKey]; ok {
		register.Port = int(pt.(float64))
	}
}

func ConfigureTracker(config map[string]interface{}) {
	FrequencyCalculationKey := "LoadParameterCalculationFrequency"
	FrequencyAveragingWidthKey := "LoadParameterHistorySize"

	if fc, ok := config[FrequencyCalculationKey]; ok {
		tracker.FrequencyCalculationPeriod = time.Duration(int64(fc.(float64)) * int64(UnitOfTime))
	}
	if fa, ok := config[FrequencyAveragingWidthKey]; ok {
		trackerLoadHistorySize = int(fa.(float64))
	}
}

func ConfigureTransmuter(config map[string]interface{}) {
	PollPeriodKey := "SwarmRestructuringFrequency"

	if pp, ok := config[PollPeriodKey]; ok {
		transmuter.PollPeriod = time.Duration(int64(pp.(float64)) * int64(UnitOfTime))
	}
}
