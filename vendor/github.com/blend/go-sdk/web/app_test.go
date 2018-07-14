package web

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/util"
)

func controllerNoOp(_ *Ctx) Result { return nil }

func TestAppNew(t *testing.T) {
	assert := assert.New(t)

	var route *Route
	app := New()
	assert.NotEmpty(app.BindAddr())
	assert.NotNil(app.state)
	assert.NotNil(app.Views())
	app.GET("/", func(c *Ctx) Result {
		route = c.Route()
		return c.Raw([]byte("ok!"))
	})

	assert.Nil(app.Mock().Get("/").Execute())
	assert.NotNil(route)
	assert.Equal("GET", route.Method)
	assert.Equal("/", route.Path)
	assert.NotNil(route.Handler)
}

func TestAppNewFromEnv(t *testing.T) {
	assert := assert.New(t)

	var route *Route
	app := NewFromEnv()
	assert.NotNil(app.state)
	assert.NotNil(app.Views())
	app.GET("/", func(c *Ctx) Result {
		route = c.Route()
		return c.Raw([]byte("ok!"))
	})

	assert.Nil(app.Mock().Get("/").Execute())
	assert.NotNil(route)
	assert.Equal("GET", route.Method)
	assert.Equal("/", route.Path)
	assert.NotNil(route.Handler)
}

func TestAppNewFromConfig(t *testing.T) {
	assert := assert.New(t)

	app := NewFromConfig(&Config{
		BindAddr: ":5555",
		Port:     5000,
		HandleMethodNotAllowed: util.OptionalBool(true),
		HandleOptions:          util.OptionalBool(true),
		RecoverPanics:          util.OptionalBool(true),
		HSTS:                   util.OptionalBool(true),
		HSTSMaxAgeSeconds:      9999,
		HSTSPreload:            util.OptionalBool(false),
		HSTSIncludeSubDomains:  util.OptionalBool(false),
		MaxHeaderBytes:         128,
		ReadHeaderTimeout:      5 * time.Second,
		ReadTimeout:            6 * time.Second,
		IdleTimeout:            7 * time.Second,
		WriteTimeout:           8 * time.Second,

		CookieName: "A GOOD ONE",

		Views: ViewCacheConfig{
			Cached: util.OptionalBool(true),
		},
	})

	assert.Equal(":5555", app.BindAddr())
	assert.True(app.HandleMethodNotAllowed())
	assert.True(app.HandleOptions())
	assert.True(app.RecoverPanics())
	assert.Equal(128, app.MaxHeaderBytes())
	assert.Equal(5*time.Second, app.ReadHeaderTimeout())
	assert.Equal(6*time.Second, app.ReadTimeout())
	assert.Equal(7*time.Second, app.IdleTimeout())
	assert.Equal(8*time.Second, app.WriteTimeout())
	assert.Equal("A GOOD ONE", app.Auth().CookieName(), "we should use the auth config for the auth manager")
	assert.True(app.Views().Cached(), "we should use the view cache config for the view cache")

	assert.True(app.HSTS())
	assert.Equal(9999, app.HSTSMaxAgeSeconds())
	assert.False(app.HSTSIncludeSubdomains())
	assert.False(app.HSTSPreload())
}

func TestAppPathParams(t *testing.T) {
	assert := assert.New(t)

	var route *Route
	var params RouteParameters
	app := New()
	app.GET("/:uuid", func(c *Ctx) Result {
		route = c.Route()
		params = c.routeParameters
		return c.Raw([]byte("ok!"))
	})

	assert.Nil(app.Mock().Get("/foo").Execute())
	assert.NotNil(route)
	assert.Equal("GET", route.Method)
	assert.Equal("/:uuid", route.Path)
	assert.NotNil(route.Handler)

	assert.NotNil(params)
	assert.NotEmpty(params)
	assert.Equal("foo", params.Get("uuid"))
}

func TestAppPathParamsForked(t *testing.T) {
	assert := assert.New(t)

	var route *Route
	var params RouteParameters
	app := New()
	app.GET("/foos/bar/:uuid", func(c *Ctx) Result {
		route = c.Route()
		params = c.routeParameters
		return c.Raw([]byte("ok!"))
	})
	app.GET("/foo/:uuid", func(c *Ctx) Result { return nil })

	assert.Nil(app.Mock().Get("/foos/bar/foo").Execute())
	assert.NotNil(route)
	assert.Equal("GET", route.Method)
	assert.Equal("/foos/bar/:uuid", route.Path)
	assert.NotNil(route.Handler)

	assert.NotNil(params)
	assert.NotEmpty(params)
	assert.Equal("foo", params.Get("uuid"))
}

func TestAppSetLogger(t *testing.T) {
	assert := assert.New(t)

	log := logger.All()
	defer log.Close()
	app := New().WithLogger(log)
	assert.NotNil(app.Logger())
	assert.True(app.Logger().Flags().All())
}

func TestAppCtx(t *testing.T) {
	assert := assert.New(t)

	app := New()

	rc, err := app.Mock().CreateCtx(nil)
	assert.Nil(err)
	assert.NotNil(rc)
	assert.Nil(rc.log)

	result := rc.Raw([]byte("foo"))
	assert.NotNil(result)

	err = result.Render(rc)
	assert.Nil(err)
	assert.NotZero(rc.Response().ContentLength())
}

func TestAppCreateStaticMountedRoute(t *testing.T) {
	assert := assert.New(t)
	app := New()

	assert.Equal("/testPath/*filepath", app.createStaticMountRoute("/testPath/*filepath"))
	assert.Equal("/testPath/*filepath", app.createStaticMountRoute("/testPath/"))
	assert.Equal("/testPath/*filepath", app.createStaticMountRoute("/testPath"))
}

func TestAppStaticRewrite(t *testing.T) {
	assert := assert.New(t)
	app := New()

	app.ServeStatic("/testPath", "_static")
	assert.NotEmpty(app.statics)
	assert.NotNil(app.statics["/testPath/*filepath"])
	assert.Nil(app.SetStaticRewriteRule("/testPath", "(.*)", func(path string, pieces ...string) string {
		return path
	}))

	assert.NotEmpty(app.statics["/testPath/*filepath"].RewriteRules())
}

func TestAppStaticRewriteBadExp(t *testing.T) {
	assert := assert.New(t)
	app := New()
	app.ServeStatic("/testPath", "_static")
	assert.NotEmpty(app.statics)
	assert.NotNil(app.statics["/testPath/*filepath"])

	err := app.SetStaticRewriteRule("/testPath", "((((", func(path string, pieces ...string) string {
		return path
	})

	assert.NotNil(err)
	assert.Empty(app.statics["/testPath/*filepath"].RewriteRules())
}

func TestAppStaticHeader(t *testing.T) {
	assert := assert.New(t)
	app := New()
	app.ServeStatic("/testPath", "_static")
	assert.NotEmpty(app.statics)
	assert.NotNil(app.statics["/testPath/*filepath"])
	assert.Nil(app.SetStaticHeader("/testPath/*filepath", "cache-control", "haha what is caching."))
	assert.NotEmpty(app.statics["/testPath/*filepath"].Headers())
}

func TestAppMiddleWarePipeline(t *testing.T) {
	assert := assert.New(t)
	app := New()

	didRun := false
	app.GET("/",
		func(r *Ctx) Result { return r.Raw([]byte("OK!")) },
		func(action Action) Action {
			didRun = true
			return action
		},
		func(action Action) Action {
			return func(r *Ctx) Result {
				return r.Raw([]byte("foo"))
			}
		},
	)

	result, err := app.Mock().WithPathf("/").Bytes()
	assert.Nil(err)
	assert.True(didRun)
	assert.Equal("foo", string(result))
}

func TestAppStatic(t *testing.T) {
	assert := assert.New(t)
	app := New()
	app.ServeStatic("/static/*filepath", "testdata")

	index, err := app.Mock().WithPathf("/static/test_file.html").Bytes()
	assert.Nil(err)
	assert.True(strings.Contains(string(index), "Test!"), string(index))
}

func TestAppStaticSingleFile(t *testing.T) {
	assert := assert.New(t)
	app := New()
	app.GET("/", func(r *Ctx) Result {
		return r.Static("testdata/test_file.html")
	})

	index, err := app.Mock().WithPathf("/").Bytes()
	assert.Nil(err)
	assert.True(strings.Contains(string(index), "Test!"), string(index))
}

func TestAppProviderMiddleware(t *testing.T) {
	assert := assert.New(t)

	var okAction = func(r *Ctx) Result {
		assert.NotNil(r.DefaultResultProvider())
		return r.Raw([]byte("OK!"))
	}

	app := New()
	app.GET("/", okAction, JSONProviderAsDefault)

	err := app.Mock().WithPathf("/").Execute()
	assert.Nil(err)
}

func TestAppProviderMiddlewareOrder(t *testing.T) {
	assert := assert.New(t)

	app := New()

	var okAction = func(r *Ctx) Result {
		assert.NotNil(r.DefaultResultProvider())
		return r.Raw([]byte("OK!"))
	}

	var dependsOnProvider = func(action Action) Action {
		return func(r *Ctx) Result {
			assert.NotNil(r.DefaultResultProvider())
			return action(r)
		}
	}

	app.GET("/", okAction, dependsOnProvider, JSONProviderAsDefault)

	err := app.Mock().WithPathf("/").Execute()
	assert.Nil(err)
}

func TestAppDefaultResultProvider(t *testing.T) {
	assert := assert.New(t)

	app := New()
	assert.Nil(app.DefaultMiddleware())

	rc := app.createCtx(nil, nil, nil, nil, nil)
	assert.NotNil(rc.view)
	assert.NotNil(rc.json)
	assert.NotNil(rc.xml)
	assert.NotNil(rc.text)
	assert.NotNil(rc.defaultResultProvider)
}

func TestAppDefaultResultProviderWithDefault(t *testing.T) {
	assert := assert.New(t)
	app := New().WithDefaultMiddleware(ViewProviderAsDefault)
	assert.NotNil(app.DefaultMiddleware())

	rc := app.createCtx(nil, nil, nil, nil, nil)
	assert.NotNil(rc.view)
	assert.NotNil(rc.json)
	assert.NotNil(rc.xml)
	assert.NotNil(rc.text)

	// this will be set to the default initially
	assert.NotNil(rc.defaultResultProvider)

	app.GET("/", func(ctx *Ctx) Result {
		assert.NotNil(ctx.DefaultResultProvider())
		_, isTyped := ctx.DefaultResultProvider().(*ViewResultProvider)
		assert.True(isTyped)
		return nil
	})
}

func TestAppDefaultResultProviderWithDefaultFromRoute(t *testing.T) {
	assert := assert.New(t)

	app := New().WithDefaultMiddleware(JSONProviderAsDefault)
	app.Views().AddLiterals(DefaultTemplateNotAuthorized)
	app.GET("/", controllerNoOp, SessionRequired, ViewProviderAsDefault)

	//somehow assert that the content type is html
	meta, err := app.Mock().WithPathf("/").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(ContentTypeHTML, meta.Headers.Get(HeaderContentType))
}

func TestAppViewResult(t *testing.T) {
	assert := assert.New(t)

	app := New()
	app.Views().AddPaths("testdata/test_file.html")
	app.GET("/", func(r *Ctx) Result {
		return r.View().View("test", "foobarbaz")
	})

	res, meta, err := app.Mock().WithPathf("/").BytesWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(res))
	assert.Equal(ContentTypeHTML, meta.Headers.Get(HeaderContentType))
	assert.Contains(string(res), "foobarbaz")
}

func TestAppWritesLogs(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	agent := logger.New().WithFlags(logger.AllFlags()).WithWriter(logger.NewTextWriter(buffer))

	app := New().WithLogger(agent)
	app.GET("/", func(r *Ctx) Result {
		return r.Raw([]byte("ok!"))
	})
	err := app.Mock().Get("/").Execute()
	assert.Nil(err)
	assert.Nil(agent.Drain())

	assert.NotZero(buffer.Len())
	assert.NotEmpty(buffer.String())
}

func TestAppBindAddr(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvironmentVariableBindAddr, ":9999")
	env.Env().Set(EnvironmentVariablePort, "1111")
	defer env.Restore()

	assert.Equal(":3333", New().WithBindAddr(":3333").BindAddr())
	assert.Equal(":2222", New().WithPort(2222).BindAddr())
	assert.Equal(":9999", New().WithBindAddrFromEnv().BindAddr())
	assert.Equal(":1111", New().WithPortFromEnv().BindAddr())
}

func TestAppNotFound(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	agent := logger.New().WithFlags(logger.AllFlags()).WithWriter(logger.NewTextWriter(buffer).WithShowHeadings(true).WithUseColor(false).WithShowTimestamp(false))
	app := New().WithLogger(agent)
	app.GET("/", func(r *Ctx) Result {
		return r.Raw([]byte("ok!"))
	})

	wg := sync.WaitGroup{}
	wg.Add(1)

	app.WithNotFoundHandler(func(r *Ctx) Result {
		defer wg.Done()
		return r.JSON().NotFound()
	})

	agent.Listen(logger.HTTPResponse, "foo", logger.NewHTTPResponseEventListener(func(wre *logger.HTTPResponseEvent) {
		assert.NotNil(wre.Request())
		assert.Empty(wre.Route())
	}))

	err := app.Mock().Get("/doesntexist").Execute()
	assert.Nil(err)
	assert.Nil(agent.Drain())
	wg.Wait()
}

func TestAppDefaultHeaders(t *testing.T) {
	assert := assert.New(t)
	app := New().WithDefaultHeader("foo", "bar").WithDefaultHeader("baz", "buzz")
	app.GET("/", func(r *Ctx) Result {
		return r.Text().Result("ok")
	})

	meta, err := app.Mock().Get("/").ExecuteWithMeta()
	assert.Nil(err)
	assert.NotEmpty(meta.Headers)
	assert.Equal("bar", meta.Headers.Get("foo"))
	assert.Equal("buzz", meta.Headers.Get("baz"))
}

func TestAppIssuesHSTSHeaders(t *testing.T) {
	assert := assert.New(t)

	app := New().WithHSTS(true).WithHSTSMaxAgeSeconds(9999).WithHSTSIncludeSubdomains(true).WithHSTSPreload(true)
	app.GET("/", func(r *Ctx) Result {
		return r.Text().Result("ok")
	})
	assert.Nil(app.SetTLSCertPair([]byte(TestTLSCert), []byte(TestTLSKey)))

	meta, err := app.Mock().Get("/").ExecuteWithMeta()
	assert.Nil(err)
	assert.NotEmpty(meta.Headers)
	assert.NotEmpty(meta.Headers.Get(HeaderStrictTransportSecurity))
	assert.Equal("max-age=9999; includeSubDomains; preload", meta.Headers.Get(HeaderStrictTransportSecurity))
}

func TestAppTLSOptions(t *testing.T) {
	assert := assert.New(t)

	app := New()
	assert.NotNil(app.SetTLSCertPair([]byte{}, []byte{}))

	app = New()
	assert.Nil(app.SetTLSCertPair([]byte(TestTLSCert), []byte(TestTLSKey)))
	assert.NotNil(app.TLSConfig())
	assert.NotNil(app.TLSConfig().Certificates)

	app = New()
	assert.NotNil(app.SetTLSClientCertPool([]byte{}))
	app = New()
	assert.Nil(app.SetTLSClientCertPool([]byte(TestTLSCert)))
	assert.NotNil(app.TLSConfig())
	assert.NotNil(app.TLSConfig().ClientCAs)
	assert.NotNil(app.TLSConfig().GetConfigForClient)

	app = New()
	app.WithTLSClientCertVerification(tls.RequireAndVerifyClientCert)
	assert.NotNil(app.TLSConfig())
	assert.Equal(tls.RequireAndVerifyClientCert, app.TLSConfig().ClientAuth)

	app = New()
	app.WithTLSClientCertVerification(tls.RequireAndVerifyClientCert)
	assert.Nil(app.SetTLSCertPair([]byte(TestTLSCert), []byte(TestTLSKey)))
	assert.NotNil(app.TLSConfig())
	assert.NotNil(app.TLSConfig().Certificates)
	assert.Equal(tls.RequireAndVerifyClientCert, app.TLSConfig().ClientAuth)
}

func TestAppViewErrorsRenderErrorView(t *testing.T) {
	assert := assert.New(t)

	app := New()
	app.Views().AddLiterals(`{{ define "malformed" }} {{ .Ctx ALSKADJALSKDJA }} {{ end }}`)
	app.GET("/", func(r *Ctx) Result {
		return r.View().View("malformed", nil)
	})

	_, err := app.Mock().Get("/").Bytes()
	assert.NotNil(err)
}

func TestAppAddsDefaultHeaders(t *testing.T) {
	assert := assert.New(t)

	app := NewFromConfig(&Config{})
	app.WithBindAddr(DefaultIntegrationBindAddr)
	assert.NotEmpty(app.DefaultHeaders())
	app.GET("/", func(r *Ctx) Result {
		return r.Text().Result("OK!")
	})
	go app.Start()
	defer app.Shutdown()
	<-app.Started()

	res, err := http.Get("http://" + app.Listener().Addr().String() + "/")
	assert.Nil(err)
	assert.NotEmpty(res.Header)
	assert.Equal(PackageName, res.Header.Get(HeaderServer))
}

func TestAppHandlesPanics(t *testing.T) {
	assert := assert.New(t)

	app := NewFromConfig(&Config{})
	app.WithBindAddr(DefaultIntegrationBindAddr)
	app.GET("/", func(r *Ctx) Result {
		panic("this is only a test")
	})

	var didRecover bool
	go func() {
		defer func() {
			if r := recover(); r != nil {
				didRecover = true
			}
		}()
		app.Start()
	}()
	defer app.Shutdown()
	<-app.Started()

	res, err := http.Get("http://" + app.Listener().Addr().String() + "/")
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, res.StatusCode)
	assert.False(didRecover)
}
