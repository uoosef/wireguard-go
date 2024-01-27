package psiphon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Psiphon-Labs/psiphon-tunnel-core/psiphon"
	"github.com/refraction-networking/conjure/pkg/station/log"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Parameters provide an easier way to modify the tunnel config at runtime.
type Parameters struct {
	// Used as the directory for the datastore, remote server list, and obfuscasted
	// server list.
	// Empty string means the default will be used (current working directory).
	// nil means the values in the config file will be used.
	// Optional, but strongly recommended.
	DataRootDirectory *string

	// Overrides config.ClientPlatform. See config.go for details.
	// nil means the value in the config file will be used.
	// Optional, but strongly recommended.
	ClientPlatform *string

	// Overrides config.NetworkID. For details see:
	// https://godoc.org/github.com/Psiphon-Labs/psiphon-tunnel-core/psiphon#NetworkIDGetter
	// nil means the value in the config file will be used. (If not set in the config,
	// an error will result.)
	// Empty string will produce an error.
	// Optional, but strongly recommended.
	NetworkID *string

	// Overrides config.EstablishTunnelTimeoutSeconds. See config.go for details.
	// nil means the EstablishTunnelTimeoutSeconds value in the config file will be used.
	// If there's no such value in the config file, the default will be used.
	// Zero means there will be no timeout.
	// Optional.
	EstablishTunnelTimeoutSeconds *int

	// EmitDiagnosticNoticesToFile indicates whether to use the rotating log file
	// facility to record diagnostic notices instead of sending diagnostic
	// notices to noticeReceiver. Has no effect unless the tunnel
	// config.EmitDiagnosticNotices flag is set.
	EmitDiagnosticNoticesToFiles bool
}

// Tunnel is the tunnel object. It can be used for stopping the tunnel and
// retrieving proxy ports.
type Tunnel struct {
	embeddedServerListWaitGroup sync.WaitGroup
	controllerWaitGroup         sync.WaitGroup
	stopController              context.CancelFunc

	// The port on which the HTTP proxy is running
	HTTPProxyPort int
	// The port on which the SOCKS proxy is running
	SOCKSProxyPort int
}

// ParametersDelta allows for fine-grained modification of parameters.Parameters.
// NOTE: Ordinary users of this library should never need this.
type ParametersDelta map[string]interface{}

// NoticeEvent represents the notices emitted by tunnel core. It will be passed to
// noticeReceiver, if supplied.
// NOTE: Ordinary users of this library should never need this.
type NoticeEvent struct {
	Data      map[string]interface{} `json:"data"`
	Type      string                 `json:"noticeType"`
	Timestamp string                 `json:"timestamp"`
}

// ErrTimeout is returned when the tunnel establishment attempt fails due to timeout
var ErrTimeout = errors.New("clientlib: tunnel establishment timeout")

// StartTunnel establishes a Psiphon tunnel. It returns an error if the establishment
// was not successful. If the returned error is nil, the returned tunnel can be used
// to find out the proxy ports and subsequently stop the tunnel.
//
// ctx may be cancelable, if the caller wants to be able to interrupt the establishment
// attempt, or context.Background().
//
// configJSON will be passed to psiphon.LoadConfig to configure the tunnel. Required.
//
// embeddedServerEntryList is the encoded embedded server entry list. It is optional.
//
// params are config values that typically need to be overridden at runtime.
//
// paramsDelta contains changes that will be applied to the Parameters.
// NOTE: Ordinary users of this library should never need this and should pass nil.
//
// noticeReceiver, if non-nil, will be called for each notice emitted by tunnel core.
// NOTE: Ordinary users of this library should never need this and should pass nil.
func StartTunnel(
	ctx context.Context,
	configJSON []byte,
	embeddedServerEntryList string,
	params Parameters,
	paramsDelta ParametersDelta,
	noticeReceiver func(NoticeEvent)) (retTunnel *Tunnel, retErr error) {

	config, err := psiphon.LoadConfig(configJSON)
	if err != nil {
		return nil, errors.New("failed to load config file")
	}

	// Use params.DataRootDirectory to set related config values.
	if params.DataRootDirectory != nil {
		config.DataRootDirectory = *params.DataRootDirectory

		// Migrate old fields
		config.MigrateDataStoreDirectory = *params.DataRootDirectory
		config.MigrateObfuscatedServerListDownloadDirectory = *params.DataRootDirectory
		config.MigrateRemoteServerListDownloadFilename = filepath.Join(*params.DataRootDirectory, "server_list_compressed")
	}

	if params.NetworkID != nil {
		config.NetworkID = *params.NetworkID
	}

	if params.ClientPlatform != nil {
		config.ClientPlatform = *params.ClientPlatform
	} // else use the value in config

	if params.EstablishTunnelTimeoutSeconds != nil {
		config.EstablishTunnelTimeoutSeconds = params.EstablishTunnelTimeoutSeconds
	} // else use the value in config

	if config.UseNoticeFiles == nil && config.EmitDiagnosticNotices && params.EmitDiagnosticNoticesToFiles {
		config.UseNoticeFiles = &psiphon.UseNoticeFiles{
			RotatingFileSize:      0,
			RotatingSyncFrequency: 0,
		}
	} // else use the value in the config

	// config.Commit must be called before calling config.SetParameters
	// or attempting to connect.
	err = config.Commit(true)
	if err != nil {
		return nil, errors.New("config.Commit failed")
	}

	// If supplied, apply the parameters delta
	if len(paramsDelta) > 0 {
		err = config.SetParameters("", false, paramsDelta)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("SetParameters failed for delta: %v", paramsDelta))
		}
	}

	// Will receive a value when the tunnel has successfully connected.
	connected := make(chan struct{}, 1)
	// Will receive a value if an error occurs during the connection sequence.
	errored := make(chan error, 1)

	// Create the tunnel object
	tunnel := new(Tunnel)

	// Set up notice handling
	psiphon.SetNoticeWriter(psiphon.NewNoticeReceiver(
		func(notice []byte) {
			var event NoticeEvent
			err := json.Unmarshal(notice, &event)
			if err != nil {
				// This is unexpected and probably indicates something fatal has occurred.
				// We'll interpret it as a connection error and abort.
				err = errors.New("failed to unmarshal notice JSON")
				select {
				case errored <- err:
				default:
				}
				return
			}

			if event.Type == "ListeningHttpProxyPort" {
				port := event.Data["port"].(float64)
				tunnel.HTTPProxyPort = int(port)
			} else if event.Type == "ListeningSocksProxyPort" {
				port := event.Data["port"].(float64)
				tunnel.SOCKSProxyPort = int(port)
			} else if event.Type == "EstablishTunnelTimeout" {
				select {
				case errored <- ErrTimeout:
				default:
				}
			} else if event.Type == "Tunnels" {
				count := event.Data["count"].(float64)
				if count > 0 {
					select {
					case connected <- struct{}{}:
					default:
					}
				}
			}

			// Some users of this package may need to add special processing of notices.
			// If the caller has requested it, we'll pass on the notices.
			if noticeReceiver != nil {
				noticeReceiver(event)
			}
		}))

	err = psiphon.OpenDataStore(config)
	if err != nil {
		return nil, errors.New("failed to open data store")
	}
	// Make sure we close the datastore in case of error
	defer func() {
		if retErr != nil {
			tunnel.controllerWaitGroup.Wait()
			tunnel.embeddedServerListWaitGroup.Wait()
			psiphon.CloseDataStore()
		}
	}()

	// Create a cancelable context that will be used for stopping the tunnel
	var controllerCtx context.Context
	controllerCtx, tunnel.stopController = context.WithCancel(ctx)

	// If specified, the embedded server list is loaded and stored. When there
	// are no server candidates at all, we wait for this import to complete
	// before starting the Psiphon controller. Otherwise, we import while
	// concurrently starting the controller to minimize delay before attempting
	// to connect to existing candidate servers.
	//
	// If the import fails, an error notice is emitted, but the controller is
	// still started: either existing candidate servers may suffice, or the
	// remote server list fetch may obtain candidate servers.
	//
	// The import will be interrupted if it's still running when the controller
	// is stopped.
	tunnel.embeddedServerListWaitGroup.Add(1)
	go func() {
		defer tunnel.embeddedServerListWaitGroup.Done()

		err := psiphon.ImportEmbeddedServerEntries(
			controllerCtx,
			config,
			"",
			embeddedServerEntryList)
		if err != nil {
			psiphon.NoticeError("error importing embedded server entry list: %s", err)
			return
		}
	}()
	if !psiphon.HasServerEntries() {
		psiphon.NoticeInfo("awaiting embedded server entry list import")
		tunnel.embeddedServerListWaitGroup.Wait()
	}

	// Create the Psiphon controller
	controller, err := psiphon.NewController(config)
	if err != nil {
		tunnel.stopController()
		tunnel.embeddedServerListWaitGroup.Wait()
		return nil, errors.New("psiphon.NewController failed")
	}

	// Begin tunnel connection
	tunnel.controllerWaitGroup.Add(1)
	go func() {
		defer tunnel.controllerWaitGroup.Done()

		// Start the tunnel. Only returns on error (or internal timeout).
		controller.Run(controllerCtx)

		// controller.Run does not exit until the goroutine that posts
		// EstablishTunnelTimeout has terminated; so, if there was a
		// EstablishTunnelTimeout event, ErrTimeout is guaranteed to be sent to
		// errord before this next error and will be the StartTunnel return value.

		var err error
		switch ctx.Err() {
		case context.DeadlineExceeded:
			err = ErrTimeout
		case context.Canceled:
			err = errors.New("StartTunnel canceled")
		default:
			err = errors.New("controller.Run exited unexpectedly")
		}
		select {
		case errored <- err:
		default:
		}
	}()

	// Wait for an active tunnel or error
	select {
	case <-connected:
		return tunnel, nil
	case err := <-errored:
		tunnel.Stop()
		if err != ErrTimeout {
			err = errors.New("tunnel start produced error")
		}
		return nil, err
	}
}

// Stop stops/disconnects/shuts down the tunnel. It is safe to call when not connected.
// Not safe to call concurrently with Start.
func (tunnel *Tunnel) Stop() {
	if tunnel.stopController == nil {
		return
	}
	tunnel.stopController()
	tunnel.controllerWaitGroup.Wait()
	tunnel.embeddedServerListWaitGroup.Wait()
	psiphon.CloseDataStore()
}

func RunPsiphon(wgBind, localSocksPort, country string) context.Context {
	ctx := context.Background()
	// Embedded configuration
	host, port, err := net.SplitHostPort(localSocksPort)
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(host, "127.0.0") {
		host = ""
	} else {
		host = "any"
	}
	configJSON := `{
		"EgressRegion": "` + country + `",
		"ListenInterface": "` + host + `",
		"LocalSocksProxyPort": ` + port + `,
		"UpstreamProxyURL": "socks5://` + wgBind + `",
		"DisableLocalHTTPProxy": true,
		"PropagationChannelId":"FFFFFFFFFFFFFFFF",
		"RemoteServerListDownloadFilename":"remote_server_list",
		"RemoteServerListSignaturePublicKey":"MIICIDANBgkqhkiG9w0BAQEFAAOCAg0AMIICCAKCAgEAt7Ls+/39r+T6zNW7GiVpJfzq/xvL9SBH5rIFnk0RXYEYavax3WS6HOD35eTAqn8AniOwiH+DOkvgSKF2caqk/y1dfq47Pdymtwzp9ikpB1C5OfAysXzBiwVJlCdajBKvBZDerV1cMvRzCKvKwRmvDmHgphQQ7WfXIGbRbmmk6opMBh3roE42KcotLFtqp0RRwLtcBRNtCdsrVsjiI1Lqz/lH+T61sGjSjQ3CHMuZYSQJZo/KrvzgQXpkaCTdbObxHqb6/+i1qaVOfEsvjoiyzTxJADvSytVtcTjijhPEV6XskJVHE1Zgl+7rATr/pDQkw6DPCNBS1+Y6fy7GstZALQXwEDN/qhQI9kWkHijT8ns+i1vGg00Mk/6J75arLhqcodWsdeG/M/moWgqQAnlZAGVtJI1OgeF5fsPpXu4kctOfuZlGjVZXQNW34aOzm8r8S0eVZitPlbhcPiR4gT/aSMz/wd8lZlzZYsje/Jr8u/YtlwjjreZrGRmG8KMOzukV3lLmMppXFMvl4bxv6YFEmIuTsOhbLTwFgh7KYNjodLj/LsqRVfwz31PgWQFTEPICV7GCvgVlPRxnofqKSjgTWI4mxDhBpVcATvaoBl1L/6WLbFvBsoAUBItWwctO2xalKxF5szhGm8lccoc5MZr8kfE0uxMgsxz4er68iCID+rsCAQM=",
		"RemoteServerListUrl":"https://s3.amazonaws.com//psiphon/web/mjr4-p23r-puwl/server_list_compressed",
		"SponsorId":"FFFFFFFFFFFFFFFF",
		"UseIndistinguishableTLS":true
	}`

	dir := "."
	ClientPlatform := "Android_4.0.4_com.example.exampleClientLibraryApp"
	network := "test"
	timeout := 60

	p := Parameters{
		DataRootDirectory:             &dir,
		ClientPlatform:                &ClientPlatform,
		NetworkID:                     &network,
		EstablishTunnelTimeoutSeconds: &timeout,
		EmitDiagnosticNoticesToFiles:  false,
	}

	log.Println("Handshaking, Please Wait...")

	var tunnel *Tunnel
	startTime := time.Now()

	go func() {
		for {
			tunnel, err = StartTunnel(ctx, []byte(configJSON), "", p, nil, nil)
			if err == nil {
				break
			}
			log.Error("Unable to start psiphon", err, "reconnecting...")
			time.Sleep(2 * time.Second)
		}

		log.Println("Psiphon started successfully on port", tunnel.SOCKSProxyPort, "handshake operation took", int64(time.Since(startTime)/time.Millisecond), "milliseconds")
	}()

	return ctx
}
