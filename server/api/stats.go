package api

import (
	"net/http"

	"github.com/gorilla/mux"
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-incident-collaboration/server/bot"
	"github.com/mattermost/mattermost-plugin-incident-collaboration/server/sqlstore"
)

type StatsHandler struct {
	pluginAPI  *pluginapi.Client
	log        bot.Logger
	statsStore *sqlstore.StatsStore
}

func NewStatsHandler(router *mux.Router, api *pluginapi.Client, log bot.Logger, statsStore *sqlstore.StatsStore) *StatsHandler {
	handler := &StatsHandler{
		pluginAPI:  api,
		log:        log,
		statsStore: statsStore,
	}

	statsRouter := router.PathPrefix("/stats").Subrouter()
	statsRouter.HandleFunc("", handler.stats).Methods(http.MethodGet)

	return handler
}

type Stats struct {
	TotalReportedIncidents                int `json:"total_reported_incidents"`
	TotalActiveIncidents                  int `json:"total_active_incidents"`
	TotalActiveParticipants               int `json:"total_active_participants"`
	AverageDurationActiveIncidentsMinutes int `json:"average_duration_active_incidents_minutes"`

	ActiveIncidents        []int `json:"active_incidents"`
	PeopleInIncidents      []int `json:"people_in_incidents"`
	AverageStartToActive   []int `json:"average_start_to_active"`
	AverageStartToResolved []int `json:"average_start_to_resolved"`
}

func (h *StatsHandler) stats(w http.ResponseWriter, r *http.Request) {
	stats := Stats{
		TotalReportedIncidents:                h.statsStore.TotalReportedIncidents(),
		TotalActiveIncidents:                  h.statsStore.TotalActiveIncidents(),
		TotalActiveParticipants:               h.statsStore.TotalActiveParticipants(),
		AverageDurationActiveIncidentsMinutes: h.statsStore.AverageDurationActiveIncidentsMinutes(),

		ActiveIncidents:        h.statsStore.ActiveIncidents(),
		PeopleInIncidents:      h.statsStore.PeopleInIncidents(),
		AverageStartToActive:   h.statsStore.AverageStartToActive(),
		AverageStartToResolved: h.statsStore.AverageStartToResolved(),
	}

	ReturnJSON(w, stats, http.StatusOK)
}
