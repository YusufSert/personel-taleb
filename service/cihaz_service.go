package service

import (
	"context"
	"errors"
	"log/slog"
	"websocketTest/repo"
)

type CihazService struct {
	repo *repo.CihazRepo
}

func NewCihazService() (*CihazService, error) {
	db, err := repo.NewRepo("user=postgres password=Banana@@ dbname=cihazlar sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &CihazService{db}, nil
}

func (s CihazService) GetAll() ([]repo.Cihaz, error) {
	cihazler, err := s.repo.GetAll(context.Background())
	if err != nil {
		return nil, err
	}
	return cihazler, nil
}

func (s CihazService) GetIndirimTalebEkraniData() ([]repo.CihazJoinTalebJoinRealiziton, error) {
	records, err := s.repo.GetCihazJoinTalebJoinRealization(context.Background())
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s CihazService) AddTaleb(ctx context.Context, user string, cihazId int) error {
	if user == "" || cihazId == 0 {
		return errors.New("bad argument error")
	}

	r, err := s.repo.GetRealizationByCihazIdAndState(ctx, cihazId, "pending", "active")
	for _, e := range r { // only one row
		if e.TalebState == "active" {
			slog.Info("indirim talebi bulundu", "id", e.Id, "state", e.TalebState)
			return errors.New("there already active discount")
		}

		slog.Info("indirim talebi bulundu", "id", e.Id, "state", e.TalebState)
		err = s.repo.AddTaleb(ctx, user, cihazId)
		if err != nil {
			return err
		}
		slog.Info("yeni taleb eklendi")

		c, err := s.repo.GetNotEvaulatedTalebCountByCihazId(ctx, cihazId)
		if err != nil {
			return err
		}
		slog.Info("cihaz taleb sayısı bulundu", "ader", c)

		if e.RequiredTalebCount <= c {
			err = s.repo.EvaluateTalebsWithCihazId(ctx, e.Id, cihazId)
			if err != nil {
				return err
			}
			slog.Info("talebs are evaluated", "rel_id", e.Id)

			e.TalebState = "active"
			err = s.repo.UpdateTalebRealization(ctx, e)
			if err != nil {
				return err
			}
			slog.Info("rel state updated", "new_state", e.TalebState)
		}
	}

	if len(r) < 1 {
		err = s.repo.AddTaleb(ctx, user, cihazId)
		if err != nil {
			return err
		}
		slog.Info("rel bulunamadı yeni taleb eklendi")
	}

	return err
}
func (s CihazService) GetAdminEkran(ctx context.Context) ([]repo.CihazWithTalebCount, error) {
	recs, err := s.repo.GetCihazWithTalebCount(ctx)
	if err != nil {
		return nil, err
	}
	return recs, nil
}
