/**
 * 后端API调用服务
 * 封装所有后端API请求
 */

import request from '@/utils/request'
import type {
  Seiyuu,
  SeiyuuGroup,
  SeiyuuListResponse,
  SeiyuuDetailResponse,
  GroupListResponse,
  SchedulerRequest,
  SchedulerResponse,
  MoegirlRawDataResponse,
  ProcessProfileRequest,
  ProcessProfileResponse,
  AdminLoginRequest,
  AdminLoginResponse,
  CreateSeiyuuRequest,
  UpdateSeiyuuRequest,
  SeiyuuRelationship,
  CreateRelationshipRequest,
  UpdateRelationshipRequest,
  GenerateRelationshipRequest,
  GenerateRelationshipResponse
} from '@/types'

// ========== 公开声优接口 ==========

/**
 * 获取已发布声优列表
 */
export function getSeiyuuList(): Promise<SeiyuuListResponse> {
  return request.get<Seiyuu[]>('/seiyuu')
}

/**
 * 获取声优详细资料
 */
export function getSeiyuuDetail(id: string): Promise<SeiyuuDetailResponse> {
  return request.get<Seiyuu>(`/seiyuu/${id}`)
}

/**
 * 获取群组配置列表
 */
export function getGroups(): Promise<GroupListResponse> {
  return request.get<SeiyuuGroup[]>('/groups')
}

// ========== AI调度器 ==========

/**
 * 智能选择声优
 */
export function selectSeiyuu(data: SchedulerRequest): Promise<SchedulerResponse> {
  return request.post<SchedulerResponse['data']>('/scheduler/select', data)
}

// ========== 管理员权限接口 ==========

/**
 * 管理员登录
 */
export function adminLogin(data: AdminLoginRequest): Promise<AdminLoginResponse> {
  return request.post<AdminLoginResponse['data']>('/admin/login', data)
}

/**
 * 获取所有声优(含待审核)
 */
export function adminGetAllSeiyuu(): Promise<SeiyuuListResponse> {
  return request.get<Seiyuu[]>('/admin/seiyuu')
}

/**
 * 创建声优
 */
export function adminCreateSeiyuu(data: CreateSeiyuuRequest): Promise<SeiyuuDetailResponse> {
  return request.post<Seiyuu>('/admin/seiyuu', data)
}

/**
 * 更新声优资料
 */
export function adminUpdateSeiyuu(id: string, data: UpdateSeiyuuRequest): Promise<SeiyuuDetailResponse> {
  return request.put<Seiyuu>(`/admin/seiyuu/${id}`, data)
}

/**
 * 删除声优
 */
export function adminDeleteSeiyuu(id: string): Promise<{ success: boolean }> {
  return request.delete(`/admin/seiyuu/${id}`)
}

/**
 * 发布声优(审核通过)
 */
export function adminPublishSeiyuu(id: string): Promise<SeiyuuDetailResponse> {
  return request.post<Seiyuu>(`/admin/publish/${id}`)
}

// ========== 管理员声优创建辅助接口 ==========

/**
 * 获取萌娘百科原始数据
 */
export function getMoegirlData(name: string): Promise<MoegirlRawDataResponse> {
  return request.get<MoegirlRawDataResponse['data']>(`/admin/moegirl/${encodeURIComponent(name)}`)
}

/**
 * AI处理资料转Markdown
 */
export function processProfile(data: ProcessProfileRequest): Promise<ProcessProfileResponse> {
  return request.post<ProcessProfileResponse['data']>('/admin/process-profile', data, {
    timeout: 120000  // 2分钟超时，因为AI处理可能需要较长时间
  })
}

// ========== 群组管理 ==========

/**
 * 获取所有群组
 */
export function adminGetAllGroups(): Promise<GroupListResponse> {
  return request.get<SeiyuuGroup[]>('/admin/groups')
}

/**
 * 创建群组
 */
export function adminCreateGroup(data: Omit<SeiyuuGroup, 'id' | 'created_at' | 'updated_at'>): Promise<{ success: boolean; data: SeiyuuGroup }> {
  return request.post<SeiyuuGroup>('/admin/groups', data)
}

/**
 * 更新群组
 */
export function adminUpdateGroup(id: string, data: Partial<SeiyuuGroup>): Promise<{ success: boolean; data: SeiyuuGroup }> {
  return request.put<SeiyuuGroup>(`/admin/groups/${id}`, data)
}

/**
 * 删除群组
 */
export function adminDeleteGroup(id: string): Promise<{ success: boolean }> {
  return request.delete(`/admin/groups/${id}`)
}

// ========== 声优关系管理 ==========

/**
 * AI生成声优关系
 */
export function adminGenerateRelationship(data: GenerateRelationshipRequest): Promise<GenerateRelationshipResponse> {
  return request.post<GenerateRelationshipResponse['data']>('/admin/relationships/generate', data)
}

/**
 * 创建声优关系
 */
export function adminCreateRelationship(data: CreateRelationshipRequest): Promise<{ success: boolean; data: SeiyuuRelationship }> {
  return request.post<SeiyuuRelationship>('/admin/relationships', data)
}

/**
 * 查询声优关系
 */
export function adminGetRelationship(seiyuuIdA?: string, seiyuuIdB?: string): Promise<{ success: boolean; data: SeiyuuRelationship[] | SeiyuuRelationship }> {
  const params = new URLSearchParams()
  if (seiyuuIdA) params.append('seiyuu_id_a', seiyuuIdA)
  if (seiyuuIdB) params.append('seiyuu_id_b', seiyuuIdB)

  const queryString = params.toString()
  const url = queryString ? `/admin/relationships?${queryString}` : '/admin/relationships'

  return request.get<SeiyuuRelationship[] | SeiyuuRelationship>(url)
}

/**
 * 获取所有声优关系
 */
export function adminGetAllRelationships(): Promise<{ success: boolean; data: SeiyuuRelationship[] }> {
  return request.get<SeiyuuRelationship[]>('/admin/relationships')
}

/**
 * 更新声优关系
 */
export function adminUpdateRelationship(id: string, data: UpdateRelationshipRequest): Promise<{ success: boolean; data: SeiyuuRelationship }> {
  return request.put<SeiyuuRelationship>(`/admin/relationships/${id}`, data)
}

/**
 * 删除声优关系
 */
export function adminDeleteRelationship(id: string): Promise<{ success: boolean }> {
  return request.delete(`/admin/relationships/${id}`)
}

/**
 * 获取声优的所有关系
 */
export function adminGetSeiyuuRelationships(seiyuuId: string): Promise<{ success: boolean; data: SeiyuuRelationship[] }> {
  return request.get<SeiyuuRelationship[]>(`/admin/seiyuu/${seiyuuId}/relationships`)
}

// ========== 系统接口 ==========

/**
 * 健康检查
 */
export function healthCheck(): Promise<{ success: boolean; message: string }> {
  return request.get('/health')
}
