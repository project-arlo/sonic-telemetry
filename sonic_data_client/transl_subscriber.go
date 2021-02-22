////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2021 Broadcom. The term Broadcom refers to Broadcom Inc. and/or //
//  its subsidiaries.                                                         //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License");           //
//  you may not use this file except in compliance with the License.          //
//  You may obtain a copy of the License at                                   //
//                                                                            //
//     http://www.apache.org/licenses/LICENSE-2.0                             //
//                                                                            //
//  Unless required by applicable law or agreed to in writing, software       //
//  distributed under the License is distributed on an "AS IS" BASIS,         //
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  //
//  See the License for the specific language governing permissions and       //
//  limitations under the License.                                            //
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

package client

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib"
	"github.com/Azure/sonic-mgmt-common/translib/path"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	spb "github.com/Azure/sonic-telemetry/proto"
	"github.com/Azure/sonic-telemetry/transl_utils"
	"github.com/Workiva/go-datastructures/queue"
	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/ygot"
)

// translSubscriber is an extension of TranslClient to service Subscribe RPC.
type translSubscriber struct {
	client      *TranslClient
	updatesOnly bool
	isHeartbeat bool
	stopOnSync  bool           // Stop upon sync message from translib
	synced      sync.WaitGroup // To signal receipt of sync message from translib
	pathMap     map[string]*gnmipb.Path
	filterFunc  func(string, *spb.Value) bool
}

// doSample invokes translib.Stream API to service SAMPLE, POLL and ONCE subscriptions.
// Timer for SAMPLE subscription should be managed outside.
func (ts *translSubscriber) doSample(p *gnmipb.Path) {
	// Temporary workaround to support legacy app implementations.
	// TODO remove once all apps have migrated to new subscribe infra
	if !path.HasWildcardKey(p) {
		ts.notifyUsingGet(p)
		return
	}

	c := ts.client
	paths := []string{c.path2URI[p]}
	req := translib.SubscribeRequest{
		Paths: paths,
		Q:     queue.NewPriorityQueue(1, false),
	}
	if c.version != nil {
		req.ClientVersion = *c.version
	}

	c.w.Add(1)
	ts.synced.Add(1)
	ts.stopOnSync = true
	go ts.processResponses(req.Q)

	err := translib.Stream(req)
	if err != nil {
		req.Q.Dispose()
		enqueFatalMsgTranslib(c, fmt.Sprintf("Subscribe operation failed with error = %v", err))
	}

	ts.synced.Wait()
}

// doOnChange handles the ON_CHANGE subscriptions through translib.Subscribe API.
// Returns only after initial updates and sync message are sent to the RPC queue.
func (ts *translSubscriber) doOnChange(stringPaths []string) {
	c := ts.client
	q := queue.NewPriorityQueue(1, false)

	req := translib.SubscribeRequest{
		Paths: stringPaths,
		Q:     q,
		Stop:  c.channel,
	}
	if c.version != nil {
		req.ClientVersion = *c.version
	}

	c.w.Add(1)
	ts.synced.Add(1)
	go ts.processResponses(q)

	_, err := translib.Subscribe(req)
	if err != nil {
		q.Dispose()
		enqueFatalMsgTranslib(c, "Subscribe operation failed with error: "+err.Error())
	}

	ts.synced.Wait()
}

// processResponses waits for SubscribeResponse messages from translib over a
// queue, formats them as spb.Value and pushes to the RPC queue.
func (ts *translSubscriber) processResponses(q *queue.PriorityQueue) {
	c := ts.client
	var syncDone bool
	defer c.w.Done()
	defer func() {
		if !syncDone {
			ts.synced.Done()
		}
	}()
	defer recoverSubscribe(c)

	for {
		items, err := q.Get(1)
		if err == queue.ErrDisposed {
			log.V(3).Info("PriorityQueue was disposed!")
			return
		}
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

			if v.SyncComplete && !syncDone {
				if ts.stopOnSync {
					log.V(6).Infof("Stopping on sync signal from translib")
					return
				}

				log.V(6).Infof("SENDING SYNC")
				enqueueSyncMessage(c)
				syncDone = true
				ts.synced.Done()
				break
			}

			//Don't send initial update with full object if user wants updates only.
			if ts.updatesOnly && !syncDone {
				log.V(1).Infof("Msg suppressed due to updates_only")
				break
			}

			ts.notify(v)

		default:
			log.V(1).Infof("Unknown data type %v for %s in queue", items[0], c)
		}
	}
}

func (ts *translSubscriber) notify(v *translib.SubscribeResponse) {
	c := ts.client

	//TODO change SubscribeResponse to return *gnmi.Path itself
	p := ts.pathMap[v.Path]
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
				return
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
			ju, err := ygotToIetfJsonValue(extraPrefix, v.Update)
			if err != nil {
				log.Error(err)
				return
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

	if ts.filterFunc != nil && !ts.filterFunc(v.Path, spbv) && !ts.isHeartbeat {
		log.V(6).Infof("Redundant message suppressed for path: %v", v.Path)
	} else {
		c.q.EnqueueItem(Value{spbv})
		log.V(6).Infof("Added spbv #%v", spbv)
	}
}

// notifyUsingGet performs translib.Get to retrieve data for a path
// and generates a notification message from it.
func (ts *translSubscriber) notifyUsingGet(p *gnmipb.Path) {
	pathStr := ts.client.path2URI[p]
	_, valTree, err := transl_utils.TranslProcessGet(pathStr, nil, ts.client.ctx, gnmipb.Encoding_PROTO)

	switch err.(type) {
	case nil: // no errors
	case tlerr.NotFoundError:
		err = nil
	default:
		enqueFatalMsgTranslib(ts.client, "DB Connection Error")
		return
	}

	resp := translib.SubscribeResponse{
		Path:      pathStr,
		Timestamp: time.Now().UnixNano(),
	}

	if valTree != nil {
		resp.Update = *valTree
	} else {
		//TODO send delete notification if applicable
		return
	}

	ts.notify(&resp)
}

// ygotToIetfJsonValue creates a gnmipb.Update object with IETF JSON value of a ygot struct.
func ygotToIetfJsonValue(targetElem *gnmipb.PathElem, yval ygot.ValidatedGoStruct) (*gnmipb.Update, error) {
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
		Elem:   elems,
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

//Compares val to cache, removes any duplicates from val,
//updates cache. Returns true if there are new values.
func filterValues(cache map[string]string, val *spb.Value, URI string, enc gnmipb.Encoding) bool {
	switch enc {
	case gnmipb.Encoding_JSON:
		fallthrough
	case gnmipb.Encoding_JSON_IETF:
		var v string
		if val.Notification != nil && len(val.Notification.Update) != 0 {
			v = string(val.Notification.Update[0].Val.GetJsonIetfVal())
		} else if val.Val != nil {
			v = string(val.Val.GetJsonIetfVal())
		}
		if len(v) != 0 && v != cache[URI] {
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
