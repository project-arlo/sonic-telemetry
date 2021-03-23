//  Package client provides a generic access layer for data available in system
package client

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib"
	"github.com/Azure/sonic-telemetry/common_utils"
	spb "github.com/Azure/sonic-telemetry/proto"
	transutil "github.com/Azure/sonic-telemetry/transl_utils"
	log "github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
	gnmi_extpb "github.com/openconfig/gnmi/proto/gnmi_ext"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TranslClient struct {
	prefix *gnmipb.Path
	/* GNMI Path to REST URL Mapping */
	path2URI map[*gnmipb.Path]string
	encoding gnmipb.Encoding
	channel  chan struct{}
	q        *LimitedQueue

	w          *sync.WaitGroup // wait for all sub go routines to finish
	ctx        context.Context //Contains Auth info and request info
	extensions []*gnmi_extpb.Extension
	version    *translib.Version // Client version; populated by parseVersion()
}

func NewTranslClient(prefix *gnmipb.Path, getpaths []*gnmipb.Path, ctx context.Context, extensions []*gnmi_extpb.Extension) (Client, error) {
	var client TranslClient
	var err error
	client.ctx = ctx
	client.prefix = prefix
	client.extensions = extensions
	if getpaths != nil {
		client.path2URI = make(map[*gnmipb.Path]string)
		/* Populate GNMI path to REST URL map. */
		err = transutil.PopulateClientPaths(prefix, getpaths, &client.path2URI)
	}

	if err != nil {
		return nil, err
	} else {
		return &client, nil
	}
}

func (c *TranslClient) Get(w *sync.WaitGroup) ([]*spb.Value, error) {
	rc, ctx := common_utils.GetContext(c.ctx)
	c.ctx = ctx
	var values []*spb.Value
	ts := time.Now()

	version := getBundleVersion(c.extensions)
	if version != nil {
		rc.BundleVersion = version
	}
	/* Iterate through all GNMI paths. */

	for gnmiPath, URIPath := range c.path2URI {
		/* Fill values for each GNMI path. */
		val, valTree, err := transutil.TranslProcessGet(URIPath, nil, c.ctx, c.encoding)

		if err != nil {
			return nil, err
		}

		v, err := buildValue(c.prefix, gnmiPath, c.encoding, val, valTree)
		if err != nil {
			return nil, err
		}
		/* Value of each path is added to spb value structure. */
		values = append(values, v)
	}

	/* The values structure at the end is returned and then updates in notitications as
	specified in the proto file in the server.go */

	log.V(6).Infof("TranslClient : Getting #%v", values)
	log.V(4).Infof("TranslClient :Get done, total time taken: %v ms", int64(time.Since(ts)/time.Millisecond))

	return values, nil
}

func (c *TranslClient) Set(delete []*gnmipb.Path, replace []*gnmipb.Update, update []*gnmipb.Update) error {
	rc, ctx := common_utils.GetContext(c.ctx)
	c.ctx = ctx
	var uri string
	version := getBundleVersion(c.extensions)
	if version != nil {
		rc.BundleVersion = version
	}

	if (len(delete) + len(replace) + len(update)) > 1 {
		return transutil.TranslProcessBulk(delete, replace, update, c.prefix, c.ctx)
	} else {
		if len(delete) == 1 {
			/* Convert the GNMI Path to URI. */
			transutil.ConvertToURI(c.prefix, delete[0], &uri)
			return transutil.TranslProcessDelete(uri, c.ctx)
		}
		if len(replace) == 1 {
			/* Convert the GNMI Path to URI. */
			transutil.ConvertToURI(c.prefix, replace[0].GetPath(), &uri)
			return transutil.TranslProcessReplace(uri, replace[0].GetVal(), c.ctx)
		}
		if len(update) == 1 {
			/* Convert the GNMI Path to URI. */
			transutil.ConvertToURI(c.prefix, update[0].GetPath(), &uri)
			return transutil.TranslProcessUpdate(uri, update[0].GetVal(), c.ctx)
		}
	}
	return nil
}

func enqueFatalMsgTranslib(c *TranslClient, msg string) {
	log.Error(msg)
	c.q.ForceEnqueueItem(Value{
		&spb.Value{
			Timestamp: time.Now().UnixNano(),
			Fatal:     msg,
		},
	})
}

func enqueueSyncMessage(c *TranslClient) {
	m := &spb.Value{
		Timestamp:    time.Now().UnixNano(),
		SyncResponse: true,
	}
	c.q.EnqueueItem(Value{m})
}

// recoverSubscribe recovers from possible panics during subscribe handling.
// It pushes a fatal message to the RPC handler's queue, which will force it to
// close the RPC with an error status. Should always be used as a deferred function.
func recoverSubscribe(c *TranslClient) {
	if r := recover(); r != nil {
		buff := make([]byte, 1<<12)
		buff = buff[:runtime.Stack(buff, false)]
		log.Error(string(buff))

		err := status.Errorf(codes.Internal, "%v", r)
		enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
	}
}

type ticker_info struct {
	t         *time.Ticker
	sub       *gnmipb.Subscription
	heartbeat bool
}

func tickerCleanup(ticker_map map[int][]*ticker_info, c *TranslClient) {
	for _, v := range ticker_map {
		for _, ti := range v {
			fmt.Println("Ticker Cleanup: ", c.path2URI[ti.sub.Path])
			ti.t.Stop()
		}
	}
}

func (c *TranslClient) StreamRun(q *LimitedQueue, stop chan struct{}, w *sync.WaitGroup, subscribe *gnmipb.SubscriptionList) {
	c.w = w

	defer c.w.Done()
	defer recoverSubscribe(c)

	c.q = q
	c.channel = stop

	if err := c.parseVersion(); err != nil {
		enqueFatalMsgTranslib(c, err.Error())
		return
	}

	ticker_map := make(map[int][]*ticker_info)

	defer tickerCleanup(ticker_map, c)
	var cases []reflect.SelectCase
	cases_map := make(map[int]int)
	var subscribe_mode gnmipb.SubscriptionMode
	stringPaths := make([]string, len(subscribe.Subscription))
	pathsMap := make(map[string]*gnmipb.Path)

	for i, sub := range subscribe.Subscription {
		p := c.path2URI[sub.Path]
		stringPaths[i] = p
		pathsMap[p] = sub.Path
	}

	ss := translib.NewSubscribeSession()
	defer translib.CloseSubscribeSession(ss)

	req := translib.IsSubscribeRequest{
		Paths:   stringPaths,
		Session: ss,
	}
	if c.version != nil {
		req.ClientVersion = *c.version
	}
	subSupport, _ := translib.IsSubscribeSupported(req)

	var onChangeSubsString []string

	valueCache := make(map[string]string)
	filterFunc := func(path string, v *spb.Value) bool {
		return filterValues(valueCache, v, path, c.encoding)
	}

	for i, sub := range subscribe.Subscription {
		log.V(6).Infof("%s %s", sub.Mode, sub.SampleInterval)

		switch sub.Mode {

		case gnmipb.SubscriptionMode_TARGET_DEFINED:

			if subSupport[i].Err == nil && subSupport[i].IsOnChangeSupported {
				if subSupport[i].PreferredType == translib.Sample {
					subscribe_mode = gnmipb.SubscriptionMode_SAMPLE
				} else if subSupport[i].PreferredType == translib.OnChange {
					subscribe_mode = gnmipb.SubscriptionMode_ON_CHANGE
				}
			} else {
				subscribe_mode = gnmipb.SubscriptionMode_SAMPLE
			}

		case gnmipb.SubscriptionMode_ON_CHANGE:
			if subSupport[i].Err == nil && subSupport[i].IsOnChangeSupported {
				if subSupport[i].MinInterval > 0 {
					subscribe_mode = gnmipb.SubscriptionMode_ON_CHANGE
				} else {
					enqueFatalMsgTranslib(c, fmt.Sprintf("Invalid subscribe path %v", stringPaths[i]))
					return
				}
			} else {
				enqueFatalMsgTranslib(c, fmt.Sprintf("ON_CHANGE Streaming mode invalid for %v", stringPaths[i]))
				return
			}
		case gnmipb.SubscriptionMode_SAMPLE:
			if subSupport[i].MinInterval > 0 {
				subscribe_mode = gnmipb.SubscriptionMode_SAMPLE
			} else {
				enqueFatalMsgTranslib(c, fmt.Sprintf("Invalid subscribe path %v", stringPaths[i]))
				return
			}
		default:
			log.V(1).Infof("Bad Subscription Mode for client %s ", c)
			enqueFatalMsgTranslib(c, fmt.Sprintf("Invalid Subscription Mode %d", sub.Mode))
			return
		}
		log.V(6).Infof("subscribe_mode:", subscribe_mode)
		if subscribe_mode == gnmipb.SubscriptionMode_SAMPLE {
			interval := int(sub.SampleInterval)
			if interval == 0 {
				interval = subSupport[i].MinInterval * int(time.Second)
			} else {
				if interval < (subSupport[i].MinInterval * int(time.Second)) {
					enqueFatalMsgTranslib(c, fmt.Sprintf("Invalid Sample Interval %ds, minimum interval is %ds", interval/int(time.Second), subSupport[i].MinInterval))
					return
				}
			}

			if !subscribe.UpdatesOnly {
				//Send initial data now so we can send sync response, unless updates_only is set.
				ts := translSubscriber{
					client:  c,
					session: ss,
					pathMap: pathsMap,
				}
				if sub.SuppressRedundant {
					ts.filterFunc = filterFunc
				}
				ts.doSample(sub.Path)
			}

			addTimer(c, ticker_map, &cases, cases_map, interval, sub, false)
			//Heartbeat intervals are valid for SAMPLE in the case suppress_redundant is specified
			if sub.SuppressRedundant && sub.HeartbeatInterval > 0 {
				if int(sub.HeartbeatInterval) < subSupport[i].MinInterval*int(time.Second) {
					enqueFatalMsgTranslib(c, fmt.Sprintf("Invalid Heartbeat Interval %ds, minimum interval is %ds", int(sub.HeartbeatInterval)/int(time.Second), subSupport[i].MinInterval))
					return
				}
				addTimer(c, ticker_map, &cases, cases_map, int(sub.HeartbeatInterval), sub, true)
			}
		} else if subscribe_mode == gnmipb.SubscriptionMode_ON_CHANGE {
			onChangeSubsString = append(onChangeSubsString, c.path2URI[sub.Path])
			if sub.HeartbeatInterval > 0 {
				if int(sub.HeartbeatInterval) < subSupport[i].MinInterval*int(time.Second) {
					enqueFatalMsgTranslib(c, fmt.Sprintf("Invalid Heartbeat Interval %ds, minimum interval is %ds", int(sub.HeartbeatInterval)/int(time.Second), subSupport[i].MinInterval))
					return
				}
				addTimer(c, ticker_map, &cases, cases_map, int(sub.HeartbeatInterval), sub, true)
			}

		}
	}

	if len(onChangeSubsString) > 0 {
		ts := translSubscriber{
			client:      c,
			session:     ss,
			pathMap:     pathsMap,
			updatesOnly: subscribe.UpdatesOnly,
		}
		ts.doOnChange(onChangeSubsString)
	} else {
		// If at least one ON_CHANGE subscription was present, then
		// ts.doOnChange() would have sent the sync message.
		// Explicitly send one here if all are SAMPLE subscriptions.
		enqueueSyncMessage(c)
	}

	cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c.channel)})

	for {
		chosen, _, ok := reflect.Select(cases)
		if !ok {
			return
		}

		for _, tick := range ticker_map[cases_map[chosen]] {
			log.V(6).Infof("tick, heartbeat: %t, path: %s\n", tick.heartbeat, c.path2URI[tick.sub.Path])
			ts := translSubscriber{
				client:      c,
				session:     ss,
				pathMap:     pathsMap,
				isHeartbeat: tick.heartbeat,
			}
			if tick.sub.SuppressRedundant {
				ts.filterFunc = filterFunc
			}
			ts.doSample(tick.sub.Path)
		}
	}
}

func addTimer(c *TranslClient, ticker_map map[int][]*ticker_info, cases *[]reflect.SelectCase, cases_map map[int]int, interval int, sub *gnmipb.Subscription, heartbeat bool) {
	//Reuse ticker for same sample intervals, otherwise create a new one.
	if ticker_map[interval] == nil {
		ticker_map[interval] = make([]*ticker_info, 1, 1)
		ticker_map[interval][0] = &ticker_info{
			t:         time.NewTicker(time.Duration(interval) * time.Nanosecond),
			sub:       sub,
			heartbeat: heartbeat,
		}
		cases_map[len(*cases)] = interval
		*cases = append(*cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ticker_map[interval][0].t.C)})
	} else {
		ticker_map[interval] = append(ticker_map[interval], &ticker_info{
			t:         ticker_map[interval][0].t,
			sub:       sub,
			heartbeat: heartbeat,
		})
	}

}

func (c *TranslClient) PollRun(q *LimitedQueue, poll chan struct{}, w *sync.WaitGroup, subscribe *gnmipb.SubscriptionList) {
	c.w = w

	defer c.w.Done()
	defer recoverSubscribe(c)

	c.q = q
	c.channel = poll
	synced := false

	if err := c.parseVersion(); err != nil {
		enqueFatalMsgTranslib(c, err.Error())
		return
	}

	for {
		_, more := <-c.channel
		if !more {
			log.V(1).Infof("%v poll channel closed, exiting pollDb routine", c)
			enqueFatalMsgTranslib(c, "")
			return
		}
		t1 := time.Now()
		for gnmiPath := range c.path2URI {
			if synced || !subscribe.UpdatesOnly {
				ts := translSubscriber{client: c}
				ts.doSample(gnmiPath)
			}
		}

		enqueueSyncMessage(c)
		synced = true
		log.V(4).Infof("Sync done, poll time taken: %v ms", int64(time.Since(t1)/time.Millisecond))
	}
}

func (c *TranslClient) OnceRun(q *LimitedQueue, once chan struct{}, w *sync.WaitGroup, subscribe *gnmipb.SubscriptionList) {
	c.w = w

	defer c.w.Done()
	defer recoverSubscribe(c)

	c.q = q
	c.channel = once

	if err := c.parseVersion(); err != nil {
		enqueFatalMsgTranslib(c, err.Error())
		return
	}

	_, more := <-c.channel
	if !more {
		log.V(1).Infof("%v once channel closed, exiting onceDb routine", c)
		enqueFatalMsgTranslib(c, "")
		return
	}
	t1 := time.Now()
	for gnmiPath := range c.path2URI {
		ts := translSubscriber{client: c}
		ts.doSample(gnmiPath)
	}

	enqueueSyncMessage(c)
	log.V(4).Infof("Sync done, once time taken: %v ms", int64(time.Since(t1)/time.Millisecond))

}

// setPrefixTarget fills prefix taregt for given Notification objects.
func setPrefixTarget(notifs []*gnmipb.Notification, target string) {
	for _, n := range notifs {
		if n.Prefix == nil {
			n.Prefix = &gnmipb.Path{Target: target}
		} else {
			n.Prefix.Target = target
		}
	}
}

//Creates a spb.Value out of data from the translib according to the requested encoding.
func buildValue(prefix *gnmipb.Path, path *gnmipb.Path, enc gnmipb.Encoding,
	typedVal *gnmipb.TypedValue, valueTree *ygot.ValidatedGoStruct) (*spb.Value, error) {

	switch enc {
	case gnmipb.Encoding_JSON:
		fallthrough
	case gnmipb.Encoding_JSON_IETF:
		if typedVal == nil {
			return nil, nil
		}
		return &spb.Value{
			Prefix:       prefix,
			Path:         path,
			Timestamp:    time.Now().UnixNano(),
			SyncResponse: false,
			Val:          typedVal,
		}, nil

	case gnmipb.Encoding_PROTO:
		if valueTree == nil {
			return nil, nil
		}

		fullPath := transutil.GnmiTranslFullPath(prefix, path)
		removeLastPathElem(fullPath)
		notifications, err := ygot.TogNMINotifications(
			*valueTree,
			time.Now().UnixNano(),
			ygot.GNMINotificationsConfig{
				UsePathElem:    true,
				PathElemPrefix: fullPath.GetElem(),
			})
		if err != nil {
			return nil, fmt.Errorf("Cannot convert OC Struct to Notifications: %s", err)
		}
		if len(notifications) != 1 {
			return nil, fmt.Errorf("YGOT returned wrong number of notifications")
		}
		if len(prefix.Target) != 0 {
			// Copy target from reqest.. ygot.TogNMINotifications does not fill it.
			setPrefixTarget(notifications, prefix.Target)
		}
		return &spb.Value{
			Notification: notifications[0],
		}, nil
	default:
		return nil, fmt.Errorf("Unsupported Encoding %s", enc)
	}

}

func (c *TranslClient) Capabilities() []gnmipb.ModelData {
	rc, ctx := common_utils.GetContext(c.ctx)
	c.ctx = ctx
	version := getBundleVersion(c.extensions)
	if version != nil {
		rc.BundleVersion = version
	}
	/* Fetch the supported models. */
	supportedModels := transutil.GetModels()
	return supportedModels
}

func (c *TranslClient) Close() error {
	return nil
}

func getBundleVersion(extensions []*gnmi_extpb.Extension) *string {
	for _, e := range extensions {
		switch v := e.Ext.(type) {
		case *gnmi_extpb.Extension_RegisteredExt:
			if v.RegisteredExt.Id == spb.BUNDLE_VERSION_EXT {
				var bv spb.BundleVersion
				proto.Unmarshal(v.RegisteredExt.Msg, &bv)
				return &bv.Version
			}

		}
	}
	return nil
}

func (c *TranslClient) parseVersion() error {
	bv := getBundleVersion(c.extensions)
	if bv == nil {
		return nil
	}
	v, err := translib.NewVersion(*bv)
	if err != nil {
		c.version = &v
		return nil
	}
	log.V(4).Infof("Failed to parse version \"%s\"; err=%v", *bv, err)
	return fmt.Errorf("Invalid bundle version: %v", *bv)
}

// Set the desired encoding for Get and Subcribe responses
func (c *TranslClient) SetEncoding(enc gnmipb.Encoding) {
	c.encoding = enc
}
