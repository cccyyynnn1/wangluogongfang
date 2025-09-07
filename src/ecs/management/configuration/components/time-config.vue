<script lang="ts">
  export default {
    name: 'TimeConfig',
  }
</script>
<script setup lang="ts">
  import { setSystemTimeApi, setNtpdateApi, getNtpdateApi, checkVerificationApi } from '@/api-ecs/system'
  import { solar2lunar } from '@/plugins/VabCalendar'
  import { useNow } from '@vueuse/core'
  import dayjs from 'dayjs'
  import { ElInput } from 'element-plus'
  import AesEncryptCBC from '~/src/utils/crypto'

  const $baseMessage: any = inject('$baseMessage')
  const weekDay = ['星期天', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  const now = useNow()
  const curTime = computed(() => dayjs(now.value).format('HH : mm : ss'))
  const curDate = computed(() => dayjs(now.value).format('YYYY年MM月DD日'))
  const curWeek = computed(() => weekDay[now.value.getDay()])
  const lunarDate = computed(() => {
    const { gzYear, IMonthCn, IDayCn } = solar2lunar(dayjs(now.value).format('YYYY-MM-DD'))
    return `${gzYear}年${IMonthCn}${IDayCn}`
  })

  const formData = reactive({
    type: 0,
    address: '',
    date: '',
    enable: false,
  })

  const submitHandle = () => {
    if (formData.type === 0 && !formData.date) {
      return $baseMessage('请选择服务器时间', 'error', 'vab-hey-message-error')
    }
    if (formData.type === 1 && !formData.address) {
      return $baseMessage('请填写NTP服务器地址', 'error', 'vab-hey-message-error')
    }
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
            const { data } = await checkVerificationApi({ password: AesEncryptCBC(password.value) })
            if (data?.status) {
              const request =
                formData.type === 0
                  ? setSystemTimeApi(formData.date)
                  : setNtpdateApi({ server: formData.address, enable: true })
              request.then(({ msg }) => $baseMessage(msg, 'success', 'vab-hey-message-success'))
            }
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

  onMounted(() => {
    getNtpdateApi().then(({ data: { enable, server } }) => {
      formData.address = server || ''
      formData.enable = enable || false
    })
  })
</script>

<template>
  <div class="time-config">
    <el-form ref="formRef" label-position="right" label-width="140px" :model="formData">
      <el-row>
        <el-col :span="5" />
        <el-col :span="12">
          <el-form-item prop="type">
            <el-radio-group v-model="formData.type">
              <el-radio :label="0">手动</el-radio>
              <el-radio :label="1">NTP时间</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item v-if="!formData.type" prop="date">
            <div class="time-container">
              <div class="time-panel">
                <div class="top"><span>当前时间</span></div>
                <div class="bottom">
                  <div class="left">{{ curTime }}</div>
                  <div class="right">
                    {{ curDate }} {{ curWeek }}
                    <br />
                    {{ lunarDate }}
                  </div>
                </div>
              </div>
              <el-date-picker
                v-model="formData.date"
                placeholder="选择日期时间"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </div>
          </el-form-item>
          <template v-else>
            <el-form-item label="服务器地址：" prop="address">
              <el-input v-model="formData.address" />
            </el-form-item>
            <el-form-item label="启用Ntp同步：" prop="enable">
              <el-switch v-model="formData.enable" active-text="开启" inactive-text="关闭" inline-prompt />
            </el-form-item>
          </template>
          <el-form-item>
            <el-button type="primary" @click="submitHandle">保存</el-button>
          </el-form-item>
        </el-col>
        <el-col :span="9" />
      </el-row>
    </el-form>
  </div>
</template>

<style lang="scss" scoped>
  .time-config {
    padding-top: 20px;

    .time-panel {
      padding: 5px 25px;
      margin-bottom: 20px;
      height: 102px;
      width: 400px;
      border: 1px solid var(--el-border-color);
      .top {
        height: 40px;
        line-height: 40px;
      }

      .bottom {
        display: flex;
        align-items: center;
        justify-content: space-between;

        .left {
          font-size: 34px;
          font-weight: bold;
          color: #303133;
        }

        .right {
          line-height: 17px;
        }
      }
    }
    .time-container {
      display: flex;
      flex-direction: column;
    }
  }
</style>
