/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package actions

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/auditcontroller/audit/logics"
	"encoding/json"
	"io/ioutil"
	"strconv"

	restful "github.com/emicklei/go-restful"
)

var procAudit *procAuditAction = &procAuditAction{}

// ObjectAction
type procAuditAction struct {
	base.BaseAction
}

func init() {
	procAudit.CreateAction()

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/proc/{owner_id}/{biz_id}/{user}", Params: nil, Handler: procAudit.AddLog})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/procs/{owner_id}/{biz_id}/{user}", Params: nil, Handler: procAudit.AddLogs})

	// set cc api resource
}

//操作日志
func (a *procAuditAction) AddLog(req *restful.Request, resp *restful.Response) {
	type paramsStruct struct {
		Content interface{}            `json:"content"`
		OpDesc  string                 `json:"op_desc"`
		OpType  auditoplog.AuditOpType `json:"op_type"`
		ProcID  int                    `json:"inst_id"`
	}
	language := util.GetActionLanguage(req)
	defErr := a.CC.Error.CreateDefaultCCErrorIf(language)

	ownerID := util.GetActionOnwerID(req)
	strAppID := req.PathParameter("biz_id")
	appID, _ := strconv.Atoi(strAppID)
	user := req.PathParameter("user")
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request boody error:%s", err.Error())
		procAudit.ResponseFailed(common.CCErrCommHTTPReadBodyFailed, defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error(), resp)
		return
	}
	params := paramsStruct{}
	err = json.Unmarshal([]byte(value), &params)
	if err != nil {
		blog.Error("json unmarshal failed,input:%v error:%v", string(value), err)
		procAudit.ResponseFailed(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	logics.DB = appAudit.CC.InstCli
	err = logics.AddLogWithStr(appID, params.ProcID, params.OpType, common.BKInnerObjIDProc, params.Content, "", params.OpDesc, ownerID, user)
	if nil != err {
		blog.Error("json unmarshal failed,input:%v error:%v", string(value), err)
		procAudit.ResponseFailed(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	} else {
		procAudit.ResponseSuccess(nil, resp)
		return
	}

}

//插入多行主机操作日志型操作
func (a *procAuditAction) AddLogs(req *restful.Request, resp *restful.Response) {
	type paramsStruct struct {
		Content []auditoplog.AuditLogContext `json:"bk_content"`
		OpDesc  string                       `json:"bk_op_desc"`
		OpType  auditoplog.AuditOpType       `json:"bk_op_type"`
	}
	language := util.GetActionLanguage(req)
	defErr := a.CC.Error.CreateDefaultCCErrorIf(language)

	ownerID := util.GetActionOnwerID(req)
	strAppID := req.PathParameter("biz_id")
	appID, _ := strconv.Atoi(strAppID)
	user := req.PathParameter("user")
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request boody error:%s", err.Error())
		procAudit.ResponseFailed(common.CCErrCommHTTPReadBodyFailed, defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error(), resp)
		return
	}
	params := paramsStruct{}
	err = json.Unmarshal([]byte(value), &params)
	if err != nil {
		blog.Error("json unmarshal failed,input:%v error:%v", string(value), err)
		procAudit.ResponseFailed(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	logics.DB = appAudit.CC.InstCli
	err = logics.AddLogMulti(appID, params.OpType, common.BKInnerObjIDProc, params.Content, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Error("json unmarshal failed,input:%v error:%v", string(value), err)
		procAudit.ResponseFailed(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	} else {
		procAudit.ResponseSuccess(nil, resp)
		return
	}

}
