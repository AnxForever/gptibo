package store

import "context"

// PruneAuth deletes only expired authentication records, in bounded batches.
// Financial and audit history never goes through this retention job.
func (s *Store) PruneAuth(ctx context.Context) error {
	for _, sql := range []string{
		"DELETE FROM auth_challenges WHERE id IN (SELECT id FROM auth_challenges WHERE expires_at<now()-interval '1 hour' ORDER BY expires_at LIMIT 1000)",
		"DELETE FROM auth_sessions WHERE token_hash IN (SELECT token_hash FROM auth_sessions WHERE expires_at<now() ORDER BY expires_at LIMIT 1000)",
		"DELETE FROM auth_limits WHERE reset_at<now() AND key IN (SELECT key FROM auth_limits WHERE reset_at<now() ORDER BY reset_at LIMIT 1000)",
	} {
		if _, e := s.DB.Exec(ctx, sql); e != nil {
			return e
		}
	}
	return nil
}

// PruneAnalytics bounds the two anonymously writable tables. Both accept public
// inserts, so without retention one client can grow them without limit; the
// daily aggregates keep their shape, only old days disappear.
func (s *Store) PruneAnalytics(ctx context.Context, visitorDays, interestDays int) error {
	for _, job := range []struct {
		sql  string
		days int
	}{
		{`DELETE FROM visitor_days WHERE (day,shop_id,visitor_hash) IN (
			SELECT day,shop_id,visitor_hash FROM visitor_days
			WHERE day < to_char((now() - make_interval(days => $1::int)) at time zone 'UTC','YYYY-MM-DD')
			LIMIT 2000)`, visitorDays},
		{`DELETE FROM launch_interest WHERE id IN (
			SELECT id FROM launch_interest
			WHERE created_at < to_char((now() - make_interval(days => $1::int)) at time zone 'UTC','YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
			LIMIT 2000)`, interestDays},
	} {
		if _, e := s.DB.Exec(ctx, job.sql, job.days); e != nil {
			return e
		}
	}
	return nil
}
