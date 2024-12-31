/*
 * Copyright (c) 2024. Devtron Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package devtronApps

import (
	"context"
	"errors"
	"fmt"
	pubsub "github.com/devtron-labs/common-lib/pubsub-lib"
	util5 "github.com/devtron-labs/common-lib/utils/k8s"
	bean3 "github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/api/bean/gitOps"
	"github.com/devtron-labs/devtron/api/helm-app/gRPC"
	client2 "github.com/devtron-labs/devtron/api/helm-app/service"
	"github.com/devtron-labs/devtron/client/argocdServer"
	bean7 "github.com/devtron-labs/devtron/client/argocdServer/bean"
	client "github.com/devtron-labs/devtron/client/events"
	gitSensorClient "github.com/devtron-labs/devtron/client/gitSensor"
	"github.com/devtron-labs/devtron/internal/middleware"
	"github.com/devtron-labs/devtron/internal/sql/models"
	repository3 "github.com/devtron-labs/devtron/internal/sql/repository"
	appRepository "github.com/devtron-labs/devtron/internal/sql/repository/app"
	"github.com/devtron-labs/devtron/internal/sql/repository/appWorkflow"
	"github.com/devtron-labs/devtron/internal/sql/repository/chartConfig"
	repository4 "github.com/devtron-labs/devtron/internal/sql/repository/dockerRegistry"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig/bean/timelineStatus"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig/bean/workflow/cdWorkflow"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/app"
	bean4 "github.com/devtron-labs/devtron/pkg/app/bean"
	"github.com/devtron-labs/devtron/pkg/app/status"
	statusBean "github.com/devtron-labs/devtron/pkg/app/status/bean"
	"github.com/devtron-labs/devtron/pkg/attributes"
	"github.com/devtron-labs/devtron/pkg/auth/user"
	bean2 "github.com/devtron-labs/devtron/pkg/bean"
	"github.com/devtron-labs/devtron/pkg/build/git/gitMaterial/read"
	pipeline2 "github.com/devtron-labs/devtron/pkg/build/pipeline"
	chartRepoRepository "github.com/devtron-labs/devtron/pkg/chartRepo/repository"
	repository2 "github.com/devtron-labs/devtron/pkg/cluster/environment/repository"
	repository5 "github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/devtron-labs/devtron/pkg/deployment/common"
	bean9 "github.com/devtron-labs/devtron/pkg/deployment/common/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/gitOps/config"
	"github.com/devtron-labs/devtron/pkg/deployment/gitOps/git"
	"github.com/devtron-labs/devtron/pkg/deployment/manifest"
	bean10 "github.com/devtron-labs/devtron/pkg/deployment/manifest/deploymentTemplate/bean"
	bean5 "github.com/devtron-labs/devtron/pkg/deployment/manifest/deploymentTemplate/chartRef/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/manifest/publish"
	"github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps/adapter"
	"github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps/helper"
	"github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps/userDeploymentRequest/service"
	clientErrors "github.com/devtron-labs/devtron/pkg/errors"
	"github.com/devtron-labs/devtron/pkg/eventProcessor/out"
	"github.com/devtron-labs/devtron/pkg/imageDigestPolicy"
	k8s2 "github.com/devtron-labs/devtron/pkg/k8s"
	"github.com/devtron-labs/devtron/pkg/pipeline"
	bean8 "github.com/devtron-labs/devtron/pkg/pipeline/bean"
	"github.com/devtron-labs/devtron/pkg/pipeline/history"
	"github.com/devtron-labs/devtron/pkg/pipeline/repository"
	"github.com/devtron-labs/devtron/pkg/pipeline/types"
	"github.com/devtron-labs/devtron/pkg/plugin"
	security2 "github.com/devtron-labs/devtron/pkg/policyGovernance/security/imageScanning"
	read2 "github.com/devtron-labs/devtron/pkg/policyGovernance/security/imageScanning/read"
	repository6 "github.com/devtron-labs/devtron/pkg/policyGovernance/security/imageScanning/repository"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/devtron-labs/devtron/pkg/variables"
	"github.com/devtron-labs/devtron/pkg/workflow/cd"
	globalUtil "github.com/devtron-labs/devtron/util"
	"github.com/devtron-labs/devtron/util/argo"
	util2 "github.com/devtron-labs/devtron/util/event"
	"github.com/devtron-labs/devtron/util/rbac"
	"github.com/go-pg/pg"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	status2 "google.golang.org/grpc/status"
	"helm.sh/helm/v3/pkg/chart"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TriggerService interface {
	TriggerPostStage(request bean.TriggerRequest) (*bean4.ManifestPushTemplate, error)
	TriggerPreStage(request bean.TriggerRequest) (*bean4.ManifestPushTemplate, error)

	TriggerAutoCDOnPreStageSuccess(triggerContext bean.TriggerContext, cdPipelineId, ciArtifactId, workflowId int, triggerdBy int32) error

	TriggerStageForBulk(triggerRequest bean.TriggerRequest) error

	ManualCdTrigger(triggerContext bean.TriggerContext, overrideRequest *bean3.ValuesOverrideRequest) (int, string, *bean4.ManifestPushTemplate, error)
	TriggerAutomaticDeployment(request bean.TriggerRequest) error

	TriggerRelease(overrideRequest *bean3.ValuesOverrideRequest, envDeploymentConfig *bean9.DeploymentConfig, ctx context.Context, triggeredAt time.Time, triggeredBy int32) (releaseNo int, manifestPushTemplate *bean4.ManifestPushTemplate, err error)
}

type TriggerServiceImpl struct {
	logger                              *zap.SugaredLogger
	cdWorkflowCommonService             cd.CdWorkflowCommonService
	gitOpsManifestPushService           publish.GitOpsPushService
	gitOpsConfigReadService             config.GitOpsConfigReadService
	argoK8sClient                       argocdServer.ArgoK8sClient
	ACDConfig                           *argocdServer.ACDConfig
	argoClientWrapperService            argocdServer.ArgoClientWrapperService
	pipelineStatusTimelineService       status.PipelineStatusTimelineService
	chartTemplateService                util.ChartTemplateService
	eventFactory                        client.EventFactory
	eventClient                         client.EventClient
	globalEnvVariables                  *globalUtil.GlobalEnvVariables
	workflowEventPublishService         out.WorkflowEventPublishService
	manifestCreationService             manifest.ManifestCreationService
	deployedConfigurationHistoryService history.DeployedConfigurationHistoryService
	argoUserService                     argo.ArgoUserService
	pipelineStageService                pipeline.PipelineStageService
	globalPluginService                 plugin.GlobalPluginService
	customTagService                    pipeline.CustomTagService
	pluginInputVariableParser           pipeline.PluginInputVariableParser
	prePostCdScriptHistoryService       history.PrePostCdScriptHistoryService
	scopedVariableManager               variables.ScopedVariableCMCSManager
	cdWorkflowService                   pipeline.WorkflowService
	imageDigestPolicyService            imageDigestPolicy.ImageDigestPolicyService
	userService                         user.UserService
	gitSensorClient                     gitSensorClient.Client
	config                              *types.CdConfig
	helmAppService                      client2.HelmAppService
	imageScanService                    security2.ImageScanService
	enforcerUtil                        rbac.EnforcerUtil
	userDeploymentRequestService        service.UserDeploymentRequestService
	helmAppClient                       gRPC.HelmAppClient //TODO refactoring: use helm app service instead
	appRepository                       appRepository.AppRepository
	ciPipelineMaterialRepository        pipelineConfig.CiPipelineMaterialRepository
	imageScanHistoryReadService         read2.ImageScanHistoryReadService
	imageScanDeployInfoService          security2.ImageScanDeployInfoService
	imageScanDeployInfoReadService      read2.ImageScanDeployInfoReadService
	pipelineRepository                  pipelineConfig.PipelineRepository
	pipelineOverrideRepository          chartConfig.PipelineOverrideRepository
	manifestPushConfigRepository        repository.ManifestPushConfigRepository
	chartRepository                     chartRepoRepository.ChartRepository
	envRepository                       repository2.EnvironmentRepository
	cdWorkflowRepository                pipelineConfig.CdWorkflowRepository
	ciWorkflowRepository                pipelineConfig.CiWorkflowRepository
	ciArtifactRepository                repository3.CiArtifactRepository
	ciTemplateService                   pipeline2.CiTemplateReadService
	gitMaterialReadService              read.GitMaterialReadService
	appLabelRepository                  pipelineConfig.AppLabelRepository
	ciPipelineRepository                pipelineConfig.CiPipelineRepository
	appWorkflowRepository               appWorkflow.AppWorkflowRepository
	dockerArtifactStoreRepository       repository4.DockerArtifactStoreRepository
	K8sUtil                             *util5.K8sServiceImpl
	transactionUtilImpl                 *sql.TransactionUtilImpl
	deploymentConfigService             common.DeploymentConfigService
	deploymentServiceTypeConfig         *globalUtil.DeploymentServiceTypeConfig
	ciCdPipelineOrchestrator            pipeline.CiCdPipelineOrchestrator
	gitOperationService                 git.GitOperationService
	attributeService                    attributes.AttributesService
	clusterRepository                   repository5.ClusterRepository
}

func NewTriggerServiceImpl(logger *zap.SugaredLogger,
	cdWorkflowCommonService cd.CdWorkflowCommonService,
	gitOpsManifestPushService publish.GitOpsPushService,
	gitOpsConfigReadService config.GitOpsConfigReadService,
	argoK8sClient argocdServer.ArgoK8sClient,
	ACDConfig *argocdServer.ACDConfig,
	argoClientWrapperService argocdServer.ArgoClientWrapperService,
	pipelineStatusTimelineService status.PipelineStatusTimelineService,
	chartTemplateService util.ChartTemplateService,
	workflowEventPublishService out.WorkflowEventPublishService,
	manifestCreationService manifest.ManifestCreationService,
	deployedConfigurationHistoryService history.DeployedConfigurationHistoryService,
	argoUserService argo.ArgoUserService,
	pipelineStageService pipeline.PipelineStageService,
	globalPluginService plugin.GlobalPluginService,
	customTagService pipeline.CustomTagService,
	pluginInputVariableParser pipeline.PluginInputVariableParser,
	prePostCdScriptHistoryService history.PrePostCdScriptHistoryService,
	scopedVariableManager variables.ScopedVariableCMCSManager,
	cdWorkflowService pipeline.WorkflowService,
	imageDigestPolicyService imageDigestPolicy.ImageDigestPolicyService,
	userService user.UserService,
	gitSensorClient gitSensorClient.Client,
	helmAppService client2.HelmAppService,
	enforcerUtil rbac.EnforcerUtil,
	userDeploymentRequestService service.UserDeploymentRequestService,
	helmAppClient gRPC.HelmAppClient,
	eventFactory client.EventFactory,
	eventClient client.EventClient,
	envVariables *globalUtil.EnvironmentVariables,
	appRepository appRepository.AppRepository,
	ciPipelineMaterialRepository pipelineConfig.CiPipelineMaterialRepository,
	imageScanHistoryReadService read2.ImageScanHistoryReadService,
	imageScanDeployInfoReadService read2.ImageScanDeployInfoReadService,
	imageScanDeployInfoService security2.ImageScanDeployInfoService,
	pipelineRepository pipelineConfig.PipelineRepository,
	pipelineOverrideRepository chartConfig.PipelineOverrideRepository,
	manifestPushConfigRepository repository.ManifestPushConfigRepository,
	chartRepository chartRepoRepository.ChartRepository,
	envRepository repository2.EnvironmentRepository,
	cdWorkflowRepository pipelineConfig.CdWorkflowRepository,
	ciWorkflowRepository pipelineConfig.CiWorkflowRepository,
	ciArtifactRepository repository3.CiArtifactRepository,
	ciTemplateService pipeline2.CiTemplateReadService,
	gitMaterialReadService read.GitMaterialReadService,
	appLabelRepository pipelineConfig.AppLabelRepository,
	ciPipelineRepository pipelineConfig.CiPipelineRepository,
	appWorkflowRepository appWorkflow.AppWorkflowRepository,
	dockerArtifactStoreRepository repository4.DockerArtifactStoreRepository,
	imageScanService security2.ImageScanService,
	K8sUtil *util5.K8sServiceImpl,
	transactionUtilImpl *sql.TransactionUtilImpl,
	deploymentConfigService common.DeploymentConfigService,
	ciCdPipelineOrchestrator pipeline.CiCdPipelineOrchestrator,
	gitOperationService git.GitOperationService,
	attributeService attributes.AttributesService,
	clusterRepository repository5.ClusterRepository,
) (*TriggerServiceImpl, error) {
	impl := &TriggerServiceImpl{
		logger:                              logger,
		cdWorkflowCommonService:             cdWorkflowCommonService,
		gitOpsManifestPushService:           gitOpsManifestPushService,
		gitOpsConfigReadService:             gitOpsConfigReadService,
		argoK8sClient:                       argoK8sClient,
		ACDConfig:                           ACDConfig,
		argoClientWrapperService:            argoClientWrapperService,
		pipelineStatusTimelineService:       pipelineStatusTimelineService,
		chartTemplateService:                chartTemplateService,
		workflowEventPublishService:         workflowEventPublishService,
		manifestCreationService:             manifestCreationService,
		deployedConfigurationHistoryService: deployedConfigurationHistoryService,
		argoUserService:                     argoUserService,
		pipelineStageService:                pipelineStageService,
		globalPluginService:                 globalPluginService,
		customTagService:                    customTagService,
		pluginInputVariableParser:           pluginInputVariableParser,
		prePostCdScriptHistoryService:       prePostCdScriptHistoryService,
		scopedVariableManager:               scopedVariableManager,
		cdWorkflowService:                   cdWorkflowService,
		imageDigestPolicyService:            imageDigestPolicyService,
		userService:                         userService,
		gitSensorClient:                     gitSensorClient,
		helmAppService:                      helmAppService,
		enforcerUtil:                        enforcerUtil,
		eventFactory:                        eventFactory,
		eventClient:                         eventClient,

		globalEnvVariables:             envVariables.GlobalEnvVariables,
		userDeploymentRequestService:   userDeploymentRequestService,
		helmAppClient:                  helmAppClient,
		appRepository:                  appRepository,
		ciPipelineMaterialRepository:   ciPipelineMaterialRepository,
		imageScanHistoryReadService:    imageScanHistoryReadService,
		imageScanDeployInfoReadService: imageScanDeployInfoReadService,
		imageScanDeployInfoService:     imageScanDeployInfoService,
		pipelineRepository:             pipelineRepository,
		pipelineOverrideRepository:     pipelineOverrideRepository,
		manifestPushConfigRepository:   manifestPushConfigRepository,
		chartRepository:                chartRepository,
		envRepository:                  envRepository,
		cdWorkflowRepository:           cdWorkflowRepository,
		ciWorkflowRepository:           ciWorkflowRepository,
		ciArtifactRepository:           ciArtifactRepository,
		ciTemplateService:              ciTemplateService,
		gitMaterialReadService:         gitMaterialReadService,
		appLabelRepository:             appLabelRepository,
		ciPipelineRepository:           ciPipelineRepository,
		appWorkflowRepository:          appWorkflowRepository,
		dockerArtifactStoreRepository:  dockerArtifactStoreRepository,

		imageScanService: imageScanService,
		K8sUtil:          K8sUtil,

		transactionUtilImpl: transactionUtilImpl,

		deploymentConfigService:     deploymentConfigService,
		deploymentServiceTypeConfig: envVariables.DeploymentServiceTypeConfig,
		ciCdPipelineOrchestrator:    ciCdPipelineOrchestrator,
		gitOperationService:         gitOperationService,
		attributeService:            attributeService,

		clusterRepository: clusterRepository,
	}
	config, err := types.GetCdConfig()
	if err != nil {
		return nil, err
	}
	impl.config = config
	return impl, nil
}

func (impl *TriggerServiceImpl) TriggerStageForBulk(triggerRequest bean.TriggerRequest) error {

	preStage, err := impl.pipelineStageService.GetCdStageByCdPipelineIdAndStageType(triggerRequest.Pipeline.Id, repository.PIPELINE_STAGE_TYPE_PRE_CD, false)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error in fetching CD pipeline stage", "cdPipelineId", triggerRequest.Pipeline.Id, "stage ", repository.PIPELINE_STAGE_TYPE_PRE_CD, "err", err)
		return err
	}

	// handle corrupt data (https://github.com/devtron-labs/devtron/issues/3826)
	err, deleted := impl.deleteCorruptedPipelineStage(preStage, triggerRequest.TriggeredBy)
	if err != nil {
		impl.logger.Errorw("error in deleteCorruptedPipelineStage ", "cdPipelineId", triggerRequest.Pipeline.Id, "err", err, "preStage", preStage, "triggeredBy", triggerRequest.TriggeredBy)
		return err
	}

	triggerRequest.TriggerContext.Context = context.Background()
	if len(triggerRequest.Pipeline.PreStageConfig) > 0 || (preStage != nil && !deleted) {
		// pre stage exists
		impl.logger.Debugw("trigger pre stage for pipeline", "artifactId", triggerRequest.Artifact.Id, "pipelineId", triggerRequest.Pipeline.Id)
		triggerRequest.RefCdWorkflowRunnerId = 0
		err = impl.preStageHandlingForTriggerStageInBulk(&triggerRequest)
		if err != nil {
			return err
		}
		_, err = impl.TriggerPreStage(triggerRequest) // TODO handle error here
		return err
	} else {
		// trigger deployment
		impl.logger.Debugw("trigger cd for pipeline", "artifactId", triggerRequest.Artifact.Id, "pipelineId", triggerRequest.Pipeline.Id)
		err = impl.TriggerAutomaticDeployment(triggerRequest)
		return err
	}
}

func (impl *TriggerServiceImpl) getCdPipelineForManualCdTrigger(ctx context.Context, pipelineId int) (*pipelineConfig.Pipeline, error) {
	_, span := otel.Tracer("TriggerService").Start(ctx, "getCdPipelineForManualCdTrigger")
	defer span.End()
	cdPipeline, err := impl.pipelineRepository.FindById(pipelineId)
	if err != nil {
		impl.logger.Errorw("manual trigger request with invalid pipelineId, ManualCdTrigger", "pipelineId", pipelineId, "err", err)
		return nil, err
	}
	//checking if namespace exist or not
	clusterIdToNsMap := map[int]string{
		cdPipeline.Environment.ClusterId: cdPipeline.Environment.Namespace,
	}

	err = impl.helmAppService.CheckIfNsExistsForClusterIds(clusterIdToNsMap)
	if err != nil {
		impl.logger.Errorw("manual trigger request with invalid namespace, ManualCdTrigger", "pipelineId", pipelineId, "err", err)
		return nil, err
	}
	return cdPipeline, nil
}

func (impl *TriggerServiceImpl) validateDeploymentTriggerRequest(ctx context.Context, runner *pipelineConfig.CdWorkflowRunner,
	cdPipeline *pipelineConfig.Pipeline, imageDigest string, envDeploymentConfig *bean9.DeploymentConfig, triggeredBy int32) error {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.validateDeploymentTriggerRequest")
	defer span.End()
	// custom GitOps repo url validation --> Start
	err := impl.handleCustomGitOpsRepoValidation(runner, cdPipeline, envDeploymentConfig, triggeredBy)
	if err != nil {
		impl.logger.Errorw("custom GitOps repository validation error, TriggerStage", "err", err)
		return err
	}
	// custom GitOps repo url validation --> Ends

	// checking vulnerability for deploying image
	vulnerabilityCheckRequest := adapter.GetVulnerabilityCheckRequest(cdPipeline, imageDigest)
	isVulnerable, err := impl.imageScanService.GetArtifactVulnerabilityStatus(newCtx, vulnerabilityCheckRequest)
	if err != nil {
		impl.logger.Errorw("error in getting Artifact vulnerability status, ManualCdTrigger", "err", err)
		return err
	}

	if isVulnerable == true {
		// if image vulnerable, update timeline status and return
		if err = impl.cdWorkflowCommonService.MarkCurrentDeploymentFailed(runner, errors.New(cdWorkflow.FOUND_VULNERABILITY), triggeredBy); err != nil {
			impl.logger.Errorw("error while updating current runner status to failed, TriggerDeployment", "wfrId", runner.Id, "err", err)
		}
		return fmt.Errorf("found vulnerability for image digest %s", imageDigest)
	}
	return nil
}

// TODO: write a wrapper to handle auto and manual trigger
func (impl *TriggerServiceImpl) ManualCdTrigger(triggerContext bean.TriggerContext, overrideRequest *bean3.ValuesOverrideRequest) (int, string, *bean4.ManifestPushTemplate, error) {

	triggerContext.TriggerType = bean.Manual
	// setting triggeredAt variable to have consistent data for various audit log places in db for deployment time
	triggeredAt := time.Now()

	releaseId := 0
	ctx := triggerContext.Context
	cdPipeline, err := impl.getCdPipelineForManualCdTrigger(ctx, overrideRequest.PipelineId)
	if err != nil {
		if overrideRequest.WfrId != 0 {
			err2 := impl.cdWorkflowCommonService.MarkDeploymentFailedForRunnerId(overrideRequest.WfrId, err, overrideRequest.UserId)
			if err2 != nil {
				impl.logger.Errorw("error while updating current runner status to failed, ManualCdTrigger", "cdWfr", overrideRequest.WfrId, "err2", err2)
			}
		}
		return 0, "", nil, err
	}
	envDeploymentConfig, err := impl.deploymentConfigService.GetAndMigrateConfigIfAbsentForDevtronApps(cdPipeline.AppId, cdPipeline.EnvironmentId)
	if err != nil {
		impl.logger.Errorw("error in fetching environment deployment config by appId and envId", "appId", cdPipeline.AppId, "envId", cdPipeline.EnvironmentId, "err", err)
		return 0, "", nil, err
	}

	adapter.SetPipelineFieldsInOverrideRequest(overrideRequest, cdPipeline, envDeploymentConfig)
	ciArtifactId := overrideRequest.CiArtifactId

	_, span := otel.Tracer("orchestrator").Start(ctx, "ciArtifactRepository.Get")
	artifact, err := impl.ciArtifactRepository.Get(ciArtifactId)
	span.End()
	if err != nil {
		impl.logger.Errorw("error in getting CiArtifact", "CiArtifactId", overrideRequest.CiArtifactId, "err", err)
		return 0, "", nil, err
	}

	if artifact.IsMigrationRequired() {
		// Migration of deprecated DataSource Type
		migrationErr := impl.ciArtifactRepository.MigrateToWebHookDataSourceType(artifact.Id)
		if migrationErr != nil {
			impl.logger.Warnw("unable to migrate deprecated DataSource", "artifactId", artifact.Id)
		}
	}

	_, imageTag, err := artifact.ExtractImageRepoAndTag()
	if err != nil {
		impl.logger.Errorw("error in getting image tag and repo", "err", err)
	}
	helmPackageName := fmt.Sprintf("%s-%s-%s", cdPipeline.App.AppName, cdPipeline.Environment.Name, imageTag)
	var manifestPushTemplate *bean4.ManifestPushTemplate

	switch overrideRequest.CdWorkflowType {
	case bean3.CD_WORKFLOW_TYPE_PRE:
		var cdWf *pipelineConfig.CdWorkflow
		if overrideRequest.CdWorkflowId == 0 {
			cdWf = &pipelineConfig.CdWorkflow{
				CiArtifactId: artifact.Id,
				PipelineId:   cdPipeline.Id,
				AuditLog:     sql.AuditLog{CreatedOn: triggeredAt, CreatedBy: 1, UpdatedOn: triggeredAt, UpdatedBy: 1},
			}
			err := impl.cdWorkflowRepository.SaveWorkFlow(ctx, cdWf)
			if err != nil {
				return 0, "", nil, err
			}
		} else {
			cdWf, err = impl.cdWorkflowRepository.FindById(overrideRequest.CdWorkflowId)
			if err != nil {
				impl.logger.Errorw("error in TriggerPreStage, ManualCdTrigger", "err", err)
				return 0, "", nil, err
			}
		}
		overrideRequest.CdWorkflowId = cdWf.Id

		_, span = otel.Tracer("orchestrator").Start(ctx, "TriggerPreStage")
		triggerRequest := bean.TriggerRequest{
			CdWf:                  cdWf,
			Artifact:              artifact,
			Pipeline:              cdPipeline,
			TriggeredBy:           overrideRequest.UserId,
			ApplyAuth:             false,
			TriggerContext:        triggerContext,
			RefCdWorkflowRunnerId: 0,
			CdWorkflowRunnerId:    overrideRequest.WfrId,
		}
		manifestPushTemplate, err = impl.TriggerPreStage(triggerRequest)
		span.End()
		if err != nil {
			impl.logger.Errorw("error in TriggerPreStage, ManualCdTrigger", "err", err)
			return 0, "", nil, err
		}
	case bean3.CD_WORKFLOW_TYPE_DEPLOY:
		if overrideRequest.DeploymentType == models.DEPLOYMENTTYPE_UNKNOWN {
			overrideRequest.DeploymentType = models.DEPLOYMENTTYPE_DEPLOY
		}

		cdWf, err := impl.cdWorkflowRepository.FindByWorkflowIdAndRunnerType(ctx, overrideRequest.CdWorkflowId, bean3.CD_WORKFLOW_TYPE_PRE)
		if err != nil && !util.IsErrNoRows(err) {
			impl.logger.Errorw("error in getting cdWorkflow, ManualCdTrigger", "CdWorkflowId", overrideRequest.CdWorkflowId, "err", err)
			return 0, "", nil, err
		}

		cdWorkflowId := cdWf.CdWorkflowId
		if cdWf.CdWorkflowId == 0 {
			cdWf := &pipelineConfig.CdWorkflow{
				CiArtifactId: overrideRequest.CiArtifactId,
				PipelineId:   overrideRequest.PipelineId,
				AuditLog:     sql.AuditLog{CreatedOn: triggeredAt, CreatedBy: overrideRequest.UserId, UpdatedOn: triggeredAt, UpdatedBy: overrideRequest.UserId},
			}
			err := impl.cdWorkflowRepository.SaveWorkFlow(ctx, cdWf)
			if err != nil {
				impl.logger.Errorw("error in creating cdWorkflow, ManualCdTrigger", "PipelineId", overrideRequest.PipelineId, "err", err)
				return 0, "", nil, err
			}
			cdWorkflowId = cdWf.Id
		}

		runner := &pipelineConfig.CdWorkflowRunner{
			Name:         cdPipeline.Name,
			WorkflowType: bean3.CD_WORKFLOW_TYPE_DEPLOY,
			ExecutorType: cdWorkflow.WORKFLOW_EXECUTOR_TYPE_AWF,
			Status:       cdWorkflow.WorkflowInitiated, //deployment Initiated for manual trigger
			TriggeredBy:  overrideRequest.UserId,
			StartedOn:    triggeredAt,
			Namespace:    impl.config.GetDefaultNamespace(),
			CdWorkflowId: cdWorkflowId,
			AuditLog:     sql.AuditLog{CreatedOn: triggeredAt, CreatedBy: overrideRequest.UserId, UpdatedOn: triggeredAt, UpdatedBy: overrideRequest.UserId},
			ReferenceId:  triggerContext.ReferenceId,
		}
		savedWfr, err := impl.cdWorkflowRepository.SaveWorkFlowRunner(runner)
		if err != nil {
			impl.logger.Errorw("err in creating cdWorkflowRunner, ManualCdTrigger", "cdWorkflowId", cdWorkflowId, "err", err)
			return 0, "", nil, err
		}
		runner.CdWorkflow = &pipelineConfig.CdWorkflow{
			Pipeline: cdPipeline,
		}
		overrideRequest.WfrId = savedWfr.Id
		overrideRequest.CdWorkflowId = cdWorkflowId
		// creating cd pipeline status timeline for deployment initialisation
		timeline := impl.pipelineStatusTimelineService.NewDevtronAppPipelineStatusTimelineDbObject(runner.Id, timelineStatus.TIMELINE_STATUS_DEPLOYMENT_INITIATED, timelineStatus.TIMELINE_DESCRIPTION_DEPLOYMENT_INITIATED, overrideRequest.UserId)
		_, span = otel.Tracer("orchestrator").Start(ctx, "cdPipelineStatusTimelineRepo.SaveTimelineForACDHelmApps")
		_, err = impl.pipelineStatusTimelineService.SaveTimelineIfNotAlreadyPresent(timeline, nil)

		span.End()
		if err != nil {
			impl.logger.Errorw("error in creating timeline status for deployment initiation, ManualCdTrigger", "err", err, "timeline", timeline)
		}
		if isNotHibernateRequest(overrideRequest.DeploymentType) {
			validationErr := impl.validateDeploymentTriggerRequest(ctx, runner, cdPipeline, artifact.ImageDigest, envDeploymentConfig, overrideRequest.UserId)
			if validationErr != nil {
				impl.logger.Errorw("validation error deployment request", "cdWfr", runner.Id, "err", validationErr)
				return 0, "", nil, validationErr
			}
		}
		// Deploy the release
		var releaseErr error
		releaseId, manifestPushTemplate, releaseErr = impl.handleCDTriggerRelease(ctx, overrideRequest, envDeploymentConfig, triggeredAt, overrideRequest.UserId)
		// if releaseErr found, then the mark current deployment Failed and return
		if releaseErr != nil {
			err := impl.cdWorkflowCommonService.MarkCurrentDeploymentFailed(runner, releaseErr, overrideRequest.UserId)
			if err != nil {
				impl.logger.Errorw("error while updating current runner status to failed", "cdWfr", runner.Id, "err", err)
			}
			return 0, "", nil, releaseErr
		}

	case bean3.CD_WORKFLOW_TYPE_POST:
		cdWfRunner, err := impl.cdWorkflowRepository.FindByWorkflowIdAndRunnerType(ctx, overrideRequest.CdWorkflowId, bean3.CD_WORKFLOW_TYPE_DEPLOY)
		if err != nil && !util.IsErrNoRows(err) {
			impl.logger.Errorw("err in getting cdWorkflowRunner, ManualCdTrigger", "cdWorkflowId", overrideRequest.CdWorkflowId, "err", err)
			return 0, "", nil, err
		}

		var cdWf *pipelineConfig.CdWorkflow
		if cdWfRunner.CdWorkflowId == 0 {
			cdWf = &pipelineConfig.CdWorkflow{
				CiArtifactId: ciArtifactId,
				PipelineId:   overrideRequest.PipelineId,
				AuditLog:     sql.AuditLog{CreatedOn: triggeredAt, CreatedBy: overrideRequest.UserId, UpdatedOn: triggeredAt, UpdatedBy: overrideRequest.UserId},
			}
			err := impl.cdWorkflowRepository.SaveWorkFlow(ctx, cdWf)
			if err != nil {
				impl.logger.Errorw("error in creating cdWorkflow, ManualCdTrigger", "CdWorkflowId", overrideRequest.CdWorkflowId, "err", err)
				return 0, "", nil, err
			}
			overrideRequest.CdWorkflowId = cdWf.Id
		} else {
			_, span = otel.Tracer("orchestrator").Start(ctx, "cdWorkflowRepository.FindById")
			cdWf, err = impl.cdWorkflowRepository.FindById(overrideRequest.CdWorkflowId)
			span.End()
			if err != nil && !util.IsErrNoRows(err) {
				impl.logger.Errorw("error in getting cdWorkflow, ManualCdTrigger", "CdWorkflowId", overrideRequest.CdWorkflowId, "err", err)
				return 0, "", nil, err
			}
		}
		_, span = otel.Tracer("orchestrator").Start(ctx, "TriggerPostStage")
		triggerRequest := bean.TriggerRequest{
			CdWf:                  cdWf,
			Pipeline:              cdPipeline,
			TriggeredBy:           overrideRequest.UserId,
			RefCdWorkflowRunnerId: 0,
			TriggerContext:        triggerContext,
			CdWorkflowRunnerId:    overrideRequest.WfrId,
		}
		manifestPushTemplate, err = impl.TriggerPostStage(triggerRequest)
		span.End()
		if err != nil {
			impl.logger.Errorw("error in TriggerPostStage, ManualCdTrigger", "CdWorkflowId", cdWf.Id, "err", err)
			return 0, "", nil, err
		}
	default:
		impl.logger.Errorw("invalid CdWorkflowType, ManualCdTrigger", "CdWorkflowType", overrideRequest.CdWorkflowType, "err", err)
		return 0, "", nil, fmt.Errorf("invalid CdWorkflowType %s for the trigger request", string(overrideRequest.CdWorkflowType))
	}
	return releaseId, helmPackageName, manifestPushTemplate, err
}

func isNotHibernateRequest(deploymentType models.DeploymentType) bool {
	return deploymentType != models.DEPLOYMENTTYPE_STOP && deploymentType != models.DEPLOYMENTTYPE_START
}

// TODO: write a wrapper to handle auto and manual trigger
func (impl *TriggerServiceImpl) TriggerAutomaticDeployment(request bean.TriggerRequest) error {
	// in case of manual trigger auth is already applied and for auto triggers there is no need for auth check here
	triggeredBy := request.TriggeredBy
	pipeline := request.Pipeline
	artifact := request.Artifact

	//setting triggeredAt variable to have consistent data for various audit log places in db for deployment time
	triggeredAt := time.Now()
	cdWf := request.CdWf
	ctx := context.Background()

	if cdWf == nil || (cdWf != nil && cdWf.CiArtifactId != artifact.Id) {
		// cdWf != nil && cdWf.CiArtifactId != artifact.Id for auto trigger case when deployment is triggered with image generated by plugin
		cdWf = &pipelineConfig.CdWorkflow{
			CiArtifactId: artifact.Id,
			PipelineId:   pipeline.Id,
			AuditLog:     sql.AuditLog{CreatedOn: triggeredAt, CreatedBy: 1, UpdatedOn: triggeredAt, UpdatedBy: 1},
		}
		err := impl.cdWorkflowRepository.SaveWorkFlow(ctx, cdWf)
		if err != nil {
			return err
		}
	}

	runner := &pipelineConfig.CdWorkflowRunner{
		Name:         pipeline.Name,
		WorkflowType: bean3.CD_WORKFLOW_TYPE_DEPLOY,
		ExecutorType: cdWorkflow.WORKFLOW_EXECUTOR_TYPE_SYSTEM,
		Status:       cdWorkflow.WorkflowInitiated, // deployment Initiated for auto trigger
		TriggeredBy:  1,
		StartedOn:    triggeredAt,
		Namespace:    impl.config.GetDefaultNamespace(),
		CdWorkflowId: cdWf.Id,
		AuditLog:     sql.AuditLog{CreatedOn: triggeredAt, CreatedBy: triggeredBy, UpdatedOn: triggeredAt, UpdatedBy: triggeredBy},
		ReferenceId:  request.TriggerContext.ReferenceId,
	}
	savedWfr, err := impl.cdWorkflowRepository.SaveWorkFlowRunner(runner)
	if err != nil {
		return err
	}
	runner.CdWorkflow = &pipelineConfig.CdWorkflow{
		Pipeline: pipeline,
	}
	// creating cd pipeline status timeline for deployment initialisation
	timeline := &pipelineConfig.PipelineStatusTimeline{
		CdWorkflowRunnerId: runner.Id,
		Status:             timelineStatus.TIMELINE_STATUS_DEPLOYMENT_INITIATED,
		StatusDetail:       "Deployment initiated successfully.",
		StatusTime:         time.Now(),
	}
	timeline.CreateAuditLog(1)
	err = impl.pipelineStatusTimelineService.SaveTimeline(timeline, nil)
	if err != nil {
		impl.logger.Errorw("error in creating timeline status for deployment initiation", "err", err, "timeline", timeline)
	}
	envDeploymentConfig, err := impl.deploymentConfigService.GetAndMigrateConfigIfAbsentForDevtronApps(pipeline.AppId, pipeline.EnvironmentId)
	if err != nil {
		impl.logger.Errorw("error in fetching environment deployment config by appId and envId", "appId", pipeline.AppId, "envId", pipeline.EnvironmentId, "err", err)
		return err
	}
	// setting triggeredBy as 1(system user) since case of auto trigger
	validationErr := impl.validateDeploymentTriggerRequest(ctx, runner, pipeline, artifact.ImageDigest, envDeploymentConfig, 1)
	if validationErr != nil {
		impl.logger.Errorw("validation error deployment request", "cdWfr", runner.Id, "err", validationErr)
		return validationErr
	}
	releaseErr := impl.TriggerCD(ctx, artifact, cdWf.Id, savedWfr.Id, pipeline, envDeploymentConfig, triggeredAt)
	// if releaseErr found, then the mark current deployment Failed and return
	if releaseErr != nil {
		err := impl.cdWorkflowCommonService.MarkCurrentDeploymentFailed(runner, releaseErr, triggeredBy)
		if err != nil {
			impl.logger.Errorw("error while updating current runner status to failed, updatePreviousDeploymentStatus", "cdWfr", runner.Id, "err", err)
		}
		return releaseErr
	}
	return nil
}

func (impl *TriggerServiceImpl) TriggerCD(ctx context.Context, artifact *repository3.CiArtifact, cdWorkflowId, wfrId int, pipeline *pipelineConfig.Pipeline, envDeploymentConfig *bean9.DeploymentConfig, triggeredAt time.Time) error {
	impl.logger.Debugw("automatic pipeline trigger attempt async", "artifactId", artifact.Id)
	err := impl.triggerReleaseAsync(ctx, artifact, cdWorkflowId, wfrId, pipeline, envDeploymentConfig, triggeredAt)
	if err != nil {
		impl.logger.Errorw("error in cd trigger", "err", err)
		return err
	}
	return err
}

func (impl *TriggerServiceImpl) triggerReleaseAsync(ctx context.Context, artifact *repository3.CiArtifact, cdWorkflowId, wfrId int, pipeline *pipelineConfig.Pipeline, envDeploymentConfig *bean9.DeploymentConfig, triggeredAt time.Time) error {
	err := impl.validateAndTrigger(ctx, pipeline, envDeploymentConfig, artifact, cdWorkflowId, wfrId, triggeredAt)
	if err != nil {
		impl.logger.Errorw("error in trigger for pipeline", "pipelineId", strconv.Itoa(pipeline.Id))
	}
	impl.logger.Debugw("trigger attempted for all pipeline ", "artifactId", artifact.Id)
	return err
}

func (impl *TriggerServiceImpl) validateAndTrigger(ctx context.Context, p *pipelineConfig.Pipeline, envDeploymentConfig *bean9.DeploymentConfig, artifact *repository3.CiArtifact, cdWorkflowId, wfrId int, triggeredAt time.Time) error {
	//TODO: verify this logic
	object := impl.enforcerUtil.GetAppRBACNameByAppId(p.AppId)
	envApp := strings.Split(object, "/")
	if len(envApp) != 2 {
		impl.logger.Error("invalid req, app and env not found from rbac")
		return errors.New("invalid req, app and env not found from rbac")
	}
	err := impl.releasePipeline(ctx, p, envDeploymentConfig, artifact, cdWorkflowId, wfrId, triggeredAt)
	return err
}

func (impl *TriggerServiceImpl) releasePipeline(ctx context.Context, pipeline *pipelineConfig.Pipeline, envDeploymentConfig *bean9.DeploymentConfig, artifact *repository3.CiArtifact, cdWorkflowId, wfrId int, triggeredAt time.Time) error {
	startTime := time.Now()
	defer func() {
		impl.logger.Debugw("auto trigger release process completed", "timeTaken", time.Since(startTime), "cdPipelineId", pipeline.Id, "artifactId", artifact.Id, "wfrId", wfrId)
	}()
	impl.logger.Debugw("auto triggering release for", "cdPipelineId", pipeline.Id, "artifactId", artifact.Id, "wfrId", wfrId)
	pipeline, err := impl.pipelineRepository.FindById(pipeline.Id)
	if err != nil {
		impl.logger.Errorw("error in fetching pipeline by pipelineId", "err", err)
		return err
	}

	request := &bean3.ValuesOverrideRequest{
		PipelineId:           pipeline.Id,
		UserId:               artifact.CreatedBy,
		CiArtifactId:         artifact.Id,
		AppId:                pipeline.AppId,
		CdWorkflowId:         cdWorkflowId,
		ForceTrigger:         true,
		DeploymentWithConfig: bean3.DEPLOYMENT_CONFIG_TYPE_LAST_SAVED,
		WfrId:                wfrId,
	}

	adapter.SetPipelineFieldsInOverrideRequest(request, pipeline, envDeploymentConfig)

	releaseCtx, err := impl.argoUserService.GetACDContext(ctx)
	if err != nil {
		impl.logger.Errorw("error in creating acd sync context", "pipelineId", pipeline.Id, "artifactId", artifact.Id, "err", err)
		return err
	}
	// setting deployedBy as 1(system user) since case of auto trigger
	id, _, err := impl.handleCDTriggerRelease(releaseCtx, request, envDeploymentConfig, triggeredAt, 1)
	if err != nil {
		impl.logger.Errorw("error in auto  cd pipeline trigger", "pipelineId", pipeline.Id, "artifactId", artifact.Id, "err", err)
	} else {
		impl.logger.Infow("pipeline successfully triggered", "cdPipelineId", pipeline.Id, "artifactId", artifact.Id, "releaseId", id)
	}
	return err
}

func (impl *TriggerServiceImpl) triggerAsyncRelease(userDeploymentRequestId int, overrideRequest *bean3.ValuesOverrideRequest, ctx context.Context, triggeredAt time.Time, deployedBy int32) (releaseNo int, manifestPushTemplate *bean4.ManifestPushTemplate, err error) {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.triggerAsyncRelease")
	defer span.End()
	// build merged values and save PCO history for the release
	valuesOverrideResponse, err := impl.manifestCreationService.GetValuesOverrideForTrigger(overrideRequest, triggeredAt, newCtx)
	// auditDeploymentTriggerHistory is performed irrespective of GetValuesOverrideForTrigger error - for auditing purposes
	historyErr := impl.auditDeploymentTriggerHistory(overrideRequest.WfrId, valuesOverrideResponse, newCtx, triggeredAt, deployedBy)
	if historyErr != nil {
		impl.logger.Errorw("error in auditing deployment trigger history", "cdWfrId", overrideRequest.WfrId, "err", err)
		return releaseNo, manifestPushTemplate, err
	}
	// handling GetValuesOverrideForTrigger error
	if err != nil {
		impl.logger.Errorw("error in fetching values for trigger", "err", err)
		return releaseNo, manifestPushTemplate, err
	}
	// asynchronous mode of Helm/ArgoCd installation starts
	return impl.workflowEventPublishService.TriggerAsyncRelease(userDeploymentRequestId, overrideRequest, valuesOverrideResponse, newCtx, deployedBy)
}

func (impl *TriggerServiceImpl) handleCDTriggerRelease(ctx context.Context, overrideRequest *bean3.ValuesOverrideRequest, envDeploymentConfig *bean9.DeploymentConfig, triggeredAt time.Time, deployedBy int32) (releaseNo int, manifestPushTemplate *bean4.ManifestPushTemplate, err error) {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.handleCDTriggerRelease")
	defer span.End()
	// Handling for auto trigger
	if overrideRequest.UserId == 0 {
		overrideRequest.UserId = deployedBy
	}
	tx, err := impl.transactionUtilImpl.StartTx()
	if err != nil {
		impl.logger.Errorw("error in starting transaction to update userDeploymentRequest", "error", err)
		return releaseNo, manifestPushTemplate, err
	}
	defer impl.transactionUtilImpl.RollbackTx(tx)
	newDeploymentRequest := adapter.NewUserDeploymentRequest(overrideRequest, triggeredAt, overrideRequest.UserId)
	// creating new user deployment request
	userDeploymentRequestId, err := impl.userDeploymentRequestService.SaveNewDeployment(newCtx, tx, newDeploymentRequest)
	if err != nil {
		impl.logger.Errorw("error in saving new userDeploymentRequest", "overrideRequest", overrideRequest, "err", err)
		return releaseNo, manifestPushTemplate, err
	}
	timeline := impl.pipelineStatusTimelineService.NewDevtronAppPipelineStatusTimelineDbObject(overrideRequest.WfrId, timelineStatus.TIMELINE_STATUS_DEPLOYMENT_REQUEST_VALIDATED, timelineStatus.TIMELINE_DESCRIPTION_DEPLOYMENT_REQUEST_VALIDATED, deployedBy)
	// creating cd pipeline status timeline for deployment trigger request validated
	_, err = impl.pipelineStatusTimelineService.SaveTimelineIfNotAlreadyPresent(timeline, tx)
	err = impl.transactionUtilImpl.CommitTx(tx)
	if err != nil {
		impl.logger.Errorw("error in committing transaction to update userDeploymentRequest", "error", err)
		return userDeploymentRequestId, manifestPushTemplate, err
	}
	isAsyncMode, err := impl.isDevtronAsyncInstallModeEnabled(overrideRequest)
	if err != nil {
		impl.logger.Errorw("error in checking async mode for devtron app", "err", err, "deploymentType", overrideRequest.DeploymentType,
			"forceSyncDeployment", overrideRequest.ForceSyncDeployment, "appId", overrideRequest.AppId, "envId", overrideRequest.EnvId)
		return userDeploymentRequestId, manifestPushTemplate, err
	}
	if isAsyncMode {
		return impl.triggerAsyncRelease(userDeploymentRequestId, overrideRequest, newCtx, triggeredAt, deployedBy)
	}
	// synchronous mode of installation starts
	return impl.TriggerRelease(overrideRequest, envDeploymentConfig, newCtx, triggeredAt, deployedBy)
}

func (impl *TriggerServiceImpl) auditDeploymentTriggerHistory(cdWfrId int, valuesOverrideResponse *app.ValuesOverrideResponse, ctx context.Context, triggeredAt time.Time, triggeredBy int32) (err error) {
	if valuesOverrideResponse.Pipeline == nil || valuesOverrideResponse.EnvOverride == nil {
		impl.logger.Warnw("unable to save histories for deployment trigger, invalid valuesOverrideResponse received", "cdWfrId", cdWfrId)
		return nil
	}
	err1 := impl.deployedConfigurationHistoryService.CreateHistoriesForDeploymentTrigger(ctx, valuesOverrideResponse.Pipeline, valuesOverrideResponse.PipelineStrategy, valuesOverrideResponse.EnvOverride, triggeredAt, triggeredBy)
	if err1 != nil {
		impl.logger.Errorw("error in saving histories for deployment trigger", "err", err1, "pipelineId", valuesOverrideResponse.Pipeline.Id, "cdWfrId", cdWfrId)
		return nil
	}
	return nil
}

// TriggerRelease will trigger Install/Upgrade request for Devtron App releases synchronously
func (impl *TriggerServiceImpl) TriggerRelease(overrideRequest *bean3.ValuesOverrideRequest, envDeploymentConfig *bean9.DeploymentConfig, ctx context.Context, triggeredAt time.Time, triggeredBy int32) (releaseNo int, manifestPushTemplate *bean4.ManifestPushTemplate, err error) {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.TriggerRelease")
	defer span.End()
	triggerEvent, skipRequest, err := impl.buildTriggerEventForOverrideRequest(overrideRequest, triggeredAt)
	if err != nil {
		return releaseNo, manifestPushTemplate, err
	}
	impl.logger.Debugw("processing TriggerRelease", "wfrId", overrideRequest.WfrId, "triggerEvent", triggerEvent)
	// request has already been served, skipping
	if skipRequest {
		impl.logger.Infow("request already served, skipping", "wfrId", overrideRequest.WfrId)
		return releaseNo, manifestPushTemplate, nil
	}
	// build merged values and save PCO history for the release
	valuesOverrideResponse, builtChartPath, err := impl.manifestCreationService.BuildManifestForTrigger(overrideRequest, triggeredAt, newCtx)

	if envDeploymentConfig == nil || (envDeploymentConfig != nil && envDeploymentConfig.Id == 0) {
		envDeploymentConfig, err1 := impl.deploymentConfigService.GetAndMigrateConfigIfAbsentForDevtronApps(overrideRequest.AppId, overrideRequest.EnvId)
		if err1 != nil {
			impl.logger.Errorw("error in getting deployment config by appId and envId", "appId", overrideRequest.AppId, "envId", overrideRequest.EnvId, "err", err1)
			return releaseNo, manifestPushTemplate, err1
		}
		valuesOverrideResponse.DeploymentConfig = envDeploymentConfig
	}
	valuesOverrideResponse.DeploymentConfig = envDeploymentConfig

	// auditDeploymentTriggerHistory is performed irrespective of BuildManifestForTrigger error - for auditing purposes
	historyErr := impl.auditDeploymentTriggerHistory(overrideRequest.WfrId, valuesOverrideResponse, newCtx, triggeredAt, triggeredBy)
	if historyErr != nil {
		impl.logger.Errorw("error in auditing deployment trigger history", "cdWfrId", overrideRequest.WfrId, "err", err)
		return releaseNo, manifestPushTemplate, err
	}
	if err != nil {
		impl.logger.Errorw("error in building merged manifest for trigger", "err", err)
		impl.manifestGenerationFailedTimelineHandling(triggerEvent, overrideRequest, err)
		return releaseNo, manifestPushTemplate, err
	}
	helmManifest, err := impl.getHelmManifestForTriggerRelease(ctx, triggerEvent, overrideRequest, valuesOverrideResponse, builtChartPath)
	if err != nil {
		impl.logger.Errorw("error, getHelmManifestForTriggerRelease", "err", err)
		return releaseNo, manifestPushTemplate, err
	}
	impl.logger.Debugw("triggering pipeline for release", "wfrId", overrideRequest.WfrId, "builtChartPath", builtChartPath)
	releaseNo, err = impl.triggerPipeline(overrideRequest, valuesOverrideResponse, builtChartPath, triggerEvent, newCtx)
	if err != nil {
		return 0, manifestPushTemplate, err
	}

	err = impl.triggerReleaseSuccessHandling(triggerEvent, overrideRequest, valuesOverrideResponse, helmManifest)
	if err != nil {
		impl.logger.Errorw("error, triggerReleaseSuccessHandling", "triggerEvent", triggerEvent, "err", err)
		return releaseNo, manifestPushTemplate, err
	}
	// creating cd pipeline status timeline for deployment triggered - for successfully triggered requests
	timeline := impl.pipelineStatusTimelineService.NewDevtronAppPipelineStatusTimelineDbObject(overrideRequest.WfrId, timelineStatus.TIMELINE_STATUS_DEPLOYMENT_TRIGGERED, timelineStatus.TIMELINE_DESCRIPTION_DEPLOYMENT_COMPLETED, overrideRequest.UserId)
	_, dbErr := impl.pipelineStatusTimelineService.SaveTimelineIfNotAlreadyPresent(timeline, nil)
	if dbErr != nil {
		impl.logger.Errorw("error in creating timeline status for deployment completed", "err", dbErr, "timeline", timeline)
	}
	impl.logger.Debugw("triggered pipeline for release successfully", "wfrId", overrideRequest.WfrId, "builtChartPath", builtChartPath)
	return releaseNo, valuesOverrideResponse.ManifestPushTemplate, nil
}

func (impl *TriggerServiceImpl) performGitOps(ctx context.Context,
	overrideRequest *bean3.ValuesOverrideRequest, valuesOverrideResponse *app.ValuesOverrideResponse,
	builtChartPath string, triggerEvent bean.TriggerEvent) error {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.performGitOps")
	defer span.End()
	// update workflow runner status, used in app workflow view
	err := impl.cdWorkflowCommonService.UpdateNonTerminalStatusInRunner(newCtx, overrideRequest.WfrId, overrideRequest.UserId, cdWorkflow.WorkflowInProgress)
	if err != nil {
		impl.logger.Errorw("error in updating the workflow runner status", "err", err)
		return err
	}
	manifestPushTemplate, err := impl.buildManifestPushTemplate(overrideRequest, valuesOverrideResponse, builtChartPath)
	if err != nil {
		impl.logger.Errorw("error in building manifest push template", "err", err)
		return err
	}
	manifestPushService := impl.getManifestPushService(triggerEvent.ManifestStorageType)
	manifestPushResponse := manifestPushService.PushChart(newCtx, manifestPushTemplate)
	if manifestPushResponse.Error != nil {
		impl.logger.Errorw("error in pushing manifest to git/helm", "err", manifestPushResponse.Error, "git_repo_url", manifestPushTemplate.RepoUrl)
		return manifestPushResponse.Error
	}
	if manifestPushResponse.IsNewGitRepoConfigured() {
		// Update GitOps repo url after repo new repo created
		valuesOverrideResponse.DeploymentConfig.RepoURL = manifestPushResponse.NewGitRepoUrl
	}
	valuesOverrideResponse.ManifestPushTemplate = manifestPushTemplate
	return nil
}

func (impl *TriggerServiceImpl) buildTriggerEventForOverrideRequest(overrideRequest *bean3.ValuesOverrideRequest, triggeredAt time.Time) (triggerEvent bean.TriggerEvent, skipRequest bool, err error) {
	triggerEvent = helper.NewTriggerEvent(overrideRequest.DeploymentAppType, triggeredAt, overrideRequest.UserId)
	request := statusBean.NewTimelineGetRequest().
		WithCdWfrId(overrideRequest.WfrId).
		ExcludingStatuses(timelineStatus.TIMELINE_STATUS_UNABLE_TO_FETCH_STATUS,
			timelineStatus.TIMELINE_STATUS_KUBECTL_APPLY_STARTED,
			timelineStatus.TIMELINE_STATUS_KUBECTL_APPLY_SYNCED)
	timelineStatuses, err := impl.pipelineStatusTimelineService.GetTimelineStatusesFor(request)
	if err != nil {
		impl.logger.Errorw("error in getting last timeline status by cdWfrId", "cdWfrId", overrideRequest.WfrId, "err", err)
		return triggerEvent, skipRequest, err
	} else if !slices.Contains(timelineStatuses, timelineStatus.TIMELINE_STATUS_DEPLOYMENT_REQUEST_VALIDATED) {
		impl.logger.Errorw("pre-condition missing: timeline for deployment request validation", "cdWfrId", overrideRequest.WfrId, "timelineStatuses", timelineStatuses)
		return triggerEvent, skipRequest, fmt.Errorf("pre-condition missing: timeline for deployment request validation")
	} else if timelineStatus.ContainsTerminalTimelineStatus(timelineStatuses) {
		impl.logger.Info("deployment is already terminated", "cdWfrId", overrideRequest.WfrId, "timelineStatuses", timelineStatuses)
		skipRequest = true
		return triggerEvent, skipRequest, nil
	} else if slices.Contains(timelineStatuses, timelineStatus.TIMELINE_STATUS_DEPLOYMENT_TRIGGERED) {
		impl.logger.Info("deployment has been performed. skipping", "cdWfrId", overrideRequest.WfrId, "timelineStatuses", timelineStatuses)
		skipRequest = true
		return triggerEvent, skipRequest, nil
	}
	if slices.Contains(timelineStatuses, timelineStatus.TIMELINE_STATUS_GIT_COMMIT) ||
		slices.Contains(timelineStatuses, timelineStatus.TIMELINE_STATUS_ARGOCD_SYNC_INITIATED) {
		// git commit has already been performed
		triggerEvent.PerformChartPush = false
	}
	if slices.Contains(timelineStatuses, timelineStatus.TIMELINE_STATUS_ARGOCD_SYNC_COMPLETED) {
		// ArgoCd sync has already been performed
		triggerEvent.DeployArgoCdApp = false
	}
	return triggerEvent, skipRequest, nil
}

func (impl *TriggerServiceImpl) triggerPipeline(overrideRequest *bean3.ValuesOverrideRequest, valuesOverrideResponse *app.ValuesOverrideResponse, builtChartPath string, triggerEvent bean.TriggerEvent, ctx context.Context) (releaseNo int, err error) {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.triggerPipeline")
	defer span.End()
	if triggerEvent.PerformChartPush {
		impl.logger.Debugw("performing chart push operation in triggerPipeline", "cdWfrId", overrideRequest.WfrId)
		err = impl.performGitOps(newCtx, overrideRequest, valuesOverrideResponse, builtChartPath, triggerEvent)
		if err != nil {
			impl.logger.Errorw("error in performing GitOps", "cdWfrId", overrideRequest.WfrId, "err", err)
			return releaseNo, err
		}
		impl.logger.Debugw("chart push operation completed successfully", "cdWfrId", overrideRequest.WfrId)
	}

	if triggerEvent.PerformDeploymentOnCluster {
		err = impl.deployApp(newCtx, overrideRequest, valuesOverrideResponse, triggerEvent)
		if err != nil {
			impl.logger.Errorw("error in deploying app", "err", err)
			return releaseNo, err
		}
	}

	go impl.writeCDTriggerEvent(overrideRequest, valuesOverrideResponse.Artifact, valuesOverrideResponse.PipelineOverride.PipelineReleaseCounter, valuesOverrideResponse.PipelineOverride.Id, overrideRequest.WfrId)

	_ = impl.markImageScanDeployed(newCtx, overrideRequest.AppId, overrideRequest.EnvId, overrideRequest.ClusterId,
		valuesOverrideResponse.Artifact.ImageDigest, valuesOverrideResponse.Artifact.ScanEnabled, valuesOverrideResponse.Artifact.Image)

	middleware.CdTriggerCounter.WithLabelValues(overrideRequest.AppName, overrideRequest.EnvName).Inc()

	// Update previous deployment runner status (in transaction): Failed
	dbErr := impl.cdWorkflowCommonService.SupersedePreviousDeployments(newCtx, overrideRequest.WfrId, overrideRequest.PipelineId, triggerEvent.TriggeredAt, overrideRequest.UserId)
	if dbErr != nil {
		impl.logger.Errorw("error while update previous cd workflow runners", "err", dbErr, "currentRunnerId", overrideRequest.WfrId, "pipelineId", overrideRequest.PipelineId)
		return releaseNo, dbErr
	}
	return valuesOverrideResponse.PipelineOverride.PipelineReleaseCounter, nil
}

func (impl *TriggerServiceImpl) buildManifestPushTemplate(overrideRequest *bean3.ValuesOverrideRequest, valuesOverrideResponse *app.ValuesOverrideResponse, builtChartPath string) (*bean4.ManifestPushTemplate, error) {

	manifestPushTemplate := &bean4.ManifestPushTemplate{
		WorkflowRunnerId:      overrideRequest.WfrId,
		AppId:                 overrideRequest.AppId,
		ChartRefId:            valuesOverrideResponse.EnvOverride.Chart.ChartRefId,
		EnvironmentId:         valuesOverrideResponse.EnvOverride.Environment.Id,
		EnvironmentName:       valuesOverrideResponse.EnvOverride.Environment.Namespace,
		UserId:                overrideRequest.UserId,
		PipelineOverrideId:    valuesOverrideResponse.PipelineOverride.Id,
		AppName:               overrideRequest.AppName,
		TargetEnvironmentName: valuesOverrideResponse.EnvOverride.TargetEnvironment,
		BuiltChartPath:        builtChartPath,
		MergedValues:          valuesOverrideResponse.MergedValues,
	}

	manifestPushConfig, err := impl.manifestPushConfigRepository.GetManifestPushConfigByAppIdAndEnvId(overrideRequest.AppId, overrideRequest.EnvId)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error in fetching manifest push config from db", "err", err)
		return manifestPushTemplate, err
	}

	if manifestPushConfig != nil && manifestPushConfig.Id != 0 {
		if manifestPushConfig.StorageType == bean2.ManifestStorageGit {
			// need to implement for git repo push
			// currently manifest push config doesn't have git push config. GitOps config is derived from charts, chart_env_config_override and chart_ref table
		} else {
			err2 := impl.buildManifestPushTemplateForNonGitStorageType(overrideRequest, valuesOverrideResponse, builtChartPath, err, manifestPushConfig, manifestPushTemplate)
			if err2 != nil {
				return manifestPushTemplate, err2
			}
		}
	} else {
		manifestPushTemplate.ChartReferenceTemplate = valuesOverrideResponse.EnvOverride.Chart.ReferenceTemplate
		manifestPushTemplate.ChartName = valuesOverrideResponse.EnvOverride.Chart.ChartName
		manifestPushTemplate.ChartVersion = valuesOverrideResponse.EnvOverride.Chart.ChartVersion
		manifestPushTemplate.ChartLocation = valuesOverrideResponse.EnvOverride.Chart.ChartLocation
		manifestPushTemplate.RepoUrl = valuesOverrideResponse.DeploymentConfig.RepoURL
		manifestPushTemplate.IsCustomGitRepository = common.IsCustomGitOpsRepo(valuesOverrideResponse.DeploymentConfig.ConfigType)
	}
	return manifestPushTemplate, nil
}

func (impl *TriggerServiceImpl) deployApp(ctx context.Context, overrideRequest *bean3.ValuesOverrideRequest, valuesOverrideResponse *app.ValuesOverrideResponse, triggerEvent bean.TriggerEvent) error {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.deployApp")
	defer span.End()
	var referenceChartByte []byte
	var err error

	if util.IsAcdApp(overrideRequest.DeploymentAppType) && triggerEvent.DeployArgoCdApp {
		err = impl.deployArgoCdApp(newCtx, overrideRequest, valuesOverrideResponse)
		if err != nil {
			impl.logger.Errorw("error in deploying app on ArgoCd", "err", err)
			return err
		}
	} else if util.IsHelmApp(overrideRequest.DeploymentAppType) {
		_, referenceChartByte, err = impl.createHelmAppForCdPipeline(newCtx, overrideRequest, valuesOverrideResponse)
		if err != nil {
			impl.logger.Errorw("error in creating or updating helm application for cd pipeline", "err", err)
			return err
		}
	}
	impl.postDeployHook(overrideRequest, valuesOverrideResponse, referenceChartByte, err)
	return nil
}

func (impl *TriggerServiceImpl) createHelmAppForCdPipeline(ctx context.Context, overrideRequest *bean3.ValuesOverrideRequest, valuesOverrideResponse *app.ValuesOverrideResponse) (bool, []byte, error) {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.createHelmAppForCdPipeline")
	defer span.End()
	pipelineModel := valuesOverrideResponse.Pipeline
	envOverride := valuesOverrideResponse.EnvOverride
	mergeAndSave := valuesOverrideResponse.MergedValues
	chartMetaData, helmRevisionHistory, releaseIdentifier, err := impl.getHelmHistoryLimitAndChartMetadataForHelmAppCreation(ctx, valuesOverrideResponse)
	if err != nil {
		impl.logger.Errorw("error, getHelmHistoryLimitAndChartMetadataForHelmAppCreation", "valuesOverrideResponse", valuesOverrideResponse, "err", err)
		return false, nil, err
	}
	referenceTemplate := envOverride.Chart.ReferenceTemplate
	referenceTemplatePath := path.Join(bean5.RefChartDirPath, referenceTemplate)
	var referenceChartByte []byte
	if util.IsHelmApp(valuesOverrideResponse.DeploymentConfig.DeploymentAppType) {
		sanitizedK8sVersion, err := impl.getSanitizedK8sVersion(referenceTemplate)
		if err != nil {
			return false, nil, err
		}
		referenceChartByte, err = impl.getReferenceChartByteForHelmTypeApp(envOverride, chartMetaData, referenceTemplatePath, overrideRequest, valuesOverrideResponse)
		if err != nil {
			impl.logger.Errorw("error, getReferenceChartByteForHelmTypeApp", "envOverride", envOverride, "err", err)
			return false, nil, err
		}
		if pipelineModel.DeploymentAppCreated {
			req := &gRPC.UpgradeReleaseRequest{
				ReleaseIdentifier: releaseIdentifier,
				ValuesYaml:        mergeAndSave,
				HistoryMax:        helmRevisionHistory,
				ChartContent:      &gRPC.ChartContent{Content: referenceChartByte},
			}
			if len(sanitizedK8sVersion) > 0 {
				req.K8SVersion = sanitizedK8sVersion
			}
			if impl.isDevtronAsyncHelmInstallModeEnabled(overrideRequest.ForceSyncDeployment) {
				req.RunInCtx = true
			}
			// For cases where helm release was not found, kubelink will install the same configuration
			updateApplicationResponse, err := impl.helmAppClient.UpdateApplication(newCtx, req)
			if err != nil {
				impl.logger.Errorw("error in updating helm application for cd pipelineModel", "err", err)
				if util.IsErrorContextCancelled(err) {
					return false, nil, cdWorkflow.ErrorDeploymentSuperseded
				} else if util.IsErrorContextDeadlineExceeded(err) {
					return false, nil, context.DeadlineExceeded
				}
				apiError := clientErrors.ConvertToApiError(err)
				if apiError != nil {
					return false, nil, apiError
				}
				return false, nil, err
			} else {
				impl.logger.Debugw("updated helm application", "response", updateApplicationResponse, "isSuccess", updateApplicationResponse.Success)
			}

		} else {

			helmResponse, err := impl.helmInstallReleaseWithCustomChart(newCtx, releaseIdentifier, referenceChartByte,
				mergeAndSave, sanitizedK8sVersion, overrideRequest.ForceSyncDeployment)

			// For connection related errors, no need to update the db
			if err != nil && strings.Contains(err.Error(), "connection error") {
				impl.logger.Errorw("error in helm install custom chart", "err", err)
				return false, nil, err
			}

			// IMP: update cd pipelineModel to mark deployment app created, even if helm install fails
			// If the helm install fails, it still creates the app in failed state, so trying to
			// re-create the app results in error from helm that cannot re-use name which is still in use
			_, pgErr := impl.updatePipeline(pipelineModel, overrideRequest.UserId)

			if err != nil {
				impl.logger.Errorw("error in helm install custom chart", "err", err)
				if pgErr != nil {
					impl.logger.Errorw("failed to update deployment app created flag in pipelineModel table", "err", err)
				}
				if util.IsErrorContextCancelled(err) {
					return false, nil, cdWorkflow.ErrorDeploymentSuperseded
				} else if util.IsErrorContextDeadlineExceeded(err) {
					return false, nil, context.DeadlineExceeded
				}
				apiError := clientErrors.ConvertToApiError(err)
				if apiError != nil {
					return false, nil, apiError
				}
				return false, nil, err
			}

			if pgErr != nil {
				impl.logger.Errorw("failed to update deployment app created flag in pipelineModel table", "err", err)
				return false, nil, err
			}

			impl.logger.Debugw("received helm release response", "helmResponse", helmResponse, "isSuccess", helmResponse.Success)
		}

		//update workflow runner status, used in app workflow view
		err = impl.cdWorkflowCommonService.UpdateNonTerminalStatusInRunner(newCtx, overrideRequest.WfrId, overrideRequest.UserId, cdWorkflow.WorkflowInProgress)
		if err != nil {
			impl.logger.Errorw("error in updating the workflow runner status, createHelmAppForCdPipeline", "err", err)
			return false, nil, err
		}
	}
	return true, referenceChartByte, nil
}

func (impl *TriggerServiceImpl) deployArgoCdApp(ctx context.Context, overrideRequest *bean3.ValuesOverrideRequest,
	valuesOverrideResponse *app.ValuesOverrideResponse) error {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.deployArgoCdApp")
	defer span.End()
	impl.logger.Debugw("new pipeline found", "pipeline", valuesOverrideResponse.Pipeline)
	name, err := impl.createArgoApplicationIfRequired(newCtx, overrideRequest.AppId, valuesOverrideResponse.EnvOverride, valuesOverrideResponse.Pipeline, overrideRequest.UserId)
	if err != nil {
		impl.logger.Errorw("acd application create error on cd trigger", "err", err, "req", overrideRequest)
		return err
	}
	impl.logger.Debugw("ArgoCd application created", "name", name)
	updateAppInArgoCd, err := impl.updateArgoPipeline(newCtx, valuesOverrideResponse.Pipeline, valuesOverrideResponse.EnvOverride, valuesOverrideResponse.DeploymentConfig)
	if err != nil {
		impl.logger.Errorw("error in updating argocd app ", "err", err)
		return err
	}
	syncTime := time.Now()
	err = impl.argoClientWrapperService.SyncArgoCDApplicationIfNeededAndRefresh(newCtx, valuesOverrideResponse.Pipeline.DeploymentAppName)
	if err != nil {
		impl.logger.Errorw("error in getting argo application with normal refresh", "argoAppName", valuesOverrideResponse.Pipeline.DeploymentAppName)
		return fmt.Errorf("%s. err: %s", bean.ARGOCD_SYNC_ERROR, util.GetClientErrorDetailedMessage(err))
	}
	if impl.ACDConfig.IsManualSyncEnabled() {
		timeline := &pipelineConfig.PipelineStatusTimeline{
			CdWorkflowRunnerId: overrideRequest.WfrId,
			StatusTime:         syncTime,
			Status:             timelineStatus.TIMELINE_STATUS_ARGOCD_SYNC_COMPLETED,
			StatusDetail:       timelineStatus.TIMELINE_DESCRIPTION_ARGOCD_SYNC_COMPLETED,
		}
		timeline.CreateAuditLog(overrideRequest.UserId)
		_, err = impl.pipelineStatusTimelineService.SaveTimelineIfNotAlreadyPresent(timeline, nil)
		if err != nil {
			impl.logger.Errorw("error in saving pipeline status timeline", "err", err)
		}
	}
	if updateAppInArgoCd {
		impl.logger.Debug("argo-cd successfully updated")
	} else {
		impl.logger.Debug("argo-cd failed to update, ignoring it")
	}
	return nil
}

// update repoUrl, revision and argo app sync mode (auto/manual) if needed
func (impl *TriggerServiceImpl) updateArgoPipeline(ctx context.Context, pipeline *pipelineConfig.Pipeline, envOverride *bean10.EnvConfigOverride, deploymentConfig *bean9.DeploymentConfig) (bool, error) {
	if ctx == nil {
		impl.logger.Errorw("err in syncing ACD, ctx is NULL", "pipelineName", pipeline.Name)
		return false, nil
	}
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.updateArgoPipeline")
	defer span.End()
	argoAppName := pipeline.DeploymentAppName
	impl.logger.Infow("received payload, updateArgoPipeline", "appId", pipeline.AppId, "pipelineName", pipeline.Name, "envId", envOverride.TargetEnvironment, "argoAppName", argoAppName)
	argoApplication, err := impl.argoClientWrapperService.GetArgoAppByName(newCtx, argoAppName)
	if err != nil {
		impl.logger.Errorw("unable to get ArgoCd app", "app", argoAppName, "pipeline", pipeline.Name, "err", err)
		return false, err
	}
	// if status, ok:=status.FromError(err);ok{
	appStatus, _ := status2.FromError(err)
	if appStatus.Code() == codes.OK {
		impl.logger.Debugw("argo app exists", "app", argoAppName, "pipeline", pipeline.Name)
		if impl.argoClientWrapperService.IsArgoAppPatchRequired(argoApplication.Spec.Source, deploymentConfig.RepoURL, envOverride.Chart.ChartLocation) {
			patchRequestDto := &bean7.ArgoCdAppPatchReqDto{
				ArgoAppName:    argoAppName,
				ChartLocation:  envOverride.Chart.ChartLocation,
				GitRepoUrl:     deploymentConfig.RepoURL,
				TargetRevision: bean7.TargetRevisionMaster,
				PatchType:      bean7.PatchTypeMerge,
			}
			url, err := impl.gitOperationService.GetRepoUrlWithUserName(deploymentConfig.RepoURL)
			if err != nil {
				return false, err
			}
			patchRequestDto.GitRepoUrl = url
			err = impl.argoClientWrapperService.PatchArgoCdApp(newCtx, patchRequestDto)
			if err != nil {
				impl.logger.Errorw("error in patching argo pipeline", "err", err, "req", patchRequestDto)
				return false, err
			}
			if deploymentConfig.RepoURL != argoApplication.Spec.Source.RepoURL {
				impl.logger.Infow("patching argo application's repo url", "argoAppName", argoAppName)
			}
			impl.logger.Debugw("pipeline update req", "res", patchRequestDto)
		} else {
			impl.logger.Debug("pipeline no need to update ")
		}
		err := impl.argoClientWrapperService.UpdateArgoCDSyncModeIfNeeded(newCtx, argoApplication)
		if err != nil {
			impl.logger.Errorw("error in updating argocd sync mode", "err", err)
			return false, err
		}
		return true, nil
	} else if appStatus.Code() == codes.NotFound {
		impl.logger.Errorw("argo app not found", "app", argoAppName, "pipeline", pipeline.Name)
		return false, nil
	} else {
		impl.logger.Errorw("err in checking application on argoCD", "err", err, "pipeline", pipeline.Name)
		return false, err
	}
}

func (impl *TriggerServiceImpl) createArgoApplicationIfRequired(ctx context.Context, appId int, envConfigOverride *bean10.EnvConfigOverride, pipeline *pipelineConfig.Pipeline, userId int32) (string, error) {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.createArgoApplicationIfRequired")
	defer span.End()
	// repo has been registered while helm create
	chart, err := impl.chartRepository.FindLatestChartForAppByAppId(appId)
	if err != nil {
		impl.logger.Errorw("no chart found ", "app", appId)
		return "", err
	}
	envModel, err := impl.envRepository.FindById(envConfigOverride.TargetEnvironment)
	if err != nil {
		return "", err
	}
	argoAppName := pipeline.DeploymentAppName
	if pipeline.DeploymentAppCreated {
		return argoAppName, nil
	} else {
		// create
		appNamespace := envConfigOverride.Namespace
		if appNamespace == "" {
			appNamespace = "default"
		}
		namespace := argocdServer.DevtronInstalationNs

		appRequest := &argocdServer.AppTemplate{
			ApplicationName: argoAppName,
			Namespace:       namespace,
			TargetNamespace: appNamespace,
			TargetServer:    envModel.Cluster.ServerUrl,
			Project:         "default",
			ValuesFile:      helper.GetValuesFileForEnv(envModel.Id),
			RepoPath:        chart.ChartLocation,
			RepoUrl:         chart.GitRepoUrl,
			AutoSyncEnabled: impl.ACDConfig.ArgoCDAutoSyncEnabled,
		}
		appRequest.RepoUrl, err = impl.gitOperationService.GetRepoUrlWithUserName(appRequest.RepoUrl)
		if err != nil {
			return "", err
		}
		argoAppName, err := impl.argoK8sClient.CreateAcdApp(newCtx, appRequest, argocdServer.ARGOCD_APPLICATION_TEMPLATE)
		if err != nil {
			return "", err
		}
		// update cd pipeline to mark deployment app created
		_, err = impl.updatePipeline(pipeline, userId)
		if err != nil {
			impl.logger.Errorw("error in update cd pipeline for deployment app created or not", "err", err)
			return "", err
		}
		return argoAppName, nil
	}
}

func (impl *TriggerServiceImpl) updatePipeline(pipeline *pipelineConfig.Pipeline, userId int32) (bool, error) {
	err := impl.pipelineRepository.SetDeploymentAppCreatedInPipeline(true, pipeline.Id, userId)
	if err != nil {
		impl.logger.Errorw("error on updating cd pipeline for setting deployment app created", "err", err)
		return false, err
	}
	return true, nil
}

// helmInstallReleaseWithCustomChart performs helm install with custom chart
func (impl *TriggerServiceImpl) helmInstallReleaseWithCustomChart(ctx context.Context, releaseIdentifier *gRPC.ReleaseIdentifier,
	referenceChartByte []byte, valuesYaml, k8sServerVersion string, forceSync bool) (*gRPC.HelmInstallCustomResponse, error) {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.helmInstallReleaseWithCustomChart")
	defer span.End()
	helmInstallRequest := gRPC.HelmInstallCustomRequest{
		ValuesYaml:        valuesYaml,
		ChartContent:      &gRPC.ChartContent{Content: referenceChartByte},
		ReleaseIdentifier: releaseIdentifier,
	}
	if len(k8sServerVersion) > 0 {
		helmInstallRequest.K8SVersion = k8sServerVersion
	}
	if impl.isDevtronAsyncHelmInstallModeEnabled(forceSync) {
		helmInstallRequest.RunInCtx = true
	}
	// Request exec
	return impl.helmAppClient.InstallReleaseWithCustomChart(newCtx, &helmInstallRequest)
}

func (impl *TriggerServiceImpl) writeCDTriggerEvent(overrideRequest *bean3.ValuesOverrideRequest, artifact *repository3.CiArtifact, releaseId, pipelineOverrideId, wfrId int) {

	event, err := impl.eventFactory.Build(util2.Trigger, &overrideRequest.PipelineId, overrideRequest.AppId, &overrideRequest.EnvId, util2.CD)
	if err != nil {
		impl.logger.Errorw("error in building cd trigger event", "cdPipelineId", overrideRequest.PipelineId, "err", err)
	}
	impl.logger.Debugw("event WriteCDTriggerEvent", "event", event)
	wfr := impl.getEnrichedWorkflowRunner(overrideRequest, artifact, wfrId)
	event = impl.eventFactory.BuildExtraCDData(event, wfr, pipelineOverrideId, bean3.CD_WORKFLOW_TYPE_DEPLOY)
	_, evtErr := impl.eventClient.WriteNotificationEvent(event)
	if evtErr != nil {
		impl.logger.Errorw("CD trigger event not sent", "error", evtErr)
	}
	deploymentEvent := app.DeploymentEvent{
		ApplicationId:      overrideRequest.AppId,
		EnvironmentId:      overrideRequest.EnvId, // check for production Environment
		ReleaseId:          releaseId,
		PipelineOverrideId: pipelineOverrideId,
		TriggerTime:        time.Now(),
		CiArtifactId:       overrideRequest.CiArtifactId,
	}

	ciPipelineMaterials, err := impl.ciPipelineMaterialRepository.GetByPipelineId(artifact.PipelineId)
	if err != nil {
		impl.logger.Errorw("error in ")
	}
	materialInfoMap, mErr := artifact.ParseMaterialInfo()
	if mErr != nil {
		impl.logger.Errorw("material info map error", mErr)
		return
	}
	for _, ciPipelineMaterial := range ciPipelineMaterials {
		hash := materialInfoMap[ciPipelineMaterial.GitMaterial.Url]
		pipelineMaterialInfo := &app.PipelineMaterialInfo{PipelineMaterialId: ciPipelineMaterial.Id, CommitHash: hash}
		deploymentEvent.PipelineMaterials = append(deploymentEvent.PipelineMaterials, pipelineMaterialInfo)
	}
	impl.logger.Infow("triggering deployment event", "event", deploymentEvent)
	err = impl.eventClient.WriteNatsEvent(pubsub.CD_SUCCESS, deploymentEvent)
	if err != nil {
		impl.logger.Errorw("error in writing cd trigger event", "err", err)
	}
}

func (impl *TriggerServiceImpl) markImageScanDeployed(ctx context.Context, appId, envId, clusterId int,
	imageDigest string, isScanEnabled bool, image string) error {
	_, span := otel.Tracer("orchestrator").Start(ctx, "TriggerServiceImpl.markImageScanDeployed")
	defer span.End()
	// TODO KB: send NATS event for self consumption
	impl.logger.Debugw("mark image scan deployed for devtron app, from cd auto or manual trigger", "imageDigest", imageDigest)
	executionHistory, err := impl.imageScanHistoryReadService.FindByImageAndDigest(imageDigest, image)
	if err != nil && !errors.Is(err, pg.ErrNoRows) {
		impl.logger.Errorw("error in fetching execution history", "err", err)
		return err
	}
	if errors.Is(err, pg.ErrNoRows) || executionHistory == nil || executionHistory.Id == 0 {
		if isScanEnabled {
			// There should ImageScanHistory for ScanEnabled artifacts
			impl.logger.Errorw("no execution history found for digest", "digest", imageDigest)
			return fmt.Errorf("no execution history found for digest - %s", imageDigest)
		} else {
			// For ScanDisabled artifacts it should be an expected condition
			impl.logger.Infow("no execution history found for digest", "digest", imageDigest)
			return nil
		}
	}
	impl.logger.Debugw("saving image_scan_deploy_info for cd auto or manual trigger", "executionHistory", executionHistory)
	var ids []int
	ids = append(ids, executionHistory.Id)

	ot, err := impl.imageScanDeployInfoReadService.FetchByAppIdAndEnvId(appId, envId, []string{repository6.ScanObjectType_APP})

	if err == pg.ErrNoRows && !isScanEnabled {
		// ignoring if no rows are found and scan is disabled
		return nil
	}

	if err != nil && err != pg.ErrNoRows {
		return err
	} else if err == pg.ErrNoRows && isScanEnabled {
		imageScanDeployInfo := &repository6.ImageScanDeployInfo{
			ImageScanExecutionHistoryId: ids,
			ScanObjectMetaId:            appId,
			ObjectType:                  repository6.ScanObjectType_APP,
			EnvId:                       envId,
			ClusterId:                   clusterId,
			AuditLog: sql.AuditLog{
				CreatedOn: time.Now(),
				CreatedBy: 1,
				UpdatedOn: time.Now(),
				UpdatedBy: 1,
			},
		}
		impl.logger.Debugw("mark image scan deployed for normal app, from cd auto or manual trigger", "imageScanDeployInfo", imageScanDeployInfo)
		err = impl.imageScanDeployInfoService.Save(imageScanDeployInfo)
		if err != nil {
			impl.logger.Errorw("error in creating deploy info", "err", err)
		}
	} else {
		// Updating Execution history for Latest Deployment to fetch out security Vulnerabilities for latest deployed info
		if isScanEnabled {
			ot.ImageScanExecutionHistoryId = ids
		} else {
			arr := []int{-1}
			ot.ImageScanExecutionHistoryId = arr
		}
		err = impl.imageScanDeployInfoService.Update(ot)
		if err != nil {
			impl.logger.Errorw("error in updating deploy info for latest deployed image", "err", err)
		}
	}
	return err
}

func (impl *TriggerServiceImpl) isDevtronAsyncHelmInstallModeEnabled(forceSync bool) bool {
	return impl.globalEnvVariables.EnableAsyncHelmInstallDevtronChart && !forceSync
}

func (impl *TriggerServiceImpl) isDevtronAsyncInstallModeEnabled(overrideRequest *bean3.ValuesOverrideRequest) (bool, error) {
	if util.IsHelmApp(overrideRequest.DeploymentAppType) {
		return impl.isDevtronAsyncHelmInstallModeEnabled(overrideRequest.ForceSyncDeployment), nil
	} else if util.IsAcdApp(overrideRequest.DeploymentAppType) {
		return impl.isDevtronAsyncArgoCdInstallModeEnabledForApp(overrideRequest.AppId,
			overrideRequest.EnvId, overrideRequest.ForceSyncDeployment)
	}
	return false, nil
}

func (impl *TriggerServiceImpl) deleteCorruptedPipelineStage(pipelineStage *repository.PipelineStage, triggeredBy int32) (error, bool) {
	if pipelineStage != nil {
		stageReq := &bean8.PipelineStageDto{
			Id:   pipelineStage.Id,
			Type: pipelineStage.Type,
		}
		err, deleted := impl.pipelineStageService.DeletePipelineStageIfReq(stageReq, triggeredBy)
		if err != nil {
			impl.logger.Errorw("error in deleting the corrupted pipeline stage", "err", err, "pipelineStageReq", stageReq)
			return err, false
		}
		return nil, deleted
	}
	return nil, false
}

func (impl *TriggerServiceImpl) handleCustomGitOpsRepoValidation(runner *pipelineConfig.CdWorkflowRunner, pipeline *pipelineConfig.Pipeline, envDeploymentConfig *bean9.DeploymentConfig, triggeredBy int32) error {
	if !util.IsAcdApp(pipeline.DeploymentAppName) {
		return nil
	}
	isGitOpsConfigured := false
	gitOpsConfig, err := impl.gitOpsConfigReadService.GetGitOpsConfigActive()
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error while getting active GitOpsConfig", "err", err)
	}
	if gitOpsConfig != nil && gitOpsConfig.Id > 0 {
		isGitOpsConfigured = true
	}
	if isGitOpsConfigured && gitOpsConfig.AllowCustomRepository {
		//chart, err := impl.chartRepository.FindLatestChartForAppByAppId(pipeline.AppId)
		//if err != nil {
		//	impl.logger.Errorw("error in fetching latest chart for app by appId", "err", err, "appId", pipeline.AppId)
		//	return err
		//}
		if gitOps.IsGitOpsRepoNotConfigured(envDeploymentConfig.RepoURL) {
			if err = impl.cdWorkflowCommonService.MarkCurrentDeploymentFailed(runner, errors.New(cdWorkflow.GITOPS_REPO_NOT_CONFIGURED), triggeredBy); err != nil {
				impl.logger.Errorw("error while updating current runner status to failed, TriggerDeployment", "wfrId", runner.Id, "err", err)
			}
			apiErr := &util.ApiError{
				HttpStatusCode:  http.StatusConflict,
				UserMessage:     cdWorkflow.GITOPS_REPO_NOT_CONFIGURED,
				InternalMessage: cdWorkflow.GITOPS_REPO_NOT_CONFIGURED,
			}
			return apiErr
		}
	}
	return nil
}

func (impl *TriggerServiceImpl) getSanitizedK8sVersion(referenceTemplate string) (string, error) {
	var sanitizedK8sVersion string
	//handle specific case for all cronjob charts from cronjob-chart_1-2-0 to cronjob-chart_1-5-0 where semverCompare
	//comparison func has wrong api version mentioned, so for already installed charts via these charts that comparison
	//is always false, handles the gh issue:- https://github.com/devtron-labs/devtron/issues/4860
	cronJobChartRegex := regexp.MustCompile(bean.CronJobChartRegexExpression)
	if cronJobChartRegex.MatchString(referenceTemplate) {
		k8sServerVersion, err := impl.K8sUtil.GetKubeVersion()
		if err != nil {
			impl.logger.Errorw("exception caught in getting k8sServerVersion", "err", err)
			return "", err
		}
		sanitizedK8sVersion = k8s2.StripPrereleaseFromK8sVersion(k8sServerVersion.String())
	}
	return sanitizedK8sVersion, nil
}

func (impl *TriggerServiceImpl) getReferenceChartByteForHelmTypeApp(envOverride *bean10.EnvConfigOverride,
	chartMetaData *chart.Metadata, referenceTemplatePath string, overrideRequest *bean3.ValuesOverrideRequest,
	valuesOverrideResponse *app.ValuesOverrideResponse) ([]byte, error) {
	referenceChartByte := envOverride.Chart.ReferenceChart
	// here updating reference chart into database.
	if len(envOverride.Chart.ReferenceChart) == 0 {
		refChartByte, err := impl.chartTemplateService.GetByteArrayRefChart(chartMetaData, referenceTemplatePath)
		if err != nil {
			impl.logger.Errorw("ref chart commit error on cd trigger", "err", err, "req", overrideRequest)
			return nil, err
		}
		ch := envOverride.Chart
		ch.ReferenceChart = refChartByte
		ch.UpdatedOn = time.Now()
		ch.UpdatedBy = overrideRequest.UserId
		err = impl.chartRepository.Update(ch)
		if err != nil {
			impl.logger.Errorw("chart update error", "err", err, "req", overrideRequest)
			return nil, err
		}
		referenceChartByte = refChartByte
	}
	var err error
	referenceChartByte, err = impl.overrideReferenceChartByteForHelmTypeApp(valuesOverrideResponse, chartMetaData, referenceTemplatePath, referenceChartByte)
	if err != nil {
		impl.logger.Errorw("ref chart commit error on cd trigger", "err", err, "req", overrideRequest)
		return nil, err
	}
	return referenceChartByte, nil
}