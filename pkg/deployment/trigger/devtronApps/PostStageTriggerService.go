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
	bean2 "github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig/bean/workflow/cdWorkflow"
	bean4 "github.com/devtron-labs/devtron/pkg/app/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps/bean"
	bean3 "github.com/devtron-labs/devtron/pkg/pipeline/bean"
	repository3 "github.com/devtron-labs/devtron/pkg/pipeline/history/repository"
	"github.com/devtron-labs/devtron/pkg/pipeline/types"
	util2 "github.com/devtron-labs/devtron/util/event"
	"time"
)

func (impl *TriggerServiceImpl) TriggerPostStage(request bean.TriggerRequest) (*bean4.ManifestPushTemplate, error) {
	request.WorkflowType = bean2.CD_WORKFLOW_TYPE_POST
	// setting triggeredAt variable to have consistent data for various audit log places in db for deployment time
	triggeredAt := time.Now()
	triggeredBy := request.TriggeredBy
	pipeline := request.Pipeline
	cdWf := request.CdWf
	ctx := context.Background() //before there was only one context. To check why here we are not using ctx from request.TriggerContext
	env, namespace, err := impl.getEnvAndNsIfRunStageInEnv(ctx, request)
	if err != nil {
		impl.logger.Errorw("error, getEnvAndNsIfRunStageInEnv", "err", err, "pipeline", pipeline, "stage", request.WorkflowType)
		return nil, nil
	}
	request.RunStageInEnvNamespace = namespace

	cdWf, runner, err := impl.createStartingWfAndRunner(request, triggeredAt)
	if err != nil {
		impl.logger.Errorw("error in creating wf starting and runner entry", "err", err, "request", request)
		return nil, err
	}
	if cdWf.CiArtifact == nil || cdWf.CiArtifact.Id == 0 {
		cdWf.CiArtifact, err = impl.ciArtifactRepository.Get(cdWf.CiArtifactId)
		if err != nil {
			impl.logger.Errorw("error fetching artifact data", "err", err)
			return nil, err
		}
	}

	// Migration of deprecated DataSource Type
	if cdWf.CiArtifact.IsMigrationRequired() {
		migrationErr := impl.ciArtifactRepository.MigrateToWebHookDataSourceType(cdWf.CiArtifact.Id)
		if migrationErr != nil {
			impl.logger.Warnw("unable to migrate deprecated DataSource", "artifactId", cdWf.CiArtifact.Id)
		}
	}

	filterEvaluationAudit, err := impl.checkFeasibilityForPostStage(pipeline, &request, env, cdWf, triggeredBy)
	if err != nil {
		impl.logger.Errorw("error, checkFeasibilityForPostStage", "err", err, "pipeline", pipeline)
		return nil, nil
	}

	envDevploymentConfig, err := impl.deploymentConfigService.GetAndMigrateConfigIfAbsentForDevtronApps(pipeline.AppId, pipeline.EnvironmentId)
	if err != nil {
		impl.logger.Errorw("error in fetching deployment config by appId and envId", "appId", pipeline.AppId, "envId", pipeline.EnvironmentId, "err", err)
		return nil, err
	}

	dbErr := impl.createAuditDataForDeploymentWindowBypass(request, runner.Id)
	if dbErr != nil {
		impl.logger.Errorw("error in creating audit data for deployment window bypass", "runnerId", runner.Id, "err", dbErr)
		// skip error for audit data creation
	}

	err = impl.handlerFilterEvaluationAudit(filterEvaluationAudit, runner)
	if err != nil {
		impl.logger.Errorw("error, handlerFilterEvaluationAudit", "err", err)
		return nil, err
	}

	// custom GitOps repo url validation --> Start
	err = impl.handleCustomGitOpsRepoValidation(runner, pipeline, envDevploymentConfig, triggeredBy)
	if err != nil {
		impl.logger.Errorw("custom GitOps repository validation error, TriggerPreStage", "err", err)
		return nil, err
	}
	// custom GitOps repo url validation --> Ends

	// checking vulnerability for the selected image
	err = impl.checkVulnerabilityStatusAndFailWfIfNeeded(ctx, cdWf.CiArtifact, pipeline, runner, triggeredBy)
	if err != nil {
		impl.logger.Errorw("error, checkVulnerabilityStatusAndFailWfIfNeeded", "err", err, "runner", runner)
		return nil, err
	}
	cdStageWorkflowRequest, err := impl.buildWFRequest(runner, cdWf, pipeline, envDevploymentConfig, triggeredBy)
	if err != nil {
		return impl.buildWfRequestErrorHandler(runner, err, triggeredBy)
	}
	cdStageWorkflowRequest.StageType = types.POST
	cdStageWorkflowRequest.Pipeline = pipeline
	cdStageWorkflowRequest.Env = env
	cdStageWorkflowRequest.Type = bean3.CD_WORKFLOW_PIPELINE_TYPE
	// handling plugin specific logic

	pluginImagePathReservationIds, err := impl.setCopyContainerImagePluginDataAndReserveImages(cdStageWorkflowRequest, pipeline.Id, types.POST, cdWf.CiArtifact)
	if err != nil {
		runner.Status = cdWorkflow.WorkflowFailed
		runner.Message = err.Error()
		_ = impl.cdWorkflowRepository.UpdateWorkFlowRunner(runner)
		return nil, err
	}

	_, jobHelmPackagePath, err := impl.cdWorkflowService.SubmitWorkflow(cdStageWorkflowRequest)
	if err != nil {
		impl.logger.Errorw("error in submitting workflow", "err", err, "workflowId", cdStageWorkflowRequest.WorkflowId, "pipeline", pipeline, "env", env)
		return nil, err
	}
	manifestPushTempate, err := impl.getManifestPushTemplateForPostStage(request, envDevploymentConfig, jobHelmPackagePath, cdStageWorkflowRequest, cdWf, runner, pipeline, triggeredBy, triggeredAt)
	if err != nil {
		impl.logger.Errorw("error in getting manifest push template", "err", err)
		return nil, err
	}
	wfr, err := impl.cdWorkflowRepository.FindByWorkflowIdAndRunnerType(context.Background(), cdWf.Id, bean2.CD_WORKFLOW_TYPE_POST)
	if err != nil {
		impl.logger.Errorw("error in getting wfr by workflowId and runnerType", "err", err, "wfId", cdWf.Id)
		return nil, err
	}
	wfr.ImagePathReservationIds = pluginImagePathReservationIds
	err = impl.cdWorkflowRepository.UpdateWorkFlowRunner(&wfr)
	if err != nil {
		impl.logger.Error("error in updating image path reservation ids in cd workflow runner", "err", "err")
	}

	event, _ := impl.eventFactory.Build(util2.Trigger, &pipeline.Id, pipeline.AppId, &pipeline.EnvironmentId, util2.CD)
	impl.logger.Debugw("event Cd Post Trigger", "event", event)
	event = impl.eventFactory.BuildExtraCDData(event, &wfr, 0, bean2.CD_WORKFLOW_TYPE_POST)
	_, evtErr := impl.eventClient.WriteNotificationEvent(event)
	if evtErr != nil {
		impl.logger.Errorw("CD trigger event not sent", "error", evtErr)
	}
	// creating cd config history entry
	err = impl.prePostCdScriptHistoryService.CreatePrePostCdScriptHistory(pipeline, nil, repository3.POST_CD_TYPE, true, triggeredBy, triggeredAt)
	if err != nil {
		impl.logger.Errorw("error in creating post cd script entry", "err", err, "pipeline", pipeline)
		return nil, err
	}
	return manifestPushTempate, nil
}

func (impl *TriggerServiceImpl) buildWfRequestErrorHandler(runner *pipelineConfig.CdWorkflowRunner, err error, triggeredBy int32) (*bean4.ManifestPushTemplate, error) {
	dbErr := impl.cdWorkflowCommonService.MarkCurrentDeploymentFailed(runner, err, triggeredBy)
	if dbErr != nil {
		impl.logger.Errorw("error while updating current runner status to failed, buildWfRequestErrorHandler", "runner", runner.Id, "err", dbErr, "releaseErr", err)
	}
	return nil, err
}