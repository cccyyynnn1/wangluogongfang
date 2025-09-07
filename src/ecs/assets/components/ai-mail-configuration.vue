<script lang="ts">
  export default {
    name: 'AiMailConfiguration',
  }
  interface AiMailConfigurationProps {
    modelValue: boolean
    retrieveIndexType?: number
  }
</script>

<script setup lang="ts">
  import {
    getInnerMailApi,
    deleteInnerMailApi,
    addInnerMailApi,
    editInnerMailApi,
    getMailConfigApi,
    updateMailConfigApi,
    deleteMailConfigApi,
    getMailAlgoConfigPageApi,
    updateMailAlgoConfigPageApi,
    getEnterpriseMailApi,
  } from '@/api-ecs/alert'
  import { MailConfigItem } from '@/types'
  import { ElMessageBox, ElInput } from 'element-plus'
  import { h } from 'vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import type { FormInstance } from 'element-plus'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  let oldMailAddress = ''
  const multipleSelection = ref<{ id: number; mailAddress: string }[]>([])
  const domainMail = reactive({
    pageNum: 1,
    pageSize: 10,
    total: 0,
    list: [] as {
      id: number
      mailAddress: string
      edit?: boolean
    }[],
  })
  const mailConfiguration = reactive({
    pageNum: 1,
    pageSize: 10,
    total: 0,
    list: [] as MailConfigItem[],
  })
  const props = withDefaults(defineProps<AiMailConfigurationProps>(), {
    modelValue: false,
    retrieveIndexType: 33,
  })
  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
  }>()
  //邮箱域名配置
  const visible = useVModel(props, 'modelValue', emits)
  const configurationVisible = ref(false)
  const configurationType = ref('domain')
  const configurationRef = ref<FormInstance>()
  const configurationForm = reactive({
    address: '',
    name: '',
    password: '',
    remark: '',
    enable: false,
  } as MailConfigItem)

  const algorithmForm = reactive({
    id: undefined,
    type: 'default',
    param: '',
  })

  const enterpriseMails = ref<{ name: string; used: boolean }[]>([])

  const handleGetInnerMail = async () => {
    const { data } = await getInnerMailApi({
      pageNum: domainMail.pageNum,
      pageSize: domainMail.pageSize,
    })
    domainMail.list = data.records
    domainMail.total = data.total
    multipleSelection.value = []
  }
  function handleSelectionChange(
    val: {
      id: number
      mailAddress: string
    }[]
  ) {
    multipleSelection.value = val
  }
  const handleEmailValidate = (val: string) => {
    const isMail = /@[\w-]+(\.[\w-]+)+$/
    const address = val
      .replace(/\s+/g, '')
      .replaceAll('，', ',')
      .split(',')
      .filter((i) => isMail.test(i))
      .join(',')
    return address
  }
  const handleEmailDomainEdit = async (data: { id: number; mailAddress: string; edit?: boolean }) => {
    const { id, mailAddress: address } = data
    const { msg } = await editInnerMailApi({ id, address })
    data.edit = false
    $baseMessage(msg, 'success', 'vab-hey-message-success')
  }
  const handleEmailDomainAdd = () => {
    const mails = ref('')
    ElMessageBox({
      title: '新增内部邮箱',
      confirmButtonText: '保存',
      customClass: 'sharedMessageBox',
      message: () =>
        h('div', null, [
          h(ElInput, {
            modelValue: mails.value,
            type: 'textarea',
            placeholder: '请输入内部邮箱域名',
            rows: 5,
            style: 'margin-top: 10px',
            resize: 'none',
            'onUpdate:modelValue': (val: string) => {
              mails.value = val
            },
          }),
        ]),
      beforeClose: async (action, instance, done) => {
        if (action === 'confirm') {
          const { msg } = await addInnerMailApi({ address: mails.value })
          domainMail.pageNum = 1
          handleGetInnerMail()
          $baseMessage(msg, 'success', 'vab-hey-message-success')
          done()
        } else {
          done()
        }
      },
    }).catch(() => {})
  }
  const handleEmailDomainDelete = (ids: number[]) => {
    $baseConfirm('确认删除当前内部邮箱么？', null, async () => {
      const { msg } = await deleteInnerMailApi({ ids, deleteAll: !ids.length })
      domainMail.pageNum === 1 ? handleGetInnerMail() : (domainMail.pageNum = 1)
      $baseMessage(msg, 'success', 'vab-hey-message-success')
    })
  }
  const handleCancelEmailDomainEdit = (data: { id: number; mailAddress: string; edit?: boolean }) => {
    data.mailAddress = oldMailAddress
    oldMailAddress = ''
    data.edit = false
  }

  //邮箱配置
  const handleEmailSave = () => {
    configurationRef.value?.validate(async (valid) => {
      if (valid) {
        const { msg } = await updateMailConfigApi(configurationForm)
        mailConfiguration.pageNum = 1
        handleGetMailConfigurationList()
        handleGetEnterpriseMail()
        configurationVisible.value = false
        $baseMessage(msg, 'success', 'vab-hey-message-success')
      }
    })
  }
  const handleEnableChange = async (row: any) => {
    Object.assign(configurationForm, row)
    const { msg } = await updateMailConfigApi(configurationForm)
    Object.assign(configurationForm, {
      id: undefined,
      address: '',
      name: '',
      password: '',
      remark: '',
      enable: false,
    })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
  }
  const handleEmailDelete = (ids: number[]) => {
    $baseConfirm('确认删除当前邮箱么？', null, async () => {
      const { msg } = await deleteMailConfigApi({ ids, deleteAll: !ids.length })
      mailConfiguration.pageNum = 1
      handleGetMailConfigurationList()
      $baseMessage(msg, 'success', 'vab-hey-message-success')
    })
  }
  const handleEmailEdit = (
    data = {
      id: undefined,
      address: '',
      name: '',
      password: '',
      remark: '',
      enable: false,
    }
  ) => {
    Object.assign(configurationForm, data)
    configurationVisible.value = true
  }
  const handleGetMailConfigurationList = async () => {
    Object.assign(configurationForm, {
      id: undefined,
      address: '',
      name: '',
      password: '',
      remark: '',
      enable: false,
    })
    const { data } = await getMailConfigApi({
      pageNum: mailConfiguration.pageNum,
      pageSize: mailConfiguration.pageSize,
    })
    mailConfiguration.list = data.records
    mailConfiguration.total = data.total
    multipleSelection.value = []
  }

  //邮箱算法

  const handleGetEnterpriseMail = async () => {
    const { data } = await getEnterpriseMailApi()
    enterpriseMails.value = data || []
  }
  const handleGetMailAlgoConfigPage = async () => {
    const { data } = await getMailAlgoConfigPageApi()
    Object.assign(algorithmForm, data)
  }
  const handleMailAlgorithmSave = async () => {
    algorithmForm.param = algorithmForm.type === 'default' ? '' : algorithmForm.param
    const { msg } = await updateMailAlgoConfigPageApi(algorithmForm)
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    handleGetMailAlgoConfigPage()
  }
  watch(
    () => configurationVisible.value,
    () => {
      if (!configurationVisible.value) {
        Object.assign(configurationForm, {
          id: undefined,
          address: '',
          name: '',
          password: '',
          remark: '',
          enable: false,
        })
      }
    }
  )

  watch(
    () => configurationType.value,
    (newType) => {
      if (newType === 'domain') {
        handleGetInnerMail()
      } else if (newType === 'email') {
        handleGetEnterpriseMail()
        handleGetMailConfigurationList()
      } else if (newType === 'algorithm') {
        handleGetMailAlgoConfigPage()
      }
    },
    {
      immediate: true,
    }
  )
</script>

<template>
  <vab-dialog v-model="visible" destroy-on-close :width="1000">
    <template #header>
      <el-select v-model="configurationType" placeholder="Select" style="width: 240px">
        <el-option label="邮箱配置" value="email" />
        <el-option label="邮箱域名配置" value="domain" />
        <el-option label="算法配置" value="algorithm" />
      </el-select>
    </template>
    <div class="emailConfigurationBox">
      <template v-if="configurationType === 'domain'">
        <el-button type="primary" @click="handleEmailDomainAdd">添加</el-button>
        <el-dropdown style="margin-left: 10px">
          <span class="el-dropdown-link">
            <el-button type="danger">
              批量删除
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleEmailDomainDelete(multipleSelection.map((item) => item.id))">
                删除选中
              </el-dropdown-item>
              <el-dropdown-item @click="handleEmailDomainDelete([])">删除所有</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-table
          class-name="mail-table"
          :data="domainMail.list"
          style="height: 450px"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="序号" type="index" width="75" />
          <el-table-column label="邮箱地址" prop="mail">
            <template #default="{ row }">
              <span v-if="!row.edit">{{ row.mailAddress }}</span>
              <el-input
                v-else
                v-model="row.mailAddress"
                placeholder="请输入邮箱地址"
                @keyup.enter="handleEmailDomainEdit(row)"
              />
            </template>
          </el-table-column>
          <el-table-column align="center" fixed="right" label="操作" width="140">
            <template #default="{ row }">
              <template v-if="row.edit">
                <el-button size="small" @click="() => handleEmailDomainEdit(row)">保存</el-button>
                <el-button size="small" @click="() => handleCancelEmailDomainEdit(row)">取消</el-button>
              </template>
              <template v-else>
                <el-button
                  size="small"
                  @click="
                    () => {
                      row.edit = true
                      oldMailAddress = row.mailAddress
                    }
                  "
                >
                  编辑
                </el-button>
                <el-button size="small" @click="() => handleEmailDomainDelete([row.id])">删除</el-button>
              </template>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="domainMail.pageNum"
          v-model:page-size="domainMail.pageSize"
          background
          layout="total, sizes, prev, pager, next, jumper"
          :page-num-sizes="[10, 20, 30]"
          :total="domainMail.total"
          @current-change="handleGetInnerMail"
          @size-change="handleGetInnerMail"
        />
      </template>
      <template v-else-if="configurationType === 'email'">
        <el-button type="primary" @click="handleEmailEdit()">添加</el-button>
        <el-dropdown style="margin-left: 10px">
          <span class="el-dropdown-link">
            <el-button type="danger">
              批量删除
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleEmailDelete(multipleSelection.map((item) => item.id))">
                删除选中
              </el-dropdown-item>
              <el-dropdown-item @click="handleEmailDelete([])">删除所有</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-table
          class-name="mail-table"
          :data="mailConfiguration.list"
          style="height: 450px"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="序号" type="index" width="75" />
          <el-table-column label="邮箱名称" prop="name" show-overflow-tooltip />
          <el-table-column label="邮箱地址" prop="address" show-overflow-tooltip />
          <el-table-column label="是否启用" prop="enable">
            <template #default="{ row }">
              <el-switch
                v-model="row.enable"
                active-text="是"
                inactive-text="否"
                inline-prompt
                @change="handleEnableChange(row)"
              />
            </template>
          </el-table-column>
          <el-table-column label="自定义参数" prop="remark" />
          <el-table-column align="center" fixed="right" label="操作" width="140">
            <template #default="{ row }">
              <el-button size="small" @click="handleEmailEdit(row)">编辑</el-button>
              <el-button size="small" @click="() => handleEmailDelete([row.id])">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="mailConfiguration.pageNum"
          v-model:page-size="mailConfiguration.pageSize"
          background
          layout="total, sizes, prev, pager, next, jumper"
          :page-num-sizes="[10, 20, 30]"
          :total="mailConfiguration.total"
          @current-change="handleGetMailConfigurationList"
          @size-change="handleGetMailConfigurationList"
        />
      </template>
      <template v-else-if="configurationType === 'algorithm'">
        <el-form
          label-position="left"
          label-width="80"
          :model="algorithmForm"
          style="width: 60%; margin: 0 auto; height: 536px"
        >
          <el-form-item label="算法类型" prop="type">
            <el-radio-group v-model="algorithmForm.type">
              <el-radio label="default" size="large">默认</el-radio>
              <el-radio label="custom" size="large">自定义</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item v-if="algorithmForm.type === 'custom'" label="算法名称" prop="param">
            <el-input v-model.trim="algorithmForm.param" />
          </el-form-item>
          <el-form-item label="&nbsp;">
            <el-button type="primary" @click="handleMailAlgorithmSave">保存</el-button>
          </el-form-item>
        </el-form>
      </template>
      <!-- 邮箱配置Form -->
      <vab-dialog v-model="configurationVisible" destroy-on-close title="邮箱配置">
        <el-form
          ref="configurationRef"
          label-position="left"
          label-width="80"
          :model="configurationForm"
          :rules="{
            name: [
              {
                required: true,
                message: '请输入邮箱名称',
                trigger: 'blur',
              },
            ],
            address: [
              {
                required: true,
                message: '请输入邮箱地址',
                trigger: 'blur',
              },
            ],
            password: [
              {
                required: true,
                message: '请输入邮箱密码',
                trigger: 'blur',
              },
            ],
          }"
          style="width: 75%; margin: 0 auto"
        >
          <el-form-item label="邮件名称" prop="name">
            <el-select
              v-model="configurationForm.name"
              placeholder="请选择邮箱"
              style="width: 100%"
              :teleported="false"
            >
              <el-option
                v-for="mail in enterpriseMails"
                :key="mail"
                :disabled="mail.used"
                :label="mail.name"
                :value="mail.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="邮件地址" prop="address">
            <el-input v-model.trim="configurationForm.address" />
          </el-form-item>
          <el-form-item label="密码信息" prop="password">
            <el-input v-model.trim="configurationForm.password" clearable show-password type="password" />
          </el-form-item>
          <el-form-item label="自定义参数" prop="remark">
            <el-input v-model.trim="configurationForm.remark" />
          </el-form-item>
          <el-form-item label="是否启用" prop="enable">
            <el-switch v-model="configurationForm.enable" active-text="是" inactive-text="否" inline-prompt />
          </el-form-item>
          <el-form-item label="&nbsp;">
            <el-button type="primary" @click="handleEmailSave">保存</el-button>
          </el-form-item>
        </el-form>
      </vab-dialog>
    </div>
  </vab-dialog>
</template>

<style scoped lang="scss">
  .emailConfigurationBox {
    .mail-table,
    .algorithm-table {
      margin-top: 10px;
    }
  }
  :deep() {
    .el-radio.el-radio--large {
      height: 30px;
    }
  }
</style>
