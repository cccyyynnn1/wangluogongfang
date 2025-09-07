<script lang="ts">
  export default {
    name: 'ThuAlertApi',
  }
</script>

<script setup lang="ts">
  import { stringToHex } from '~/src/utils/text'
  import AttackHighlight from '@/components/attack-highlight.vue'
  import { PayloadDivider } from '~/src/hooks/useHttpTelegramMatch'
  import { removeLinerBeaksAndReturns } from '@/utils/text'
  const props = defineProps<{
    value: any
    type: string
  }>()
  const magicIcon = require('@/assets/mofabang.svg')
  const showRequest = ref(false)
  const showResponse = ref(false)
  const codeMode = ref()
  const magicRequestVal = ref('')
  const magicResponseVal = ref('')

  const allData = ref()
  const loading = ref(true)
  const requestList = ref([
    { Host: 'requestHost', label: '' },
    { Accept: 'requestAccept', label: '' },
    { Cookie: 'requestCookie', label: '' },
    { Authorization: 'requestProxyAuthorization', label: '' },
    { 'User-Agent': 'requestUserAgent', label: '' },
    { ContentType: 'requestContentType', label: '' },
    { Referer: 'requestReferer', label: '' },
    { 'Content-Length': 'requestContentLength', label: '' },
    { 'Accept-Language': 'requestAcceptLanguage', label: '' },
    { Connection: 'requestConnection', label: '' },
    { 'X-Forwarded-For': 'requestXForwardedFor', label: '' },
    { Via: 'requestVia', label: '' },
    { 'Accept-Charset': 'requestAcceptCharset', label: '' },
    { Rang: 'requestRang', label: '' },
    { 'If-Rang': 'requestIfRang', label: '' },
    { 'Accept-Encoding': 'requestAcceptEncoding', label: '' },
  ])

  const responseList = ref([
    { Server: 'responseServer', label: '' },
    { 'Content-Type': 'responseContentType', label: '' },
    { 'Content-Length': 'responseContentLength', label: '' },
    { Connection: 'responseConnection', label: '' },
    { 'Set-Cookie': 'responseSetCookie', label: '' },
    { Authorization: 'responseProxyAuthorization', label: '' },
    { 'Content-Encoding': 'responseContentEncoding', label: '' },
    { Authorization: 'responseWwwAuthorization', label: '' },
    { 'Content-Disposition': 'responseContentDisposition', label: '' },
    { Via: 'responseVia', label: '' },
    { Location: 'responseLocation', label: '' },
  ])
  const attack_req_ref = ref<InstanceType<typeof AttackHighlight> | null>(null)
  const attack_res_ref = ref<InstanceType<typeof AttackHighlight> | null>(null)
  // 将数据转化成相应格式
  const changeData = (list: any) => {
    list.forEach((item: any) => {
      Object.keys(item).forEach((td) => {
        if (td !== 'label') {
          if (allData.value[item[td]]) {
            item['label'] = `${td}: ${allData.value[item[td]]}`
          }
        }
      })
    })
  }

  const copyData = reactive({
    req: '',
    res: '',
  })
  const needToFormet = inject('isNeedToFormet')

  const getCopyData = () => {
    // @ts-ignore
    if (needToFormet.value) {
      let str = ''
      requestList.value.forEach((item: any) => {
        if (item.label) {
          str += `${item.label}\n`
        }
      })
      copyData.req = `${allData.value.requestMethod} ${allData.value.requestUrl}\n${str}${removeLinerBeaksAndReturns(
        allData.value.requestOtherData
      )}${PayloadDivider}${allData.value.requestPayload || ''}`
      let str1 = ''
      responseList.value.forEach((item: any) => {
        if (item.label) {
          str1 += `${item.label}\n`
        }
      })
      copyData.res = `HTTP/${allData.value.responseVersion} ${
        allData.value.responseStatusCode
      }\n${str1}${removeLinerBeaksAndReturns(allData.value.responseOtherData)}${PayloadDivider}${
        allData.value.responsePayload || ''
      }`
    } else {
      const splitStringWithReduce = (str: string, n: number) =>
        [...str].reduce(
          (acc, char, index) => (index % n ? (acc[acc.length - 1] += char) : acc.push(char), acc),
          [] as string[]
        )
      const _res = splitStringWithReduce(atob(props.value.packet) || '', 16)
      if (props.type === 'TUH') copyData.res = _res.join('\n')
      copyData.req = stringToHex(_res.join(''), '', 1, '', 16)
    }
  }
  const toMagic = (type: 'req' | 'res') => {
    if (type === 'req') {
      attack_req_ref.value?.setConvertsType('magic')
    } else {
      attack_res_ref.value?.setConvertsType('magic')
    }
    codeMode.value = ''
  }
  const handleChange = (value: any) => {
    attack_req_ref.value?.setConvertsType(value)
    attack_res_ref.value?.setConvertsType(value)
  }
  watch(
    () => props.value,
    () => {
      handleChange('default')
      codeMode.value = ''
      allData.value = props.value
      changeData(requestList.value)
      changeData(responseList.value)
      showRequest.value = false
      showResponse.value = false
      getCopyData()
    },
    { deep: true }
  )
  onMounted(() => {
    allData.value = props.value
    changeData(requestList.value)
    changeData(responseList.value)
    loading.value = false
    getCopyData()
  })
</script>

<template>
  <div v-if="!loading">
    <div class="hostTable_title marginB15">
      {{ needToFormet ? '请求和响应' : 'Payload' }}
      <div class="fomart-code">
        <el-select v-model="codeMode" placeholder="请选择编码模式" style="width: 140px" @change="handleChange">
          <el-option key="default" label="恢复原始编码" value="default" />
          <el-option key="GB2312" label="GB2312" value="gb2312" />
          <el-option key="GBK" label="GBK" value="gbk" />
        </el-select>
      </div>
    </div>
    <div class="requests-and-responses">
      <div class="requests line" :style="{ width: needToFormet ? '50%' : '69%' }">
        <div class="top">
          <span style="color: #4b4764; font-weight: 500">{{ needToFormet ? '请求' : 'HEX' }}</span>
          <el-tooltip content="魔法棒" effect="light" placement="top" :show-arrow="false">
            <el-image
              v-if="needToFormet"
              :src="magicIcon"
              style="width: 16px; height: 16px; cursor: pointer"
              @click="toMagic('req')"
            />
          </el-tooltip>
        </div>
        <div class="content">
          <attack-highlight
            v-if="needToFormet"
            ref="attack_req_ref"
            :attack-http-msg="copyData.req"
            attack-type="request"
          />
          <div v-else class="diyLog" v-html="copyData.req"></div>
        </div>
      </div>
      <div class="requests" :style="{ width: needToFormet ? '50%' : '31%' }">
        <div class="top">
          <span style="color: #4b4764; font-weight: 500">{{ needToFormet ? '响应' : 'ASCII' }}</span>
          <el-tooltip content="魔法棒" effect="light" placement="top" :show-arrow="false">
            <el-image :src="magicIcon" style="width: 16px; height: 16px; cursor: pointer" @click="toMagic('res')" />
          </el-tooltip>
        </div>
        <div class="content" :class="{ ascii: !needToFormet }">
          <attack-highlight ref="attack_res_ref" :attack-http-msg="copyData.res" attack-type="response" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .hostTable_title {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .requests-and-responses {
    height: calc(100vh - 20px);
    width: 100%;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));

    display: flex;

    .line {
      border-right: 1px solid var(--el-border-color-lighter);
    }
    .requests {
      width: 50%;

      .top {
        padding: 0 20px;
        height: 40px;
        background: #f5f7fa;
        border-bottom: 1px solid var(--el-border-color-lighter);
        display: flex;
        align-items: center;
        justify-content: space-between;
      }
      .content {
        height: calc(100vh - 65px);
        overflow-y: auto;
        &.ascii {
          :deep() {
            .cm-scroller {
              letter-spacing: 4px;
              line-height: 2;
            }
          }
        }
        :deep() {
          .attack-highlight {
            background-color: transparent;
            border-radius: 0;
            border: 0;
            height: 100%;
            // background-color: red;
            // overflow-y: auto;
            .cm-scroller {
              min-height: 300px;
              .cm-content {
                width: 100%;
                white-space: break-spaces;
                word-wrap: break-word;
              }
            }
          }
        }
        &::-webkit-scrollbar {
          width: 0;
          height: 0;
        }
        .content-item {
          padding: 20px 20px 0;
          line-height: 20px;
          overflow: hidden;
          word-break: break-all;
          word-wrap: break-word;
          // box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color-lighter)) inset;
          // border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
          transition: var(--el-transition-box-shadow);
          white-space: pre-wrap;
          margin: 0;
          &:hover {
            cursor: text;
          }
        }
        :deep() {
          .el-textarea {
            margin: 0;
            border: none !important;
            background-color: #fff;
            .el-textarea__inner {
              padding: 15px 20px;
              background-color: #fff;
              border: none !important;
              box-shadow: none;
            }
          }
        }
      }
    }
  }
  .marginB15 {
    margin-bottom: 15px !important;
  }
  .my-item {
    background: #fff !important;
    padding: 20px;
    line-height: 20px;
    overflow: hidden;
    word-break: break-all;
    word-wrap: break-word;
    border: none !important;
    box-shadow: none !important;
    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
    transition: var(--el-transition-box-shadow);
    white-space: pre-wrap;
    margin: 0 !important;
    &:hover {
      cursor: text;
    }
  }
</style>
