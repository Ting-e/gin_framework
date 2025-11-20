/*
 * @Author: your name
 * @Date: 2022-03-10 16:54:52
 * @LastEditTime: 2022-07-07 16:42:09
 * @LastEditors: 三丰 1209490572@qq.com
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: \newUser：12007\api\api.go
 */
package api

import "project/web/handler/model"

type APIService interface {

	//测试接口
	GetDatas(req *model.GetDatasReq) *model.GetDatasResp
}
