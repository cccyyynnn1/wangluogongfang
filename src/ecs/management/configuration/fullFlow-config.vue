<script lang="ts">
  export default {
    name: 'FullFlowOnfig',
  }
</script>

<script setup lang="ts">
  import UpdateDialog from './components/update-dialog.vue'
  import { useTableCopy } from '@/utils'
  import { checkVerificationApi, systemUpgrade } from '@/api-ecs/system'
  import { getSystemConfigApi, setSystemConfigApi } from '~/src/api-ecs/system'
  import AesEncryptCBC from '~/src/utils/crypto'
  const $baseMessage: any = inject('$baseMessage')
  const listLoading = ref(false) // 是否加载
  const isShowDialog = ref(false) // 是否显示弹框
  const mode = ref('info') // 弹框显示类型
  const description = ref('')
  const varKeyId = 88

  const varVal = reactive({
    updateSite: 'https://172.17.0.1',
    interval: 'day',
    timeUnit: 1,
    enable: false,
  })

  // 表格数据
  const listDate = ref()
  const currVersion = ref()
  // 关闭弹窗
  const closeDialog = (updateVal?: any) => {
    if (mode.value === 'updateConfig' && updateVal) {
      saveVar(updateVal)
    }
    getData()
    isShowDialog.value = false
  }

  // 版本说明
  const showInfo = (val: string) => {
    mode.value = 'info'
    description.value = val
    isShowDialog.value = true
  }

  // 升级配置
  const updateConfig = () => {
    mode.value = 'updateConfig'
    isShowDialog.value = true
  }

  // 升级
  const update = () => {
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
            const { msg } = await checkVerificationApi({ password: AesEncryptCBC(password.value) })
            mode.value = 'update'
            isShowDialog.value = true
            done()
          } catch (error) {
            instance.confirmButtonLoading = false
            done()
          }
        } else {
          done()
        }
      },
    }).catch(() => {})
  }

  const getData = async () => {
    const { data } = await systemUpgrade()
    const cur = (data || [])[0]
    currVersion.value = cur
    listDate.value = data || []
  }

  const getVar = async () => {
    const { data: config } = await getSystemConfigApi({ keys: 'sys_update_config' })
    try {
      const _val = JSON.parse(config.sys_update_config.value)
      varVal.enable = _val?.enable || false
      varVal.updateSite = _val?.updateSite || ''
      varVal.timeUnit = _val?.timeUnit || 1
      varVal.interval = _val?.interval || ''
    } catch (error) {
      console.error(error)
    }
  }

  const saveVar = (saveVar?: any) => {
    const _varVal = JSON.stringify(saveVar)
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
            const { msg } = await setSystemConfigApi(
              {
                id: varKeyId,
                value: _varVal,
              },
              { password: AesEncryptCBC(password.value) }
            )
            $baseMessage('更新成功', 'success', 'vab-hey-message-success')
            done()
          } catch (error) {
            instance.confirmButtonLoading = false
          }
        } else {
          done()
        }
      },
    }).catch(() => {})
    if (saveVar) {
      getVar()
    }
  }

  onMounted(() => {
    getData()
    getVar()
  })
</script>

<template>
  <div class="update-config">
    <el-card>
      <el-tabs :model-value="'authorization'">
        <el-tab-pane label="全流量升级管理" name="authorization">
          <el-row>
            <el-col :span="12">
              <span class="my-tip-1">当前系统版本</span>
              <span class="my-tip-2">{{ currVersion?.version }}</span>
            </el-col>
            <el-col :span="12">
              <span class="my-tip-1">更新时间</span>
              <span class="my-tip-2">{{ currVersion?.date }}</span>
            </el-col>
          </el-row>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    <el-card>
      <div class="top-bar">
        <div class="left">升级历史记录</div>
        <div class="right">
          <el-switch
            v-model="varVal.enable"
            class="auto-update"
            inactive-text="自动升级"
            @change="() => saveVar(varVal)"
          />
          <el-button type="primary" @click="updateConfig">升级配置</el-button>
          <el-button type="primary" @click="update">升级</el-button>
        </div>
      </div>
      <el-table
        v-loading="listLoading"
        class="my-table"
        :data="listDate"
        default-expand-all
        row-key="id"
        :tree-props="{ children: 'children' }"
        @cell-contextmenu="useTableCopy"
      >
        <el-table-column :formatter="() => '系统'" label="版本类型" show-overflow-tooltip />
        <el-table-column label="更新版本" prop="version" show-overflow-tooltip />
        <el-table-column label="更新时间" prop="date" show-overflow-tooltip />
        <el-table-column label="更新结果" prop="updateResult" show-overflow-tooltip>
          <el-button plain type="primary">成功</el-button>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="{ row }">
            <el-button size="small" @click="() => showInfo(row.description)">版本说明</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty class="vab-data-empty" description="暂无数据" />
        </template>
      </el-table>
    </el-card>
    <update-dialog
      v-if="isShowDialog"
      :description="description"
      :is-show-dialog="isShowDialog"
      :mode="mode"
      :var-val="varVal"
      @on-close-event="closeDialog"
    />
  </div>
</template>

<style lang="scss" scoped>
  .update-config {
    :deep() {
      .el-card__body {
        padding: 24px;
      }
      .el-scrollbar {
        height: calc(100vh - 295px);
      }
    }

    .my-tip-1 {
      margin-right: 50px;
      color: #303133;
    }

    .my-tip-2 {
      color: #303133;
      opacity: 0.5;
    }

    .my-bottom {
      margin-top: 20px;
    }

    .my-table {
      margin-top: 20px;
    }

    .auto-update {
      margin-right: 20px;
    }

    .top-bar {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .left {
        color: #178bf6;
        font-weight: 500;
      }
    }
  }
</style>
