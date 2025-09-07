<script lang="ts">
  export default {
    name: 'AlertAiAnalyse',
  }
</script>

<script setup lang="ts">
  import { useAiAnalyseHistoryStore } from '@/store/modules/ai-analyse'
  import { contentType } from '@/config'
  import { fetchEventSource } from '@microsoft/fetch-event-source'
  import { formatNstime } from '@/utils/time'
  import 'highlight.js/styles/default.css'
  import { useUserStore } from '@/store/modules/user'
  import { getAiQuickInputApi } from '@/api-ecs/public'
  import AiChatAnswer from '@/components/ai-chat-answer.vue'
  const abortController = ref<AbortController>()
  const isGenerating = ref(false)
  const isComposition = ref(false)
  const quickInputs = ref<string[]>([])
  const router = useRouter()
  const aiAnalyseStore = useAiAnalyseHistoryStore()
  const props = defineProps<{
    question?: string
    collapse?: boolean
    alarmData?: any
    moduleEnable?: boolean
  }>()
  const queryStr = window.location.hash.split('#/asset-access-insights')[1] || ''
  const routeParams = new URLSearchParams(queryStr)
  let resizeObserverDom: Element | undefined
  const analyseRef = ref()
  const showHelp = ref(!aiAnalyseStore.aiAnalyseHistory.size)
  const showSend = ref(false)
  // 用户咨询的问题
  const questionStr = ref(props?.question || '')
  // AI研判
  const analyseAnswer = ref('')
  // AI问答
  const questionList = ref<{ role: 'assistant' | 'user'; content: string; status?: 'loading' | 'done' | 'cancel' }[]>(
    []
  )
  const { start, stop } = useTimeoutFn(() => {
    handleAiAnalyse()
  }, 2000)
  const userStore = useUserStore()
  const { token = '' } = userStore
  const setQuestionStr = (question: string) => {
    questionStr.value = question
  }

  const request_header: Record<string, string> = window.location.href.includes('asset-access-insights')
    ? {
        authId: routeParams.get('authId') || '',
        timestamp: routeParams.get('timestamp') || '',
        token: routeParams.get('token') || '',
      }
    : { Authorization: `Bearer ${token}`, 'Content-Type': contentType, Cookies: `EcsSessionId=${token}` }
  const handleGetAiChatStrem = async (data: { q: string }) => {
    abortController.value = new AbortController()

    const fetchTimeOut = () =>
      new Promise((resolve, reject) => {
        setTimeout(() => {
          reject('请求超时')
        }, 80000)
        fetchEventSource('/v3/ecsPlatform/gpt/attackHttp', {
          method: 'POST',
          headers: request_header,
          openWhenHidden: true,
          body: JSON.stringify(data),
          signal: abortController.value!.signal,
          onmessage(event) {
            const message_data = JSON.parse(event.data).msg
            showSend.value = false
            analyseAnswer.value = message_data.replaceAll('*', '')
            resolve('success')
          },
          onclose() {
            questionList.value = [
              {
                role: 'user',
                content: data.q,
              },
              {
                role: 'assistant',
                content: analyseAnswer.value,
              },
            ]
            aiAnalyseStore.setAnalyseHistory(props.alarmData?.id, {
              assistant: analyseAnswer.value,
              user: data.q,
            })
            isGenerating.value = false
          },
        })
      })
    fetchTimeOut().catch((err) => {
      abortController.value?.abort()
      showSend.value = false
      isGenerating.value = false
      analyseAnswer.value = '请求超时，请重新请求'
    })
  }
  const handleGetAiAnalyse = async () => {
    const {
      threatName,
      sourceData,
      id,
      startTimeNs,
      victimIp = '',
      targetPort = '',
      sourcePort = '',
      clientIp = '',
    } = props.alarmData
    const str = `告警时间:${
      formatNstime(startTimeNs) || ''
    }\n源IP:${clientIp}\n源端口:${sourcePort}\n目的IP:${victimIp}\n目的端口:${targetPort}\n威胁名称:${threatName}\n`
    const _sourceData = JSON.parse(sourceData)
    const analyseStr =
      _sourceData.event_type === 'api-inspection'
        ? `${str}请求体:${_sourceData.http.request}\n响应体:${_sourceData.http.response}`
        : `${str}请求头:${_sourceData?.http?.request_header || ''}\n响应头:${
            _sourceData?.http?.response_header || ''
          }\n请求体:${_sourceData?.http?.http_request_body || ''}\n响应体:${
            _sourceData?.http?.http_response_body || ''
          }`

    isGenerating.value = true
    handleGetAiChatStrem({ q: analyseStr })
  }
  // 获取Chart对话
  const hanldeGetAiChart = async () => {
    const userStore = useUserStore()
    const { token } = userStore
    const { id } = props.alarmData
    const history = questionList.value.toSpliced(-1).map(({ role, content }) => ({ role, content }))
    const [a, b, ...c] = history
    const questionHistoryStr = questionStr.value
    let aiContent = ''
    questionStr.value = ''
    isGenerating.value = true
    const fetchTimeOut = () =>
      new Promise((resolve, reject) => {
        setTimeout(() => {
          reject('请求超时')
        }, 80000)
        fetchEventSource('/v3/ecsPlatform/gpt/chat', {
          method: 'POST',
          headers: request_header,
          openWhenHidden: true,
          body: JSON.stringify({
            q: questionHistoryStr,
            history: history.length <= 20 ? history : [a, b, ...c.slice(-18)],
          }),
          onmessage(event) {
            const length = questionList.value.length - 1
            showSend.value = false
            const message_data = JSON.parse(event.data)
            const text = message_data.msg
            if (aiContent) {
              const chat = questionList.value[length]
              chat!.content = text
            } else {
              questionList.value.push({
                role: 'assistant',
                content: text,
              })
            }
            aiContent = text
            resolve('success')
          },
          onclose() {
            const length = questionList.value.length - 1
            const chat = questionList.value[length]
            chat!.status = 'done'
            isGenerating.value = false
            showSend.value = false
            aiAnalyseStore.setAnalyseHistory(id, { assistant: aiContent, user: questionHistoryStr })
          },
        })
      })

    fetchTimeOut().catch((err) => {
      abortController.value?.abort()
      showSend.value = false
      questionList.value.push({
        role: 'assistant',
        content: '请求超时，请重新请求',
      })
      const length = questionList.value.length - 1
      const chat = questionList.value[length]
      chat!.status = 'done'
      isGenerating.value = false
    })
  }
  // 发送AI Chart问答
  const handleSendQuestion = () => {
    if (!questionStr.value) return
    questionList.value.push({
      role: 'user',
      content: questionStr.value,
    })
    showSend.value = true
    hanldeGetAiChart()
  }
  // Ai研判（刚进页面判断）
  const handleAiAnalyse = () => {
    if (!props.moduleEnable) return
    const histories = aiAnalyseStore.getCurrentAlarmHistory(props.alarmData.id) || []
    showSend.value = true
    // 当能获取到历史
    if (histories.length) {
      const arr = histories
        .map(({ assistant, user }) => [
          {
            role: 'user',
            content: user,
          },
          {
            role: 'assistant',
            content: assistant,
          },
        ])
        .flat(2)
      questionList.value = arr as { role: 'assistant' | 'user'; content: string }[]
      analyseAnswer.value = histories[0].assistant
      showSend.value = false
    } else {
      // 否则就去AI调研
      handleGetAiAnalyse()
    }
  }
  const handleToModuleConfig = () => {
    const resolveRouter = router.resolve({
      path: '/managements/configuration/systemConfig',
      query: { params: 'modelConfig' },
    })
    window.open(resolveRouter.href, '_blank')
  }
  const handleAiReload = () => {
    questionList.value.pop()
    const length = questionList.value.length - 1
    const chat = questionList.value[length]
    questionStr.value = chat!.content
    const histories = aiAnalyseStore.getCurrentAlarmHistory(props.alarmData.id) || []
    histories.pop()
    showSend.value = true
    hanldeGetAiChart()
  }
  const handleCompositionStart = () => {
    isComposition.value = true
  }
  const handleCompositionEnd = () => {
    isComposition.value = false
  }
  const handleAnalyseSearchKeydown = (event: KeyboardEvent) => {
    if (event.code === 'Enter' && !event.shiftKey) {
      event.preventDefault()
      if (!questionStr.value || isGenerating.value || isComposition.value) return
      handleSendQuestion()
    }
  }

  const hanldeGetAiQuickInput = async () => {
    const { attackIp, threatName, victimIp, host } = props.alarmData
    const params = {
      attackIp: attackIp,
      victimIp: victimIp,
      host: host,
      threatName: threatName,
    }
    const { data } = await getAiQuickInputApi(params)
    quickInputs.value = Array.isArray(data) ? data.sort((a, b) => a.length - b.length) : []
  }

  const resizeObserver = new ResizeObserver((entries) => {
    analyseRef.value?.scrollTo({
      top: analyseRef.value.scrollHeight,
      behavior: 'smooth',
    })
  })
  watch(
    () => props.alarmData,
    (val) => {
      if (val.id) {
        abortController.value?.abort()
        isGenerating.value = false
        questionList.value = []
        questionStr.value = ''
        analyseAnswer.value = ''
        showSend.value = true
        start()
        hanldeGetAiQuickInput()
      }
    },
    {
      deep: true,
      immediate: true,
    }
  )
  watch(
    () => props.collapse,
    () => {
      if (props.collapse) {
        nextTick(() => {
          resizeObserverDom = analyseRef.value.getElementsByClassName('analyse-lists')[0]
          resizeObserver.observe(resizeObserverDom as Element)
        })
      }
    },
    {
      deep: true,
      immediate: true,
    }
  )
  onUnmounted(() => {
    if (resizeObserverDom) resizeObserver.unobserve(resizeObserverDom as Element)
    stop()
  })
</script>

<template>
  <div class="alert-ai-analyse">
    <!-- 展开时显示 -->
    <template v-if="collapse">
      <div ref="analyseRef" class="ai-analyse-box">
        <div class="analyse-lists">
          <!-- 帮助提示 -->
          <div v-if="showHelp" class="analyse answer help">
            HI，我是您的AI助理，我能帮助您做告警研判、调查取证，我还知道如何防范各类攻击以及漏洞修复的方法，同时我还能协助您写各类报告。
            您有什么需要我做的，尽管吩咐我就好了。
            <el-image class="aiAvatar" :src="require('@/assets/alert_images/ai_avatar.svg')" />
          </div>
          <!-- AI研判 -->
          <div v-if="analyseAnswer" class="analyse answer analyseAnswer">
            <ai-chat-answer :answer-content="analyseAnswer" role="assistant" />
            <el-image class="aiAvatar" :src="require('@/assets/alert_images/ai_avatar.svg')" />
          </div>
          <!-- 对话信息 -->
          <template v-for="(question, index) in questionList" :key="index">
            <div v-if="index > 1" class="analyse" :class="[question.role === 'assistant' ? 'answer' : 'question']">
              <template v-if="question.role === 'assistant'">
                <ai-chat-answer :answer-content="question.content" role="assistant" />
                <el-image class="aiAvatar" :src="require('@/assets/alert_images/ai_avatar.svg')" />
                <span
                  v-if="index === questionList.length - 1 && question?.status === 'done'"
                  class="ai-reload"
                  @click="handleAiReload"
                >
                  <el-icon>
                    <Refresh />
                  </el-icon>
                  重新回答
                </span>
              </template>
              <template v-else>
                {{ question.content }}
              </template>
            </div>
          </template>
          <!-- 发送中 -->
          <div v-if="showSend" class="analyse answer">
            <div style="width: 50px">
              <div class="loader"></div>
            </div>
            <el-image class="aiAvatar" :src="require('@/assets/alert_images/ai_avatar.svg')" />
          </div>
        </div>
      </div>
      <div class="analyse-search">
        <div class="tip">
          <svg
            height="20px"
            style="vertical-align: middle; margin-right: 10px"
            version="1.1"
            viewBox="0 0 20 20"
            width="20px"
          >
            <g id="告警大模型" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
              <g id="大模型告警详情1-更多问题" fill-rule="nonzero" transform="translate(-687, -632)">
                <g id="新增宽带配置" transform="translate(220, 0)">
                  <g id="编组-20" transform="translate(346, 76)">
                    <g id="编组-21" transform="translate(121, 445)">
                      <g id="编组-5" transform="translate(0, 107)">
                        <g id="wen-2" transform="translate(0, 4)">
                          <path
                            id="形状"
                            d="M18.2182409,15.679278 C19.3850847,13.9982538 20,12.042359 20,9.99965807 C20,4.48553499 15.5144052,0 10.0000684,0 C4.48573158,0 0,4.48617325 0,10.0003419 C0,15.5130289 4.48570882,19.999225 10.0000456,20 C11.6066812,20 13.1972909,19.6072299 14.6216677,18.8621293 L18.3057799,19.9681088 C18.5691945,20.0469136 18.8498434,19.9753123 19.0421783,19.7822326 C19.2346499,19.5897684 19.306254,19.306942 19.2281073,19.0458399 L18.2182409,15.679278 Z M6.39496971,14.6633367 C6.38765616,14.8666893 6.32584569,15.0339122 6.20967479,15.1647324 C6.09336743,15.2808715 5.9443401,15.3426796 5.76278382,15.3499929 C5.58109108,15.3426796 5.43203647,15.28461 5.31589285,15.1756204 C5.19958549,15.0374869 5.14151367,14.8667166 5.14151367,14.6633639 L5.14151367,7.73144645 C5.14151367,7.53527067 5.19601056,7.38267432 5.30497709,7.27368476 C5.42846157,7.15754564 5.5810638,7.09931234 5.76275651,7.09931234 C5.9371357,7.09931234 6.08242438,7.15380712 6.19875902,7.2627967 C6.32952975,7.37909955 6.39494241,7.53527064 6.39494241,7.73147373 L6.39494243,14.6633367 L6.39496971,14.6633367 Z M7.1470379,7.29548815 C6.96534516,7.26653523 6.80542937,7.16112041 6.66745424,6.97940744 C6.34047282,6.57256576 6.05333391,6.26737309 5.80639228,6.06388403 C5.63913563,5.92594153 5.54105759,5.76960671 5.51207627,5.595207 C5.50476271,5.42814785 5.55928691,5.26823825 5.67559427,5.11564193 C5.79173789,4.96304561 5.93716303,4.88319996 6.11159678,4.87585938 C6.28597596,4.86139657 6.46038243,4.92320462 6.63476158,5.06114713 C7.49208493,5.75869136 7.91002185,6.30740506 7.88821763,6.70690619 C7.86644071,6.88861916 7.77922382,7.03764072 7.62664886,7.15377984 C7.46670577,7.27011 7.30678995,7.31726422 7.1470379,7.29548815 Z M8.20431056,12.8758693 C7.94986433,12.8687198 7.75368094,12.774193 7.61570582,12.59248 C7.44132663,12.3673786 7.35771196,12.1311436 7.36502551,11.8840478 L7.36502551,9.24643985 C7.36502551,9.02128387 7.44490154,8.82136956 7.60481736,8.64699716 C7.77919654,8.47993801 7.97898211,8.39629929 8.20428328,8.39629929 L8.84733035,8.39629929 L11.0708422,8.39629929 C11.3613923,8.38914976 11.6049228,8.47636324 11.8011334,8.6578852 C11.9681991,8.81048152 12.0553886,9.00665733 12.0627022,9.24643985 L12.0627022,11.8840478 C12.0553886,12.1456064 11.9499697,12.3782394 11.7466365,12.581592 C11.5431123,12.7777678 11.3177838,12.8758693 11.0708422,12.8758693 L8.20431056,12.8758693 Z M13.7630765,14.8922312 C13.3415919,15.2990729 12.4769278,15.4116372 11.1689748,15.230088 C10.958137,15.1573646 10.8093007,15.0410617 10.7220838,14.8813158 C10.6566712,14.7214062 10.6493576,14.5543471 10.7003069,14.3799747 C10.8746861,13.9729693 11.1725497,13.8459148 11.5940616,13.9984839 C12.1608348,14.0710435 12.5314247,14.0530059 12.7058311,13.9439891 C12.9527728,13.7767662 13.0763937,13.5151803 13.0763937,13.1592313 L13.0763937,7.34998293 C13.0690801,7.11020038 13.0072697,6.92491263 12.8911261,6.79411969 C12.7747914,6.67798057 12.5930987,6.61974727 12.3461297,6.61974727 L8.90185455,6.61974727 C8.68389422,6.61259773 8.50948775,6.55076238 8.37868974,6.43445952 C8.26969594,6.31097983 8.2114604,6.16198558 8.20431056,5.98758587 C8.21877393,5.79872336 8.2805844,5.64972909 8.38960548,5.54071222 C8.52752605,5.4245731 8.69835758,5.36631252 8.90185455,5.36631252 L12.8693219,5.36631252 C13.3925139,5.35916298 13.7666787,5.48619014 13.9920072,5.74780335 C14.2171173,6.00223972 14.3299043,6.35816139 14.3299043,6.81592308 L14.3299043,13.6170203 C14.3299043,14.074782 14.1408162,14.4998796 13.7630765,14.8922312 Z"
                            fill="#9F90F9"
                          />
                          <path
                            id="路径"
                            d="M10.7688755,9.70814132 L8.80316002,9.70814132 C8.74632409,9.70814132 8.69899639,9.72827696 8.66112948,9.76859875 C8.61688426,9.81568285 8.59793894,9.86937793 8.60429356,9.92986058 L8.60429356,11.401277 C8.60429356,11.4483611 8.62321515,11.49542 8.66112948,11.542378 C8.69899639,11.5894369 8.74634779,11.6129032 8.80316002,11.6129032 L10.7688755,11.6129032 C10.825664,11.6129032 10.8730391,11.5894369 10.910906,11.542378 C10.9424183,11.4954452 10.9613636,11.4483864 10.9677419,11.401277 L10.9677419,9.92986058 C10.9677419,9.87614028 10.9487966,9.82577592 10.910906,9.77869182 C10.8603062,9.73178429 10.8129547,9.70814132 10.7688755,9.70814132 Z"
                            fill="#A08FF9"
                          />
                        </g>
                      </g>
                    </g>
                  </g>
                </g>
              </g>
            </g>
          </svg>
          <span v-for="(question, index) in quickInputs.slice(0, 3)" :key="index" @click="setQuestionStr(question)">
            {{ question }}
          </span>
          <el-popover placement="top-end" :show-arrow="false" :teleported="false" trigger="hover" width="max-content">
            <template #reference>
              <span class="exchange">
                <el-icon style="vertical-align: 5px"><MoreFilled /></el-icon>
              </span>
            </template>
            <template #default>
              <div class="quickAction" style="margin: 8px">
                <div>
                  <svg
                    height="20px"
                    style="vertical-align: bottom; margin-right: 6px"
                    version="1.1"
                    viewBox="0 0 20 20"
                    width="20px"
                  >
                    <g id="告警大模型" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
                      <g id="大模型告警详情1-更多问题" fill-rule="nonzero" transform="translate(-687, -632)">
                        <g id="新增宽带配置" transform="translate(220, 0)">
                          <g id="编组-20" transform="translate(346, 76)">
                            <g id="编组-21" transform="translate(121, 445)">
                              <g id="编组-5" transform="translate(0, 107)">
                                <g id="wen-2" transform="translate(0, 4)">
                                  <path
                                    id="形状"
                                    d="M18.2182409,15.679278 C19.3850847,13.9982538 20,12.042359 20,9.99965807 C20,4.48553499 15.5144052,0 10.0000684,0 C4.48573158,0 0,4.48617325 0,10.0003419 C0,15.5130289 4.48570882,19.999225 10.0000456,20 C11.6066812,20 13.1972909,19.6072299 14.6216677,18.8621293 L18.3057799,19.9681088 C18.5691945,20.0469136 18.8498434,19.9753123 19.0421783,19.7822326 C19.2346499,19.5897684 19.306254,19.306942 19.2281073,19.0458399 L18.2182409,15.679278 Z M6.39496971,14.6633367 C6.38765616,14.8666893 6.32584569,15.0339122 6.20967479,15.1647324 C6.09336743,15.2808715 5.9443401,15.3426796 5.76278382,15.3499929 C5.58109108,15.3426796 5.43203647,15.28461 5.31589285,15.1756204 C5.19958549,15.0374869 5.14151367,14.8667166 5.14151367,14.6633639 L5.14151367,7.73144645 C5.14151367,7.53527067 5.19601056,7.38267432 5.30497709,7.27368476 C5.42846157,7.15754564 5.5810638,7.09931234 5.76275651,7.09931234 C5.9371357,7.09931234 6.08242438,7.15380712 6.19875902,7.2627967 C6.32952975,7.37909955 6.39494241,7.53527064 6.39494241,7.73147373 L6.39494243,14.6633367 L6.39496971,14.6633367 Z M7.1470379,7.29548815 C6.96534516,7.26653523 6.80542937,7.16112041 6.66745424,6.97940744 C6.34047282,6.57256576 6.05333391,6.26737309 5.80639228,6.06388403 C5.63913563,5.92594153 5.54105759,5.76960671 5.51207627,5.595207 C5.50476271,5.42814785 5.55928691,5.26823825 5.67559427,5.11564193 C5.79173789,4.96304561 5.93716303,4.88319996 6.11159678,4.87585938 C6.28597596,4.86139657 6.46038243,4.92320462 6.63476158,5.06114713 C7.49208493,5.75869136 7.91002185,6.30740506 7.88821763,6.70690619 C7.86644071,6.88861916 7.77922382,7.03764072 7.62664886,7.15377984 C7.46670577,7.27011 7.30678995,7.31726422 7.1470379,7.29548815 Z M8.20431056,12.8758693 C7.94986433,12.8687198 7.75368094,12.774193 7.61570582,12.59248 C7.44132663,12.3673786 7.35771196,12.1311436 7.36502551,11.8840478 L7.36502551,9.24643985 C7.36502551,9.02128387 7.44490154,8.82136956 7.60481736,8.64699716 C7.77919654,8.47993801 7.97898211,8.39629929 8.20428328,8.39629929 L8.84733035,8.39629929 L11.0708422,8.39629929 C11.3613923,8.38914976 11.6049228,8.47636324 11.8011334,8.6578852 C11.9681991,8.81048152 12.0553886,9.00665733 12.0627022,9.24643985 L12.0627022,11.8840478 C12.0553886,12.1456064 11.9499697,12.3782394 11.7466365,12.581592 C11.5431123,12.7777678 11.3177838,12.8758693 11.0708422,12.8758693 L8.20431056,12.8758693 Z M13.7630765,14.8922312 C13.3415919,15.2990729 12.4769278,15.4116372 11.1689748,15.230088 C10.958137,15.1573646 10.8093007,15.0410617 10.7220838,14.8813158 C10.6566712,14.7214062 10.6493576,14.5543471 10.7003069,14.3799747 C10.8746861,13.9729693 11.1725497,13.8459148 11.5940616,13.9984839 C12.1608348,14.0710435 12.5314247,14.0530059 12.7058311,13.9439891 C12.9527728,13.7767662 13.0763937,13.5151803 13.0763937,13.1592313 L13.0763937,7.34998293 C13.0690801,7.11020038 13.0072697,6.92491263 12.8911261,6.79411969 C12.7747914,6.67798057 12.5930987,6.61974727 12.3461297,6.61974727 L8.90185455,6.61974727 C8.68389422,6.61259773 8.50948775,6.55076238 8.37868974,6.43445952 C8.26969594,6.31097983 8.2114604,6.16198558 8.20431056,5.98758587 C8.21877393,5.79872336 8.2805844,5.64972909 8.38960548,5.54071222 C8.52752605,5.4245731 8.69835758,5.36631252 8.90185455,5.36631252 L12.8693219,5.36631252 C13.3925139,5.35916298 13.7666787,5.48619014 13.9920072,5.74780335 C14.2171173,6.00223972 14.3299043,6.35816139 14.3299043,6.81592308 L14.3299043,13.6170203 C14.3299043,14.074782 14.1408162,14.4998796 13.7630765,14.8922312 Z"
                                    fill="#9F90F9"
                                  />
                                  <path
                                    id="路径"
                                    d="M10.7688755,9.70814132 L8.80316002,9.70814132 C8.74632409,9.70814132 8.69899639,9.72827696 8.66112948,9.76859875 C8.61688426,9.81568285 8.59793894,9.86937793 8.60429356,9.92986058 L8.60429356,11.401277 C8.60429356,11.4483611 8.62321515,11.49542 8.66112948,11.542378 C8.69899639,11.5894369 8.74634779,11.6129032 8.80316002,11.6129032 L10.7688755,11.6129032 C10.825664,11.6129032 10.8730391,11.5894369 10.910906,11.542378 C10.9424183,11.4954452 10.9613636,11.4483864 10.9677419,11.401277 L10.9677419,9.92986058 C10.9677419,9.87614028 10.9487966,9.82577592 10.910906,9.77869182 C10.8603062,9.73178429 10.8129547,9.70814132 10.7688755,9.70814132 Z"
                                    fill="#A08FF9"
                                  />
                                </g>
                              </g>
                            </g>
                          </g>
                        </g>
                      </g>
                    </g>
                  </svg>
                  <span class="title">更多问题</span>
                </div>
                <ul class="quickActionBox">
                  <li
                    v-for="(question, index) in quickInputs.slice(4)"
                    :key="index"
                    class="quickActionItem"
                    @click="setQuestionStr(question)"
                  >
                    {{ question }}
                  </li>
                </ul>
              </div>
            </template>
          </el-popover>
        </div>
        <div class="search-box">
          <el-input
            v-model="questionStr"
            :autosize="{ minRows: 1, maxRows: 4 }"
            class="text"
            clearable
            placeholder="按回车发送消息，按shift+回车换行"
            resize="none"
            type="textarea"
            @compositionend="handleCompositionEnd"
            @compositionstart="handleCompositionStart"
            @keydown="handleAnalyseSearchKeydown"
          />
          <el-button
            class="search-btn"
            :disabled="!questionStr || isGenerating"
            size="small"
            type="primary"
            @click="handleSendQuestion"
          >
            <svg height="40px" version="1.1" viewBox="0 0 40 40" width="40px">
              <title>发送</title>
              <g id="告警大模型" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
                <g id="大模型告警详情1-原始-默认" transform="translate(-1521, -577)">
                  <g id="新增宽带配置" transform="translate(220, 0)">
                    <g id="编组-20" transform="translate(346, 76)">
                      <g id="编组-21" transform="translate(121, 445)">
                        <g id="编组-16" transform="translate(0, 52)">
                          <g id="编组-23" transform="translate(834, 4)">
                            <rect id="矩形" fill="#6954F0" height="40" rx="20" width="40" x="0" y="0" />
                            <g
                              id="发送-(21)"
                              fill="#FFFFFF"
                              fill-rule="nonzero"
                              stroke="#6954F0"
                              stroke-width="0.5"
                              transform="translate(14, 12)"
                            >
                              <path
                                id="路径"
                                d="M9.53913708,7.4349797 L7.52251433,4.9842245 L7.52251433,15.5583793 C7.52251433,16.361568 6.90518084,17 6.12322508,17 C5.34126932,17 4.72393582,16.361568 4.72393582,15.5583793 L4.72393582,4.86065701 L2.41922411,7.45557428 C1.90477953,8.03222256 1.01993486,8.07341172 0.443756929,7.53795261 C-0.111843216,7.00249349 -0.152998783,6.07573732 0.361445796,5.49908904 L4.66220247,0.67995697 C5.4853138,-0.246799197 6.88460305,-0.226204615 7.66655881,0.741740715 L11.6586487,5.5402782 C12.1730933,6.13752106 12.0907822,7.06427723 11.5146043,7.57914177 C10.9384263,8.09400631 10.0535817,8.03222256 9.53913708,7.4349797 Z"
                              />
                            </g>
                          </g>
                        </g>
                      </g>
                    </g>
                  </g>
                </g>
              </g>
            </svg>
          </el-button>
        </div>
      </div>
    </template>
    <!-- 关闭时显示 -->
    <template v-else>
      <div class="analyse-lists" style="margin-top: 32px">
        <div class="analyse answer analyseAnswer">
          <template v-if="!analyseAnswer">
            <div v-if="moduleEnable" style="width: 50px">
              <div class="loader"></div>
            </div>
            <div v-else style="position: relative">
              <h4 style="margin-bottom: 0">HI，我是您的AI助理。</h4>
              <p style="margin-right: 150px; margin-bottom: 0">
                我能帮助您做告警研判、调查取证，我还知道如何防范各类攻击以及漏洞修复的方法，同时我还能协助您写各类报告。
              </p>
              <span class="opneMoudle" @click="handleToModuleConfig">立即开启</span>
            </div>
          </template>
          <template v-else><ai-chat-answer :answer-content="analyseAnswer" role="assistant" /></template>
          <el-image class="aiAvatar" :src="require('@/assets/alert_images/ai_avatar.svg')" />
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped lang="scss">
  .alert-ai-analyse {
    width: 100%;
    height: 100%;
    // padding-bottom: 35px;
    display: flex;
    flex-direction: column;
    background: url('@/assets/alert_images/analyse.png') no-repeat right 0px/140px 140px,
      linear-gradient(180deg, #efedff 0%, #f1eaff 32%, #ecf2ff 54%, #f8f8ff 100%);
    .ai-analyse-box {
      overflow-y: auto;
      scrollbar-width: none;
      flex: 1;
    }
    .analyse-lists {
      flex: 1;
      // width: 880px;
      margin: 20px 60px 0;
      .analyse {
        font-weight: 400;
        font-size: 14px;
        color: #4a4759;
        line-height: 24px;
        background-color: #fff;
        padding: 10px 20px;
        border-radius: 8px;
        position: relative;
        transition: all 0.3s ease;
        max-width: calc(100% - 80px);
        width: fit-content;
        word-wrap: break-word;
        &::after {
          // content: '|';
          // opacity: 0;
          // display: inline-block;
          // animation: cursor-blink 1s;
        }
        &:not(:last-child) {
          margin-bottom: 14px;
        }
        .ai-reload {
          position: absolute;
          right: -80px;
          bottom: 5px;
          display: flex;
          align-items: center;
          font-size: 13px;
          color: #888b91;
          cursor: pointer;
          :deep(.el-icon) {
            color: #888b91;
            margin-right: 5px;
          }
        }
        &.answer {
          margin-right: auto;
          &.help {
            width: 520px;
          }
          .aiAvatar {
            left: -38px;
          }
        }
        &.question {
          background: #e3dfff;
          color: #4a4759;
          margin-left: auto;
          .userAvatar {
            right: -38px;
          }
        }
        &.analyseAnswer {
          white-space: normal;
        }
        .aiAvatar,
        .userAvatar {
          width: 28px;
          height: 28px;
          position: absolute;
          top: 7px;
        }
        .opneMoudle {
          position: absolute;
          right: 23px;
          top: 22px;
          width: 88px;
          height: 32px;
          background: #ffffff;
          line-height: 32px;
          border-radius: 6px;
          border: 1px solid var(--el-color-primary);
          color: var(--el-color-primary);
          display: block;
          text-align: center;
          font-weight: 500;
          font-size: 14px;
          cursor: pointer;
        }
      }
    }
    .analyse-search {
      min-height: 100px;
      margin-block: 20px 45px;
      display: flex;
      flex-direction: column;
      justify-content: end;
      .search-box {
        position: relative;
        margin: 0 60px;
        border-radius: 27px;
        display: flex;
        padding: 4px;
        background-color: #fff;
        align-items: end;
        border: var(--el-border);
        &:has(.el-textarea__inner:focus) {
          border: 1px solid var(--el-color-primary);
        }
        :deep(.el-textarea__inner) {
          box-shadow: none;
          margin-bottom: 5px;
        }

        .el-textarea {
          margin: 0 10px;
        }
        .search-btn {
          width: 40px;
          height: 40px;
          background: #6954f0;
          border-radius: 100%;
          // position: absolute;
          // right: 5px;
          // bottom: 5px;
          &.is-disabled {
            filter: grayscale(50%);
            -webkit-filter: grayscale(10%);
            opacity: 0.5;
          }
        }
      }
      .tip {
        margin: 0 60px 15px;
        font-weight: 400;
        font-size: 13px;
        color: #888b91;
        span:not(.title) {
          height: 28px;
          border-radius: 6px;
          border: 1px solid #dedcee;
          color: #494758;
          margin-right: 8px;
          line-height: 28px;
          background: #ffffff;
          padding: 5px;
          cursor: pointer;
        }
        .title {
          font-weight: 500;
          font-size: 15px;
          color: #4a4759;
        }
        .quickActionBox {
          padding-left: 0;
          margin-block: 15px 10px;
          max-height: min(448px, 40vh);
          overflow-y: auto;
          &:last-child {
            margin-bottom: 0;
          }
        }
        .quickActionItem {
          height: 30px;
          line-height: 30px;
          background: #f5f4ff;
          border-radius: 4px;
          font-weight: 400;
          font-size: 13px;
          color: #494758;
          padding-inline: 10px;
          margin-bottom: 8px;
          cursor: pointer;
          &:hover {
            color: var(--el-color-primary);
          }
        }
        .exchange {
          float: right;
          cursor: pointer;
          display: inline-block;
          color: var(--el-color-primary);
          padding: 7px !important;
        }
      }
    }
    .typing {
      border-right: 0.1em solid;
      width: 9em; /*宽度为“字数 + em”*/
      white-space: nowrap;
      overflow: hidden;
      animation: typing 5s steps(9, end), /*步数为字数*/ blink-caret 0.5s step-end infinite alternate;
    }
    @keyframes typing {
      from {
        width: 0;
      }
    }
    @keyframes blink-caret {
      50% {
        border-color: transparent;
      }
    }
    .loader {
      color: #6954f0;
      width: 4px;
      margin: 10px auto;
      aspect-ratio: 1;
      border-radius: 50%;
      box-shadow: 19px 0 0 4px, 34px 0 0 3px, 50px 0 0 0;
      transform: translateX(-38px);
      animation: l21 0.3s infinite alternate linear;
    }

    @keyframes l21 {
      50% {
        box-shadow: 19px 0 0 3px, 34px 0 0 4px, 50px 0 0 3px;
      }
      100% {
        box-shadow: 19px 0 0 0, 34px 0 0 3px, 50px 0 0 4px;
      }
    }
    @keyframes cursor-blink {
      50% {
        opacity: 1;
      }
    }
  }
</style>
