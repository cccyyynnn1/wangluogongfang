<script lang="ts">
  export default {
    name: 'ParamsConfig',
  }
</script>

<script setup lang="ts">
  import { Search } from '@element-plus/icons-vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { ElMessageBox, ElInput } from 'element-plus'
  import { getAllCustomConfigApi, getParamsByBelongApi, updateParamsByBelongApi } from '~/src/api-ecs/system'
  import AesEncryptCBC from '~/src/utils/crypto'
  const $baseMessage: any = inject('$baseMessage')

  const types = ref()
  const currType = ref('')
  const fieldDetailVisible = ref(false)
  const editForm = ref()
  const tableData = ref()
  const handleGetAllCustomConfigApi = async () => {
    const { data } = await getAllCustomConfigApi()
    types.value = data || []
    currType.value = types.value[0] || ''
  }
  const searchStr = ref('')
  const handelGetParamsByBelong = async () => {
    const { data } = await getParamsByBelongApi(currType.value, searchStr.value)
    tableData.value = data || []
  }
  const handleEditParam = (data: {
    belong: string
    createTime: number
    id: number
    param: string
    remark: string
    value: string
  }) => {
    editForm.value = data
    fieldDetailVisible.value = true
  }
  const handleEditParamSubmit = () => {
    if (!editForm.value?.value) return $baseMessage('请填写参数值！', 'error', 'vab-hey-message-error')
    const { belong, param, remark, value } = editForm.value
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
            const { code } = await updateParamsByBelongApi({
              custom: JSON.stringify({
                belong,
                param,
                remark,
                value,
              }),
              password: AesEncryptCBC(password.value),
            })
            fieldDetailVisible.value = false
            instance.confirmButtonLoading = false
            $baseMessage('编辑成功', 'success', 'vab-hey-message-success')
            handelGetParamsByBelong()
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
  watch(
    () => currType.value,
    () => {
      handelGetParamsByBelong()
    }
  )

  onMounted(() => {
    handleGetAllCustomConfigApi()
  })
</script>

<template>
  <div class="params-config-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>统一参数配置</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-input
          v-model="searchStr"
          class="search"
          clearable
          placeholder="请输入检索条件"
          style="width: 300px"
          @clear="handelGetParamsByBelong"
        />
        <el-button :icon="Search" style="margin-left: -2px" type="primary" @click="handelGetParamsByBelong">
          检索
        </el-button>
      </vab-query-form-right-panel>
    </vab-query-form>
    <div class="config-content">
      <div class="config-nav">
        <div
          v-for="type in types"
          :key="type"
          class="config-nav-item"
          :class="{ select: currType === type }"
          @click="currType = type"
        >
          {{ type }}
        </div>
      </div>
      <div class="config-list">
        <el-table :data="tableData">
          <el-table-column label="序号" type="index" width="80" />
          <el-table-column label="参数名称" prop="param" />
          <el-table-column
            :formatter="({ value, desensitization }) => (desensitization ? '**********' : value)"
            label="值"
            prop="value"
          />
          <el-table-column label="备注" prop="remark" />
          <el-table-column label="操作" width="80">
            <template #default="{ row }">
              <el-button size="small" @click="() => handleEditParam(row)">编辑</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
    <vab-dialog v-model="fieldDetailVisible" destroy-on-close title="编辑" width="680px">
      <div>
        <el-form ref="formRef" label-position="top">
          <el-form-item label="参数名" prop="param">
            <el-input v-model="editForm.param" :disabled="true" />
          </el-form-item>
          <el-form-item label="值" prop="value">
            <el-input v-model="editForm.value" :type="editForm?.desensitization ? 'password' : 'text'" />
          </el-form-item>
          <el-form-item label="备注" prop="remark">
            <el-input v-model="editForm.remark" type="textarea" />
          </el-form-item>
          <el-form-item label="">
            <el-button class="history-add" type="primary" @click="handleEditParamSubmit">确定</el-button>
            <el-button class="history-add" plain @click="fieldDetailVisible = false">取消</el-button>
          </el-form-item>
        </el-form>
      </div>
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss">
  .params-config-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    :deep {
      .el-form-item__content {
        justify-content: end;
        .el-button {
          width: 68px;
          height: 34px;
          border-radius: 6px;
          margin-top: 40px;
        }
      }
      .el-dialog__footer:empty {
        display: none;
      }
    }
    h3 {
      margin-block: 0;
    }
    .config-content {
      flex: 1;
      display: flex;
      margin-top: 8px;
      .config-nav {
        width: 216px;
        height: 100%;
        padding: 14px 0;
        border-radius: 16px;
        border: 1px solid #f1f0ff;
        .config-nav-item {
          width: 184px;
          height: 40px;
          border-radius: 4px;
          font-weight: 600;
          font-size: 15px;
          line-height: 40px;
          margin: 0 auto 4px;
          text-indent: 10px;
          cursor: pointer;
          &.select,
          &:hover {
            background: #efeefe;
            color: #5236ff;
          }
        }
      }
      .config-list {
        flex: 1;
        height: 100%;
        :deep() {
          .el-table {
            width: calc(100% - 20px);
            height: calc(100vh - 94px);
            margin-left: auto;
            --el-table-row-hover-bg-color: transparent;
          }
          .el-table--striped .el-table__body tr.el-table__row--striped td.el-table__cell {
            background-color: #f8f9ff;
          }
          .el-table .el-table__header th {
            background-color: #f2f2ff !important;
            .cell {
              font-weight: 500 !important;
              font-size: 14px !important;
              color: #918da5 !important;
            }
          }
        }
      }
    }
  }
</style>
