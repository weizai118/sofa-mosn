/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alipay/sofa-mosn/pkg/api/v2"
	"github.com/alipay/sofa-mosn/pkg/log"
	"github.com/alipay/sofa-mosn/pkg/tls/certtool"
	"github.com/alipay/sofa-mosn/pkg/types"
)

type MockListener struct {
	net.Listener
	Mng types.TLSContextManager
}

func (ln MockListener) Accept() (net.Conn, error) {
	conn, err := ln.Listener.Accept()
	if err != nil {
		return conn, err
	}
	return ln.Mng.Conn(conn), nil
}

type MockServer struct {
	Mng    types.TLSContextManager
	Addr   string
	server *http.Server
}

func (s *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "mock server")
}
func (s *MockServer) GoListenAndServe(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)
	server := &http.Server{
		Handler: mux,
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Errorf("listen failed %v", err)
		return
	}
	s.Addr = ln.Addr().String()
	s.server = server
	listener := MockListener{ln, s.Mng}
	go server.Serve(listener)
}
func (s *MockServer) Close() {
	s.server.Close()
}

type certInfo struct {
	CommonName string
	Curve      string
	DNS        string
}

func (c *certInfo) CreateCertConfig() (*v2.TLSConfig, error) {
	priv, err := certtool.GeneratePrivateKey(c.Curve)
	if err != nil {
		return nil, fmt.Errorf("generate key failed %v", err)
	}
	var dns []string
	if c.DNS != "" {
		dns = append(dns, c.DNS)
	}
	tmpl, err := certtool.CreateTemplate(c.CommonName, false, dns)
	if err != nil {
		return nil, fmt.Errorf("generate certificate template failed %v", err)
	}
	cert, err := certtool.SignCertificate(tmpl, priv)
	if err != nil {
		return nil, fmt.Errorf("sign certificate failed %v", err)
	}
	return &v2.TLSConfig{
		Status:     true,
		CACert:     certtool.GetRootCA().CertPem,
		CertChain:  cert.CertPem,
		PrivateKey: cert.KeyPem,
	}, nil
}

// TestServerContextManagerWithMultipleCert tests the contextManager's core logic
// make three certificates with different dns and common name
// test context manager can find correct certificate for different client
func TestServerContextManagerWithMultipleCert(t *testing.T) {
	var filterChains []v2.FilterChain
	testCases := []struct {
		Info *certInfo
		Addr string
	}{
		{Info: &certInfo{"Cert1", "RSA", "www.example.com"}, Addr: "www.example.com"},
		{Info: &certInfo{"Cert2", "RSA", "*.example.com"}, Addr: "test.example.com"},
		{Info: &certInfo{"Cert3", "P256", "*.com"}, Addr: "www.foo.com"},
	}
	for i, tc := range testCases {
		cfg, err := tc.Info.CreateCertConfig()
		if err != nil {
			t.Errorf("#%d %v", i, err)
			return
		}
		fc := v2.FilterChain{
			TLS: *cfg,
		}
		filterChains = append(filterChains, fc)
	}
	lc := &v2.ListenerConfig{
		FilterChains: filterChains,
	}
	ctxMng, err := NewTLSServerContextManager(lc, nil, log.StartLogger)
	if err != nil {
		t.Errorf("create context manager failed %v", err)
		return
	}
	server := MockServer{
		Mng: ctxMng,
	}
	server.GoListenAndServe(t)
	defer server.Close()
	time.Sleep(time.Second) //wait server start
	// request with different "servername"
	// context manager just find a certificate to response
	// the certificate may be not match the client
	url := server.Addr
	for i, tc := range testCases {
		trans := &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName:         tc.Addr,
				InsecureSkipVerify: true,
			},
		}
		client := &http.Client{Transport: trans}
		resp, err := client.Get("https://" + url)
		if err != nil {
			t.Errorf("#%d request server error %v", i, err)
			continue
		}
		serverCN := resp.TLS.PeerCertificates[0].Subject.CommonName
		if serverCN != tc.Info.CommonName {
			t.Errorf("#%d expected request server config %s , but got %s", i, tc.Info.CommonName, serverCN)
		}
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	// request a unknown server name, return the first certificate
	trans := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName:         "www.example.net",
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: trans}
	resp, err := client.Get("https://" + url)
	if err != nil {
		t.Errorf("request server error %v", err)
		return
	}
	defer resp.Body.Close()
	serverCN := resp.TLS.PeerCertificates[0].Subject.CommonName
	expected := testCases[0].Info.CommonName
	if serverCN != expected {
		t.Errorf("expected request server config  %s , but got %s", expected, serverCN)
	}
	ioutil.ReadAll(resp.Body)
}

// TestVerifyClient tests a client must have certificate to server
func TestVerifyClient(t *testing.T) {
	info := &certInfo{
		CommonName: "test",
		Curve:      "P256",
	}
	cfg, err := info.CreateCertConfig()
	if err != nil {
		t.Error(err)
		return
	}
	cfg.VerifyClient = true
	filterChains := []v2.FilterChain{
		{
			TLS: *cfg,
		},
	}
	lc := &v2.ListenerConfig{
		FilterChains: filterChains,
	}
	ctxMng, err := NewTLSServerContextManager(lc, nil, log.StartLogger)
	if err != nil {
		t.Errorf("create context manager failed %v", err)
		return
	}
	server := MockServer{
		Mng: ctxMng,
	}
	server.GoListenAndServe(t)
	defer server.Close()
	time.Sleep(time.Second) //wait server start
	clientConfigs := []*v2.TLSConfig{
		// Verify Server
		{
			Status:     true,
			CACert:     cfg.CACert,
			CertChain:  cfg.CertChain,
			PrivateKey: cfg.PrivateKey,
		},
		// Skip Verify Server
		{
			Status:       true,
			CertChain:    cfg.CertChain,
			PrivateKey:   cfg.PrivateKey,
			InsecureSkip: true,
		},
	}
	for i, cfg := range clientConfigs {
		cltMng, err := NewTLSClientContextManager(cfg, nil)
		if err != nil {
			t.Errorf("#%d create client context manager failed %v", i, err)
			continue
		}
		trans := &http.Transport{
			TLSClientConfig: cltMng.Config(),
		}
		client := &http.Client{Transport: trans}
		resp, err := client.Get("https://" + server.Addr)
		if err != nil {
			t.Errorf("#%d request server error %v", i, err)
			continue
		}
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	trans := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: trans}
	resp, err := client.Get("https://" + server.Addr)
	// expected bad certificate
	if err == nil {
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		t.Errorf("server should verify client certificate")
		return
	}
}

// TestInspector tests context manager support both tls and non-tls
func TestInspector(t *testing.T) {
	info := &certInfo{
		CommonName: "test",
		Curve:      "P256",
	}
	cfg, err := info.CreateCertConfig()
	if err != nil {
		t.Error(err)
		return
	}
	cfg.VerifyClient = true
	filterChains := []v2.FilterChain{
		{
			TLS: *cfg,
		},
	}
	lc := &v2.ListenerConfig{
		Inspector:    true,
		FilterChains: filterChains,
	}
	ctxMng, err := NewTLSServerContextManager(lc, nil, log.StartLogger)
	if err != nil {
		t.Errorf("create context manager failed %v", err)
		return
	}
	server := MockServer{
		Mng: ctxMng,
	}
	server.GoListenAndServe(t)
	defer server.Close()
	time.Sleep(time.Second) //wait server start
	testCases := []string{
		"http://" + server.Addr,
		"https://" + server.Addr,
	}
	cltMng, err := NewTLSClientContextManager(&v2.TLSConfig{
		Status:     true,
		CACert:     cfg.CACert,
		CertChain:  cfg.CertChain,
		PrivateKey: cfg.PrivateKey,
	}, nil)
	if err != nil {
		t.Errorf("create client context manager failed %v", err)
		return
	}
	trans := &http.Transport{
		TLSClientConfig: cltMng.Config(),
	}
	client := &http.Client{Transport: trans}
	for i, tc := range testCases {
		resp, err := client.Get(tc)
		if err != nil {
			t.Errorf("#%d request server error %v", i, err)
			return
		}
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
}

// test ConfigHooks
// define VerifyPeerCertificate, verify common name instead of san, ignore keyusage
type testConfigHooks struct {
	defaultConfigHooks
	Name           string
	Root           *x509.CertPool
	PassCommonName string
}

// over write
func (hook *testConfigHooks) VerifyPeerCertificate() func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	return hook.verifyPeerCertificate
}
func (hook *testConfigHooks) GetX509Pool(caIndex string) (*x509.CertPool, error) {
	return hook.Root, nil
}

// verifiedChains is always nil
func (hook *testConfigHooks) verifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	var certs []*x509.Certificate
	for _, asn1Data := range rawCerts {
		cert, err := x509.ParseCertificate(asn1Data)
		if err != nil {
			return err
		}
		certs = append(certs, cert)
	}
	intermediates := x509.NewCertPool()
	for _, cert := range certs[1:] {
		intermediates.AddCert(cert)
	}
	opts := x509.VerifyOptions{
		Roots:         hook.Root,
		Intermediates: intermediates,
	}
	leaf := certs[0]
	_, err := leaf.Verify(opts)
	if err != nil {
		return err
	}
	if leaf.Subject.CommonName != hook.PassCommonName {
		return errors.New("common name miss match")
	}
	return nil
}

func pass(resp *http.Response, err error) bool {
	if err != nil {
		return false
	}
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return true
}
func fail(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return false
}

const testType = "test"

type testConfigHooksFactory struct{}

func (f *testConfigHooksFactory) CreateConfigHooks(config map[string]interface{}) ConfigHooks {
	c := make(map[string]string)
	for k, v := range config {
		if s, ok := v.(string); ok {
			c[strings.ToLower(k)] = s
		}
	}
	root := certtool.GetRootCA()
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(root.CertPem))
	return &testConfigHooks{
		defaultConfigHooks: defaultConfigHooks{},
		Name:               c["name"],
		PassCommonName:     c["cn"],
		Root:               pool,
	}
}

// TestTLSExtensionsVerifyClient tests server allow request with certificate's common name is client only
func TestTLSExtensionsVerifyClient(t *testing.T) {
	// Server
	extendVerify := map[string]interface{}{
		"name": "server",
		"cn":   "client",
	}
	serverInfo := &certInfo{
		CommonName: extendVerify["name"].(string),
		Curve:      "RSA",
	}
	serverConfig, err := serverInfo.CreateCertConfig()
	if err != nil {
		t.Errorf("create server certificate error %v", err)
		return
	}
	serverConfig.VerifyClient = true
	serverConfig.Type = testType
	serverConfig.ExtendVerify = extendVerify
	filterChains := []v2.FilterChain{
		{
			TLS: *serverConfig,
		},
	}
	lc := &v2.ListenerConfig{
		FilterChains: filterChains,
	}
	ctxMng, err := NewTLSServerContextManager(lc, nil, log.StartLogger)
	if err != nil {
		t.Errorf("create context manager failed %v", err)
		return
	}
	server := MockServer{
		Mng: ctxMng,
	}
	server.GoListenAndServe(t)
	defer server.Close()
	time.Sleep(time.Second) //wait server start
	testCases := []struct {
		Info *certInfo
		Pass func(resp *http.Response, err error) bool
	}{
		{
			Info: &certInfo{
				CommonName: extendVerify["cn"].(string),
				Curve:      serverInfo.Curve,
			},
			Pass: pass,
		},
		{
			Info: &certInfo{
				CommonName: "invalid client",
				Curve:      serverInfo.Curve,
			},
			Pass: fail,
		},
	}
	url := "https://" + server.Addr
	for i, tc := range testCases {
		cfg, err := tc.Info.CreateCertConfig()
		if err != nil {
			t.Errorf("#%d create client certificate error %v", i, err)
			continue
		}
		cltMng, err := NewTLSClientContextManager(cfg, nil)
		if err != nil {
			t.Errorf("#%d create client context manager failed %v", i, err)
			continue
		}
		trans := &http.Transport{
			TLSClientConfig: cltMng.Config(),
		}
		client := &http.Client{Transport: trans}
		resp, err := client.Get(url)
		if !tc.Pass(resp, err) {
			t.Errorf("#%d verify failed", i)
		}
	}
}

// TestTestTLSExtensionsVerifyServer tests client accept server response with cerificate's common name is server only
func TestTestTLSExtensionsVerifyServer(t *testing.T) {
	extendVerify := map[string]interface{}{
		"name": "client",
		"cn":   "server",
	}
	clientInfo := &certInfo{
		CommonName: extendVerify["name"].(string),
		Curve:      "RSA",
	}
	clientConfig, err := clientInfo.CreateCertConfig()
	if err != nil {
		t.Errorf("create client certificate error %v", err)
		return
	}
	clientConfig.Type = testType
	clientConfig.ExtendVerify = extendVerify
	cltMng, err := NewTLSClientContextManager(clientConfig, nil)
	if err != nil {
		t.Errorf("create client context manager failed %v", err)
		return
	}
	testCases := []struct {
		Info *certInfo
		Pass func(resp *http.Response, err error) bool
	}{
		{
			Info: &certInfo{
				CommonName: extendVerify["cn"].(string),
				Curve:      clientInfo.Curve,
				DNS:        "www.pass.com",
			},
			Pass: pass,
		},
		{
			Info: &certInfo{
				CommonName: "invalid server",
				Curve:      clientInfo.Curve,
				DNS:        "www.fail.com",
			},
			Pass: fail,
		},
	}
	var filterChains []v2.FilterChain
	for i, tc := range testCases {
		cfg, err := tc.Info.CreateCertConfig()
		if err != nil {
			t.Errorf("#%d %v", i, err)
			return
		}
		fc := v2.FilterChain{
			TLS: *cfg,
		}
		filterChains = append(filterChains, fc)
	}
	lc := &v2.ListenerConfig{
		FilterChains: filterChains,
	}
	ctxMng, err := NewTLSServerContextManager(lc, nil, log.StartLogger)
	if err != nil {
		t.Errorf("create context manager failed %v", err)
		return
	}
	server := MockServer{
		Mng: ctxMng,
	}
	server.GoListenAndServe(t)
	defer server.Close()
	time.Sleep(time.Second) //wait server start
	url := "https://" + server.Addr
	for i, tc := range testCases {
		cfg := cltMng.Config()
		cfg.ServerName = tc.Info.DNS
		trans := &http.Transport{
			TLSClientConfig: cfg,
		}
		client := &http.Client{Transport: trans}
		resp, err := client.Get(url)
		if !tc.Pass(resp, err) {
			t.Errorf("#%d verify failed", i)
		}
	}
	// insecure skip will skip even if it is registered
	skipConfig := &v2.TLSConfig{
		Status:       true,
		Type:         clientConfig.Type,
		CACert:       clientConfig.CACert,
		CertChain:    clientConfig.CertChain,
		PrivateKey:   clientConfig.PrivateKey,
		InsecureSkip: true,
	}
	skipMng, err := NewTLSClientContextManager(skipConfig, nil)
	if err != nil {
		t.Errorf("create client context manager failed %v", err)
		return
	}
	for i, tc := range testCases {
		cfg := skipMng.Config()
		cfg.ServerName = tc.Info.DNS
		trans := &http.Transport{
			TLSClientConfig: cfg,
		}
		client := &http.Client{Transport: trans}
		resp, err := client.Get(url)
		// ignore the case, must be pass
		if !pass(resp, err) {
			t.Errorf("#%d skip verify failed", i)
		}
	}
}

func TestMain(m *testing.M) {
	Register(testType, &testConfigHooksFactory{})
	m.Run()
}
