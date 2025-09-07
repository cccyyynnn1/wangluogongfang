<script lang="ts">
  export default {
    name: 'AiChatAnswer',
  }
</script>

<script setup lang="ts">
  import hljs from 'highlight.js/lib/core'
  import python from 'highlight.js/lib/languages/python'
  import javascript from 'highlight.js/lib/languages/javascript'
  import go from 'highlight.js/lib/languages/go'
  import java from 'highlight.js/lib/languages/java'
  import plaintext from 'highlight.js/lib/languages/plaintext'
  import html from 'highlight.js/lib/languages/xml'
  hljs.registerLanguage('javascript', javascript)
  hljs.registerLanguage('python', python)
  hljs.registerLanguage('go', go)
  hljs.registerLanguage('java', java)
  hljs.registerLanguage('plaintext', plaintext)
  hljs.registerLanguage('html', html)
  const props = defineProps<{
    answerContent?: string
    role: 'assistant' | 'user'
  }>()
  const { marked } = require('marked')
  // 创建自定义渲染器
  const renderer = new marked.Renderer()
  renderer.link = (href: string, title: string, text: string) => text
  renderer.code = (code: string, infostring: string) => {
    const language = hljs.getLanguage(infostring) ? infostring : 'plaintext'
    const hljsVal = hljs.highlight(code, { language, ignoreIllegals: true }).value
    return `<pre><code class="hljs language-${language}">${hljsVal}</code></pre>`
  }
  const startTag = `<think>\n`
  const endTag = `\n</think>`

  // 使用自定义的 renderer 选项
  const markedOptions = {
    renderer: renderer,
    gfm: true,
    headerIds: false,
    mangle: false,
  }
  const thinkText = ref('')
  const answerText = ref('')
  const showThink = ref(true)
  const getHtmlContentFromMarked = (content: string) => {
    return marked(content, markedOptions)
  }

  const changeVisible = () => {
    showThink.value = !showThink.value
  }
  const contentProcessorHandler = (content: string) => {
    const thinkIndex = content.indexOf(startTag)
    const thinkEndIndex = content.indexOf(endTag)
    if (thinkIndex !== -1) {
      thinkText.value = content.slice(thinkIndex + startTag.length, thinkEndIndex)
    }
    const answer = content.slice(thinkEndIndex + endTag.length)
    answerText.value = thinkIndex === -1 ? content : thinkEndIndex !== -1 ? answer : ''
    showThink.value = answerText.value.length > 0 ? false : true
  }
  watch(
    () => props.answerContent,
    (newVal) => {
      contentProcessorHandler(newVal || '')
    },
    {
      immediate: true,
    }
  )
</script>

<template>
  <div v-if="role === 'assistant'" class="answer-container">
    <div class="answer-title" @click="changeVisible">
      已深度思考
      <el-icon size="14" style="margin-left: 5px">
        <ArrowDownBold v-if="!showThink" />
        <ArrowUpBold v-else />
      </el-icon>
    </div>
    <div v-if="showThink" class="answer-think" v-html="getHtmlContentFromMarked(thinkText)" />
    <div v-if="answerText" class="answer-content" v-html="getHtmlContentFromMarked(answerText)" />
  </div>
  <template v-else>{{ answerContent }}</template>
</template>

<style scoped lang="scss">
  .answer-container {
    .answer-title {
      font-weight: 500;
      font-size: 14px;
      color: #2b2742;
      cursor: pointer;
      margin-bottom: 1em;
      display: flex;
      align-items: center;
    }
    .answer-think {
      border-left: 1px solid #dbd9e7;
      padding-left: 12px;
      margin-bottom: 1em;
      font-size: 14px;
      line-height: 26px;
      color: #8b8b8b;
      ::v-deep * {
        margin-block: 0.5em;
      }
    }
    .answer-content {
      background: #f1eeff;
      padding: 14px 20px;
      border-radius: 8px;
    }
  }
</style>
