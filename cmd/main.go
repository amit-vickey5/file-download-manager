package main

import (
	"context"
	"fmt"
	"github.com/amit/file-download-manager/internal/apiservice"
	"github.com/amit/file-download-manager/internal/boot"
	"github.com/amit/file-download-manager/internal/health"
	"github.com/amit/file-download-manager/internal/user"
	twirpServer "github.com/amit/file-download-manager/rpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twitchtv/twirp"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	authUsernameCtxKey = "authUser"
	authSecretKeyCtxKey = "authKey"
)

func main() {
	fmt.Println("initializing components...")
	initErr := boot.Init()
	if initErr != nil {
		fmt.Println("initializing components error :: ", initErr)
		return
	} else {
		fmt.Println("components initialized successfully...")
	}

	server := &apiservice.Server{}
	authHook := AuthenticationHook()
	twirpHandler := twirpServer.NewFileDownloadManagerServer(server, authHook)
	httpHandler := AddBasicAuthHeadersToHandler(twirpHandler)

	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), httpHandler)
	mux.Handle("/status", &health.StatusHandler{})
	mux.Handle("/metrics", promhttp.Handler())

	listener, err := net.Listen("tcp4", boot.Config.CoreConfig.Port)
	if err != nil {
		fmt.Println("failed to listen :: ", err)
		return
	}
	httpServer := http.Server{
		Handler: mux,
	}

	c := make(chan os.Signal, 1)
	// accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// or SIGTERM. SIGKILL, SIGQUIT will not be caught.
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := httpServer.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				fmt.Println("failed to start server :: ", err)
			}
		}
	}()
	fmt.Println("server started successfully...")

	// Block until signal is received.
	<-c
	// send unhealthy status to the healthcheck
	/*fmt.Println("marking server unhealthy...")
	healthCore.MarkUnhealthy()*/
	// wait for ShutdownDelay seconds
	time.Sleep(time.Duration(boot.Config.CoreConfig.ShutdownDelay) * time.Second)

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(boot.Config.CoreConfig.ShutdownTimeout)*time.Second)
	defer cancel()
	fmt.Println("shutting down server...")
	err = httpServer.Shutdown(ctxWithTimeout)
	if err != nil {
		fmt.Println("failed to shutdown server ::", err)
	}
	fmt.Println("server shut down successfully...")
}

func AuthenticationHook() *twirp.ServerHooks {
	hooks := &twirp.ServerHooks{}
	hooks.RequestRouted = func(ctx context.Context) (context.Context, error) {
		username := ctx.Value(authUsernameCtxKey).(string)
		suppliedSecretKey := ctx.Value(authSecretKeyCtxKey).(string)
		userDetails, err := getUserDetails(ctx, username)
		if err != nil {
			if err.Error() == "no_user_found_for_username" {
				return ctx, twirp.NewError(twirp.Unauthenticated, "no_user_found_for_username")
			} else {
				return ctx, twirp.NewError(twirp.Internal, err.Error())
			}
		}
		if userDetails.SecretKey != suppliedSecretKey {
			return ctx, twirp.NewError(twirp.Unauthenticated, "incorrect password")
		}
		ctx = context.WithValue(ctx, user.USER_ID_CONTEXT_KEY, userDetails.Id)
		return ctx, nil
	}
	return hooks
}

func getUserDetails(ctx context.Context, username string) (user.User, error) {
	userDetails, userErr := user.GetUserDetails(ctx, username)
	if userErr != nil {
		return user.User{}, userErr
	}
	return userDetails, nil
}

func AddBasicAuthHeadersToHandler(base http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		authUsername, authSecretKey, _ := request.BasicAuth()
		ctx = context.WithValue(ctx, authUsernameCtxKey, authUsername)
		ctx = context.WithValue(ctx, authSecretKeyCtxKey, authSecretKey)
		request = request.WithContext(ctx)
		base.ServeHTTP(responseWriter, request)
	})
}

