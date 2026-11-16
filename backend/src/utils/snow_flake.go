package utils

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	SF                     *SnowFlake
	ErrClockMovedBackwards = errors.New("clock moved backwards, refusing to generate id")
	ErrMachineIDOverflow   = errors.New("machineID must be between 0 and 1023")
	ErrSequenceOverflow    = errors.New("sequence must be between 0 and 4095")
)

const (
	// 2026-01-01 00:00:00  UTC milli
	epoch int64 = 1767196800000
	mask  int64 = 0x7BCBCBCBCBCBCBCA
)

type SnowFlakeInfo struct {
	Timestamp time.Time // 到毫秒
	MachineID int64     // 0-1023
	Sequence  int64     // 0-4095
	ID        int64     // 原始id
}

type SnowFlake struct {
	machineID int64
	sequence  int64
	lastStamp int64
	mutex     sync.Mutex
}

func NewSnowFlake(machineID int64) (*SnowFlake, error) {
	if machineID < 0 || machineID > 1023 {
		return nil, ErrMachineIDOverflow
	}
	return &SnowFlake{
		machineID: machineID,
		sequence:  0,
		lastStamp: -1,
	}, nil
}

func (sf *SnowFlake) NextID() (int64, error) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	timestamp := time.Now().UnixMilli()

	// 时钟回拨
	if timestamp < sf.lastStamp {
		return 0, ErrClockMovedBackwards
	}

	if timestamp == sf.lastStamp {
		sf.sequence = (sf.sequence + 1) & 0xFFF
		if sf.sequence == 0 {
			for timestamp <= sf.lastStamp {
				timestamp = time.Now().UnixMilli() // 下一毫秒
			}
		}
	} else {
		sf.sequence = 0
	}
	sf.lastStamp = timestamp

	return ((timestamp - epoch) << 22) | (sf.machineID << 12) | sf.sequence, nil
}

func (sf *SnowFlake) Encode(v int64) string {
	obfuscated := v ^ mask

	return strconv.FormatInt(obfuscated, 16)
}

func (sf *SnowFlake) Decode(v string) (int64, error) {
	obfuscated, err := strconv.ParseInt(v, 16, 64)
	if err != nil {
		return 0, err
	}

	return obfuscated ^ mask, nil
}

func (sf *SnowFlake) ParseInfo(id int64) *SnowFlakeInfo {
	timestamp := (id >> 22) + epoch
	machineID := (id >> 12) & 0x3FF
	sequence := id & 0xFFF

	return &SnowFlakeInfo{
		Timestamp: time.UnixMilli(timestamp),
		MachineID: machineID,
		Sequence:  sequence,
		ID:        id,
	}
}

func (sf *SnowFlake) Rebuild(info *SnowFlakeInfo) (int64, error) {
	timestamp := info.Timestamp.UnixMilli() - epoch
	if timestamp < 0 {
		return 0, fmt.Errorf("timestamp is before epoch")
	}

	if info.MachineID < 0 || info.MachineID > 1023 {
		return 0, ErrMachineIDOverflow
	}

	if info.Sequence < 0 || info.Sequence > 4095 {
		return 0, ErrSequenceOverflow
	}

	id := (timestamp << 22) | (info.MachineID << 12) | info.Sequence

	return id, nil
}
