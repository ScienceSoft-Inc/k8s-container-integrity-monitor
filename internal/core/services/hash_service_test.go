package services

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/k8s-container-integrity-monitor/internal/core/models"
	mock_ports "github.com/k8s-container-integrity-monitor/internal/core/ports/mocks"
	"github.com/k8s-container-integrity-monitor/pkg/api"
	"github.com/k8s-container-integrity-monitor/pkg/hasher"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewHashService(t *testing.T) {
	alg := "SHA256"
	logger := logrus.New()
	expected := HashService{}

	c := gomock.NewController(t)
	defer c.Finish()

	repo := mock_ports.NewMockIHashRepository(c)
	h, err := hasher.NewHashSum(alg)
	if err != nil {
		require.Error(t, err)
	}
	hashService := HashService{
		hashRepository: repo,
		hasher:         h,
		alg:            alg,
		logger:         logger,
	}
	assert.NotEqual(t, expected, hashService, "they should not be equal")
}

func TestCreateHash(t *testing.T) {
	testTable := []struct {
		name         string
		alg          string
		path         string
		mockBehavior func(s *mock_ports.MockIHashService, path string)
		expected     api.HashData
	}{
		{
			name: "exist path",
			alg:  "SHA256",
			path: "../h/h1/test.txt",
			mockBehavior: func(s *mock_ports.MockIHashService, path string) {
				s.EXPECT().CreateHash(path).Return(api.HashData{
					Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					FileName:     "test.txt",
					FullFilePath: "../h/h1/test.txt",
					Algorithm:    "SHA256",
				})

			},
			expected: api.HashData{
				Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				FileName:     "test.txt",
				FullFilePath: "../h/h1/test.txt",
				Algorithm:    "SHA256",
			},
		},
		{
			name: "not exist path",
			alg:  "SHA256",
			path: "/test.txx",
			mockBehavior: func(s *mock_ports.MockIHashService, path string) {
				s.EXPECT().CreateHash(path).Return(api.HashData{})
			},
			expected: api.HashData{
				Hash:         "",
				FileName:     "",
				FullFilePath: "",
				Algorithm:    "",
			},
		},
		{
			name: "error in a hash sum",
			alg:  "SHA256",
			path: "/test.txx",
			mockBehavior: func(s *mock_ports.MockIHashService, path string) {
				s.EXPECT().CreateHash(path).Return(api.HashData{})
			},
			expected: api.HashData{},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_ports.NewMockIHashService(c)
			testCase.mockBehavior(service, testCase.path)

			file, err := os.Open(testCase.path)
			if err != nil {
				require.Error(t, err)
			}
			defer file.Close()

			result := service.CreateHash(testCase.path)

			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestSaveHashData(t *testing.T) {
	type mockBehavior func(r *mock_ports.MockIHashRepository, ctx context.Context, allHashData []api.HashData, deploymentData models.DeploymentData)
	testTable := []struct {
		name           string
		alg            string
		allHashData    []api.HashData
		deploymentData models.DeploymentData
		mockBehavior   mockBehavior
		expected       error
	}{
		{
			name: "exist path",
			alg:  "SHA256",
			allHashData: []api.HashData{{
				Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				FileName:     "new",
				FullFilePath: "../h/h1/new",
				Algorithm:    "SHA256",
			}},
			deploymentData: models.DeploymentData{
				NameDeployment: "nginx",
				Image:          "nginx:latest",
				NamePod:        "nginx-deploy",
				Timestamp:      "01.01.2022 00:00",
			},
			mockBehavior: func(r *mock_ports.MockIHashRepository, ctx context.Context, allHashData []api.HashData, deploymentData models.DeploymentData) {
				r.EXPECT().SaveHashData(ctx, allHashData, deploymentData).Return(nil)
			},
			expected: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			l := logrus.New()
			ctx := context.Background()
			repo := mock_ports.NewMockIHashRepository(c)
			h, err := hasher.NewHashSum(testCase.alg)
			if err != nil {
				assert.Error(t, err)
			}

			hashService := HashService{
				hashRepository: repo,
				hasher:         h,
				alg:            testCase.alg,
				logger:         l,
			}
			testCase.mockBehavior(repo, ctx, testCase.allHashData, testCase.deploymentData)

			err = hashService.hashRepository.SaveHashData(ctx, testCase.allHashData, testCase.deploymentData)
			if testCase.expected != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetHashData(t *testing.T) {
	type mockBehavior func(r *mock_ports.MockIHashRepository, ctx context.Context, dirFiles, alg string, deploymentData models.DeploymentData)
	testTable := []struct {
		name           string
		alg            string
		dirFiles       string
		expected       []models.HashDataFromDB
		deploymentData models.DeploymentData
		mockBehavior   mockBehavior
		expectedErr    bool
	}{
		{
			name:     "exist path",
			alg:      "SHA256",
			dirFiles: "../h/h1/new",
			expected: []models.HashDataFromDB{{
				ID:           1,
				Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				FileName:     "new",
				FullFilePath: "../h/h1/new",
				Algorithm:    "SHA256",
			}},
			deploymentData: models.DeploymentData{
				NameDeployment: "nginx",
				Image:          "nginx:latest",
				NamePod:        "nginx-deploy",
				Timestamp:      "01.01.2022 00:00",
			},
			mockBehavior: func(r *mock_ports.MockIHashRepository, ctx context.Context, dirFiles, alg string, deploymentData models.DeploymentData) {
				r.EXPECT().GetHashData(ctx, dirFiles, alg, deploymentData).Return([]models.HashDataFromDB{
					{
						ID:           1,
						Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						FileName:     "new",
						FullFilePath: "../h/h1/new",
						Algorithm:    "SHA256",
					},
				}, nil)
			},
			expectedErr: false,
		},
		{
			name:     "not exist path",
			alg:      "SHA256",
			dirFiles: "../h/h1/new",
			expected: []models.HashDataFromDB{{
				ID:           1,
				Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				FileName:     "new",
				FullFilePath: "../h/h1/new",
				Algorithm:    "SHA256",
			}},
			deploymentData: models.DeploymentData{
				NameDeployment: "nginx",
				Image:          "nginx:latest",
				NamePod:        "nginx-deploy",
				Timestamp:      "01.01.2022 00:00",
			},
			mockBehavior: func(r *mock_ports.MockIHashRepository, ctx context.Context, dirFiles, alg string, deploymentData models.DeploymentData) {
				r.EXPECT().GetHashData(ctx, dirFiles, alg, deploymentData).Return([]models.HashDataFromDB{}, errors.New("hash service didn't get data"))
			},
			expectedErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			l := logrus.New()
			ctx := context.Background()
			repo := mock_ports.NewMockIHashRepository(c)
			h, err := hasher.NewHashSum(testCase.alg)
			if err != nil {
				assert.Error(t, err)
			}

			hashService := HashService{
				hashRepository: repo,
				hasher:         h,
				alg:            testCase.alg,
				logger:         l,
			}
			testCase.mockBehavior(repo, ctx, testCase.dirFiles, testCase.alg, testCase.deploymentData)

			data, err := hashService.hashRepository.GetHashData(ctx, testCase.dirFiles, testCase.alg, testCase.deploymentData)

			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, data)
			}
		})
	}
}
