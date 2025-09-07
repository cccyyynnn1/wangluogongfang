<script lang="ts">
  export default {
    name: 'VabPage', // es分页器
  }
</script>

<script setup lang="ts">
  import dayjs from 'dayjs'

  const props = defineProps<{
    currentPage: number
    endTime: string
    pageSize: number
    total: number
    indexType: number
    list: any[]
  }>()

  const emits = defineEmits<{
    (e: 'current-change', data: number): void
    (e: 'size-change', data: number): void
  }>()

  const threshold = 1000
  /**
   * @description '触发即将获取下一个max-size（threshold）的数组'
   */
  let triggerArr: number[] = []
  let curIndex = 0
  const pageNum = ref(1)
  const pageSize = ref(20)
  let endTimeList: any[] = [] // 跳页的endTime数组
  let rangeField = '' // 取的row数据的字段

  watch(
    () => props.currentPage,
    () => {
      pageNum.value = props.currentPage
    }
  )

  onMounted(() => {
    pageNum.value = props.currentPage
    pageSize.value = props.pageSize
    // 初始化字段
  })
  const lastPageBtnNode = ref()
  // 控制翻页器的显示样式
  const getLastPageBtnNode = () => {
    // initData(1)
    const nodes = document.querySelectorAll('.el-pager li')
    if (nodes.length > 5) {
      lastPageBtnNode.value = nodes[nodes.length - 1] as HTMLElement
      lastPageBtnNode.value.setAttribute('style', 'display: none')
      if (pageNum.value >= Math.ceil(props.total / pageSize.value) - 5) {
        lastPageBtnNode.value?.setAttribute('style', 'display: block')
      } else {
        lastPageBtnNode.value?.setAttribute('style', 'display: none')
      }
    }
  }

  const initData = (num: number) => {
    triggerArr = []
    for (let index = 5; index > 0; index--) {
      const pageCount = threshold / pageSize.value - index + num
      triggerArr.push(pageCount)
    }
  }

  /**
   * @description '清空跳页的endTime数组'
   */
  const clearEndTimeList = () => {
    endTimeList = []
  }

  /**
   * @description '得到实时的endTime'
   */
  const getRealTimeEndTime = () => {
    recordingTime()
    let realTimeEndTime = undefined
    const curTimeObj = endTimeList[curIndex]
    let _value: any = undefined
    if (curTimeObj) {
      const _name = Object.keys(curTimeObj)[0]
      _value = curTimeObj[_name]
    }
    realTimeEndTime = _value || props.endTime
    return realTimeEndTime
  }
  /**
   * @description '获取实时的pageNum'
   */
  const getRealTimePageNum = () => {
    recordingTime()
    let realTimePageNum = 1
    const curTimeObj = endTimeList[curIndex]
    let _page = 0
    if (curTimeObj) {
      const _name = Object.keys(curTimeObj)[0]
      _page = pageNum.value - +_name + 1
    }
    realTimePageNum = _page
    return realTimePageNum
  }

  /**
   * @description '每次翻页重新记录endTime'
   */
  const recordingTime = () => {
    initField()
    const currentChangeFalg = ref(false)
    if (pageNum.value == 1) {
      endTimeList[0] = { [pageNum.value]: props.endTime }
      initData(1)
    }
    const max = triggerArr[4] // 当前阈值内最后一页
    const min = triggerArr[0] // 当前阈值内倒数第五页
    // 当前页码小于最小触发条件
    if (pageNum.value < min) {
      if (currentChangeFalg && curIndex != 0) {
        // 不在当前阈值（threshold）内，获取上一阈值（threshold）数据
        if (pageNum.value <= max - threshold / pageSize.value) {
          curIndex--
        }
        currentChangeFalg.value = false
        const curTimeObj = endTimeList[curIndex]
        if (curTimeObj) {
          const _name = Object.keys(curTimeObj)[0]
          initData(+_name)
        }
        // 这里是重新判断其是否在最后5页
        if (min <= pageNum.value && pageNum.value <= max) {
          currentChangeFalg.value = true
          const item = props.list[props.list.length - 1]
          const time = item[rangeField]
          // 记录触发页码
          endTimeList[curIndex + 1] = { [pageNum.value + 1]: dayjs(+time / 1000000).format('YYYY-MM-DD HH:mm:ss') }
        }
      }
      // 当前页码满足触发条件在最后5页
    } else if (min <= pageNum.value && pageNum.value <= max) {
      currentChangeFalg.value = true
      const item = props.list[props.list.length - 1]
      const time = item[rangeField]
      // 记录触发页码
      endTimeList[curIndex + 1] = { [pageNum.value + 1]: dayjs(+time / 1000000).format('YYYY-MM-DD HH:mm:ss') }
      // 当前页码大于触发条件
    } else if (pageNum.value > max) {
      // 获取下一阈值（threshold）数据
      curIndex++
      const curTimeObj = endTimeList[curIndex]
      if (curTimeObj) {
        const _name = Object.keys(curTimeObj)[0]
        initData(+_name)
      }
    }
  }

  const MAIL_RANGE_FIELD = 'startTimeNs'
  const HTTP_RANGE_FIELD = 'requestTimeNs'
  const OTHER_RANGE_FIELD = 'startTimeNs'
  const EVETN_RANGE_FIELD = 'flowBeginTimeNs'
  const initField = () => {
    if (props.indexType) {
      switch (props.indexType) {
        case 1:
          rangeField = HTTP_RANGE_FIELD
          break
        case 22:
          rangeField = HTTP_RANGE_FIELD
          break
        case 6:
          rangeField = MAIL_RANGE_FIELD
          break
        case 11:
          rangeField = EVETN_RANGE_FIELD
          break
        default:
          rangeField = OTHER_RANGE_FIELD
          break
      }
    }
  }

  defineExpose({
    clearEndTimeList,
    // recordingTime,
    getRealTimeEndTime,
    getRealTimePageNum,
  })

  watch(
    () => props.total,
    () => {
      nextTick(() => {
        getLastPageBtnNode()
      })
    }
  )

  // 页容量改变
  const handleSizeChange = (val: number) => {
    pageSize.value = val
    pageNum.value = 1
    getLastPageBtnNode()
    emits('size-change', pageSize.value)
  }

  // 页面改变
  const handleCurrentChange = (val: number) => {
    pageNum.value = val
    getLastPageBtnNode()
    emits('current-change', pageNum.value)
  }
</script>

<template>
  <el-pagination
    v-model:current-page="pageNum"
    v-model:page-size="pageSize"
    background
    layout="total, sizes, prev, pager, next"
    :page-sizes="[20, 50, 100]"
    :total="total"
    @current-change="handleCurrentChange"
    @size-change="handleSizeChange"
  />
</template>

<style scoped lang="scss"></style>
