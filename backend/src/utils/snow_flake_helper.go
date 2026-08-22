package utils

import (
	"errors"
	"time"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

func ShutdownSnowFlake(db *gorm.DB) {
	if stopHearBeat != nil {
		close(stopHearBeat)
	}

	releaseMachineID(db, currentMID)
}

func initSnowFlake(db *gorm.DB) error {
	var initErr error

	sfOnce.Do(func() {
		var machineID int64
		var err error

		// 1.todo 从环境变量获取

		// 2.todo 从k8s中 主机名获取

		// 3.todo 从redis中获取

		// 4.从数据库中获取
		machineID, err = getMachineIDFromeDB(db)
		if err == nil {
			goto CREATE_INSTANCE
		}

		initErr = errors.New("分配machine id 失败")
		return

	CREATE_INSTANCE:
		SF, err = NewSnowFlake(machineID)
		if err != nil {
			initErr = err
			return
		}

		currentMID = machineID
		stopHearBeat = make(chan struct{})
		go startHeartBeat(db, machineID, stopHearBeat)
	})

	return initErr
}

// 从数据库中获取一个可用的 machine ID
func getMachineIDFromeDB(db *gorm.DB) (int64, error) {
	var machineID model.MachineIdPool

	// 捡漏
	// status = 0 or status = 1 and last_heartbeat < Now()-INTERVAL 1 MINUTE limit 1
	result := db.Where("status = ?", 0).
		Or("status = ? and last_heartbeat < Now() - INTERVAL 1 MINUTE", 1).
		Take(&machineID)

	if result.Error == nil {
		LogInfo("存在漏的machine id: %d, 进行复用", machineID.ID)
		// status = 1, ip = ip
		newMid := model.MachineIdPool{
			ID:            machineID.ID,
			IP:            "unknown",
			Status:        1,
			LastHeartbeat: time.Now(),
		}
		db.Model(&newMid).Updates(newMid)

		return machineID.ID, nil
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		LogInfo("不存在漏的machine id, 进行新增")
	} else {
		return 0, result.Error
	}

	machineID.Status = 1
	machineID.IP = "unknown"
	machineID.LastHeartbeat = time.Now()
	result = db.Create(&machineID)
	if result.Error != nil {
		return 0, result.Error
	}

	return machineID.ID, nil
}

// 30 更新一次 机器id拥有状态
func startHeartBeat(db *gorm.DB, mid int64, stopHearBeat chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result := db.Model(&model.MachineIdPool{}).
				Where("id = ?", mid).
				Update("last_heartbeat", time.Now())
			if result.Error != nil {
				LogError("更新机器心跳失败: %v", result.Error)
			}
		case <-stopHearBeat:
			return
		}
	}
}

func releaseMachineID(db *gorm.DB, machineID int64) {
	db.Model(&model.MachineIdPool{}).
		Update("status", 0).
		Where("id = ?", machineID)
}
