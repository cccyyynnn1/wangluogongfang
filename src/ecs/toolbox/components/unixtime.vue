<script lang="ts">
  export default {
    name: 'Unixtime',
  }
</script>

<script setup lang="ts">
  import { formatTime } from '@/utils/time'
  import { useClipboard } from '@vueuse/core'
  const timeDurationList = [
    { lable: '秒', val: 's' },
    { lable: '毫秒', val: 'ms' },
    { lable: '纳秒', val: 'ns' },
  ]
  const timeRange: {
    [key: string]: number
  } = {
    s: 1000,
    ms: 1,
    ns: 1 / 1000,
  }
  const timeRange2: {
    [key: string]: number
  } = {
    s: 1000,
    ms: 1,
    ns: 1 / 1000,
  }
  const { now, pause, resume } = useNow({ controls: true, interval: 1000 })

  const { copy, copied } = useClipboard({ legacy: true })

  const cur_time = ref()
  const cur_time2 = ref()
  const cur_time3 = ref()
  const format_time = ref()
  const format_time2 = ref()
  const format_time3 = ref()
  const timeDuration = ref('ms')
  const timeDuration2 = ref('ms')
  const timeDuration3 = ref('ms')

  const getMsTime = (val: string) => {
    return new Date(val)
  }

  const formDate = () => {
    if (!cur_time.value) return (format_time.value = '')
    const val = timeDuration.value === 'ms' ? cur_time.value : cur_time.value * timeRange[timeDuration.value]
    format_time.value = formatTime(val)
  }
  const formDate2 = () => {
    if (!cur_time2.value) return (format_time2.value = '')
    const msTime = getMsTime(cur_time2.value).getTime()
    const val = timeDuration2.value === 'ms' ? msTime : msTime / timeRange2[timeDuration2.value]
    format_time2.value = Math.round(val)
  }
  const formDate3 = () => {
    if (!cur_time3.value) return (format_time3.value = '')
    const msTime = getMsTime(cur_time3.value).getTime()
    const val = timeDuration3.value === 'ms' ? msTime : msTime / timeRange2[timeDuration3.value]
    format_time3.value = Math.round(val)
  }
</script>

<template>
  <div class="unixtime">
    <el-row>
      <el-col :span="12">
        <div class="unixtime-tip">
          现在的Unix时间戳(Unix timestamp)是：
          <span>{{ now.getTime() }}</span>
          <el-link :underline="false" @click="copy(now.getTime().toString())">
            {{ copied ? '已复制' : '复制' }}
          </el-link>
        </div>
      </el-col>
      <el-col :span="12">
        <div style="text-align: end">
          <el-button type="primary" @click="resume">开始</el-button>
          <el-button type="primary" @click="pause">暂停</el-button>
        </div>
      </el-col>
    </el-row>
    <el-form class="unixtime-container" label-width="250px">
      <el-form-item label="Unix时间戳（毫秒）：">
        <el-input v-model="cur_time" :style="{ width: '210px', marginRight: '20px' }" />
        <el-select v-model="timeDuration" :style="{ width: '80px', marginRight: '20px' }">
          <el-option v-for="item in timeDurationList" :key="item.val" :label="item.lable" :value="item.val" />
        </el-select>
        <el-button plain :style="{ marginRight: '20px', width: '145px' }" type="primary" @click="formDate">
          转换
        </el-button>
        <el-input v-model="format_time" :style="{ width: '210px', marginRight: '20px' }" />
      </el-form-item>
      <el-form-item label="时间（年/月/日 时:分:秒）：">
        <el-input v-model="cur_time2" :style="{ width: '310px', marginRight: '20px' }" />
        <el-button plain :style="{ marginRight: '20px' }" type="primary" @click="formDate2">转换成Unix时间戳</el-button>
        <el-input v-model="format_time2" :style="{ width: '210px', marginRight: '20px' }" />
        <el-select v-model="timeDuration2" :style="{ width: '80px' }">
          <el-option v-for="item in timeDurationList" :key="item.val" :label="item.lable" :value="item.val" />
        </el-select>
      </el-form-item>
      <el-form-item label="时间：">
        <el-date-picker
          v-model="cur_time3"
          placeholder="请选择时间"
          :style="{ width: '310px', marginRight: '20px' }"
          type="datetime"
        />
        <el-button plain :style="{ marginRight: '20px' }" type="primary" @click="formDate3">转换成Unix时间戳</el-button>
        <el-input v-model="format_time3" :style="{ width: '210px', marginRight: '20px' }" />
        <el-select v-model="timeDuration3" :style="{ width: '80px' }">
          <el-option v-for="item in timeDurationList" :key="item.val" :label="item.lable" :value="item.val" />
        </el-select>
      </el-form-item>
    </el-form>
  </div>
</template>

<style scoped lang="scss">
  .unixtime {
    &-tip {
      display: flex;
      font-size: 14px;
      font-weight: 500;
      color: #303133;
      align-items: center;
      span {
        display: inline-block;
        width: 150px;
        text-align: center;
        font-size: 16px;
        color: #f23a3a;
      }
    }
    &-container {
      margin-top: 30px;
      height: 475px;
      padding: 50px 0;
      background: #f8fbff;
      :deep(.el-input__wrapper) {
        width: 100%;
      }
    }
  }
</style>
