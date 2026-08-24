// Copyright 2015 Matthew Holt
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package certmagic

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestCleanStorageContextCanceled(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "certmagic_clean_test*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storage := &FileStorage{Path: tmpDir}
	logger := zap.NewNop()

	opts := CleanStorageOptions{
		Logger:                 logger,
		Interval:               1 * time.Hour,
		OCSPStaples:            true,
		ExpiredCerts:           true,
		ExpiredCertGracePeriod: 24 * time.Hour,
	}

	// 1. Test with already-canceled context
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err = CleanStorage(canceledCtx, storage, opts)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled error, got: %v", err)
	}

	// 2. Test normal clean storage runs without error
	err = CleanStorage(context.Background(), storage, opts)
	if err != nil {
		t.Fatalf("CleanStorage failed with valid context: %v", err)
	}
}
