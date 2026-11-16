package types

import "time"

// layout rfc3339
type RedisTime struct {
	time.Time
}

func (rt RedisTime) MarshalBinary() (data []byte, err error) {
	if rt.IsZero() {
		return nil, nil
	}

	return []byte(rt.Format(time.RFC3339)), nil
}

func (rt *RedisTime) UnmarshalBinary(b []byte) error {
	s := string(b)
	if s == "" {
		rt.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	rt.Time = t
	return nil
}
