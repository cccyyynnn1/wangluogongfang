<script lang="ts">
  export default {
    name: 'AbnormalLoginRules',
  }
</script>

<script setup lang="ts">
  import { useVModel } from '@vueuse/core'
  import { requireRules } from '~/src/utils/rules'
  import GLOBALCOUNTRYCITY from '@/assets/addressJson.js'
  import { AbnormalLandingAlarmSaveEntityType } from '@/types/alert'
  import { FormInstance } from 'element-plus'
  import { AbnormalLandingAlarmSaveEntityApi, AbnormalLandingAlarmUpdateEntityApi } from '@/api-ecs/alert'
  import dayjs from 'dayjs'
  import _lodash from 'lodash'
  import { log } from 'console'
  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    visible: boolean
    ruleData?: any
  }>()

  const emit = defineEmits<{
    (e: 'update:visible', visible: boolean): void
  }>()

  type optionType = { value: string; label: string }

  type address = {
    id: string
    pid: string
    path: string
    level: string
    name: string
    name_en: string
    name_pinyin: string
    code: string
    childrens: address[]
  }

  const asset_options = ref<optionType[]>([])

  const provinces_options = ref<address[]>([])

  const city_options = ref<address[]>([])

  const countryRef = ref()
  const provinceRef = ref()
  const cityRef = ref()

  const list = ref<any[]>([])

  const dialogVisible = useVModel(props, 'visible', emit)

  const formData = reactive<AbnormalLandingAlarmSaveEntityType>({
    name: '',
    serverIp: '',
    assetArray: '',
    timeStatus: '每天',
    incident: '',
    startTime: '00:00',
    endTime: '23:59',
    country: '',
    province: '',
    city: '',
  })
  // IPV4
  const ipv4Regex =
    /^(?!0)(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])$/
  // IPV6
  const ipv6Regex =
    /^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$/
  const validatePass = (rule: any, value: string, callback: any) => {
    if (value === '') {
      callback(new Error('该选项为必填项'))
    } else {
      const arr = value.split(',')
      const flag = arr.every((item: any) => {
        return ipv4Regex.test(item) || ipv6Regex.test(item)
      })
      if (!flag) {
        callback(new Error('输入不合法'))
      } else {
        if (!formRef.value) return
        // formRef.value.validateField('serverIp', () => null)
        callback()
      }
    }
  }
  // 表单数据校验规则
  const rules = reactive({
    name: requireRules,
    serverIp: [
      { validator: validatePass, trigger: 'blur' },
      { required: true, message: '该选项为必填项' },
    ],
  })

  const handleChageCountry = () => {
    formData.province = ''
    formData.city = ''
    provinces_options.value = []
    city_options.value = []
  }
  const handleChageProvince = () => {
    formData.city = ''
    city_options.value = []
  }
  watch(
    () => formData.country,
    () => {
      GLOBALCOUNTRYCITY.forEach((item) => {
        if (item.name == formData.country) {
          provinces_options.value = item.childrens
        }
      })
    },
    { immediate: true }
  )
  watch(
    () => formData.province,
    () => {
      provinces_options.value.forEach((item) => {
        if (item.name == formData.province) {
          city_options.value = item.childrens
        }
      })
    },
    { immediate: true }
  )
  const handleChageCity = () => {}

  // const dateOption = ref<null | string>('everyDay')

  const timeRange = ref()
  const timeRangePicker = ref()
  const formRef = ref<FormInstance>()
  const visiblePicker = ref(false)

  const pickerNode = ref()
  const node = ref()
  const optionsClick = () => {
    visiblePicker.value = true
    setTimeout(() => {
      pickerNode.value.focus()
      if (visiblePicker.value) {
        const picker = document.querySelector('.myself-date-picker')
        node.value = picker?.querySelector('.el-date-picker')
        const pickerWarp = document.querySelector('.picker-warp')
        pickerWarp?.appendChild(node.value)
        picker?.parentNode?.removeChild(picker)
        timeRangePicker.value = _lodash.cloneDeep(timeRange.value)
      }
    }, 0)
  }

  const handleDel = (item: any) => {
    timeRangePicker.value = timeRangePicker.value.filter((td: any) => {
      return td != item
    })
  }

  const clearData = () => {
    formData.name = ''
    formData.serverIp = ''
    formData.id = undefined
    formData.startTime = ''
    formData.endTime = ''
    formData.country = ''
    formData.province = ''
    formData.city = ''
    timeRange.value = null
    timeRangePicker.value = null
    formData.timeStatus = ''
    list.value = []
  }

  const handleConfirm = () => {
    visiblePicker.value = false
    timeRange.value = _lodash.cloneDeep(timeRangePicker.value)
  }

  const formatData = (date: Date) => {
    if (date) {
      return dayjs(date).format('YYYY-MM-DD')
    }
  }

  const submit = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate(async (valid) => {
      if (valid) {
        try {
          countryRef.value.blur()
          provinceRef.value.blur()
          cityRef.value.blur()
          const url = !formData.id ? AbnormalLandingAlarmSaveEntityApi : AbnormalLandingAlarmUpdateEntityApi
          if (formData.timeStatus == '自定义' && timeRange.value.length > 0) {
            const arr: string[] = []
            timeRange.value.forEach((_: any, index: number) => {
              arr.push(dayjs(timeRange.value[index]).format('YYYY-MM-DD'))
            })
            formData.timeStatus = arr.join(',')
          }
          const { msg } = await url({ ...formData, incident: list.value.join(',') })
          clearData()
          $baseMessage(msg, 'success', 'vab-hey-message-success')
          dialogVisible.value = false
        } catch (e) {
          console.log(e)
        }
      } else {
        return false
      }
    })
  }

  watchEffect(() => {
    if (dialogVisible.value) {
      if (props.ruleData) {
        formData.id = props.ruleData.id
        for (const key in formData) {
          if (key == 'incident') {
            formData.incident = props.ruleData.incident
            if (formData.incident) {
              list.value = formData.incident.split(',')
            }
          } else if (
            key == 'timeStatus' &&
            props.ruleData.timeStatus &&
            props.ruleData.timeStatus !== '每天' &&
            props.ruleData.timeStatus !== '工作日'
          ) {
            timeRange.value = _lodash.cloneDeep(props.ruleData.timeStatus).split(',')
            formData.timeStatus = '自定义'
          } else {
            // @ts-ignore
            formData[key] = props.ruleData[key]
          }
        }
      } else {
        clearData()
      }
    }
  })
</script>

<template>
  <el-dialog v-model="dialogVisible" :title="`${ruleData ? '编辑' : '新增'}规则`" width="765px">
    <el-form ref="formRef" label-position="top" label-width="120px" :model="formData" :rules="rules">
      <el-form-item label="规则名称" prop="name">
        <el-input v-model="formData.name" maxlength="30" placeholder="请输入" show-word-limit type="text" />
      </el-form-item>
      <el-form-item label="目的IP" prop="serverIp">
        <el-input v-model="formData.serverIp" placeholder="请输入（输入多个时用','隔开）" />
      </el-form-item>
      <!-- <el-form-item label="资产组" prop="assetArray">
        <el-select v-model="formData.assetArray" placeholder="请选择" style="width: 100%">
          <el-option v-for="item in asset_options" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item> -->
      <el-form-item label="时间">
        <el-time-select
          v-model="formData.startTime"
          class="mr-4"
          end="24:00"
          :max-time="formData.endTime"
          placeholder="开始时间"
          start="00:00"
          step="00:15"
          style="width: 160px"
        />
        <span style="margin: 0 11px">至</span>
        <el-time-select
          v-model="formData.endTime"
          end="23:59"
          :min-time="formData.startTime"
          placeholder="结束时间"
          start="00:00"
          step="00:15"
          style="width: 160px"
        />
        <el-select ref="selectNode" v-model="formData.timeStatus" style="width: 110px; margin-left: 11px">
          <el-option label="每天" value="每天" @click.stop.prevent="" />
          <el-option label="工作日" value="工作日" />
          <el-option label="自定义" value="自定义" @click="optionsClick" @click.stop.prevent="" />
        </el-select>
      </el-form-item>
      <el-form-item label="地区">
        <el-row :gutter="10" style="width: 100%; margin: 0">
          <el-col :span="8" style="padding-left: 0">
            <el-select
              ref="countryRef"
              v-model="formData.country"
              filterable
              placeholder="请选择"
              style="width: 100%"
              @change="handleChageCountry"
            >
              <el-option v-for="item in GLOBALCOUNTRYCITY" :key="item.name_en" :label="item.name" :value="item.name" />
            </el-select>
          </el-col>
          <el-col :span="8">
            <el-select
              ref="provinceRef"
              v-model="formData.province"
              filterable
              placeholder="请选择"
              style="width: 100%"
              @change="handleChageProvince"
            >
              <el-option v-for="item in provinces_options" :key="item.name_en" :label="item.name" :value="item.name" />
            </el-select>
          </el-col>
          <el-col :span="8" style="padding-right: 0">
            <el-select
              ref="cityRef"
              v-model="formData.city"
              filterable
              placeholder="请选择"
              style="width: 100%"
              @change="handleChageCity"
            >
              <el-option v-for="item in city_options" :key="item.name_en" :label="item.name" :value="item.name" />
            </el-select>
          </el-col>
        </el-row>
      </el-form-item>
      <el-form-item label="适配事件" prop="incident">
        <el-checkbox-group v-model="list">
          <el-checkbox label="HTTP登录" value="Value A" />
          <el-checkbox label="数据库登录" value="Value B" />
          <el-checkbox label="RDP登录" value="Value C" />
          <el-checkbox label="SSH登录" value="Value disabled" />
          <el-checkbox label="邮箱登录" value="Value selected and disabled" />
        </el-checkbox-group>
      </el-form-item>
    </el-form>
    <el-dialog
      v-model="visiblePicker"
      class="myself-dialog"
      :close-on-click-modal="false"
      :destroy-on-close="true"
      title="自定义时间"
      width="715"
    >
      <el-date-picker
        v-if="visiblePicker"
        ref="pickerNode"
        v-model="timeRangePicker"
        :close-on-click-modal="false"
        placeholder="选择日期"
        popper-class="myself-date-picker"
        type="dates"
      />
      <div class="picker-warp">
        <div class="content">
          <div class="content-title">已选日期</div>
          <div class="items">
            <div v-for="(item, index) in timeRangePicker" :key="index">
              <div class="item">
                {{ formatData(item) }}
                <el-icon style="color: #9b9b9b; cursor: pointer" @click="handleDel(item)"><Close /></el-icon>
              </div>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button type="primary" @click="handleConfirm">确认</el-button>
        <el-button @click="visiblePicker = false">取消</el-button>
      </template>
    </el-dialog>
    <template #footer>
      <el-button type="primary" @click="submit(formRef)">确认</el-button>
      <el-button @click="dialogVisible = false">取消</el-button>
    </template>
  </el-dialog>
</template>

<style lang="scss">
  .myself-dialog {
    .el-input__wrapper {
      display: none;
    }
    .picker-warp {
      height: 445px;
      margin-top: -30px;
      position: relative;
      .content {
        padding: 20px;
        position: absolute;
        top: -0px;
        height: 445px;
        width: 334px;
        right: 0;
        border-radius: 6px;
        border: 1px solid #eeedfe;

        .content-title {
          font-size: 14px;
          color: #9d9baa;
          position: absolute;
        }
        .items {
          margin-top: 15px;
          height: 400px;
          overflow-y: auto;
          &::-webkit-scrollbar {
            width: 0;
            height: 0;
          }
        }
        .item {
          width: 295px;
          height: 34px;
          background: #f3f2ff;
          border-radius: 6px;
          padding: 0 15px;
          margin-top: 10px;
          color: #6954f0;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }
      }
    }
    .el-date-picker {
      border: 1px solid #eeedfe;
      border-radius: 6px;
      height: 445px;
    }
    .el-picker-panel__footer {
      display: none;
    }
  }
  .myself-date-picker {
    box-shadow: inherit !important;
    .el-date-picker {
      border: 1px solid var(--el-border-color);
    }
    .el-popper__arrow {
      display: none;
    }
  }
</style>
