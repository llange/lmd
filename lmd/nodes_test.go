package main

import (
	"testing"
)

func TestNodeManager(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping nodes test in short mode")
	}
	extraConfig := `
		Listen = ['test.sock', 'http://127.0.0.1:8901']
		Nodes = ['http://127.0.0.1:8901', 'http://127.0.0.2:8902']
	`
	peer := StartTestPeerExtra(4, 10, 10, extraConfig)
	PauseTestPeers(peer)

	if nodeAccessor == nil {
		t.Fatalf("nodeAccessor should not be nil")
	}
	if err := assertEq(nodeAccessor.IsClustered(), true); err != nil {
		t.Fatalf("Nodes are not clustered")
	}
	if nodeAccessor.thisNode == nil {
		t.Fatalf("thisNode should not be nil")
	}
	if !(nodeAccessor.thisNode.String() != "") {
		t.Fatalf("got a name")
	}

	// test host request
	res, err := peer.QueryString("GET hosts\nColumns: name peer_key state\n\n")
	if err != nil {
		t.Fatal(err)
	}

	if err = assertEq(40, len(res)); err != nil {
		t.Error(err)
	}

	// test host stats request
	res, err = peer.QueryString("GET hosts\nStats: name !=\nStats: avg latency\nStats: sum latency\n\n")
	if err != nil {
		t.Fatal(err)
	}
	if err = assertEq(40, int(res[0][0].(float64))); err != nil {
		t.Error(err)
	}
	if err = assertEq(0.24065613746999998, res[0][1]); err != nil {
		t.Error(err)
	}
	if err = assertEq(9.6262454988, res[0][2]); err != nil {
		t.Error(err)
	}

	// test host grouped stats request
	res, err = peer.QueryString("GET hosts\nColumns: name alias\nStats: name !=\nStats: avg latency\nStats: sum latency\n\n")
	if err != nil {
		t.Fatal(err)
	}
	if err = assertEq("testhost_1", res[0][0]); err != nil {
		t.Error(err)
	}
	if err = assertEq("tomcat", res[0][1]); err != nil {
		t.Error(err)
	}
	if err = assertEq(4.0, res[0][2]); err != nil {
		t.Error(err)
	}
	if err = assertEq(0.24065613747, res[0][3]); err != nil {
		t.Error(err)
	}

	// test host empty stats request
	res, err = peer.QueryString("GET hosts\nFilter: check_type = 15\nStats: sum percent_state_change\nStats: min percent_state_change\n\n")
	if err != nil {
		t.Fatal(err)
	}
	if err = assertEq(1, len(res)); err != nil {
		t.Fatal(err)
	}
	if err = assertEq(float64(0), res[0][0]); err != nil {
		t.Error(err)
	}

	if err := StopTestPeer(peer); err != nil {
		panic(err.Error())
	}
}
