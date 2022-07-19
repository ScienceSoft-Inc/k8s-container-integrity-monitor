package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/k8s-container-integrity-monitor/internal/core/models"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KuberService struct {
	logger *logrus.Logger
}

// NewHashService creates a new struct HashService
func NewKuberService(logger *logrus.Logger) *KuberService {
	return &KuberService{
		logger: logger,
	}
}

func (ks *KuberService) GetDataFromK8sAPI() (*models.DataFromK8sAPI, error) {
	kuberData, err := ks.ConnectionToK8sAPI()
	if err != nil {
		ks.logger.Error("can't connection to K8sAPI: %s", err)
		return nil, err
	}
	deploymentData, err := ks.GetDataFromDeployment(kuberData)
	if err != nil {
		ks.logger.Error("error get data from kuberAPI %s", err)
		return nil, err
	}

	configData, err := ks.GetDataFromConfigMap(kuberData, deploymentData.LabelMainProcessName)
	if err != nil {
		ks.logger.Error("err while getting data from configMap K8sAPI %s", err)
		return &models.DataFromK8sAPI{}, err
	}

	dataFromK8sAPI := &models.DataFromK8sAPI{
		KuberData:      kuberData,
		DeploymentData: deploymentData,
		ConfigMapData:  configData,
	}

	return dataFromK8sAPI, nil
}

func (ks *KuberService) ConnectionToK8sAPI() (*models.KuberData, error) {
	ks.logger.Info("### 🌀 Attempting to use in cluster config")
	config, err := rest.InClusterConfig()
	if err != nil {
		ks.logger.Error(err)
		return nil, err
	}

	ks.logger.Info("### 💻 Connecting to Kubernetes API, using host: ", config.Host)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		ks.logger.Error(err)
		return nil, err
	}

	namespaceBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		ks.logger.Error(err)
		return nil, err
	}
	namespace := string(namespaceBytes)

	podName := os.Getenv("POD_NAME")

	targetName := func(podName string) string {
		elements := strings.Split(podName, "-")
		newElements := elements[:len(elements)-2]
		return strings.Join(newElements, "-")
	}(podName)
	if targetName == "" {
		ks.logger.Fatalln("### 💥 Env var DEPLOYMENT_NAME was not set")
	}
	targetType := os.Getenv("DEPLOYMENT_TYPE")
	kuberData := &models.KuberData{
		Clientset:  clientset,
		Namespace:  namespace,
		TargetName: targetName,
		TargetType: targetType,
	}
	return kuberData, nil
}
func (ks *KuberService) GetDataFromConfigMap(kuberData *models.KuberData, label string) (*models.ConfigMapData, error) {
	cm, err := kuberData.Clientset.CoreV1().ConfigMaps(kuberData.Namespace).Get(context.Background(), "hasher-config", metav1.GetOptions{})
	if err != nil {
		ks.logger.Error("err while getting data from configMap kuberAPI ", err)
		return nil, err
	}

	var configMapData models.ConfigMapData
	valuesEnv := make(map[string]string)
	for key, value := range cm.Data {
		if key == label {
			envs := strings.Split(strings.TrimSpace(value), "\n")
			for _, subStr := range envs {
				valuesEnvs := strings.Split(strings.TrimSpace(subStr), "=")
				valuesEnv[valuesEnvs[0]] = valuesEnvs[1]
			}
		}
	}

	if value, ok := valuesEnv["PID_NAME"]; ok {
		configMapData.ProcName = value
	}
	if value, ok := valuesEnv["MOUNT_PATH"]; ok {
		configMapData.MountPath = value
	}

	return &configMapData, err
}

func (ks *KuberService) GetDataFromDeployment(kuberData *models.KuberData) (*models.DeploymentData, error) {
	allDeploymentData, err := kuberData.Clientset.AppsV1().Deployments(kuberData.Namespace).Get(context.Background(), kuberData.TargetName, metav1.GetOptions{})
	if err != nil {
		ks.logger.Error("err while getting data from kuberAPI ", err)
		return nil, err
	}

	deploymentData := models.DeploymentData{}
	deploymentData.NameDeployment = kuberData.TargetName
	deploymentData.Timestamp = fmt.Sprintf("%v", allDeploymentData.CreationTimestamp)
	deploymentData.NamePod = os.Getenv("POD_NAME")

	for _, v := range allDeploymentData.Spec.Template.Spec.Containers {
		deploymentData.Image = v.Image
	}

	for label, value := range allDeploymentData.Spec.Template.Labels {
		if label == "hasher-webhook-process-name" {
			deploymentData.LabelMainProcessName = value
		}
	}

	return &deploymentData, nil
}

func (ks *KuberService) RolloutDeployment(kuberData *models.KuberData) error {
	patchData := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))
	_, err := kuberData.Clientset.AppsV1().Deployments(kuberData.Namespace).Patch(context.Background(), kuberData.TargetName, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	if err != nil {
		ks.logger.Printf("### 👎 Warning: Failed to patch %v, restart failed: %v", kuberData.TargetType, err)
		return err
	} else {
		ks.logger.Printf("### ✅ Target %v, named %v was restarted!", kuberData.TargetType, kuberData.TargetName)
	}
	return nil
}
