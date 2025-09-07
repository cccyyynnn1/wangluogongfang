<script lang="ts">
  export default {
    name: 'FomartCode', // 常见编码转换
  }
</script>

<script setup lang="ts">
  import { ref } from 'vue'

  import useCodec from '@/hooks/useCodec'

  const codeMode = ref('')

  const props = defineProps<{
    target: any
    // isRes: boolean
    mode?: undefined
  }>()

  const emit = defineEmits<{
    (e: 'on-change-event', val: any): void
  }>()

  const nowMode = ref()
  const resvalue = ref()
  const reqvalue = ref()

  onMounted(() => {
    setTimeout(() => {
      nowMode.value = props.mode
      resvalue.value = props.target.res
      reqvalue.value = props.target.req
    }, 100)
  })

  watch(
    () => codeMode.value,
    (newVal, oldval) => {
      let res1 = undefined
      let res2 = undefined
      let getRes = undefined
      let getReq = undefined
      getRes = resvalue.value || JSON.parse(JSON.stringify(resvalue.value))
      getReq = reqvalue.value || JSON.parse(JSON.stringify(reqvalue.value))
      try {
        switch (newVal) {
          case 'GB2312':
            res1 = useCodec({ data: getRes, formatting: 'GB2312', target: 1 })
            res2 = useCodec({ data: getReq, formatting: 'GB2312', target: 1 })
            break
          case 'GBK':
            res1 = useCodec({ data: getRes, formatting: 'GBK', target: 1 })
            res2 = useCodec({ data: getReq, formatting: 'GBK', target: 1 })
            break

          default:
            res1 = getRes
            res2 = getReq
            break
        }
        // switch (oldval) {
        //   case 'ASCII':
        //     getRes = useCodec({ data: resvalue.value, formatting: 'ASCII', target: 1 })
        //     getReq = useCodec({ data: reqvalue.value, formatting: 'ASCII', target: 1 })
        //     switch (newVal) {
        //       case 'UTF-8':
        //         res1 = useCodec({ data: getRes, formatting: 'UTF-8', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'UTF-8', target: 0 })
        //         break
        //       case 'GB2312':
        //         res1 = useCodec({ data: getRes, formatting: 'GB2312', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'GB2312', target: 0 })
        //         break
        //       case 'GBK':
        //         res1 = useCodec({ data: getRes, formatting: 'GBK', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'GBK', target: 0 })
        //         break
        //       default:
        //         break
        //     }
        //     break
        //   case 'UTF-8':
        //     getRes = useCodec({ data: resvalue.value, formatting: 'UTF-8', target: 1 })
        //     getReq = useCodec({ data: reqvalue.value, formatting: 'UTF-8', target: 1 })
        //     switch (newVal) {
        //       case 'ASCII':
        //         res1 = useCodec({ data: getRes, formatting: 'ASCII', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'ASCII', target: 0 })
        //         break
        //       case 'GB2312':
        //         res1 = useCodec({ data: getRes, formatting: 'GB2312', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'GB2312', target: 0 })
        //         break
        //       case 'GBK':
        //         res1 = useCodec({ data: getRes, formatting: 'GBK', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'GBK', target: 0 })
        //         break
        //       default:
        //         break
        //     }
        //     break
        //   case 'GB2312':
        //     getRes = useCodec({ data: resvalue.value, formatting: 'GB2312', target: 1 })
        //     getReq = useCodec({ data: reqvalue.value, formatting: 'GB2312', target: 1 })
        //     switch (newVal) {
        //       case 'UTF-8':
        //         res1 = useCodec({ data: getRes, formatting: 'UTF-8', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'UTF-8', target: 0 })
        //         break
        //       case 'ASCII':
        //         res1 = useCodec({ data: getRes, formatting: 'ASCII', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'ASCII', target: 0 })
        //         break
        //       case 'GBK':
        //         res1 = useCodec({ data: getRes, formatting: 'GBK', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'GBK', target: 0 })
        //         break
        //       default:
        //         break
        //     }
        //     break
        //   case 'GBK':
        //     getRes = useCodec({ data: resvalue.value, formatting: 'GBK', target: 1 })
        //     getReq = useCodec({ data: reqvalue.value, formatting: 'GBK', target: 1 })

        //     switch (newVal) {
        //       case 'UTF-8':
        //         res1 = useCodec({ data: getRes, formatting: 'UTF-8', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'UTF-8', target: 0 })
        //         break
        //       case 'ASCII':
        //         res1 = useCodec({ data: getRes, formatting: 'ASCII', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'ASCII', target: 0 })
        //         break
        //       case 'GB2312':
        //         res1 = useCodec({ data: getRes, formatting: 'GB2312', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'GB2312', target: 0 })
        //         break
        //       default:
        //         break
        //     }
        //     break

        //   default:
        //     switch (newVal) {
        //       case 'ASCII':
        //         res1 = useCodec({ data: getRes, formatting: 'ASCII', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'ASCII', target: 0 })
        //         break
        //       case 'GB2312':
        //         res1 = useCodec({ data: getRes, formatting: 'GB2312', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'GB2312', target: 0 })
        //         break
        //       case 'GBK':
        //         res1 = useCodec({ data: getRes, formatting: 'GBK', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'GBK', target: 0 })
        //         break
        //       case 'UTF-8':
        //         res1 = useCodec({ data: getRes, formatting: 'UTF-8', target: 0 })
        //         res2 = useCodec({ data: getReq, formatting: 'UTF-8', target: 0 })
        //         break
        //       default:
        //         res1 = getRes
        //         res2 = getReq
        //         break
        //     }
        //     break
        // }
        // emit('on-change-event', { res: res1, req: res2 })
      } catch {
        console.log('error')
      }
    }
  )

  watch(
    () => props.target,
    () => {
      resvalue.value = props.target.res && JSON.parse(JSON.stringify(props.target.res))
      reqvalue.value = props.target.req && JSON.parse(JSON.stringify(props.target.req))
    },
    { deep: true }
  )
</script>

<template>
  <div class="fomart-code">
    <el-select v-model="codeMode" placeholder="请选择编码模式" style="width: 140px">
      <!-- <el-option key="ASCII" label="ASCII" value="ASCII" />
      <el-option key="UTF-8" label="UTF-8" value="UTF-8" /> -->
      <el-option key="GB2312" label="GB2312" value="GB2312" />
      <el-option key="GBK" label="GBK" value="GBK" />
    </el-select>
  </div>
</template>

<style scoped lang="scss"></style>
