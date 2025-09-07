<script lang="ts">
  export default {
    name: 'DateInfo', // 日历信息维护
  }
</script>

<script setup lang="ts">
  import { downloadTextFile } from '~/src/utils/download'
  import { CalendartemplateApi, FindAllCalendarApi, CalendarImportApi } from '~/src/api-ecs/system'
  import { requireRules } from '~/src/utils/rules'
  import { uuid } from '~/src/utils'
  import UploadDate from './upload-date.vue'

  // 表单数据
  type objType = {
    [key: string]: any
  }
  const formData = reactive<objType>({
    logType: 1,
    // logType: ['攻击日志'],
    addr: '',
    post: '',
    enable: true,
    isSave: false,
  })
  const rules = reactive({
    addr: requireRules,
    post: requireRules,
    // logType: requireRules,
  })

  const showUpload = ref(false)

  const handleDownload = async () => {
    const res = await CalendartemplateApi()
    downloadTextFile(res, '日历导入模板')
    ElMessage({ message: '下载成功', type: 'success' })
  }

  onMounted(() => {
    initData()
  })

  const currentYear = ref()
  const currentMonth = ref()
  const yearOption = ref([])
  const monthOption = ref([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12])
  const dayOption = ref(['日', '一', '二', '三', '四', '五', '六'])
  const allDate = ref<any[]>([])
  const curDate = ref()

  const initData = async () => {
    const currentTime = new Date()
    currentYear.value = currentTime.getFullYear()
    currentMonth.value = currentTime.getMonth() + 1
  }

  const initCurDate = () => {
    curDate.value = []
    // 某年某月数据
    curDate.value = allDate.value?.filter((item: any) => {
      return item.month == currentMonth.value && item.annual == currentYear.value
    })
    if (curDate.value?.length == 0) return
    // 某年某月1号是星期几
    const dayOfWeek = new Date(currentYear.value, currentMonth.value - 1, 1).getDay()
    const prex = dayOfWeek - 0
    const next = 35 - curDate.value.length - prex
    const prexMouth = currentMonth.value == 1 ? 12 : currentMonth.value - 1
    const nextMouth = currentMonth.value == 12 ? 1 : currentMonth.value + 1
    for (let index = 0; index < prex; index++) {
      const obj = {
        annual: prexMouth < currentMonth.value ? currentYear.value : currentYear.value - 1,
        id: uuid(),
        month: prexMouth,
        name: `${prexMouth}-${31 - index}`,
        status: 2,
        value: '',
      }
      curDate.value.unshift(obj)
    }
    for (let index = 0; index < next; index++) {
      const obj = {
        annual: nextMouth > currentMonth.value ? currentYear.value : currentYear.value + 1,
        id: uuid(),
        month: nextMouth,
        name: `${nextMouth}-${index + 1}`,
        status: 2,
        value: '',
      }
      curDate.value.push(obj)
    }
  }

  const skipToToday = () => {
    const currentTime = new Date()
    currentYear.value = currentTime.getFullYear()
    currentMonth.value = currentTime.getMonth() + 1
  }

  const formatName = (val: string) => {
    const arr = val.split('-')
    return arr[1]
  }

  const handleUpload = () => {
    showUpload.value = true
  }

  watch(
    () => currentMonth.value,
    () => {
      if (allDate.value?.length == 0) return
      initCurDate()
    }
  )
  watch(
    () => currentYear.value,
    async () => {
      const { data } = await FindAllCalendarApi({ year: currentYear.value })
      yearOption.value = data?.years
      allDate.value = data?.data
      if (allDate.value?.length == 0) return
      initCurDate()
    }
  )

  const handleReflash = async () => {
    const { data } = await FindAllCalendarApi({ year: currentYear.value })
    yearOption.value = data?.years
    allDate.value = data?.data
    if (allDate.value?.length == 0) return
    initCurDate()
  }
</script>

<template>
  <div class="date-info">
    <div class="content">
      <el-row>
        <el-col :span="5" />
        <el-col :span="19">
          <el-form ref="formRef" label-position="right" label-width="150px" :model="formData" :rules="rules">
            <el-form-item label="下载模版：">
              <el-button style="background-color: #6954f0" type="primary" @click="handleDownload">下载</el-button>
            </el-form-item>
            <el-form-item label="上传日历：">
              <el-button style="background-color: #6954f0" type="primary" @click="handleUpload">上传</el-button>
            </el-form-item>
            <el-form-item label=" ">
              <div class="calendar">
                <div class="tools">
                  <div class="tools-left">
                    <el-select v-model="currentYear" style="width: 90px; margin-right: 10px">
                      <el-option v-for="(item, index) in yearOption" :key="index" :label="item" :value="item" />
                    </el-select>
                    <el-select v-model="currentMonth" style="width: 90px">
                      <el-option v-for="(item, index) in monthOption" :key="index" :label="item" :value="item" />
                    </el-select>
                  </div>
                  <div class="tools-right">
                    <el-button style="border-radius: 4px" @click="skipToToday">今天</el-button>
                  </div>
                </div>
                <div class="calendar-content">
                  <div v-for="(item, index) in dayOption" :key="index" class="calendar-item calendar-week">
                    {{ item }}
                  </div>
                  <div
                    v-for="item in curDate"
                    :key="item.id"
                    class="calendar-item normal"
                    :class="{ 'rest-day': item.status == 0, 'work-day': item.status == 1 }"
                  >
                    {{ formatName(item.name) }}
                    <div
                      v-if="item.status != 2"
                      class="calendar-tips"
                      :class="{ 'rest-calendar-tips': item.status == 0, 'work-calendar-tips': item.status == 1 }"
                    >
                      <span v-if="item.status == 1">班</span>
                      <span v-if="item.status == 0">休</span>
                    </div>
                  </div>
                </div>
              </div>
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
    </div>
    <UploadDate
      v-model:visible="showUpload"
      :import-assets-fnc="CalendarImportApi"
      :title="'上传'"
      @reflash="handleReflash"
    />
  </div>
</template>

<style scoped lang="scss">
  .date-info {
    .content {
      margin-top: 40px;
    }
    .calendar {
      width: 620px;
      height: 498px;
      background: fff;
      border-radius: 16px;
      border: 1px solid #eae7f5;
      overflow: hidden;
      .tools {
        background: #fafaff;
        height: 49px;
        display: flex;
        align-items: center;
        padding: 0 30px;
        justify-content: space-between;
      }
      .calendar-content {
        height: calc(100% - 49px);
        border-top: 1px solid #eae7f5;
        display: flex;
        // align-items: center;
        justify-content: space-around;
        padding: 20px 30px 30px;
        flex-wrap: wrap;
        overflow: hidden;
      }
      .calendar-item {
        position: relative;
        width: 66px;
        height: 56px;
        display: flex;
        align-items: center;
        margin-right: 8px;
        justify-content: center;
        color: #3c394f;
        font-size: 16px;
        &:nth-child(7n) {
          margin-right: 0px;
        }
      }
      .calendar-week {
        height: 32px !important;
      }
      .calendar-tips {
        position: absolute;
        top: 0;
        right: 0;
        width: 18px;
        height: 18px;
        border-radius: 4px;
        color: #fff;
        display: flex;
        align-items: center;
        font-size: 12px;
        justify-content: center;
      }
      .work-calendar-tips {
        background: #7e79a1;
      }
      .rest-calendar-tips {
        background: #f54646;
      }
      .rest-day {
        background: #fff2f2;
        border-radius: 6px;
        color: #f54646 !important;
      }
      .work-day {
        background: #f7f5ff;
        border-radius: 6px;
        color: #3c394f !important;
      }
      .normal {
        font-size: 24px;
        color: #bdbcc3;
      }
    }
  }
</style>
