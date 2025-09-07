<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import EditBasicConfiguration from './components/edit-basic-configuration.vue'
  import { FormInstance } from 'element-plus'
  import { initialSystemApi, updateSystemDynamic } from '@/api-ecs/system'
  import { getKeyVarInfoApi } from '@/api-ecs/retrieve'
  import { handleTabs } from '~/src/utils/routes'
  import AesEncryptCBC from '~/src/utils/crypto'
  const $baseMessage: any = inject('$baseMessage')
  const initializeVisible = ref(false)
  const initializeFormRef = ref<FormInstance>()
  const mode = ref('yunche_ras')
  const showPage = ref(false)
  const queryData = reactive({
    ip: '',
    user: '',
    password: '',
  })
  const varInfo = reactive({
    yunche_ras: {
      value: {
        addr: '',
        account: '',
        password: '',
        key: '',
      },
      key: '',
      name: '',
      id: '' as string | number,
    },
    ecs_kafka: {
      value: [
        {
          addr: '',
          account: '',
          password: '',
          key: '',
        },
      ],
      key: '',
      name: '',
      id: '' as string | number,
    },
    ecs_storage: {
      value: {
        addr: '',
        account: '',
        password: '',
        key: '',
      },
      key: '',
      name: '',
      id: '' as string | number,
    },
    ecs_neo4j: {
      value: {
        addr: '',
        account: '',
        password: '',
        key: '',
      },
      key: '',
      name: '',
      id: '' as string | number,
    },
  })

  const rules = {
    ip: [{ required: true, message: '请输入流量平台IP', trigger: 'change' }],
    user: [{ required: true, message: '请输入账号', trigger: 'change' }],
    password: [{ required: true, message: '请输入密码', trigger: 'change' }],
  }
  // 获取平台地址（PCAP保存使用）
  async function handleQueryByKey() {
    const { data } = await getKeyVarInfoApi()
    Object.keys(data).forEach((item) => {
      data[item].value = JSON.parse(data[item].value) || JSON.parse(data[item].value)
      // @ts-ignore
      varInfo[item] = data[item]
    })
  }
  function initializeHandle() {
    initializeVisible.value = true
    queryData.ip = varInfo.yunche_ras.value.addr
  }

  // 初始化
  function initialize() {
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
            if (!initializeFormRef.value) return
            initializeFormRef!.value.validate(async (valid) => {
              if (valid) {
                const { msg } = await initialSystemApi(queryData, { password: AesEncryptCBC(password.value) })
                $baseMessage(msg, 'success', 'vab-hey-message-success')
                initializeVisible.value = false
              }
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

  async function submit() {
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
            const values = JSON.stringify(varInfo.yunche_ras.value)
            const { msg } = await updateSystemDynamic(
              { ...varInfo, value: values },
              { password: AesEncryptCBC(password.value) }
            )
            $baseMessage(msg, 'success', 'vab-hey-message-success')
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

  const handleEdit = (val: string) => {
    mode.value = val
    showPage.value = true
  }

  const closeEvent = () => {
    showPage.value = false
  }

  onMounted(() => {
    handleQueryByKey()
  })
</script>

<script lang="ts">
  export default {
    name: 'EquipmentManagement',
  }
</script>

<template>
  <div class="article-container">
    <div class="content">
      <vab-query-form>
        <vab-query-form-left-panel :span="12">
          <h3>基础配置</h3>
        </vab-query-form-left-panel>
      </vab-query-form>
      <el-row style="margin: 100px auto">
        <el-col :offset="6" :span="14">
          <el-form class="my-form" label-position="right" label-width="100px" @submit.prevent>
            <el-form-item label="数据节点">
              <el-input
                v-model="varInfo.yunche_ras.value.addr"
                clearable
                :disabled="true"
                placeholder="请输入流量平台IP"
                style="width: 60%"
              />
              <div style="min-width: 100px">
                <el-button link style="margin-left: 15px" type="primary" @click="handleEdit('yunche_ras')">
                  编辑
                </el-button>
                <el-button link type="primary" @click="initializeHandle">初始化</el-button>
              </div>
            </el-form-item>
            <el-form-item label="KAFKA">
              <el-input
                v-model="varInfo.ecs_kafka.value[0].addr"
                clearable
                :disabled="true"
                placeholder="请输入流量平台IP"
                style="width: 60%"
              />
              <el-button link style="margin-left: 15px" type="primary" @click="handleEdit('ecs_kafka')">编辑</el-button>
            </el-form-item>
            <el-form-item label="存储库">
              <el-input
                v-model="varInfo.ecs_storage.value.addr"
                clearable
                :disabled="true"
                placeholder="请输入流量平台IP"
                style="width: 60%"
              />
              <el-button link style="margin-left: 15px" type="primary" @click="handleEdit('ecs_storage')">
                编辑
              </el-button>
            </el-form-item>
            <el-form-item label="图数据库">
              <el-input
                v-model="varInfo.ecs_neo4j.value.addr"
                clearable
                :disabled="true"
                placeholder="请输入流量平台IP"
                style="width: 60%"
              />
              <el-button link style="margin-left: 15px" type="primary" @click="handleEdit('ecs_neo4j')">编辑</el-button>
            </el-form-item>
            <el-form-item>
              <!-- <el-button type="primary" @click="submit">保存</el-button> -->
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
    </div>
    <vab-dialog v-model="initializeVisible" destroy-on-close title="内网IP范围" width="600px">
      <el-form ref="initializeFormRef" inline label-width="110px" :model="queryData" :rules="rules" @submit.prevent>
        <el-form-item label="流量平台IP" prop="ip">
          <el-input v-model="queryData.ip" clearable placeholder="请输入流量平台IP" />
        </el-form-item>
        <el-form-item label="账号" prop="user">
          <el-input v-model="queryData.user" clearable placeholder="请输入账号" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="queryData.password" clearable placeholder="请输入密码" type="password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-space style="width: 100%; justify-content: center; margin-bottom: 20px">
          <el-button type="primary" @click="initialize">初始化</el-button>
        </el-space>
      </template>
    </vab-dialog>
    <EditBasicConfiguration
      v-if="showPage"
      :current-data="varInfo"
      :mode="mode"
      :show-page="showPage"
      @on-close-event="closeEvent"
      @on-reflash="handleQueryByKey"
    />
  </div>
</template>
<style>
  .need-password-message-box .el-message-box__header {
    border-bottom: 1px solid var(--el-border-color);
  }
  .need-password-message-box .el-message-box__btns {
    flex-direction: row-reverse !important;
    justify-content: normal !important;
    .el-button {
      margin-left: 10px;
    }
  }
</style>
<style lang="scss" scoped>
  .content {
    h3 {
      margin-block: 0 0.5em;
    }
    .my-form {
      margin: 0px auto;
    }
    // display: flex;
    // align-items: center;
    // justify-content: center;
    // flex-direction: column;
    // .btn {
    //   width: 100%;
    //   display: flex;
    //   align-items: center;
    //   justify-content: center;
    // }
  }
  :deep() {
    .el-form-item__content {
      // margin-left: 20px;
      .el-input {
        width: 400px;
      }
    }
  }
</style>
