package services

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"testing"

	"github.com/golang/mock/gomock"
	mock_ports "github.com/k8s-container-integrity-monitor/internal/core/ports/mocks"
	"github.com/k8s-container-integrity-monitor/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestStartGetHashData(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIAppService, ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal)
	ctx := context.Background()
	var jobs chan string
	var results chan api.HashData

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	testTable := []struct {
		name         string
		flagName     string
		expectedErr  bool
		mockBehavior mockBehavior
	}{
		{
			name:        "valid data",
			flagName:    "d",
			expectedErr: false,
			mockBehavior: func(s *mock_ports.MockIAppService, ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) {
				s.EXPECT().StartGetHashData(ctx, flagName, jobs, results, sig).Return(nil)
			},
		},
		{
			name:        "error while saving to database",
			flagName:    "d",
			expectedErr: true,
			mockBehavior: func(s *mock_ports.MockIAppService, ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) {
				s.EXPECT().StartGetHashData(ctx, flagName, jobs, results, sig).Return(errors.New("error save hash data to database"))
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			s := mock_ports.NewMockIAppService(c)

			testCase.mockBehavior(s, ctx, testCase.flagName, jobs, results, sig)

			err := s.Start(ctx, testCase.flagName, jobs, results, sig)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStartCheckHashData(t *testing.T) {
	type mockBehavior func(s *mock_ports.MockIAppService, ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal)
	ctx := context.Background()
	var jobs chan string
	var results chan api.HashData

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	testTable := []struct {
		name         string
		flagName     string
		expectedErr  bool
		mockBehavior mockBehavior
	}{
		{
			name:        "valid data",
			flagName:    "d",
			expectedErr: false,
			mockBehavior: func(s *mock_ports.MockIAppService, ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) {
				s.EXPECT().StartCheckHashData(ctx, flagName, jobs, results, sig).Return(nil)
			},
		},
		{
			name:        "error while getting to database",
			flagName:    "d",
			expectedErr: true,
			mockBehavior: func(s *mock_ports.MockIAppService, ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) {
				s.EXPECT().StartCheckHashData(ctx, flagName, jobs, results, sig).Return(errors.New("error getting hash data from database"))
			},
		},
		{
			name:        "error while changing data to database",
			flagName:    "d",
			expectedErr: true,
			mockBehavior: func(s *mock_ports.MockIAppService, ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) {
				s.EXPECT().StartCheckHashData(ctx, flagName, jobs, results, sig).Return(errors.New("error match data currently and data from database"))
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			s := mock_ports.NewMockIAppService(c)

			testCase.mockBehavior(s, ctx, testCase.flagName, jobs, results, sig)

			err := s.Check(ctx, testCase.flagName, jobs, results, sig)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
