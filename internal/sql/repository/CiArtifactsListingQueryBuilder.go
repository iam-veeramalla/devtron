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

package repository

import (
	"fmt"
	"github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/internal/sql/repository/helper"
)

const EmptyLikeRegex = "%%"

func BuildQueryForParentTypeCIOrWebhook(listingFilterOpts bean.ArtifactsListFilterOptions) string {
	commonPaginatedQueryPart := fmt.Sprintf(" cia.image LIKE '%v'", listingFilterOpts.SearchString)
	orderByClause := " ORDER BY cia.id DESC"
	limitOffsetQueryPart := fmt.Sprintf(" LIMIT %v OFFSET %v", listingFilterOpts.Limit, listingFilterOpts.Offset)
	finalQuery := ""
	if listingFilterOpts.ParentStageType == bean.CI_WORKFLOW_TYPE {
		selectQuery := " SELECT cia.* "
		remainingQuery := " FROM ci_artifact cia" +
			" INNER JOIN ci_pipeline cp ON (cp.id=cia.pipeline_id or (cp.id=cia.component_id and cia.data_source='post_ci' ) )" +
			" INNER JOIN pipeline p ON (p.ci_pipeline_id = cp.id and p.id=%v )" +
			" WHERE "
		remainingQuery = fmt.Sprintf(remainingQuery, listingFilterOpts.PipelineId)
		if len(listingFilterOpts.ExcludeArtifactIds) > 0 {
			remainingQuery += fmt.Sprintf("cia.id NOT IN (%s) AND ", helper.GetCommaSepratedString(listingFilterOpts.ExcludeArtifactIds))
		}

		countQuery := " SELECT count(cia.id)  as total_count"
		totalCountQuery := countQuery + remainingQuery + commonPaginatedQueryPart
		selectQuery = fmt.Sprintf("%s,(%s) ", selectQuery, totalCountQuery)
		finalQuery = selectQuery + remainingQuery + commonPaginatedQueryPart + orderByClause + limitOffsetQueryPart
	} else if listingFilterOpts.ParentStageType == bean.WEBHOOK_WORKFLOW_TYPE {
		selectQuery := " SELECT cia.* "
		remainingQuery := " FROM ci_artifact cia " +
			" WHERE cia.external_ci_pipeline_id = %v AND "
		remainingQuery = fmt.Sprintf(remainingQuery, listingFilterOpts.ParentId)
		if len(listingFilterOpts.ExcludeArtifactIds) > 0 {
			remainingQuery += fmt.Sprintf("cia.id NOT IN (%s) AND ", helper.GetCommaSepratedString(listingFilterOpts.ExcludeArtifactIds))
		}

		countQuery := " SELECT count(cia.id)  as total_count"
		totalCountQuery := countQuery + remainingQuery + commonPaginatedQueryPart
		selectQuery = fmt.Sprintf("%s,(%s) ", selectQuery, totalCountQuery)
		finalQuery = selectQuery + remainingQuery + commonPaginatedQueryPart + orderByClause + limitOffsetQueryPart

	}
	return finalQuery
}

func BuildQueryForArtifactsForCdStage(listingFilterOptions bean.ArtifactsListFilterOptions) string {
	// expected result -> will fetch all successfully deployed  artifacts ar parent stage plus its own stage. Along with this it will
	// also fetch all artifacts generated by plugin at pre_cd or post_cd process (will use data_source in ci artifact table for this)

	if listingFilterOptions.UseCdStageQueryV2 {
		return buildQueryForArtifactsForCdStageV2(listingFilterOptions)
	}

	commonQuery := " from ci_artifact LEFT JOIN cd_workflow ON ci_artifact.id = cd_workflow.ci_artifact_id" +
		" LEFT JOIN cd_workflow_runner ON cd_workflow_runner.cd_workflow_id=cd_workflow.id " +
		" Where (((cd_workflow_runner.id in (select MAX(cd_workflow_runner.id) OVER (PARTITION BY cd_workflow.ci_artifact_id) FROM cd_workflow_runner inner join cd_workflow on cd_workflow.id=cd_workflow_runner.cd_workflow_id))" +
		" AND ((cd_workflow.pipeline_id= %v and cd_workflow_runner.workflow_type = '%v' ) OR (cd_workflow.pipeline_id = %v AND cd_workflow_runner.workflow_type = '%v' AND cd_workflow_runner.status IN ('Healthy','Succeeded') )))" +
		" OR (ci_artifact.component_id = %v  and ci_artifact.data_source= '%v' ))" +
		" AND (ci_artifact.image LIKE '%v' )"

	commonQuery = fmt.Sprintf(commonQuery, listingFilterOptions.PipelineId, listingFilterOptions.StageType, listingFilterOptions.ParentId, listingFilterOptions.ParentStageType, listingFilterOptions.ParentId, listingFilterOptions.PluginStage, listingFilterOptions.SearchString)
	if len(listingFilterOptions.ExcludeArtifactIds) > 0 {
		commonQuery = commonQuery + fmt.Sprintf(" AND ( ci_artifact.id NOT IN (%v))", helper.GetCommaSepratedString(listingFilterOptions.ExcludeArtifactIds))
	}

	totalCountQuery := "SELECT COUNT(DISTINCT ci_artifact.id) as total_count " + commonQuery
	selectQuery := fmt.Sprintf("SELECT DISTINCT(ci_artifact.id) , (%v) ", totalCountQuery)
	//GroupByQuery := " GROUP BY cia.id "
	limitOffSetQuery := fmt.Sprintf(" order by ci_artifact.id desc LIMIT %v OFFSET %v", listingFilterOptions.Limit, listingFilterOptions.Offset)

	//finalQuery := selectQuery + commonQuery + GroupByQuery + limitOffSetQuery
	finalQuery := selectQuery + commonQuery + limitOffSetQuery
	return finalQuery
}

func buildQueryForArtifactsForCdStageV2(listingFilterOptions bean.ArtifactsListFilterOptions) string {
	whereCondition := fmt.Sprintf(" WHERE (id IN ("+
		" SELECT DISTINCT(cd_workflow.ci_artifact_id) as ci_artifact_id "+
		" FROM cd_workflow_runner"+
		" INNER JOIN cd_workflow ON cd_workflow.id = cd_workflow_runner.cd_workflow_id "+
		" AND (cd_workflow.pipeline_id = %d OR cd_workflow.pipeline_id = %d)"+
		"    WHERE ("+
		"            (cd_workflow.pipeline_id = %d AND cd_workflow_runner.workflow_type = '%s')"+
		"            OR"+
		"            (cd_workflow.pipeline_id = %d"+
		"                AND cd_workflow_runner.workflow_type = '%s'"+
		"                AND cd_workflow_runner.status IN ('Healthy','Succeeded')"+
		"           )"+
		"      )   ) ", listingFilterOptions.PipelineId, listingFilterOptions.ParentId, listingFilterOptions.PipelineId, listingFilterOptions.StageType, listingFilterOptions.ParentId, listingFilterOptions.ParentStageType)

	whereCondition = fmt.Sprintf(" %s OR (ci_artifact.component_id = %d  AND ci_artifact.data_source= '%s' ))", whereCondition, listingFilterOptions.ParentId, listingFilterOptions.PluginStage)
	if listingFilterOptions.SearchString != EmptyLikeRegex {
		whereCondition = whereCondition + fmt.Sprintf(" AND ci_artifact.image LIKE '%s' ", listingFilterOptions.SearchString)
	}
	if len(listingFilterOptions.ExcludeArtifactIds) > 0 {
		whereCondition = whereCondition + fmt.Sprintf(" AND ( ci_artifact.id NOT IN (%s))", helper.GetCommaSepratedString(listingFilterOptions.ExcludeArtifactIds))
	}

	selectQuery := fmt.Sprintf(" SELECT ci_artifact.* ,COUNT(id) OVER() AS total_count " +
		" FROM ci_artifact")
	ordeyByAndPaginated := fmt.Sprintf(" ORDER BY id DESC LIMIT %d OFFSET %d ", listingFilterOptions.Limit, listingFilterOptions.Offset)
	finalQuery := selectQuery + whereCondition + ordeyByAndPaginated
	return finalQuery
}

func BuildQueryForArtifactsForRollback(listingFilterOptions bean.ArtifactsListFilterOptions) string {
	commonQuery := " FROM cd_workflow_runner cdwr " +
		" INNER JOIN cd_workflow cdw ON cdw.id=cdwr.cd_workflow_id " +
		" INNER JOIN ci_artifact cia ON cia.id=cdw.ci_artifact_id " +
		" WHERE cdw.pipeline_id=%v AND cdwr.workflow_type = '%v' "

	commonQuery = fmt.Sprintf(commonQuery, listingFilterOptions.PipelineId, listingFilterOptions.StageType)
	if listingFilterOptions.SearchString != EmptyLikeRegex {
		commonQuery += fmt.Sprintf(" AND cia.image LIKE '%v' ", listingFilterOptions.SearchString)
	}
	if len(listingFilterOptions.ExcludeWfrIds) > 0 {
		commonQuery = fmt.Sprintf(" %s AND cdwr.id NOT IN (%s)", commonQuery, helper.GetCommaSepratedString(listingFilterOptions.ExcludeWfrIds))
	}
	totalCountQuery := " SELECT COUNT(cia.id) as total_count " + commonQuery
	orderByQuery := " ORDER BY cdwr.id DESC "
	limitOffsetQuery := fmt.Sprintf("LIMIT %v OFFSET %v", listingFilterOptions.Limit, listingFilterOptions.Offset)
	finalQuery := fmt.Sprintf(" SELECT cdwr.id as cd_workflow_runner_id,cdwr.triggered_by,cdwr.started_on,cia.*,(%s) ", totalCountQuery) + commonQuery + orderByQuery + limitOffsetQuery
	return finalQuery
}