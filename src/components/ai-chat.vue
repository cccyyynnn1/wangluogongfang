<script lang="ts">
  export default {
    name: 'AiChat',
  }
</script>

<script setup lang="ts">
  import { Plus } from '@element-plus/icons-vue'
  import { getSystemConfigApi } from '~/src/api-ecs/system'
  import { useUserStore } from '@/store/modules/user'
  import { fetchEventSource } from '@microsoft/fetch-event-source'
  import { contentType } from '@/config'
  import { aiSessionChatsApi, aiSessionContentApi, deleteAiSessionContentApi } from '../api-ecs/public'
  import 'highlight.js/styles/default.css'
  import { getAiQuickInputApi } from '@/api-ecs/public'
  import AiChatAnswer from './ai-chat-answer.vue'
  let isDragging = false
  const $baseConfirm: any = inject('$baseConfirm')
  const aiChatRef = ref()
  const chatSzieBtn = ref()
  const { y: chatSzieBtnHeight } = useDraggable(chatSzieBtn)
  const { y: aiChatRefTop } = useDraggable(aiChatRef, {
    preventDefault: true,
    initialValue: { y: innerHeight - 150, x: 0 },
    onStart() {
      aiChatRef.value.style.transition = 'none'
    },
    onEnd() {
      aiChatRef.value.style.transition = 'all 0.3s'
      setTimeout(() => {
        isDragging = false
      }, 500)
    },
    onMove(position) {
      isDragging = true
      aiChatRefTop.value = position.y < 0 ? 0 : Math.min(position.y, innerHeight - 80)
    },
  })
  const aiDrawerVisible = ref(false)
  const drawerHeight = ref(innerHeight * 0.55)
  const moduleEnable = ref(false)
  const router = useRouter()
  const isGenerating = ref(false)
  const showSend = ref(false)
  const questionRef = ref()
  const scrollState = ref()
  const quickInputs = ref<string[]>([])
  const isComposition = ref(false)
  let isOpen = false
  const chatPromptRef = ref<HTMLDivElement>()
  const historyRef = ref<HTMLDivElement>()
  //查询会话历史列表参数
  const historyListParams = reactive({
    pageNum: 1,
    pageSize: 50,
  })
  // 会话历史列表
  const historyList = ref<
    {
      id: number
      title: string
    }[]
  >([])
  const historyTotal = ref(0)
  const historyKey = ref()
  // 用户咨询的问题
  const questionStr = ref('')
  // AI问答
  const questionList = ref<
    { role: 'assistant' | 'user'; content: string; status?: 'loading' | 'done' | 'cancel'; id?: number }[]
  >([])
  const abortController = ref<AbortController>()
  const analyseRef = ref<Element>()

  // 设置AI助理窗口样式
  const setDrawerCss = (css: 'none' | 'all .3s') => {
    const el = document.querySelector('.ai-chat.el-drawer') as HTMLDivElement
    nextTick(() => {
      el.style.transition = css
    })
  }
  //打开AI助理
  const handleOpenAiBox = () => {
    if (isDragging) return
    getModelStatus()
  }
  //关闭AI助理
  const handleCloseAiBox = () => {
    setDrawerCss('all .3s')
    aiDrawerVisible.value = false
  }
  //获取AI模型开启状态
  const getModelStatus = async () => {
    const {
      data: { module_config },
    } = await getSystemConfigApi({ keys: 'module_config' })
    moduleEnable.value = JSON.parse(module_config.value).isEnable || false
    if (moduleEnable.value) {
      hanldeGetAiChatList()
      aiDrawerVisible.value = true
    } else {
      if (isOpen) return
      if (!moduleEnable.value) {
        isOpen = true
        $baseConfirm(
          'AI大模型暂未开启，是否开启？',
          null,
          async (action: any) => {
            const resolveRouter = router.resolve({
              path: '/managements/configuration/systemConfig',
              query: { params: 'modelConfig' },
            })
            window.open(resolveRouter.href, '_blank')
          },
          () => {
            isOpen = false
          }
        )
      }
    }
  }
  //根据ID获取当前对话
  const getCurrChatHistory = async (id: number) => {
    historyKey.value = id
    abortController.value?.abort()
    showSend.value = false
    isGenerating.value = false
    questionStr.value = ''
    const { data } = await aiSessionContentApi(id)
    questionList.value = data
  }
  // 获取Chart对话
  const hanldeGetAiChart = async (reloadAnswer: boolean) => {
    if (!questionStr.value) return
    //回答问题时当抽屉高度太小的时候展开
    if (drawerHeight.value < 300) {
      setDrawerCss('all .3s')
      drawerHeight.value = innerHeight * 0.55
      historyRef.value!.style.width = 'auto'
      chatPromptRef.value!.style.marginTop = '20px'
    }
    if (!reloadAnswer) {
      questionList.value.push({
        role: 'user',
        content: questionStr.value,
      })
    }
    isGenerating.value = true
    showSend.value = true
    abortController.value = new AbortController()
    const userStore = useUserStore()
    const { token } = userStore
    const history = questionList.value.toSpliced(-1).map(({ role, content }) => ({ role, content }))
    const questionHistoryStr = questionStr.value
    questionStr.value = ''
    let aiContent = ''
    let _chatId: number

    const fetchTimeOut = () =>
      new Promise((resolve, reject) => {
        setTimeout(() => {
          reject('请求超时')
        }, 80000)
        fetchEventSource('/v3/ecsPlatform/gpt/chatHistory', {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': contentType,
            Cookies: `EcsSessionId=${token}`,
          },
          openWhenHidden: true,
          signal: abortController.value!.signal,
          body: JSON.stringify({
            historyId: historyKey.value,
            content: {
              q: questionHistoryStr,
              history: history.slice(-20),
            },
            reAnswer: reloadAnswer,
          }),
          onmessage(event) {
            const length = questionList.value.length - 1
            showSend.value = false
            const message_data = JSON.parse(event.data) as {
              msg: string
              historyId: number
            }
            const text = message_data.msg
            _chatId = +message_data.historyId
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
            showSend.value = false
            const chat = questionList.value[length]
            chat!.status = 'done'
            isGenerating.value = false
            if (historyKey.value === _chatId) return
            historyKey.value = _chatId
            historyList.value.unshift({
              id: historyKey.value,
              title: questionHistoryStr,
            })
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
  //获取会话历史列表
  const hanldeGetAiChatList = async (isAdd = false) => {
    const { data } = await aiSessionChatsApi(historyListParams)
    const arr = data.records || []
    if (isAdd) {
      historyList.value.push(...arr)
    } else {
      historyList.value = arr
    }
    historyTotal.value = data.total || 0
  }
  //重新回答
  const handleAiReload = () => {
    questionList.value.pop()
    const length = questionList.value.length - 1
    const chat = questionList.value[length]
    questionStr.value = chat!.content
    showSend.value = true
    hanldeGetAiChart(true)
  }
  // 删除会话历史列表
  const handleDeleteAiChat = async (ids: number[], deleteAll = false) => {
    await deleteAiSessionContentApi({
      deleteAll,
      ids,
    })
    if (ids.includes(historyKey.value)) newChat()
    hanldeGetAiChatList()
  }

  // 删除所有会话历史列表
  const handleDeleteAllAiChat = () => {
    $baseConfirm('是否要删除所有AI会话？', null, async () => {
      const ids = historyList.value.map((i) => i.id)
      handleDeleteAiChat(ids, true)
    })
  }

  const newChat = () => {
    historyKey.value = undefined
    showSend.value = false
    isGenerating.value = false
    questionStr.value = ''
    questionList.value = []
    questionRef.value?.focus()
  }
  const resizeObserver = new ResizeObserver(() => {
    if (isGenerating.value === false) return
    const parentNode = analyseRef.value?.parentElement
    parentNode?.scrollTo({
      top: parentNode.scrollHeight,
      behavior: 'smooth',
    })
  })
  const handleAnalyseSearchKeydown = (event: KeyboardEvent) => {
    if (event.code === 'Enter' && !event.shiftKey) {
      event.preventDefault()
      if (!questionStr.value || isGenerating.value || isComposition.value) return
      hanldeGetAiChart(false)
    }
  }
  const handleCompositionStart = (val: any) => {
    isComposition.value = true
  }
  const handleCompositionEnd = (val: any) => {
    isComposition.value = false
  }
  const hanldeGetAiQuickInput = async () => {
    const params = {
      attackIp: '',
      victimIp: '',
      host: '',
      threatName: '',
    }
    const { data } = await getAiQuickInputApi(params)
    quickInputs.value = Array.isArray(data) ? data.sort((a, b) => a.length - b.length) : []
  }
  watch(
    () => aiDrawerVisible.value,
    () => {
      if (aiDrawerVisible.value) {
        nextTick(() => {
          const scrollDom = document.querySelector('ul.history-list') as HTMLElement
          const { arrivedState } = useScroll(scrollDom)
          scrollState.value = arrivedState
          resizeObserver.observe(analyseRef.value!)
        })
        setTimeout(() => {
          questionRef.value?.focus()
        })
      } else {
        newChat()
        resizeObserver.unobserve(analyseRef.value!)
        abortController.value?.abort()
      }
    }
  )
  watch(
    () => chatSzieBtnHeight.value,
    (moveY) => {
      const maxHeight = innerHeight
      const height = innerHeight - moveY
      if (height < 168) {
        chatPromptRef.value!.setAttribute('style', 'margin-top:0')
        historyRef.value!.setAttribute('style', 'width:0')
        return
      }
      chatPromptRef.value!.removeAttribute('style')
      historyRef.value!.setAttribute('style', `height:${drawerHeight.value - 40}px`)
      drawerHeight.value = height > maxHeight ? maxHeight : height
      setDrawerCss('none')
    }
  )
  watch(
    () => scrollState.value,
    () => {
      if (scrollState.value.bottom) {
        if (historyTotal.value > historyList.value.length) {
          historyListParams.pageNum++
          hanldeGetAiChatList(true)
        }
      }
    },
    { deep: true }
  )
  onMounted(() => {
    document.addEventListener(
      'click',
      (e) => {
        const target = e.target as HTMLElement
        const activeEl = document.activeElement as HTMLElement
        if (target.tagName === 'TEXTAREA') {
          if (!target.parentElement?.classList.contains('text')) {
            questionRef.value?.blur()
            target?.focus()
            return
          }
          if (target !== activeEl) {
            setTimeout(() => {
              activeEl?.blur()
            }, 0)
          }
          setTimeout(() => {
            questionRef.value?.focus()
          }, 0)
          return
        }
        if (!aiDrawerVisible.value || target === activeEl || !target.classList.contains('cm-content')) {
          questionRef.value?.blur()
          target?.focus()
        }
      },
      true
    )
    hanldeGetAiQuickInput()
  })
</script>

<template>
  <div class="ai-chat-box">
    <!-- AI窗口显示按钮 -->
    <div
      ref="aiChatRef"
      class="ai-chat-avatar"
      :style="{
        transform: `translateY(${aiChatRefTop}px)`,
      }"
    >
      <el-image :src="require('@/assets/alert_images/ai_avatar.svg')" @click.self="handleOpenAiBox" />
    </div>

    <el-drawer
      v-model="aiDrawerVisible"
      class="ai-chat"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      direction="btt"
      :size="drawerHeight"
      :with-header="false"
    >
      <!-- AI问答历史记录窗口 -->
      <div ref="historyRef" class="ai-history" :style="{ height: `${drawerHeight - 40}px` }">
        <div class="history-title">
          <img class="aiLogo" :src="require('@/assets/ai-logg.svg')" />
          <el-button link type="primary" @click="handleDeleteAllAiChat">一键清空</el-button>
        </div>
        <ul class="history-list">
          <template v-if="historyList.length">
            <li
              v-for="chat in historyList"
              :key="chat.id"
              :class="{ active: chat.id === historyKey }"
              @click.self="() => getCurrChatHistory(chat.id)"
            >
              {{ chat.title }}
              <div class="deleteIcon">
                <el-icon @click="() => handleDeleteAiChat([chat.id])"><Delete /></el-icon>
              </div>
            </li>
          </template>
          <li v-else class="history-empty">暂无历史记录</li>
        </ul>
        <el-button class="history-add" :disabled="isGenerating" :icon="Plus" type="primary" @click="newChat">
          新建对话
        </el-button>
      </div>
      <!-- AI问答窗口 -->
      <div class="ai-chat" :style="{ height: `${drawerHeight - 40}px` }">
        <div class="ai-chat-list-box">
          <div ref="analyseRef" class="ai-chat-list">
            <div v-if="questionList.length === 0" class="chat-item help">
              HI，我是您的AI助理，
              我能帮助您做告警研判、调查取证，我还知道如何防范各类攻击以及漏洞修复的方法，同时我还能协助您写各类报告。您有什么需要我做的，尽管吩咐我就好了。
            </div>
            <div
              v-for="(question, index) in questionList"
              :key="question.id"
              class="chat-item"
              :class="{ answer: question.role === 'assistant', question: question.role === 'user' }"
            >
              <ai-chat-answer :answer-content="question.content" :role="question.role" />
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
            </div>
            <!-- 发送中 -->
            <div v-if="showSend" class="chat-item answer">
              <div style="width: 50px">
                <div class="loader"></div>
              </div>
            </div>
          </div>
        </div>
        <div ref="chatPromptRef" class="ai-chat-prompt">
          <div class="tip">
            <svg
              height="20px"
              style="vertical-align: sub; margin-right: 10px"
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
            <span
              v-for="(question, index) in quickInputs.slice(0, 5)"
              :key="index"
              @click="() => (questionStr = question)"
            >
              {{ question }}
            </span>
            <el-popover placement="top-end" :show-arrow="false" :teleported="false" trigger="hover" width="max-content">
              <template #reference>
                <span class="exchange">
                  <el-icon style="vertical-align: -2px; margin-left: -1px"><MoreFilled /></el-icon>
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
                      v-for="(question, index) in quickInputs.slice(5)"
                      :key="index"
                      class="quickActionItem"
                      @click="() => (questionStr = question)"
                    >
                      {{ question }}
                    </li>
                  </ul>
                </div>
              </template>
            </el-popover>
          </div>
          <div class="send-box">
            <el-input
              ref="questionRef"
              v-model="questionStr"
              :autosize="{ minRows: 1, maxRows: 7 }"
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
              class="send-btn"
              :disabled="!questionStr || isGenerating"
              size="small"
              type="primary"
              @click="() => hanldeGetAiChart(false)"
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
      </div>
      <!-- AI窗口拖拽按钮 -->
      <div ref="chatSzieBtn" class="ai-chat-size">
        <svg height="16px" version="1.1" viewBox="0 0 30 16" width="30px">
          <g id="告警大模型" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
            <g id="全局AI大模型2" transform="translate(-967, -364)">
              <g id="ladong" transform="translate(982, 372) rotate(90) translate(-982, -372)translate(974, 357)">
                <rect id="矩形" fill="#FFFFFF" height="29" rx="2" stroke="#E6E6F1" width="15" x="0.5" y="0.5" />
                <path
                  id="形状"
                  d="M7,21 C7,20.4477153 6.55228475,20 6,20 C5.44771525,20 5,20.4477153 5,21 C5,21.5522847 5.44771525,22 6,22 C6.53653388,21.9637079 6.9637079,21.5365339 7,21 L7,21 Z M7,17 C7,16.4477153 6.55228475,16 6,16 C5.44771525,16 5,16.4477153 5,17 C5,17.5522847 5.44771525,18 6,18 C6.53653388,17.9637079 6.9637079,17.5365339 7,17 L7,17 Z M7,13 C7,12.4477153 6.55228475,12 6,12 C5.44771525,12 5,12.4477153 5,13 C5,13.5522847 5.44771525,14 6,14 C6.53653388,13.9637079 6.9637079,13.5365339 7,13 L7,13 Z M7,9 C7,8.44771525 6.55228475,8 6,8 C5.44771525,8 5,8.44771525 5,9 C5,9.55228475 5.44771525,10 6,10 C6.53653388,9.9637079 6.9637079,9.53653388 7,9 L7,9 Z M11,21 C11,20.4477153 10.5522847,20 10,20 C9.44771525,20 9,20.4477153 9,21 C9,21.5522847 9.44771525,22 10,22 C10.5365339,21.9637079 10.9637079,21.5365339 11,21 L11,21 Z M11,17 C11,16.4477153 10.5522847,16 10,16 C9.44771525,16 9,16.4477153 9,17 C9,17.5522847 9.44771525,18 10,18 C10.5365339,17.9637079 10.9637079,17.5365339 11,17 L11,17 Z M11,13 C11,12.4477153 10.5522847,12 10,12 C9.44771525,12 9,12.4477153 9,13 C9,13.5522847 9.44771525,14 10,14 C10.5365339,13.9637079 10.9637079,13.5365339 11,13 L11,13 Z M11,9 C11,8.44771525 10.5522847,8 10,8 C9.44771525,8 9,8.44771525 9,9 C9,9.55228475 9.44771525,10 10,10 C10.5365339,9.9637079 10.9637079,9.53653388 11,9 L11,9 Z"
                  fill="#9793B0"
                  fill-rule="nonzero"
                />
              </g>
            </g>
          </g>
        </svg>
      </div>
      <!-- AI窗口关闭按钮 -->
      <el-icon class="closeBtn" @click="handleCloseAiBox">
        <CircleCloseFilled />
      </el-icon>
    </el-drawer>
  </div>
</template>

<style scoped lang="scss">
  .ai-chat-box {
    z-index: 999;
    .ai-history {
      width: 310px;
      height: 100%;
      text-indent: 10px;
      margin-inline: 20px;
      user-select: none;
      overflow: hidden;
      .history-title {
        height: 50px;
        font-weight: 500;
        font-size: 18px;
        color: #1e1842;
        line-height: 40px;
        display: flex;
        position: relative;
        align-items: center;
        justify-content: space-between;
      }
      .history-list {
        font-weight: 400;
        font-size: 15px;
        color: #4a4759;
        padding-inline-start: 0;
        height: calc(100% - 130px);
        overflow-y: auto;
        scrollbar-width: none;
        position: relative;
        width: 100%;
        li {
          height: 44px;
          line-height: 44px;
          border-radius: 6px;
          cursor: pointer;
          position: relative;
          word-break: keep-all;
          text-overflow: ellipsis;
          overflow: hidden;
          margin-bottom: 2px;
          &.active {
            background: #e3dfff;
            font-weight: 500;
            color: #6954f0;
          }
          &.history-empty {
            position: absolute;
            top: 50%;
            margin-top: -22px;
            text-align: center;
            width: inherit;
            text-indent: 0;
          }
          &:not(.history-empty) {
            &:hover {
              background: #e3dfff;
              .deleteIcon {
                display: flex;
              }
            }
          }
          .history-text {
            padding-right: 10px;
            display: inline-block;
            word-break: break-all;
            text-overflow: ellipsis;
            overflow: hidden;
            height: 44px;
          }
          .deleteIcon {
            width: 28px;
            height: 44px;
            color: #6954f0;
            position: absolute;
            right: 0;
            top: 0;
            justify-content: center;
            align-items: center;
            display: none;
            background: #e3dfff;
          }
        }
      }
      .history-add {
        width: 130px;
        height: 40px;
        background: #6954f0;
        border-radius: 6px;
        margin-left: 70px;
      }
    }
    .ai-chat {
      flex: 1;
      height: 100%;
      background: #ffffff;
      border-radius: 18px;
      max-height: 100vh;
      display: flex;
      padding: 20px 50px;
      flex-direction: column;
      container-type: inline-size;
      .ai-chat-list-box {
        height: calc(100% - 100px);
        overflow-y: auto;
        scrollbar-width: none;
        margin-left: 10px;
      }
      .ai-chat-list {
        width: 100%;
        display: flex;
        flex-direction: column;
        .chat-item {
          font-weight: 400;
          font-size: 14px;
          color: #2b2742;
          line-height: 24px;
          border-radius: 8px;
          background: #f5f4ff;
          max-width: 80%;
          word-wrap: break-word;
          padding: 10px 20px;
          border-radius: 8px;
          position: relative;
          transition: all 0.3s ease;
          clear: both;
          resize: horizontal;
          &:not(:last-child) {
            margin-bottom: 14px;
          }
          .ai-reload {
            position: absolute;
            right: -60px;
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
          &.answer,
          &.help {
            margin-left: 36px;
            &::after {
              content: ' ';
              display: block;
              width: 30px;
              height: 30px;
              background: url('@/assets/alert_images/ai_avatar.svg') no-repeat center center / 30px 30px;
              position: absolute;
              top: 7px;
              left: -36px;
            }
          }
          &.help {
            width: 720px;
          }
          &.answer {
            max-width: 80%;
            background: none;
          }
          &.question {
            align-self: flex-end;
            background: #d4cdff;
          }
        }
      }
      .ai-chat-prompt {
        width: 100%;
        margin-top: 20px;
        transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
        .send-box {
          position: relative;
          margin: 0 auto;
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
          .send-btn {
            width: 40px;
            height: 40px;
            background: #6954f0;
            border-radius: 100%;
            &.is-disabled {
              filter: grayscale(50%);
              -webkit-filter: grayscale(10%);
              opacity: 0.5;
            }
          }
        }
        .tip {
          margin: 0 auto 10px;
          font-weight: 400;
          font-size: 13px;
          color: #888b91;
          span:not(.title) {
            height: 28px;
            text-align: center;
            line-height: 28px;
            color: #494758;
            margin-right: 8px;
            background: #ffffff;
            border-radius: 6px;
            padding-inline: 8px;
            border: 1px solid #dedcee;
            cursor: pointer;
            display: inline-block;
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
            color: var(--el-color-primary);
            width: 28px;
            height: 28px;
            display: inline-block;
          }
        }
      }
    }
    .ai-chat-avatar {
      position: fixed;
      width: 52px;
      height: 70px;
      right: -24px;
      top: 0;
      overflow: visible;
      cursor: pointer;
      z-index: 9999;
      user-select: none;
      transform: translateY(calc(100vh - 150px));
      transition: all 0.3s;
      &::after {
        content: 'AI智能体';
        display: block;
        position: absolute;
        bottom: 2px;
        left: -2px;
        width: 56px;
        height: 24px;
        background: #695fa9;
        border-radius: 12px;
        font-weight: 500;
        font-size: 12px;
        color: #ffffff;
        text-align: center;
        line-height: 24px;
      }
      &::before {
        content: ' ';
        display: block;
        inset: -10px -15px;
        border-radius: 100%;
        position: absolute;
      }
      &:hover {
        right: 10px;
      }
    }
    .ai-chat-size {
      width: 30px;
      height: 16px;
      cursor: ns-resize;
      position: absolute;
      top: -8px;
      left: 50%;
      transform: translateX(-50%);
      &::after {
        content: ' ';
        display: block;
        position: absolute;
        inset: -10px -20px;
      }
    }
    .closeBtn {
      position: absolute;
      top: 5px;
      right: 5px;
      width: 28px;
      height: 28px;
      font-size: 28px;
      color: #857eae;
      cursor: pointer;
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
    :deep() {
      .ai-chat.el-drawer {
        background-color: transparent;
        box-shadow: none;
      }
      .el-image__inner {
        width: 52px;
        height: 52px;
      }
      .aiLogo {
        width: 110px;
        height: 28px;
      }
      .el-drawer__body {
        padding: 10px 10px 10px 0;
        background: linear-gradient(180deg, #efedff 0%, #f1eaff 32%, #ecf2ff 54%, #e9ebff 100%);
        margin-top: 15px;
        position: relative;
        overflow: visible;
        display: flex;
        box-shadow: 0px 12px 20px 0px rgba(0, 0, 0, 0.15), 0px 8px 20px rgba(0, 0, 0, 0.4);
      }
      .el-overlay {
        background: none !important;
        backdrop-filter: none;
        z-index: 9999999 !important;
        pointer-events: none !important;
        & > div {
          pointer-events: all;
        }
      }
    }
  }
</style>
