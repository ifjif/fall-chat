package utils

import "zzc/fall-chat/src/model"

func InjectID4Model(sf *SnowFlake, m *model.Model) error {
	id, err := sf.NextID()
	if err != nil {
		return err
	}
	m.ID = uint64(id)
	m.PublicID = sf.Encode(id)

	return nil
}
