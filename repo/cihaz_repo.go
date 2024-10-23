package repo

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
)

type CihazRepo struct {
	*sql.DB
}

type Cihaz struct {
	Id    int    `json:"id"`
	Marka string `json:"marka"`
	Model string `json:"model"`
	Fiyat int    `json:"fiyat"`
	Stock int    `json:"stock"`
}
type CihazJoinTalebJoinRealiziton struct {
	Cihaz
	TalebCount         int            `json:"taleb_count"`
	RequiredTalebCount sql.NullInt32  `json:"required_taleb_count"`
	TalebState         sql.NullString `json:"taleb_state"`
}

type TalebRealization struct {
	Id                 int    `json:"id"`
	CihazId            int    `json:"cihaz_id"`
	TalebState         string `json:"taleb_state"`
	RequiredTalebCount int    `json:"required_taleb_count"`
}

type CihazWithTalebCount struct {
	Cihaz
	TalebCount sql.NullInt32 `json:"taleb_count"`
}

func NewRepo(dbStr string) (*CihazRepo, error) {
	db, err := sql.Open("postgres", dbStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &CihazRepo{db}, nil
}

func (r CihazRepo) GetAll(ctx context.Context) ([]Cihaz, error) {
	q := `SELECT * FROM personel_cihazlar`
	rows, err := r.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var c Cihaz
	var cihazlar []Cihaz

	for rows.Next() {
		err = rows.Scan(&c.Id, &c.Marka, &c.Model, &c.Fiyat, &c.Stock)
		if err != nil {
			slog.Warn("error scaning row", "err", err)
		}
		cihazlar = append(cihazlar, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cihazlar, nil
}

func (r CihazRepo) GetCihazJoinTalebJoinRealization(ctx context.Context) ([]CihazJoinTalebJoinRealiziton, error) {
	q := `select pc.id, pc.marka, pc.model, pc.fiyat,
        pc.stock, count(pt.id) as taleb_count, ptr.gerekli_taleb_count,
       ptr.taleb_state from personel_cihazlar pc left join
        (select id,  cihaz_id from personel_taleb where evaulated = 0)pt
        on(pc.id = pt.cihaz_id) left join (select * from personeal_indirim_taleb_realization
        where taleb_state IN ('pending', 'active')) ptr on(ptr.cihaz_id = pc.id)
        group by pc.id, ptr.cihaz_id, ptr.taleb_state, ptr.gerekli_taleb_count order by pc.id`

	rows, err := r.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var record CihazJoinTalebJoinRealiziton
	var list []CihazJoinTalebJoinRealiziton

	for rows.Next() {
		err = rows.Scan(&record.Id, &record.Marka, &record.Model,
			&record.Fiyat, &record.Stock, &record.TalebCount,
			&record.RequiredTalebCount, &record.TalebState)
		if err != nil {
			slog.Warn("error scaning row", "err", err)
		}
		list = append(list, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r CihazRepo) AddTaleb(ctx context.Context, user string, cihazId int) error {
	q := `INSERT INTO personel_taleb(cihaz_id, evaulated, user_name) values($1, $2, $3)`
	_, err := r.ExecContext(ctx, q, cihazId, 0, user)
	return err
}

func (r CihazRepo) GetRealizationByCihazIdAndState(ctx context.Context, cihazId int, states ...string) ([]TalebRealization, error) {
	q := `SELECT * from personeal_indirim_taleb_realization where cihaz_id = $1 and taleb_state IN ('pending', 'active')`

	rows, err := r.DB.QueryContext(ctx, q, cihazId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var record TalebRealization
	var list []TalebRealization
	for rows.Next() {
		err = rows.Scan(&record.Id, &record.CihazId, &record.TalebState, &record.RequiredTalebCount)
		if err != nil {
			slog.Warn("error scaning row", "err", err)
		}
		list = append(list, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r CihazRepo) GetNotEvaulatedTalebCountByCihazId(ctx context.Context, cihazId int) (int, error) {
	q := `select count(*) from personel_taleb where evaulated = 0 and cihaz_id = $1 group by cihaz_id;`

	var c int
	row := r.DB.QueryRowContext(ctx, q, cihazId)
	err := row.Scan(&c)
	if err != nil {
		return -1, err
	}

	if err = row.Err(); err != nil {
		return -1, err
	}

	return c, nil
}
func (r CihazRepo) EvaluateTalebsWithCihazId(ctx context.Context, realizationId, cihazId int) error {
	q := `UPDATE personel_taleb set evaulated = 1, taleb_realization_id = $1 where cihaz_id = $2 and evaulated = 0`

	_, err := r.ExecContext(ctx, q, realizationId, cihazId)
	return err
}
func (r CihazRepo) UpdateTalebRealization(ctx context.Context, rel TalebRealization) error {
	q := `UPDATE personeal_indirim_taleb_realization
    set cihaz_id = $1, taleb_state = $2, gerekli_taleb_count = $3
    where id = $4`

	_, err := r.ExecContext(ctx, q, rel.CihazId, rel.TalebState, rel.RequiredTalebCount, rel.Id)
	return err
}
func (r CihazRepo) GetCihazWithTalebCount(ctx context.Context) ([]CihazWithTalebCount, error) {
	q := `select pc.id, pc.marka, pc.model, pc.fiyat, pc.stock, count(pt.id) from personel_cihazlar pc
    left join (select id, cihaz_id from personel_taleb where evaulated = 0) pt
    on(pc.id = pt.cihaz_id) group by pc.id order by pc.id`

	rows, err := r.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rec CihazWithTalebCount
	var list []CihazWithTalebCount
	for rows.Next() {
		err = rows.Scan(&rec.Id, &rec.Marka, &rec.Model, &rec.Fiyat, &rec.Stock, &rec.TalebCount)
		if err != nil {
			slog.Error("couldnt scan", "err", err)
		}
		list = append(list, rec)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}
