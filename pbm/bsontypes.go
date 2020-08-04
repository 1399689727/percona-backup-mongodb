package pbm

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OpTime struct {
	TS   primitive.Timestamp `bson:"ts" json:"ts"`
	Term int64               `bson:"t" json:"t"`
}

// IsMasterLastWrite represents the last write to the MongoDB server
type IsMasterLastWrite struct {
	OpTime            OpTime    `bson:"opTime"`
	LastWriteDate     time.Time `bson:"lastWriteDate"`
	MajorityOpTime    OpTime    `bson:"majorityOpTime"`
	MajorityWriteDate time.Time `bson:"majorityWriteDate"`
}

// IsMaster represents the document returned by db.runCommand( { isMaster: 1 } )
/**

mongos
{
        "ismaster" : true,
        "msg" : "isdbgrid",
        "maxBsonObjectSize" : 16777216,
        "maxMessageSizeBytes" : 48000000,
        "maxWriteBatchSize" : 100000,
        "localTime" : ISODate("2020-07-30T07:07:55.987Z"),
        "logicalSessionTimeoutMinutes" : 30,
        "connectionId" : 1015391,
        "maxWireVersion" : 8,
        "minWireVersion" : 0,
        "ok" : 1,
        "operationTime" : Timestamp(1596092874, 2),
        "$clusterTime" : {
                "clusterTime" : Timestamp(1596092874, 2),
                "signature" : {
                        "hash" : BinData(0,"6QngecKY9HQ8WyK+TOO+b3e3kcw="),
                        "keyId" : NumberLong("6836349335383310337")
                }
        }
}

replicaSet
{
        "hosts" : [
                "web-db01cn-p001.pek4.wecash.net:27017",
                "10.40.9.157:27017",
                "10.40.9.188:27017"
        ],
        "setName" : "haha01",
        "setVersion" : 10,
        "ismaster" : true,
        "secondary" : false,
        "primary" : "web-db01cn-p001.pek4.wecash.net:27017",
        "me" : "web-db01cn-p001.pek4.wecash.net:27017",
        "electionId" : ObjectId("7fffffff000000000000001b"),
        "lastWrite" : {
                "opTime" : {
                        "ts" : Timestamp(1596092708, 1),
                        "t" : NumberLong(27)
                },
                "lastWriteDate" : ISODate("2020-07-30T07:05:08Z"),
                "majorityOpTime" : {
                        "ts" : Timestamp(1596092708, 1),
                        "t" : NumberLong(27)
                },
                "majorityWriteDate" : ISODate("2020-07-30T07:05:08Z")
        },
        "maxBsonObjectSize" : 16777216,
        "maxMessageSizeBytes" : 48000000,
        "maxWriteBatchSize" : 100000,
        "localTime" : ISODate("2020-07-30T07:05:09.372Z"),
        "logicalSessionTimeoutMinutes" : 30,
        "minWireVersion" : 0,
        "maxWireVersion" : 7,
        "readOnly" : false,
        "ok" : 1,
        "operationTime" : Timestamp(1596092708, 1),
        "$clusterTime" : {
                "clusterTime" : Timestamp(1596092708, 1),
                "signature" : {
                        "hash" : BinData(0,"0n0sEhSuzpiVvtae62884V9o5zc="),
                        "keyId" : NumberLong("6820270674339168258")
                }
        }
}
 */
type IsMaster struct {
	Hosts                        []string           `bson:"hosts,omitempty"`
	Msg                          string             `bson:"msg"`
	MaxBsonObjectSise            int64              `bson:"maxBsonObjectSize"`
	MaxMessageSizeBytes          int64              `bson:"maxMessageSizeBytes"`
	MaxWriteBatchSize            int64              `bson:"maxWriteBatchSize"`
	LocalTime                    time.Time          `bson:"localTime"`
	LogicalSessionTimeoutMinutes int64              `bson:"logicalSessionTimeoutMinutes"`
	MaxWireVersion               int64              `bson:"maxWireVersion"`
	MinWireVersion               int64              `bson:"minWireVersion"`
	OK                           int                `bson:"ok"`
	SetName                      string             `bson:"setName,omitempty"`
	Primary                      string             `bson:"primary,omitempty"`
	SetVersion                   int32              `bson:"setVersion,omitempty"`
	IsMaster                     bool               `bson:"ismaster"`
	Secondary                    bool               `bson:"secondary,omitempty"`
	Hidden                       bool               `bson:"hidden,omitempty"`
	ConfigSvr                    int                `bson:"configsvr,omitempty"`
	Me                           string             `bson:"me"`
	LastWrite                    IsMasterLastWrite  `bson:"lastWrite"`
	ClusterTime                  *ClusterTime       `bson:"$clusterTime,omitempty"`
	ConfigServerState            *ConfigServerState `bson:"$configServerState,omitempty"`
	// GleStats                     *GleStats            `bson:"$gleStats,omitempty"`
	OperationTime *primitive.Timestamp `bson:"operationTime,omitempty"`
}

// IsSharded returns true is replset is part sharded cluster
//  "configsvr" : 2 就是 config server
func (im *IsMaster) IsSharded() bool {
	return im.SetName != "" && (im.ConfigServerState != nil || im.ConfigSvr == 2)
}

// IsLeader returns true if node can act as backup leader (it's configsrv or non shareded rs)
func (im *IsMaster) IsLeader() bool {
	return !im.IsSharded() || im.ReplsetRole() == ReplRoleConfigSrv
}

// ReplsetRole returns replset role in sharded clister
func (im *IsMaster) ReplsetRole() ReplRole {
	switch {
	case im.ConfigSvr == 2:
		return ReplRoleConfigSrv
	case im.ConfigServerState != nil:
		return ReplRoleShard
	default:
		return ReplRoleUnknown
	}
}

// IsStandalone returns true if node is not a part of replica set
func (im *IsMaster) IsStandalone() bool {
	return im.SetName == ""
}

type ClusterTime struct {
	ClusterTime primitive.Timestamp `bson:"clusterTime"`
	Signature   struct {
		Hash  primitive.Binary `bson:"hash"`
		KeyID int64            `bson:"keyId"`
	} `bson:"signature"`
}

type ConfigServerState struct {
	OpTime *OpTime `bson:"opTime"`
}

type Operation string

const (
	OperationInsert  Operation = "i"
	OperationNoop    Operation = "n"
	OperationUpdate  Operation = "u"
	OperationDelete  Operation = "d"
	OperationCommand Operation = "c"
)

type NodeHealth int

const (
	NodeHealthDown NodeHealth = iota
	NodeHealthUp
)

type NodeState int

const (
	NodeStateStartup NodeState = iota
	NodeStatePrimary
	NodeStateSecondary
	NodeStateRecovering
	NodeStateStartup2
	NodeStateUnknown
	NodeStateArbiter
	NodeStateDown
	NodeStateRollback
	NodeStateRemoved
)

type StatusOpTimes struct {
	LastCommittedOpTime       *OpTime `bson:"lastCommittedOpTime" json:"lastCommittedOpTime"`
	ReadConcernMajorityOpTime *OpTime `bson:"readConcernMajorityOpTime" json:"readConcernMajorityOpTime"`
	AppliedOpTime             *OpTime `bson:"appliedOpTime" json:"appliedOpTime"`
	DurableOptime             *OpTime `bson:"durableOpTime" json:"durableOpTime"`
}

type NodeStatus struct {
	ID                int                 `bson:"_id" json:"_id"`
	Name              string              `bson:"name" json:"name"`
	Health            NodeHealth          `bson:"health" json:"health"`
	State             NodeState           `bson:"state" json:"state"`
	StateStr          string              `bson:"stateStr" json:"stateStr"`
	Uptime            int64               `bson:"uptime" json:"uptime"`
	Optime            *OpTime             `bson:"optime" json:"optime"`
	OptimeDate        time.Time           `bson:"optimeDate" json:"optimeDate"`
	ConfigVersion     int                 `bson:"configVersion" json:"configVersion"`
	ElectionTime      primitive.Timestamp `bson:"electionTime,omitempty" json:"electionTime,omitempty"`
	ElectionDate      time.Time           `bson:"electionDate,omitempty" json:"electionDate,omitempty"`
	InfoMessage       string              `bson:"infoMessage,omitempty" json:"infoMessage,omitempty"`
	OptimeDurable     *OpTime             `bson:"optimeDurable,omitempty" json:"optimeDurable,omitempty"`
	OptimeDurableDate time.Time           `bson:"optimeDurableDate,omitempty" json:"optimeDurableDate,omitempty"`
	LastHeartbeat     time.Time           `bson:"lastHeartbeat,omitempty" json:"lastHeartbeat,omitempty"`
	LastHeartbeatRecv time.Time           `bson:"lastHeartbeatRecv,omitempty" json:"lastHeartbeatRecv,omitempty"`
	PingMs            int64               `bson:"pingMs,omitempty" json:"pingMs,omitempty"`
	Self              bool                `bson:"self,omitempty" json:"self,omitempty"`
	SyncingTo         string              `bson:"syncingTo,omitempty" json:"syncingTo,omitempty"`
}

type ReplsetStatus struct {
	Set                     string             `bson:"set" json:"set"`
	Date                    time.Time          `bson:"date" json:"date"`
	MyState                 NodeState          `bson:"myState" json:"myState"`
	Members                 []NodeStatus       `bson:"members" json:"members"`
	Term                    int64              `bson:"term,omitempty" json:"term,omitempty"`
	HeartbeatIntervalMillis int64              `bson:"heartbeatIntervalMillis,omitempty" json:"heartbeatIntervalMillis,omitempty"`
	Optimes                 *StatusOpTimes     `bson:"optimes,omitempty" json:"optimes,omitempty"`
	Errmsg                  string             `bson:"errmsg,omitempty" json:"errmsg,omitempty"`
	Ok                      int                `bson:"ok" json:"ok"`
	ClusterTime             *ClusterTime       `bson:"$clusterTime,omitempty" json:"$clusterTime,omitempty"`
	ConfigServerState       *ConfigServerState `bson:"$configServerState,omitempty" json:"$configServerState,omitempty"`
	// GleStats                *GleStats            `bson:"$gleStats,omitempty" json:"$gleStats,omitempty"`
	OperationTime *primitive.Timestamp `bson:"operationTime,omitempty" json:"operationTime,omitempty"`
}

// Shard represent config.shard https://docs.mongodb.com/manual/reference/config-database/#config.shards
type Shard struct {
	ID   string `bson:"_id"`
	Host string `bson:"host"`
}

type ConnectionStatus struct {
	AuthInfo AuthInfo `bson:"authInfo" json:"authInfo"`
}

type AuthInfo struct {
	Users     []AuthUser      `bson:"authenticatedUsers" json:"authenticatedUsers"`
	UserRoles []AuthUserRoles `bson:"authenticatedUserRoles" json:"authenticatedUserRoles"`
}

type AuthUser struct {
	User string `bson:"user" json:"user"`
	DB   string `bson:"db" json:"db"`
}
type AuthUserRoles struct {
	Role string `bson:"role" json:"role"`
	DB   string `bson:"db" json:"db"`
}
