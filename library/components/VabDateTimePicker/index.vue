<script lang="ts">
  export default {
    name: 'VabDatePicker',
  }
</script>

<script setup lang="ts">
  import { useVModel } from '@vueuse/core'
  import locale from 'ant-design-vue/es/date-picker/locale/zh_CN'
  import { DatePicker } from 'ant-design-vue'
  import dayjs from 'dayjs'
  import 'dayjs/locale/zh-cn'
  dayjs.locale('zh-cn')

  interface Props {
    modelValue?: [string, string] | [dayjs.Dayjs, dayjs.Dayjs]
    /**
     * 是否显示清除按钮
     */
    allowClear?: boolean
  }
  const props = withDefaults(defineProps<Props>(), {
    modelValue: undefined,
    /**
     * 是否显示清除按钮
     */
    allowClear: false,
  })
  const emit = defineEmits(['update:modelValue', 'change'])
  const datePickerVal = useVModel(props, 'modelValue', emit)
  const box = ref()
  const handleChange = (dates: [dayjs.Dayjs, dayjs.Dayjs] | [string, string], dateStrings: [string, string]) => {
    emit('change', dates)
  }
</script>

<template>
  <div ref="box" style="width: 100%">
    <DatePicker.RangePicker
      v-model:value="datePickerVal"
      :allow-clear="allowClear"
      v-bind="$attrs"
      :get-popup-container="() => box"
      :locale="locale"
      show-time
      value-format="YYYY-MM-DD HH:mm:ss"
      @change="handleChange"
    />
  </div>
</template>

<style scoped lang="scss"></style>
