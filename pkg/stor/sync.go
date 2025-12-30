// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package stor

import (
	"time"

	"gorm.io/gorm"
)

// SyncInfo stores the last synchronization timestamp
type SyncInfo struct {
	gorm.Model
	LastSync time.Time `json:"last_sync"`
}

// GetLastSyncTime retrieves the last synchronization time
func (s *Store) GetLastSyncTime() (time.Time, error) {
	var syncInfo SyncInfo
	err := s.db.Order("last_sync desc").First(&syncInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No previous sync, return zero time
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	return syncInfo.LastSync, nil
}

// SaveSyncTime saves a new synchronization time
func (s *Store) SaveSyncTime(syncTime time.Time) error {
	syncInfo := SyncInfo{
		LastSync: syncTime,
	}
	return s.db.Create(&syncInfo).Error
}
