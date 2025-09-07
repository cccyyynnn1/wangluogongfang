import { Ref } from 'vue'

import VueEvent from '@/data/event'

import { useCopy } from '@/utils'

const useContextMenu = (): {
  x: Ref<number>
  y: Ref<number>
  isShow: Ref<boolean>
  list: Ref<{ label: string; callback: Function; value: string }[]>
} => {
  const list = ref([{ label: '复制', callback: useCopy, value: '' }])

  const element = ref()

  const x = ref(0) // 横坐标

  const y = ref(0) // 纵坐标

  const rowValue = ref()

  const isShow = ref(false) // 是否显示

  const getCellData = (data: any) => {
    // if (element.value && element.value != data.cell) {
    //   reset()
    // }W
    element.value = data.cell
    if (data.list) {
      list.value = data.list
    } else {
      list.value = [{ label: '复制', callback: useCopy, value: '' }]
      // @ts-ignore
      list.value[0].value = data.value
    }
    handleContext(data.event)
    element.value?.addEventListener('contextmenu', handleContext)
    window.addEventListener('click', handleClose, true)
    window.addEventListener('contextmenu', handleClose, true)
  }

  // 激活菜单事件
  const handleContext = (e: any) => {
    e.preventDefault()
    e.stopPropagation()
    isShow.value = true
    x.value = e.clientX + 20
    y.value = e.clientY - 15
  }

  const handleClose = () => {
    isShow.value = false
  }

  // 监听事件
  VueEvent.on('useContextMenu', getCellData)

  onMounted(() => {
    nextTick(() => {
      element?.value?.addEventListener('contextmenu', handleContext)
      window.addEventListener('click', handleClose, true)
      window.addEventListener('contextmenu', handleClose, true)
    })
  })

  watch(
    () => isShow.value,
    (newVal) => {
      if (!newVal) {
        reset()
      }
    }
  )

  const reset = () => {
    element?.value?.removeEventListener('contextmenu', handleContext)
    window.removeEventListener('click', handleClose, true)
    window.removeEventListener('contextmenu', handleClose, true)
  }

  return { x, y, isShow, list }
}

export default useContextMenu
