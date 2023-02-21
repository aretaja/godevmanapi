package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aretaja/godevmanapi/config"
	_ "github.com/aretaja/godevmanapi/docs"
	"github.com/aretaja/godevmanapi/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	Conf    *config.Configuration
	Router  *chi.Mux
	Handler *handlers.Handler
	Version string
}

func (a *App) Initialize() {
	// Config
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	a.Conf = c

	// Router instance
	a.Router = chi.NewRouter()
	a.initializeMiddleware()

	// Handler instance
	a.Handler = new(handlers.Handler)
	err = a.Handler.Initialize(a.Conf.DbURL, a.Conf.Salt)
	if err != nil {
		log.Fatal(err)
	}

	// Handler instance
	a.initializeRoutes()
}

// Midleware activation
func (a *App) initializeMiddleware() {
	r := a.Router

	// httplog configuration
	logConf := httplog.Options{
		JSON:    true,
		Concise: false,
	}

	if os.Getenv("GODEVMANAPI_LOGPLAIN") != "" {
		logConf.JSON = false
		logConf.Concise = true
	}

	logger := httplog.NewLogger("godevmanapi", logConf)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

}

// Route definitions
func (a *App) initializeRoutes() {
	r := a.Router

	// Welcome
	r.Get("/", a.Handler.Hello)

	// Version
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		handlers.VersionSwagger() // Prevent function not used warning
		handlers.RespondJSON(w, r, http.StatusOK, handlers.StatusResponse{
			Code:    strconv.Itoa(http.StatusOK),
			Message: a.Version,
		})
	})

	// Routes for "/archived/interfaces" resource
	r.Route("/archived/interfaces", func(r chi.Router) {
		// Filter parameters:
		//   ifindex_f, hostname_f, host_ip4_f, host_ip6_f, descr_f, alias_f, mac_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetArchivedInterfaces)
		r.Get("/count", a.Handler.CountArchivedInterfaces)
		r.Post("/", a.Handler.CreateArchivedInterface)

		// Subroutes
		r.Route("/{ifa_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetArchivedInterface)
			r.Put("/", a.Handler.UpdateArchivedInterface)
			r.Delete("/", a.Handler.DeleteArchivedInterface)
		})
	})

	// Routes for "/archived/subinterfaces" resource
	r.Route("/archived/subinterfaces", func(r chi.Router) {
		// Filter parameters:
		//   ifindex_f, hostname_f, host_ip4_f, host_ip6_f, descr_f, alias_f, mac_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetArchivedSubinterfaces)
		r.Get("/count", a.Handler.CountArchivedSubinterfaces)
		r.Post("/", a.Handler.CreateArchivedSubinterface)

		// Subroutes
		r.Route("/{ifa_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetArchivedSubinterface)
			r.Put("/", a.Handler.UpdateArchivedSubinterface)
			r.Delete("/", a.Handler.DeleteArchivedSubinterface)
		})
	})
	// Routes for "/config" resource
	// Routes for "/config/credentials" resource
	r.Route("/config/credentials", func(r chi.Router) {
		// Filter parameters:
		//   label_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetCredentials)
		r.Get("/count", a.Handler.CountCredentials)
		r.Post("/", a.Handler.CreateCredential)

		// Subroutes
		r.Route("/{cred_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetCredential)
			r.Put("/", a.Handler.UpdateCredential)
			r.Delete("/", a.Handler.DeleteCredential)
		})
	})

	// Routes for "/connections" resource
	r.Route("/connections", func(r chi.Router) {
		// Filter parameters:
		//   hint_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetConnections)
		r.Get("/count", a.Handler.CountConnections)
		r.Post("/", a.Handler.CreateConnection)

		// Subroutes
		r.Route("/{con_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetConnection)
			r.Put("/", a.Handler.UpdateConnection)
			r.Delete("/", a.Handler.DeleteConnection)
			r.Get("/capacity", a.Handler.GetConnectionConCapacitiy)
			r.Get("/class", a.Handler.GetConnectionConClass)
			r.Get("/provider", a.Handler.GetConnectionConProvider)
			// r.Get("/site", a.Handler.GetConnectionSite)
			r.Get("/type", a.Handler.GetConnectionConType)
			// r.Get("/interfaces", a.Handler.GetConnectionInterfaces)
		})
	})

	// Routes for "/connections/capacities" resource
	r.Route("/connections/capacities", func(r chi.Router) {
		// Filter parameters:
		//   descr_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetConCapacities)
		r.Get("/count", a.Handler.CountConCapacities)
		r.Post("/", a.Handler.CreateConCapacity)

		// Subroutes
		r.Route("/{con_cap_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetConCapacity)
			r.Put("/", a.Handler.UpdateConCapacity)
			r.Delete("/", a.Handler.DeleteConCapacity)
			r.Get("/connections", a.Handler.GetConCapacityConnections)
		})
	})

	// Routes for "/connections/classes" resource
	r.Route("/connections/classes", func(r chi.Router) {
		// Filter parameters:
		//   descr_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetConClasses)
		r.Get("/count", a.Handler.CountConClasses)
		r.Post("/", a.Handler.CreateConClass)

		// Subroutes
		r.Route("/{con_class_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetConClass)
			r.Put("/", a.Handler.UpdateConClass)
			r.Delete("/", a.Handler.DeleteConClass)
			r.Get("/connections", a.Handler.GetConClassConnections)
		})
	})

	// Routes for "/connections/providers" resource
	r.Route("/connections/providers", func(r chi.Router) {
		// Filter parameters:
		//   descr_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetConProviders)
		r.Get("/count", a.Handler.CountConProviders)
		r.Post("/", a.Handler.CreateConProvider)

		// Subroutes
		r.Route("/{con_prov_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetConProvider)
			r.Put("/", a.Handler.UpdateConProvider)
			r.Delete("/", a.Handler.DeleteConProvider)
			r.Get("/connections", a.Handler.GetConProviderConnections)
		})
	})

	// Routes for "/connections/types" resource
	r.Route("/connections/types", func(r chi.Router) {
		// Filter parameters:
		//   descr_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetConTypes)
		r.Get("/count", a.Handler.CountConTypes)
		r.Post("/", a.Handler.CreateConType)

		// Subroutes
		r.Route("/{con_type_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetConType)
			r.Put("/", a.Handler.UpdateConType)
			r.Delete("/", a.Handler.DeleteConType)
			r.Get("/connections", a.Handler.GetConTypeConnections)
		})
	})

	// Routes for "/devices" resource
	r.Route("/devices", func(r chi.Router) {
		// Filter parameters:
		//   sys_id_f, host_name_f, sw_version_f, notes_f, name_f, ip4_addr_f, ip6_addr_f
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetDevices)
		r.Get("/count", a.Handler.CountDevices)
		r.Post("/", a.Handler.CreateDevice)

		// Subroutes
		r.Route("/{dev_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetDevice)
			r.Put("/", a.Handler.UpdateDevice)
			r.Delete("/", a.Handler.DeleteDevice)
			r.Get("/childs", a.Handler.GetDeviceChilds)
			r.Get("/credentials", a.Handler.GetDeviceDeviceCredentials)
			r.Get("/domain", a.Handler.GetDeviceDeviceDomain)
			r.Get("/entities", a.Handler.GetDeviceEntities)
			r.Get("/extensions", a.Handler.GetDeviceDeviceExtensions)
			r.Get("/interfaces", a.Handler.GetDeviceInterfaces)
			r.Get("/ip_interfaces", a.Handler.GetDeviceIpInterfaces)
			r.Get("/licenses", a.Handler.GetDeviceDeviceLicenses)
			r.Get("/ospf_nbrs", a.Handler.GetDeviceOspfNbrs)
			r.Get("/parent", a.Handler.GetDeviceParent)
			r.Get("/peer_xconnects", a.Handler.GetDevicePeerXconnects)
			r.Get("/rl_nbrs", a.Handler.GetDeviceRlNbrs)
			r.Get("/site", a.Handler.GetDeviceSite)
			r.Get("/snmp_credentials_main", a.Handler.GetDeviceSnmpCredentialsMain)
			r.Get("/snmp_credentials_ro", a.Handler.GetDeviceSnmpCredentialsRo)
			r.Get("/state", a.Handler.GetDeviceDeviceState)
			r.Get("/type", a.Handler.GetDeviceDeviceType)
			r.Get("/vlans", a.Handler.GetDeviceVlans)
			r.Get("/xconnects", a.Handler.GetDeviceXconnects)
		})
	})

	// Routes for "/devices/classes" resource
	r.Route("/devices/classes", func(r chi.Router) {
		// Filter parameters:
		//   descr_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetDeviceClasses)
		r.Get("/count", a.Handler.CountDeviceClasses)
		r.Post("/", a.Handler.CreateDeviceClass)

		// Subroutes
		r.Route("/{class_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetDeviceClass)
			r.Put("/", a.Handler.UpdateDeviceClass)
			r.Delete("/", a.Handler.DeleteDeviceClass)
			r.Get("/types", a.Handler.GetDeviceClassTypes)
		})
	})

	// Routes for "/devices/credentials" resource
	r.Route("/devices/credentials", func(r chi.Router) {
		// Filter parameters:
		//   username_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetDeviceCredentials)
		r.Get("/count", a.Handler.CountDeviceCredentials)
		r.Post("/", a.Handler.CreateDeviceCredential)

		// Subroutes
		r.Route("/{cred_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetDeviceCredential)
			r.Put("/", a.Handler.UpdateDeviceCredential)
			r.Delete("/", a.Handler.DeleteDeviceCredential)
		})
	})

	// Routes for "/devices/domains" resource
	r.Route("/devices/domains", func(r chi.Router) {
		// Filter parameters:
		//   descr_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetDeviceDomains)
		r.Get("/count", a.Handler.CountDeviceDomains)
		r.Post("/", a.Handler.CreateDeviceDomain)

		// Subroutes
		r.Route("/{dom_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetDeviceDomain)
			r.Put("/", a.Handler.UpdateDeviceDomain)
			r.Delete("/", a.Handler.DeleteDeviceDomain)
			r.Get("/devices", a.Handler.GetDeviceDomainDevices)
		})
	})

	// Routes for "/devices/snmp_credentials" resource
	r.Route("/devices/snmp_credentials", func(r chi.Router) {
		// Filter parameters:
		//   label_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetSnmpCredentials)
		r.Get("/count", a.Handler.CountSnmpCredentials)
		r.Post("/", a.Handler.CreateSnmpCredential)

		// Subroutes
		r.Route("/{snmp_cred_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetSnmpCredential)
			r.Put("/", a.Handler.UpdateSnmpCredential)
			r.Delete("/", a.Handler.DeleteSnmpCredential)
			r.Get("/main_devices", a.Handler.GetSnmpCredentialsMainDevices)
			r.Get("/ro_devices", a.Handler.GetSnmpCredentialsRoDevices)
		})
	})

	// Routes for "/devices/types" resource
	r.Route("/devices/types", func(r chi.Router) {
		// Filter parameters:
		//   sys_id_f, manufacturer_f, model_f
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetDeviceTypes)
		r.Get("/count", a.Handler.CountDeviceTypes)
		r.Post("/", a.Handler.CreateDeviceType)

		// Subroutes
		r.Route("/{sys_id:[\\w-\\.]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetDeviceType)
			r.Put("/", a.Handler.UpdateDeviceType)
			r.Delete("/", a.Handler.DeleteDeviceType)
			r.Get("/class", a.Handler.GetDeviceTypeClass)
			r.Get("/devices", a.Handler.GetDeviceTypeDevices)
		})
	})

	// Routes for "/entities" resource
	r.Route("/entities", func(r chi.Router) {
		// Filter parameters:
		//   slot_f, descr_f, model_f, hw_product_f, hw_revision_f, serial_nr_f, sw_product_f, sw_revision_f, manufacturer_f
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetEntities)
		r.Get("/count", a.Handler.CountEntities)
		r.Post("/", a.Handler.CreateEntity)

		// Subroutes
		r.Route("/{ent_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetEntity)
			r.Put("/", a.Handler.UpdateEntity)
			r.Delete("/", a.Handler.DeleteEntity)
			r.Get("/childs", a.Handler.GetEntityChilds)
			r.Get("/device", a.Handler.GetEntityDevice)
			r.Get("/parent", a.Handler.GetEntityParent)
			r.Get("/entity_phy_indexes", a.Handler.GetEntityEntityPhyIndexes)
			r.Get("/interfaces", a.Handler.GetEntityInterfaces)
			r.Get("/rl_nbrs", a.Handler.GetEntityRlfNbrs)
		})
	})

	// Routes for "/entities/custom_entities" resource
	r.Route("/entities/custom_entities", func(r chi.Router) {
		// Filter parameters:
		//   serial_nr_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetCustomEntities)
		r.Get("/count", a.Handler.CountCustomEntities)
		r.Post("/", a.Handler.CreateCustomEntity)

		// Subroutes
		r.Route("/{cent_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetCustomEntity)
			r.Put("/", a.Handler.UpdateCustomEntity)
			r.Delete("/", a.Handler.DeleteCustomEntity)
		})
	})

	// Routes for "/interfaces" resource
	r.Route("/interfaces", func(r chi.Router) {
		// Filter parameters:
		//   ifindex_f, descr_f, alias_f, oper_f, adm_f, speed_f, minspeed_f, type_enum_f, mac_f
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetInterfaces)
		r.Get("/count", a.Handler.CountInterfaces)
		r.Post("/", a.Handler.CreateInterface)

		// Subroutes
		r.Route("/{if_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetInterface)
			r.Put("/", a.Handler.UpdateInterface)
			r.Delete("/", a.Handler.DeleteInterface)
		})
	})

	// Routes for "/interfaces/bw_stat" resource
	r.Route("/interfaces/bw_stats", func(r chi.Router) {
		// Filter parameters:
		//   if_group_f,
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetIntBwStats)
		r.Get("/count", a.Handler.CountIntBwStats)
		r.Post("/", a.Handler.CreateIntBwStat)

		// Subroutes
		r.Route("/{bw_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetIntBwStat)
			r.Put("/", a.Handler.UpdateIntBwStat)
			r.Delete("/", a.Handler.DeleteIntBwStat)
			r.Get("/interfaces", a.Handler.GetIntBwStatInterface)
		})
	})

	// Routes for "/sites" resource

	// Routes for "/sites/countries" resource
	r.Route("/sites/countries", func(r chi.Router) {
		// Filter parameters:
		//   descr_f, code_f
		//   updated_ge, updated_le, created_ge, created_le
		// Pagination parameters:
		//   count(100), start(0).
		//   Uses default if not set.
		r.Get("/", a.Handler.GetCountries)
		r.Get("/count", a.Handler.CountCountries)
		r.Post("/", a.Handler.CreateCountry)

		// Subroutes
		r.Route("/{country_id:[0-9]+}", func(r chi.Router) {
			r.Get("/", a.Handler.GetCountry)
			r.Put("/", a.Handler.UpdateCountry)
			r.Delete("/", a.Handler.DeleteCountry)
			// r.Get("/sites", a.Handler.GetCountrySites)
		})
	})

	// Swagger
	r.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(
			httpSwagger.DocExpansion("none"),
		))
	})

	// Custom 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondError(w, r, http.StatusNotFound, "Route does not exist")
	})

	// Custom 405 handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondError(w, r, http.StatusMethodNotAllowed, "Method is not valid")
	})

}

func (a *App) Run() {
	fmt.Printf("Starting up on:%s\n", a.Conf.ApiListen)
	err := http.ListenAndServe(a.Conf.ApiListen, a.Router)
	if err != nil {
		log.Printf("Failed to launch api server:%+v\n", err)
	}
}
