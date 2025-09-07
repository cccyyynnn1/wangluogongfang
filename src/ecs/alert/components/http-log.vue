<!-- eslint-disable no-control-regex -->
<!-- eslint-disable no-control-regex -->
<script lang="ts">
  export default {
    name: 'HttpLog',
  }
</script>

<script setup lang="ts">
  import { formatStrToJson } from '~/src/utils/transition'
  import { stringToHex } from '~/src/utils/text'
  import useHttpTelegramMatch, { PayloadDivider } from '~/src/hooks/useHttpTelegramMatch'
  import AttackHighlight from '@/components/attack-highlight.vue'
  import { Ref } from 'vue'
  import { removeLinerBeaksAndReturns } from '@/utils/text'
  const magicIcon = require('@/assets/mofabang.svg')
  const matchHttpLog = ref<any>(undefined)
  const props = defineProps<{
    value: any
    type: string
  }>()

  const copyData = reactive({
    req: '',
    res: '',
  })
  const attack_req_ref = ref<InstanceType<typeof AttackHighlight> | null>(null)
  const attack_res_ref = ref<InstanceType<typeof AttackHighlight> | null>(null)
  const codeMode = ref()
  const needToFormet = inject('isNeedToFormet') as Ref<boolean>
  const httpTelegramMatchData = useHttpTelegramMatch(matchHttpLog)
  const attack_term = computed(() => [...new Set(Object.values(httpTelegramMatchData).flat(2).filter(Boolean))])
  const getCopyData = () => {
    // @ts-ignore
    if (needToFormet.value) {
      if (props.type === 'TUH' && props.value.http) {
        if (props.value.event_type === 'api-inspection') {
          copyData.req = `${props.value.http.request || ''}`
          copyData.res = `${props.value.http.response || ''}`
        } else {
          copyData.req = `${removeLinerBeaksAndReturns(
            props.value?.http?.request_header
          )}${PayloadDivider}${removeLinerBeaksAndReturns(props.value.http?.http_request_body)}`
          copyData.res = `${removeLinerBeaksAndReturns(
            props.value?.http?.response_header
          )}${PayloadDivider}${removeLinerBeaksAndReturns(props.value.http?.http_response_body)}`
        }
      } else if (props.type === 'Threatbook—http') {
        copyData.req = `${props.value?.net?.http?.method} ${
          props.value?.net?.http?.protocol
        }\n${removeLinerBeaksAndReturns(
          props.value?.net?.http?.reqs_header
        )}${PayloadDivider}${removeLinerBeaksAndReturns(props.value.http.reqs_body)}`
        copyData.res = `${props.value?.net?.http?.reqs_line}}\n${removeLinerBeaksAndReturns(
          props.value?.net?.http?.resp_header
        )}${PayloadDivider}${removeLinerBeaksAndReturns(props.value?.net?.http?.resp_body)}`
      } else if (props.type === 'Skyeye—http') {
        copyData.req = `${removeLinerBeaksAndReturns(
          props.value?.payloa?.req_header
        )}${PayloadDivider}${removeLinerBeaksAndReturns(props.value?.payloa?.req_body)}`
        copyData.res = `${removeLinerBeaksAndReturns(
          props.value?.payloa?.rsp_header
        )}${PayloadDivider}${removeLinerBeaksAndReturns(props.value?.payloa?.rsp_body)}`
      }
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
  const handleChange = (value: any) => {
    attack_req_ref.value?.setConvertsType(value)
    attack_res_ref.value?.setConvertsType(value)
  }
  const toMagic = (type: 'req' | 'res') => {
    if (type === 'req') {
      attack_req_ref.value?.setConvertsType('magic')
    } else {
      attack_res_ref.value?.setConvertsType('magic')
    }
    codeMode.value = ''
  }
  onMounted(() => {
    matchHttpLog.value = props.value?.http
    getCopyData()
  })
  watch(
    () => props.value,
    () => {
      handleChange('default')
      codeMode.value = ''
      matchHttpLog.value = props.value?.http
      getCopyData()
    },
    { deep: true }
  )
</script>

<template>
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
          :attack-term="attack_term"
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
        <attack-highlight
          ref="attack_res_ref"
          :attack-http-msg="copyData.res"
          :attack-term="attack_term"
          attack-type="response"
        />
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
    // height: calc(100vh - 20px);
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
  .diyLog {
    padding: 15px 20px 0;
    // width: 500px;
    white-space: pre-wrap;
    word-break: break-all;
    word-wrap: break-word;
  }
</style>
