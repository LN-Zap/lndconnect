package main

func ExampleLocalhost() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Version:      1,
		},
	}
	displayLink(c)

	// Output: lndconnect://127.0.0.1:0?cert=MIICiDCCAi-gAwIBAgIQdo5v0QBXHnji4hRaeeMjNDAKBggqhkjOPQQDAjBHMR8wHQYDVQQKExZsbmQgYXV0b2dlbmVyYXRlZCBjZXJ0MSQwIgYDVQQDExtKdXN0dXNzLU1hY0Jvb2stUHJvLTMubG9jYWwwHhcNMTgwODIzMDU1ODEwWhcNMTkxMDE4MDU1ODEwWjBHMR8wHQYDVQQKExZsbmQgYXV0b2dlbmVyYXRlZCBjZXJ0MSQwIgYDVQQDExtKdXN0dXNzLU1hY0Jvb2stUHJvLTMubG9jYWwwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASFhRm-w_T10PoKtg4lm9hBNJjJD473fkzHwPUFwy91vTrQSf7543j2JrgFo8mbTV0VtpgqkfK1IMVKMLrF21xio4H8MIH5MA4GA1UdDwEB_wQEAwICpDAPBgNVHRMBAf8EBTADAQH_MIHVBgNVHREEgc0wgcqCG0p1c3R1c3MtTWFjQm9vay1Qcm8tMy5sb2NhbIIJbG9jYWxob3N0ggR1bml4ggp1bml4cGFja2V0hwR_AAABhxAAAAAAAAAAAAAAAAAAAAABhxD-gAAAAAAAAAAAAAAAAAABhxD-gAAAAAAAAAwlc9Zck7bDhwTAqAEEhxD-gAAAAAAAABiNp__-GxXGhxD-gAAAAAAAAKWJ5tliDORjhwQKDwAChxD-gAAAAAAAAG6Wz__-3atFhxD92tDQyv4TAQAAAAAAABAAMAoGCCqGSM49BAMCA0cAMEQCIA9O9xtazmdxCKj0MfbFHVBq5I7JMnOFPpwRPJXQfrYaAiBd5NyJQCwlSx5ECnPOH5sRpv26T8aUcXbmynx9CoDufA&macaroon=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA
}

func ExampleQuery() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Query:     arrayFlags{"test=abc"},
			Version:      1,
		},
	}
	displayLink(c)

	// Output: lndconnect://127.0.0.1:0?cert=MIICiDCCAi-gAwIBAgIQdo5v0QBXHnji4hRaeeMjNDAKBggqhkjOPQQDAjBHMR8wHQYDVQQKExZsbmQgYXV0b2dlbmVyYXRlZCBjZXJ0MSQwIgYDVQQDExtKdXN0dXNzLU1hY0Jvb2stUHJvLTMubG9jYWwwHhcNMTgwODIzMDU1ODEwWhcNMTkxMDE4MDU1ODEwWjBHMR8wHQYDVQQKExZsbmQgYXV0b2dlbmVyYXRlZCBjZXJ0MSQwIgYDVQQDExtKdXN0dXNzLU1hY0Jvb2stUHJvLTMubG9jYWwwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASFhRm-w_T10PoKtg4lm9hBNJjJD473fkzHwPUFwy91vTrQSf7543j2JrgFo8mbTV0VtpgqkfK1IMVKMLrF21xio4H8MIH5MA4GA1UdDwEB_wQEAwICpDAPBgNVHRMBAf8EBTADAQH_MIHVBgNVHREEgc0wgcqCG0p1c3R1c3MtTWFjQm9vay1Qcm8tMy5sb2NhbIIJbG9jYWxob3N0ggR1bml4ggp1bml4cGFja2V0hwR_AAABhxAAAAAAAAAAAAAAAAAAAAABhxD-gAAAAAAAAAAAAAAAAAABhxD-gAAAAAAAAAwlc9Zck7bDhwTAqAEEhxD-gAAAAAAAABiNp__-GxXGhxD-gAAAAAAAAKWJ5tliDORjhwQKDwAChxD-gAAAAAAAAG6Wz__-3atFhxD92tDQyv4TAQAAAAAAABAAMAoGCCqGSM49BAMCA0cAMEQCIA9O9xtazmdxCKj0MfbFHVBq5I7JMnOFPpwRPJXQfrYaAiBd5NyJQCwlSx5ECnPOH5sRpv26T8aUcXbmynx9CoDufA&macaroon=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA&test=abc
}

func ExamplePort() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Port:      123,
			Version:      1,
		},
	}
	displayLink(c)

	// Output: lndconnect://127.0.0.1:123?cert=MIICiDCCAi-gAwIBAgIQdo5v0QBXHnji4hRaeeMjNDAKBggqhkjOPQQDAjBHMR8wHQYDVQQKExZsbmQgYXV0b2dlbmVyYXRlZCBjZXJ0MSQwIgYDVQQDExtKdXN0dXNzLU1hY0Jvb2stUHJvLTMubG9jYWwwHhcNMTgwODIzMDU1ODEwWhcNMTkxMDE4MDU1ODEwWjBHMR8wHQYDVQQKExZsbmQgYXV0b2dlbmVyYXRlZCBjZXJ0MSQwIgYDVQQDExtKdXN0dXNzLU1hY0Jvb2stUHJvLTMubG9jYWwwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASFhRm-w_T10PoKtg4lm9hBNJjJD473fkzHwPUFwy91vTrQSf7543j2JrgFo8mbTV0VtpgqkfK1IMVKMLrF21xio4H8MIH5MA4GA1UdDwEB_wQEAwICpDAPBgNVHRMBAf8EBTADAQH_MIHVBgNVHREEgc0wgcqCG0p1c3R1c3MtTWFjQm9vay1Qcm8tMy5sb2NhbIIJbG9jYWxob3N0ggR1bml4ggp1bml4cGFja2V0hwR_AAABhxAAAAAAAAAAAAAAAAAAAAABhxD-gAAAAAAAAAAAAAAAAAABhxD-gAAAAAAAAAwlc9Zck7bDhwTAqAEEhxD-gAAAAAAAABiNp__-GxXGhxD-gAAAAAAAAKWJ5tliDORjhwQKDwAChxD-gAAAAAAAAG6Wz__-3atFhxD92tDQyv4TAQAAAAAAABAAMAoGCCqGSM49BAMCA0cAMEQCIA9O9xtazmdxCKj0MfbFHVBq5I7JMnOFPpwRPJXQfrYaAiBd5NyJQCwlSx5ECnPOH5sRpv26T8aUcXbmynx9CoDufA&macaroon=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA
}

func ExampleHost() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Host:      "1.2.3.4",
			Version:      1,
		},
	}
	displayLink(c)

	// Output: lndconnect://1.2.3.4:0?cert=MIICiDCCAi-gAwIBAgIQdo5v0QBXHnji4hRaeeMjNDAKBggqhkjOPQQDAjBHMR8wHQYDVQQKExZsbmQgYXV0b2dlbmVyYXRlZCBjZXJ0MSQwIgYDVQQDExtKdXN0dXNzLU1hY0Jvb2stUHJvLTMubG9jYWwwHhcNMTgwODIzMDU1ODEwWhcNMTkxMDE4MDU1ODEwWjBHMR8wHQYDVQQKExZsbmQgYXV0b2dlbmVyYXRlZCBjZXJ0MSQwIgYDVQQDExtKdXN0dXNzLU1hY0Jvb2stUHJvLTMubG9jYWwwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASFhRm-w_T10PoKtg4lm9hBNJjJD473fkzHwPUFwy91vTrQSf7543j2JrgFo8mbTV0VtpgqkfK1IMVKMLrF21xio4H8MIH5MA4GA1UdDwEB_wQEAwICpDAPBgNVHRMBAf8EBTADAQH_MIHVBgNVHREEgc0wgcqCG0p1c3R1c3MtTWFjQm9vay1Qcm8tMy5sb2NhbIIJbG9jYWxob3N0ggR1bml4ggp1bml4cGFja2V0hwR_AAABhxAAAAAAAAAAAAAAAAAAAAABhxD-gAAAAAAAAAAAAAAAAAABhxD-gAAAAAAAAAwlc9Zck7bDhwTAqAEEhxD-gAAAAAAAABiNp__-GxXGhxD-gAAAAAAAAKWJ5tliDORjhwQKDwAChxD-gAAAAAAAAG6Wz__-3atFhxD92tDQyv4TAQAAAAAAABAAMAoGCCqGSM49BAMCA0cAMEQCIA9O9xtazmdxCKj0MfbFHVBq5I7JMnOFPpwRPJXQfrYaAiBd5NyJQCwlSx5ECnPOH5sRpv26T8aUcXbmynx9CoDufA&macaroon=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA
}

func ExampleNoCert() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Host:      "1.2.3.4",
			NoCert:    true,
			Version:      1,
		},
	}
	displayLink(c)

	// Output: lndconnect://1.2.3.4:0?macaroon=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA
}

func ExampleLocalhostV2() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Version:      2,
		},
	}
	displayLink(c)

	// Output: lndconn://127.0.0.1:0?c=somNeor-nJza1UFoBULswg6lRo3f1_-euBY7zV2DIUk&m=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA&v=2
}

func ExampleQueryV2() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Query:     arrayFlags{"test=abc"},
			Version:      2,
		},
	}
	displayLink(c)

	// Output: lndconn://127.0.0.1:0?c=somNeor-nJza1UFoBULswg6lRo3f1_-euBY7zV2DIUk&m=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA&test=abc&v=2
}

func ExamplePortV2() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Port:      123,
			Version:      2,
		},
	}
	displayLink(c)

	// Output: lndconn://127.0.0.1:123?c=somNeor-nJza1UFoBULswg6lRo3f1_-euBY7zV2DIUk&m=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA&v=2
}

func ExampleHostV2() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Host:      "1.2.3.4",
			Version:      2,
		},
	}
	displayLink(c)

	// Output: lndconn://1.2.3.4:0?c=somNeor-nJza1UFoBULswg6lRo3f1_-euBY7zV2DIUk&m=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA&v=2
}

func ExampleNoCertV2() {
	c := &config{
		TLSCertPath:  "testdata/tls.cert",
		AdminMacPath: "testdata/admin.macaroon",
		LndConnect: &lndConnectConfig{
			Localhost: true,
			Url:       true,
			Host:      "1.2.3.4",
			NoCert:    true,
			Version:      2,
		},
	}
	displayLink(c)

	// Output: lndconn://1.2.3.4:0?m=AgEDbG5kArsBAwoQ3_I9f6kgSE6aUPd85lWpOBIBMBoWCgdhZGRyZXNzEgRyZWFkEgV3cml0ZRoTCgRpbmZvEgRyZWFkEgV3cml0ZRoXCghpbnZvaWNlcxIEcmVhZBIFd3JpdGUaFgoHbWVzc2FnZRIEcmVhZBIFd3JpdGUaFwoIb2ZmY2hhaW4SBHJlYWQSBXdyaXRlGhYKB29uY2hhaW4SBHJlYWQSBXdyaXRlGhQKBXBlZXJzEgRyZWFkEgV3cml0ZQAABiAiUTBv3Eh6iDbdjmXCfNxp4HBEcOYNzXhrm-ncLHf5jA&v=2
}