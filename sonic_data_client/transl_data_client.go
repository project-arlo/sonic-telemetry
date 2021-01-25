//  Package client provides a generic access layer for data available in system
package client

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/Azure/sonic-telemetry/common_utils"
	spb "github.com/Azure/sonic-telemetry/proto"
	transutil "github.com/Azure/sonic-telemetry/transl_utils"
	"github.com/Workiva/go-datastructures/queue"
	log "github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
	gnmi_extpb "github.com/openconfig/gnmi/proto/gnmi_ext"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DELETE  int = 0
	REPLACE int = 1
	UPDATE  int = 2
)

type TranslClient struct {
	prefix *gnmipb.Path
	/* GNMI Path to REST URL Mapping */
	path2URI map[*gnmipb.Path]string
	encoding gnmipb.Encoding
	channel  chan struct{}
	q        *LimitedQueue

	synced     sync.WaitGroup  // Control when to send gNMI sync_response
	w          *sync.WaitGroup // wait for all sub go routines to finish
	mu         sync.RWMutex    // Mutex for data protection among routines for transl_client
	ctx        context.Context //Contains Auth info and request info
	extensions []*gnmi_extpb.Extension
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
	rc, ctx := common_utils.GetContext(c.ctx)
	c.ctx = ctx
	c.w = w

	defer func() {
		if r := recover(); r != nil {
			err := status.Errorf(codes.Internal, "%v", r)
			enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
		}
	}()

	defer c.w.Done()

	c.q = q
	c.channel = stop
	version := getBundleVersion(c.extensions)
	if version != nil {
		rc.BundleVersion = version
	}

	ticker_map := make(map[int][]*ticker_info)

	defer tickerCleanup(ticker_map, c)
	var cases []reflect.SelectCase
	cases_map := make(map[int]int)
	var subscribe_mode gnmipb.SubscriptionMode
	stringPaths := make([]string, len(subscribe.Subscription))
	for i, sub := range subscribe.Subscription {
		stringPaths[i] = c.path2URI[sub.Path]
	}
	req := translib.IsSubscribeRequest{Paths: stringPaths}
	subSupport, _ := translib.IsSubscribeSupported(req)
	var onChangeSubsString []string
	var onChangeSubsgNMI []*gnmipb.Path
	onChangeMap := make(map[string]*gnmipb.Path)
	valueCache := make(map[string]string)

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
				val, valTree, err := transutil.TranslProcessGet(c.path2URI[sub.Path], nil, c.ctx, c.encoding)
				if err != nil {
					switch err.(type) {
					case tlerr.NotFoundError:
						log.V(1).Infof("Subscribe Path Resource Not Found, ignoring: %v", c.path2URI[sub.Path])
					default:
						enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
						return
					}

				}
				v, err := buildValue(c.prefix, sub.Path, c.encoding, val, valTree)
				if err != nil {
					enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe value failed to build =%v", err.Error()))
					return
				}
				c.q.EnqueueItem(Value{v})
				log.V(6).Infof("Added spbv #%v", v)
				filterValues(valueCache, v, c.path2URI[sub.Path], c.encoding)

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
			onChangeSubsgNMI = append(onChangeSubsgNMI, sub.Path)
			onChangeMap[c.path2URI[sub.Path]] = sub.Path
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
		c.w.Add(1)
		c.synced.Add(1)
		go TranslSubscribe(onChangeSubsgNMI, onChangeSubsString, onChangeMap, c, subscribe.UpdatesOnly)

	}
	// Wait until all data values corresponding to the path(s) specified
	// in the SubscriptionList has been transmitted at least once
	c.synced.Wait()
	spbs := &spb.Value{
		Timestamp:    time.Now().UnixNano(),
		SyncResponse: true,
	}
	c.q.EnqueueItem(Value{spbs})
	cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c.channel)})

	for {
		chosen, _, ok := reflect.Select(cases)

		if !ok {
			return
		}

		for _, tick := range ticker_map[cases_map[chosen]] {
			log.V(6).Infof("tick, heartbeat: %t, path: %s\n", tick.heartbeat, c.path2URI[tick.sub.Path])
			val, valTree, err := transutil.TranslProcessGet(c.path2URI[tick.sub.Path], nil, c.ctx, c.encoding)
			if err != nil {
				switch err.(type) {
				case tlerr.NotFoundError:
					log.V(1).Infof("Subscribe Path Resource Not Found, ignoring: %v", c.path2URI[tick.sub.Path])
				default:
					enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
					return
				}

			}
			v, err := buildValue(c.prefix, tick.sub.Path, c.encoding, val, valTree)
			if err != nil {
				enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe value failed to build =%v", err.Error()))
				return
			}

			if tick.sub.SuppressRedundant && !filterValues(valueCache, v, c.path2URI[tick.sub.Path], c.encoding) && !tick.heartbeat {
				log.V(6).Infof("Redundant Message Suppressed #%v", string(val.GetJsonIetfVal()))
			} else {
				c.q.EnqueueItem(Value{v})
				log.V(6).Infof("Added spbv #%v", v)
			}
		}
	}
}

//Compares val to cache, removes any duplicates from val,
//updates cache. Returns true if there are new values.
func filterValues(cache map[string]string, val *spb.Value, URI string, enc gnmipb.Encoding) bool {
	switch enc {
	case gnmipb.Encoding_JSON:
		fallthrough
	case gnmipb.Encoding_JSON_IETF:
		if v := string(val.GetVal().GetJsonIetfVal()); v != cache[URI] {
			cache[URI] = v
			return true
		}
		return false
	case gnmipb.Encoding_PROTO:
		prefix, err := ygot.PathToString(val.Notification.GetPrefix())
		if err != nil {
			log.V(4).Infof("Failed to stringify prefix: #%v", err)
		}

		//Remove delete messages from cache.
		for _, d := range val.Notification.GetDelete() {
			path, err := ygot.PathToString(d)
			if err != nil {
				log.V(4).Infof("Failed to stringify delete path: #%v", err)
				continue
			}
			delete(cache, prefix+path)
		}

		//Only forward values that have changed from the cache.
		var filteredUpdates []*gnmipb.Update
		for _, u := range val.Notification.GetUpdate() {
			path, err := ygot.PathToString(u.GetPath())
			if err != nil {
				log.V(4).Infof("Failed to stringify update path: #%v", err)
				continue
			}
			leaf := proto.MarshalTextString(u.GetVal())
			if cache[prefix+path] != leaf {
				cache[prefix+path] = leaf
				filteredUpdates = append(filteredUpdates, u)
			}

		}
		if len(filteredUpdates) == 0 && len(val.Notification.GetDelete()) == 0 {
			return false
		}
		val.GetNotification().Update = filteredUpdates
		return true
	default:
		return false
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

func TranslSubscribe(gnmiPaths []*gnmipb.Path, stringPaths []string, pathMap map[string]*gnmipb.Path, c *TranslClient, updates_only bool) {
	var sync_done bool
	defer func() {
		if r := recover(); r != nil {
			err := status.Errorf(codes.Internal, "%v", r)
			enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
		}
		if !sync_done {
			// We reach here (without sync_done) only on error. RPC stream
			// would have been notified through enqueFatalMsgTranslib call.
			// But StreamRun would still be waiting on c.synced!! Wake him up.
			c.synced.Done()
		}
	}()
	defer c.w.Done()
	rc, ctx := common_utils.GetContext(c.ctx)
	c.ctx = ctx
	q := queue.NewPriorityQueue(1, false)

	log.V(6).Infof("Received Encoding:", c.encoding)
	req := translib.SubscribeRequest{Paths: stringPaths, Q: q, Stop: c.channel}
	if rc.BundleVersion != nil {
		nver, err := translib.NewVersion(*rc.BundleVersion)
		if err != nil {
			log.V(2).Infof("Subscribe operation failed with error =%v", err.Error())
			enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
			return
		}
		req.ClientVersion = nver
	}

	_, err := translib.Subscribe(req)
	if err != nil {
		enqueFatalMsgTranslib(c, "Subscribe operation failed with error: "+err.Error())
		return
	}

	for {
		items, err := q.Get(1)
		if err != nil {
			log.V(1).Infof("%v", err)
			enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
			return
		}
		switch v := items[0].(type) {
		case *translib.SubscribeResponse:

			if v.IsTerminated {
				//DB Connection or other backend error
				enqueFatalMsgTranslib(c, "DB Connection Error")
				close(c.channel)
				return
			}

			if v.SyncComplete && !sync_done {
				log.V(6).Infof("SENDING SYNC")
				c.synced.Done()
				sync_done = true
				break
			}

			//TODO change SubscribeResponse to return *gnmi.Path itself
			p := pathMap[v.Path]
			if p == nil {
				p, _ = ygot.StringToStructuredPath(v.Path)
				p.Target = c.prefix.Target
				p.Origin = c.prefix.Origin
			} else {
				p = copyPath(p)
			}

			n := gnmipb.Notification{
				Prefix:    p,
				Timestamp: v.Timestamp,
			}

			var extraPrefix *gnmipb.PathElem

			if c.encoding == gnmipb.Encoding_PROTO {
				if strSliceContains(v.Delete, "") {
					extraPrefix = removeLastPathElem(p)
				}
				if v.Update != nil {
					ju, err := ygot.TogNMINotifications(v.Update, v.Timestamp,
						ygot.GNMINotificationsConfig{UsePathElem: true, PathElemPrefix: p.GetElem()})
					if err != nil {
						log.Error(err)
						break
					}
					n.Update = ju[0].Update
					if extraPrefix != nil {
						for _, u := range n.Update {
							insertFirstPathElem(u.Path, extraPrefix)
						}
					}
				}
			} else {
				extraPrefix = removeLastPathElem(p)
				if v.Update != nil {
					ju, err := translToIetfJsonValue(extraPrefix, v.Update)
					if err != nil {
						log.Error(err)
						break
					}
					n.Update = append(n.Update, ju)
				}
			}

			for _, del := range v.Delete {
				p, err := ygot.StringToStructuredPath(del)
				if err != nil {
					log.Errorf("Invalid path \"%s\"; err=%v", del, err)
					break
				}
				insertFirstPathElem(p, extraPrefix)
				n.Delete = append(n.Delete, p)
			}

			spbv := &spb.Value{
				Notification: &n,
			}

			//Don't send initial update with full object if user wants updates only.
			if updates_only && !sync_done {
				log.V(1).Infof("Msg suppressed due to updates_only")
			} else {
				c.q.EnqueueItem(Value{spbv})
			}

			log.V(6).Infof("Added spbv #%v", spbv)

		default:
			log.V(1).Infof("Unknown data type %v for %s in queue", items[0], c)
		}
	}
}

func translToIetfJsonValue(targetElem *gnmipb.PathElem, yval ygot.ValidatedGoStruct) (*gnmipb.Update, error) {
	jv, err := ygot.EmitJSON(yval, &ygot.EmitJSONConfig{
		Format:         ygot.RFC7951,
		SkipValidation: true,
	})
	if err != nil {
		return nil, fmt.Errorf("EmitJSON failed; err=%v", err)
	}

	p := &gnmipb.Path{}

	if targetElem != nil {
		p.Elem = append(p.Elem, targetElem)
		// Translib always returns ygot struct for the SubscribeResponse.Path.
		// IETF JSON payload should being with the target container/list name of
		// this path; i.e, payload for /a/b should be {"b":{***}}. But ygot.EmitJSON
		// dumps only the subtree data. Adding the top node here.
		buff := new(strings.Builder)
		fmt.Fprintf(buff, "{\"%s\":", targetElem.Name)
		if len(targetElem.Key) != 0 {
			buff.WriteByte('[')
			buff.WriteString(jv)
			buff.WriteByte(']')
		} else {
			buff.WriteString(jv)
		}
		buff.WriteByte('}')
		jv = buff.String()
	}

	return &gnmipb.Update{
		Path: p,
		Val: &gnmipb.TypedValue{
			Value: &gnmipb.TypedValue_JsonIetfVal{
				JsonIetfVal: []byte(jv),
			}},
	}, nil
}

// copyPath creates a copy of p whose elements can be added/removed
// without affecting p. Returne value is not a full clone.
func copyPath(p *gnmipb.Path) *gnmipb.Path {
	elems := make([]*gnmipb.PathElem, len(p.Elem))
	copy(elems, p.Elem)
	return &gnmipb.Path{
		Elem: elems,
		Target: p.Target,
		Origin: p.Origin,
	}
}

// insertFirstPathElem inserts newElem at the beginning of path p.
func insertFirstPathElem(p *gnmipb.Path, newElem *gnmipb.PathElem) {
	if newElem != nil {
		ne := make([]*gnmipb.PathElem, 0, len(p.Elem)+1)
		ne = append(ne, newElem)
		p.Elem = append(ne, p.Elem...)
	}
}

// removeLastPathElem removes the last PathElem from the path p.
// Returns the removed element.
func removeLastPathElem(p *gnmipb.Path) *gnmipb.PathElem {
	k := len(p.Elem) - 1
	if k < 0 {
		return nil
	}
	last := p.Elem[k]
	p.Elem = p.Elem[:k]
	return last
}

func strSliceContains(ss []string, v string) bool {
	for _, s := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func (c *TranslClient) PollRun(q *LimitedQueue, poll chan struct{}, w *sync.WaitGroup, subscribe *gnmipb.SubscriptionList) {
	rc, ctx := common_utils.GetContext(c.ctx)
	c.ctx = ctx
	c.w = w

	defer func() {
		if r := recover(); r != nil {
			err := status.Errorf(codes.Internal, "%v", r)
			enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
		}
	}()

	defer c.w.Done()
	c.q = q
	c.channel = poll
	version := getBundleVersion(c.extensions)
	if version != nil {
		rc.BundleVersion = version
	}
	synced := false
	for {
		_, more := <-c.channel
		if !more {
			log.V(1).Infof("%v poll channel closed, exiting pollDb routine", c)
			enqueFatalMsgTranslib(c, "")
			return
		}
		t1 := time.Now()
		for gnmiPath, URIPath := range c.path2URI {
			if synced || !subscribe.UpdatesOnly {

				val, valTree, err := transutil.TranslProcessGet(URIPath, nil, c.ctx, c.encoding)
				if err != nil {
					switch err.(type) {
					case tlerr.NotFoundError:
						log.V(1).Infof("Subscribe Path Resource Not Found, ignoring: %v", URIPath)
					default:
						enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
						return
					}

				}
				v, err := buildValue(c.prefix, gnmiPath, c.encoding, val, valTree)
				if err != nil {
					enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe value failed to build =%v", err.Error()))
					return
				}
				if v != nil {
					c.q.EnqueueItem(Value{v})
					log.V(6).Infof("Added spbv #%v", v)
				}

			}
		}

		c.q.EnqueueItem(Value{
			&spb.Value{
				Timestamp:    time.Now().UnixNano(),
				SyncResponse: true,
			},
		})
		synced = true
		log.V(4).Infof("Sync done, poll time taken: %v ms", int64(time.Since(t1)/time.Millisecond))
	}
}

func (c *TranslClient) OnceRun(q *LimitedQueue, once chan struct{}, w *sync.WaitGroup, subscribe *gnmipb.SubscriptionList) {
	rc, ctx := common_utils.GetContext(c.ctx)
	c.ctx = ctx
	c.w = w

	defer func() {
		if r := recover(); r != nil {
			err := status.Errorf(codes.Internal, "%v", r)
			enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
		}
	}()

	defer c.w.Done()
	c.q = q
	c.channel = once
	version := getBundleVersion(c.extensions)
	if version != nil {
		rc.BundleVersion = version
	}

	_, more := <-c.channel
	if !more {
		log.V(1).Infof("%v once channel closed, exiting onceDb routine", c)
		enqueFatalMsgTranslib(c, "")
		return
	}
	t1 := time.Now()
	for gnmiPath, URIPath := range c.path2URI {

		val, valTree, err := transutil.TranslProcessGet(URIPath, nil, c.ctx, c.encoding)
		if err != nil {
			switch err.(type) {
			case tlerr.NotFoundError:
				log.V(1).Infof("Subscribe Path Resource Not Found, ignoring: %v", URIPath)
			default:
				enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error =%v", err.Error()))
				return
			}
		}
		if !subscribe.UpdatesOnly {
			v, err := buildValue(c.prefix, gnmiPath, c.encoding, val, valTree)
			if err != nil {
				enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe value failed to build =%v", err.Error()))
				return
			}
			if v != nil {
				c.q.EnqueueItem(Value{v})
				log.V(6).Infof("Added spbv #%v", v)
			}
		}

	}

	c.q.EnqueueItem(Value{
		&spb.Value{
			Timestamp:    time.Now().UnixNano(),
			SyncResponse: true,
		},
	})
	log.V(4).Infof("Sync done, once time taken: %v ms", int64(time.Since(t1)/time.Millisecond))

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
		notifications, err := ygot.TogNMINotifications(*valueTree, time.Now().UnixNano(), ygot.GNMINotificationsConfig{UsePathElem: true, PathElemPrefix: path.GetElem()})
		if err != nil {
			return nil, fmt.Errorf("Cannot convert OC Struct to Notifications: %s", err)
		}
		if len(notifications) != 1 {
			return nil, fmt.Errorf("YGOT returned wrong number of notifications")
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

// Set the desired encoding for Get and Subcribe responses
func (c *TranslClient) SetEncoding(enc gnmipb.Encoding) {
	c.encoding = enc
}
