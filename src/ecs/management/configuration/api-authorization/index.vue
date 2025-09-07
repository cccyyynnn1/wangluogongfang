<script lang="ts">
  export default {
    name: 'ApiAuthorization',
  }
</script>

<script setup lang="ts">
  import {
    getAuthTokenPageApi,
    getApiDataTreeApi,
    updateAuthTokenApi,
    deleteAuthTokenApi,
    getAuthInfoApi,
  } from '@/api-ecs/system'
  import type { FormInstance, ElTree } from 'element-plus'
  import { AuthToken, AuthTokenPage } from '@/types'
  import { formatTime } from '@/utils/time'
  import AesEncryptCBC from '~/src/utils/crypto'

  let timer: number

  const spacer = h(ElDivider, { direction: 'vertical' })
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')

  const showAuthInfoBtnRef = ref()
  const authInfoRef = ref()
  const authInfoVisible = ref(false)
  const showEdit = ref(false)

  const multipleSelection = ref<number[]>([])
  const editAuthFormRef = ref<FormInstance>()
  const editAuthForm = reactive<AuthToken>({
    id: undefined,
    apiIds: [],
    description: '',
    name: '',
  })
  const authInfo = ref<{
    authId: string
    secret: string
    timestamp: number
    token: string
  }>()
  const rules = reactive({
    apiIds: [{ required: true, message: '请选择授权', trigger: 'change' }],
    name: [
      {
        required: true,
        message: '请输入权限名称',
        trigger: 'change',
      },
    ],
    description: [
      {
        required: true,
        message: '请输入权限描述',
        trigger: 'change',
      },
    ],
  })
  const apiTree = ref<InstanceType<typeof ElTree>>()
  const authApiData = ref()
  const authList = ref<AuthTokenPage[]>([])
  const listTotal = ref(0)
  const queryData = reactive({
    pageNum: 1,
    pageSize: 10,
  })

  const setSelectRows = (val: AuthTokenPage[]) => {
    multipleSelection.value = val.map((i) => i.id) as number[]
  }
  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)
  const handleDeleteAuth = async (ids: number[], deleteAll = false) => {
    const password = ref('')
    ElMessageBox({
      title: '提示',
      showCancelButton: true,
      customClass: 'need-password-message-box',
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      customStyle: {
        maxWidth: '500px',
      },
      message: () =>
        h('div', null, [
          h('div', { style: 'margin: 15px 0 5px 0;color:#55585b' }, '请输入敏感操作密码:(通过验证后方可进行操作)'),
          h(ElInput, {
            type: 'password',
            modelValue: password.value,
            placeholder: '请输入敏感操作密码',
            showPassword: true,
            style: 'margin-block: 10px',
            'onUpdate:modelValue': (val: string) => {
              password.value = val
            },
          }),
        ]),

      beforeClose: async (action, instance, done) => {
        if (action === 'confirm') {
          instance.confirmButtonLoading = true
          try {
            $baseConfirm(!deleteAll ? '你确定要删除当前项吗' : '你确定要删除所有数据吗', null, async () => {
              const { msg } = await deleteAuthTokenApi({ ids, deleteAll }, { password: AesEncryptCBC(password.value) })
              $baseMessage(msg, 'success', 'vab-hey-message-success')
              hanldeGetAuthTokenPage()
            })
            done()
          } catch (error) {
            instance.confirmButtonLoading = false
          }
        } else {
          done()
        }
      },
    }).catch(() => {})
  }
  const handlerInfo = (row: any) => {
    Object.assign(editAuthForm, row)
    editAuthForm.apiIds = row.apiIds.split(',')
    showEdit.value = true
  }
  const hanldeGetAuthTokenPage = async () => {
    const {
      data: { records, total },
    } = await getAuthTokenPageApi(queryData)
    authList.value = records
    listTotal.value = total
  }
  const handleGetAuthApiData = async () => {
    const { data } = await getApiDataTreeApi()
    authApiData.value = data || undefined
  }
  const handleTreeChange = () => {
    editAuthForm.apiIds = apiTree.value?.getCheckedKeys().filter(Boolean) as number[]
  }
  const handleFormatterTime = (row: AuthTokenPage, column: any, val: number) => {
    return formatTime(val)
  }
  const onClickOutside = (event: MouseEvent) => {
    const dom = event.target as HTMLDivElement
    if (dom.innerText === '授权信息') return
    authInfoVisible.value = false
  }
  const handleGetAuthInfo = async (row: AuthTokenPage, event: MouseEvent) => {
    authInfo.value = undefined
    showAuthInfoBtnRef.value = event.target
    authInfoVisible.value = true
    const { data } = await getAuthInfoApi({ id: row.id as number })
    authInfo.value = data
  }

  const submitForm = async () => {
    await editAuthFormRef.value?.validate(async (valid) => {
      if (!valid) return
      const password = ref('')
      ElMessageBox({
        title: '提示',
        showCancelButton: true,
        customClass: 'need-password-message-box',
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        customStyle: {
          maxWidth: '500px',
        },
        message: () =>
          h('div', null, [
            h('div', { style: 'margin: 15px 0 5px 0;color:#55585b' }, '请输入敏感操作密码:(通过验证后方可进行操作)'),
            h(ElInput, {
              type: 'password',
              modelValue: password.value,
              placeholder: '请输入敏感操作密码',
              showPassword: true,
              style: 'margin-block: 10px',
              'onUpdate:modelValue': (val: string) => {
                password.value = val
              },
            }),
          ]),

        beforeClose: async (action, instance, done) => {
          if (action === 'confirm') {
            instance.confirmButtonLoading = true
            try {
              const { id, name, description, apiIds } = editAuthForm
              const { msg } = await updateAuthTokenApi(
                { id, description, name, apiIds },
                { password: AesEncryptCBC(password.value) }
              )
              $baseMessage(msg, 'success', 'vab-hey-message-success')
              showEdit.value = false
              hanldeGetAuthTokenPage()
              done()
            } catch (error) {
              instance.confirmButtonLoading = false
            }
          } else {
            done()
          }
        },
      }).catch(() => {})
    })
  }
  const closeAuthInfo = () => {
    timer = window.setTimeout(() => {
      authInfoVisible.value = false
    }, 1500)
  }
  const showAuthInfo = () => {
    clearInterval(timer)
  }
  const renderContent = (
    h: any,
    {
      data,
    }: {
      data: any
    }
  ) => h('span', null, data.methodRemark || data.controllerRemark)

  watch(showEdit, () => {
    if (!showEdit.value) {
      multipleSelection.value = []
      Object.assign(editAuthForm, { id: undefined, apiIds: [], description: '', name: '' })
    }
  })
  onMounted(async () => {
    hanldeGetAuthTokenPage()
    handleGetAuthApiData()
  })
</script>

<template>
  <div class="web-authorization-container" @click="onClickOutside">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>API授权管理</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-button color="#6954F0" style="margin-right: 10px" @click="showEdit = true">新增</el-button>
        <el-dropdown>
          <span class="el-dropdown-link">
            <el-button type="danger">
              批量删除
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleDeleteAuth(multipleSelection)">删除选中</el-dropdown-item>
              <el-dropdown-item @click="handleDeleteAuth([], true)">删除所有</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-table :data="authList" style="margin-top: 8px; height: calc(100vh - 120px)" @selection-change="setSelectRows">
      <el-table-column type="selection" />
      <el-table-column label="序号" width="60">
        <template #default="{ $index }">
          <span>{{ curIndex + $index }}</span>
        </template>
      </el-table-column>
      <el-table-column label="授权ID" prop="authId" />
      <el-table-column label="名称" prop="name" />
      <el-table-column label="描述" prop="description" />
      <el-table-column :formatter="handleFormatterTime" label="创建时间" prop="createTime" />
      <el-table-column :formatter="handleFormatterTime" label="修改时间" prop="updateTime" />
      <el-table-column fixed="right" label="操作" width="220">
        <template #default="{ row }">
          <el-button class="row_action" size="small" @click="handlerInfo(row)">详情</el-button>
          <el-button class="row_action" size="small" @click="handleGetAuthInfo(row, $event)">授权信息</el-button>
          <el-button class="row_action" :disabled="row.buildIn === 1" size="small" @click="handleDeleteAuth([row.id])">
            删除
          </el-button>
        </template>
      </el-table-column>
      <template #empty><el-empty /></template>
    </el-table>
    <el-pagination
      v-model:current-page="queryData.pageNum"
      v-model:page-size="queryData.pageSize"
      background
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[10, 20, 50, 100]"
      :total="listTotal"
      @current-change="hanldeGetAuthTokenPage"
      @size-change="hanldeGetAuthTokenPage"
    />
    <template v-if="showEdit">
      <div class="authorization-info">
        <el-space :spacer="spacer">
          <el-button :auto-insert-space="false" icon="Back" @click="showEdit = false">返回</el-button>
          <h3>{{ editAuthForm.id ? '编辑' : '新增' }}授权</h3>
        </el-space>
        <el-form ref="editAuthFormRef" label-width="auto" :model="editAuthForm" :rules="rules">
          <el-form-item label="名称" prop="name">
            <el-input v-model="editAuthForm.name" />
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="editAuthForm.description" />
          </el-form-item>
          <el-form-item label="授权权限" prop="apiIds">
            <el-tree
              ref="apiTree"
              :data="authApiData"
              :default-checked-keys="editAuthForm.apiIds"
              node-key="id"
              :render-content="renderContent"
              show-checkbox
              @check="handleTreeChange"
            />
          </el-form-item>
          <el-form-item label="&nbsp;">
            <el-button type="primary" @click="submitForm">确认</el-button>
            <el-button @click="showEdit = false">取消</el-button>
          </el-form-item>
        </el-form>
      </div>
    </template>
    <el-popover
      ref="authInfoRef"
      placement="left"
      :popper-style="{ width: '356px' }"
      :show-arrow="false"
      trigger="click"
      :virtual-ref="showAuthInfoBtnRef"
      virtual-triggering
      :visible="authInfoVisible"
    >
      <div v-loading="!authInfo" class="authInfo" @mouseleave="closeAuthInfo" @mousemove="showAuthInfo">
        <h3>
          <vab-icon icon="bookmark-3-fill" />
          查看授权信息
        </h3>
        <div class="info-content">
          <div>
            <span class="label">authId：</span>
            <span class="text">{{ authInfo?.authId }}</span>
          </div>
          <div>
            <span class="label">timestamp：</span>
            <span class="text">{{ authInfo?.timestamp }}</span>
          </div>
          <div>
            <span class="label">token：</span>
            <span class="text">{{ authInfo?.token }}</span>
          </div>
        </div>
      </div>
    </el-popover>
  </div>
</template>

<style scoped lang="scss">
  .authInfo {
    h3 {
      font-weight: 500;
      font-size: 14px;
      color: #1e1842;
      line-height: 28px;
      margin-block: 0;
    }
    .info-content {
      width: 328px;
      background: #f7f6ff;
      border-radius: 6px;
      border: 1px solid #f1f0ff;
      padding: 14px;
      margin-top: 10px;
      div {
        line-height: 24px;
        display: flex;
        .label {
          font-weight: 500;
          font-size: 14px;
          color: #4a4759;
        }
        .text {
          flex: 1;
          font-weight: 400;
          font-size: 14px;
          color: #494758;
          display: block;
        }
      }
    }
  }
  .web-authorization-container {
    position: relative;
    h3 {
      margin-block: 0;
    }
    .authorization-info {
      position: absolute;
      inset: 0;
      z-index: 10;
      background-color: #fff;
      padding: 10px 20px;
      .el-form {
        width: 50%;
        margin: 0 auto;
      }
      .el-tree {
        width: 100%;
        height: calc(100vh - 220px);
        overflow-y: auto;
      }
    }
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 180px);
      }
    }
  }
</style>
