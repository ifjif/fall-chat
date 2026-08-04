package utils

import (
	"testing"
)

func TestSnowFlake(t *testing.T) {
	sf, err := NewSnowFlake(1)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	id, err := sf.NextID()
	if err != nil {
		t.Errorf("Err: %v", err)
	}

	t.Logf("id: %d\n", id)
	obfuscaked := sf.Encode(id)
	t.Logf("混淆后: %s\n", obfuscaked)
	real, err := sf.Decode(obfuscaked)
	if err != nil {
		t.Errorf("Err: %v", err)
	}
	t.Logf("解码后： %d\n", real)

	if id != real {
		t.Errorf("Err: %s 解码后不是 %d, got = %d\n", obfuscaked, id, real)
	}

	info := sf.ParseInfo(id)

	t.Logf("id: %d\n", info.ID)
	t.Logf("time: %s\n", info.Timestamp)
	t.Logf("machine id: %d\n", info.MachineID)
	t.Logf("sequence: %d\n", info.Sequence)

	rid, err := sf.Rebuild(info)
	if err != nil {
		t.Errorf("Err: %v\n", err)
	}

	t.Logf("根据info重建后的id: %d", rid)
	if id != rid {
		t.Errorf("Err: 重建后不是 %d, got = %d\n", id, real)
	}

	// info.Timestamp = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	// info.MachineID = -1
	// info.MachineID = 1024
	// info.Sequence = -1
	info.Sequence = 4096
	_, err = sf.Rebuild(info)
	if err != nil {
		t.Errorf("Err: %v\n", err)
	}
}
