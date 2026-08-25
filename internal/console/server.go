package console

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"venueops/internal/audit"
	"venueops/internal/booth"
	"venueops/internal/camera"
	"venueops/internal/gate"
	"venueops/internal/guard"
	"venueops/internal/hall"
	"venueops/internal/light"
	"venueops/internal/ns"
	"venueops/internal/quota"
	"venueops/internal/screen"
	"venueops/internal/store"
	"venueops/internal/zone"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	store    *store.Store
	nsSvc    *ns.Service
	hallSvc  *hall.Service
	zoneSvc  *zone.Service
	quotaSvc *quota.Service
	boothSvc *booth.Service
	gateSvc  *gate.Service
	lightSvc *light.Service
	screenSvc *screen.Service
	cameraSvc *camera.Service
	guardSvc *guard.Service
	auditSvc *audit.Service
	router   chi.Router
	pages    map[string]*template.Template
}

func NewServer(st *store.Store, webDir string) (*Server, error) {
	auditSvc := audit.NewService(st)
	nsSvc := ns.NewService(st)
	hallSvc := hall.NewService(st)
	zoneSvc := zone.NewService(st)
	quotaSvc := quota.NewService(st)
	boothSvc := booth.NewService(st, quotaSvc)
	gateSvc := gate.NewService(st, auditSvc)
	lightSvc := light.NewService(st, auditSvc)
	screenSvc := screen.NewService(st)
	lightSvc.SetHandover(screenSvc.HandoverToEmergency)
	var _ light.HandoverFunc = screenSvc.HandoverToEmergency
	cameraSvc := camera.NewService(st, hallSvc)
	guardSvc := guard.NewService(st, zoneSvc, hallSvc, auditSvc)
	var _ zone.SceneProvider = lightSvc.SceneFor

	pages, err := loadPages(webDir)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:     st,
		nsSvc:     nsSvc,
		hallSvc:   hallSvc,
		zoneSvc:   zoneSvc,
		quotaSvc:  quotaSvc,
		boothSvc:  boothSvc,
		gateSvc:   gateSvc,
		lightSvc:  lightSvc,
		screenSvc: screenSvc,
		cameraSvc: cameraSvc,
		guardSvc:  guardSvc,
		auditSvc:  auditSvc,
		pages:     pages,
	}
	server.routes()
	return server, nil
}

func (s *Server) Router() http.Handler {
	return s.router
}

func loadPages(webDir string) (map[string]*template.Template, error) {
	base := filepath.Join(webDir, "base.html")
	names := []string{"halls", "screens", "lights", "guards"}
	pages := map[string]*template.Template{}
	for _, name := range names {
		path := filepath.Join(webDir, name+".html")
		if _, err := os.Stat(path); err != nil {
			return nil, err
		}
		tmpl, err := template.ParseFiles(base, path)
		if err != nil {
			return nil, err
		}
		pages[name] = tmpl
	}
	return pages, nil
}
