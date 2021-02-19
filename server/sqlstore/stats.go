package sqlstore

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-plugin-incident-collaboration/server/bot"
	"github.com/mattermost/mattermost-server/v5/model"
)

type StatsStore struct {
	pluginAPI PluginAPIClient
	log       bot.Logger
	store     *SQLStore
}

func NewStatsStore(pluginAPI PluginAPIClient, log bot.Logger, sqlStore *SQLStore) *StatsStore {
	return &StatsStore{
		pluginAPI: pluginAPI,
		log:       log,
		store:     sqlStore,
	}
}

func (s *StatsStore) GetTotalActiveIncidents() int {
	query := s.store.builder.
		Select("COUNT(ID)").
		From("IR_Incident").
		Where("EndAt = 0")

	var total int
	if err := s.store.getBuilder(s.store.db, &total, query); err != nil {
		// TODO: Error properly
		return -1
	}

	return total
}

func (s *StatsStore) GetTotalActiveParticipants() int {
	query := s.store.builder.
		Select("COUNT(DISTINCT cm.UserId)").
		From("ChannelMembers as cm").
		Join("IR_Incident AS i ON i.ChannelId = cm.ChannelId").
		Where("i.EndAt = 0")

	var total int
	if err := s.store.getBuilder(s.store.db, &total, query); err != nil {
		// TODO: Error properly
		return -1
	}

	return total
}

func (s *StatsStore) GetAverageDurationActiveIncidentsMinutes() int {
	query := s.store.builder.
		Select("AVG(c.CreateAt)").
		From("IR_Incident AS i").
		Join("Channels AS c ON (c.Id = i.ChannelId)").
		Where("i.EndAt = 0")

	var averageCreateAt float64
	if err := s.store.getBuilder(s.store.db, &averageCreateAt, query); err != nil {
		// TODO: Error properly
		return -1
	}

	return int((float64(model.GetMillis()) - averageCreateAt) / 60000)
}

type PlaybookUse struct {
	Name    string `json:"name"`
	NumUses int    `json:"num_uses"`
}

func (s *StatsStore) GetPlaybookUses() []PlaybookUse {
	query := s.store.builder.
		Select("pb.Title as Name, count(i.PlaybookID) as NumUses").
		From("IR_Incident as i").
		Join("IR_Playbook as pb ON i.PlaybookID = pb.ID").
		GroupBy("pb.Title")

	var uses []PlaybookUse
	if err := s.store.selectBuilder(s.store.db, &uses, query); err != nil {
		// TODO: Error properly
		fmt.Println(err)
		return []PlaybookUse{}
	}

	return uses
}

func (s *StatsStore) GetActiveIncidentsOverTime() []int {
	now := model.GetMillis()
	dayInMS := int64(86400)

	windowStart := now - dayInMS
	windowEnd := now
	numIncidents := []int{}
	for i := int64(0); i < 14; i++ {
		query := s.store.builder.
			Select("COUNT(i.Id)").
			From("IR_Incident as i").
			Join("Channels AS c ON (c.Id = i.ChannelId)").
			Where(sq.Expr(`c.CreateAt < ? AND i.EndAt > ?`, windowEnd-(i*dayInMS), windowStart-(i*dayInMS)))

		var numActiveOnDay int
		if err := s.store.getBuilder(s.store.db, &numActiveOnDay, query); err != nil {
			// TODO: Error properly
			fmt.Println(err)
			return numIncidents
		}

		numIncidents = append(numIncidents, numActiveOnDay)
	}

	return numIncidents
}
