<script lang="ts">
  export default {
    name: 'D3ForceGraph',
  }
</script>

<script setup lang="ts">
  import * as d3 from 'd3'
  import install from './d3-plugins/d3-context-menu'
  install(d3) // 为d3注册右键菜单插件
  interface ContextMenu {
    title: string
    action: (elm: any, nodeData: any) => void
    disabled?: boolean
  }

  interface Props {
    graphData: {
      nodes: any[]
      links: any[]
    }
    isTraceability?: boolean
    contextMenu?: ContextMenu[]
    /**
     * @param disableZoomKeys 禁止缩放的kes Event.type 如：'wheel' (滚轮)
     */
    disableKeys?: string[]
  }
  const symbolSize = 60

  const props = defineProps<Props>()
  const emits = defineEmits<{
    (e: 'click', data: any): void
    (e: 'graph-move'): void
  }>()
  let my_svg: any = null
  let my_tooltip: any = null
  let transformData = {
    k: 1,
    x: 0,
    y: 0,
  }
  const d3_force_ref = ref()
  const d3_instance = ref()
  const d3_links = ref<any[]>([])
  const d3_nodes = ref<any[]>([])
  const svg_zoom = d3
    .zoom()
    .scaleExtent([0.5, 1.5])
    .on('zoom', function (this: any, event: any) {
      transformData = event.transform
      emits('graph-move')
      d3.select('.group').attr('transform', d3.zoomTransform(my_svg.node()))
    })
    .filter((event: any) => {
      const filterKeys = props.disableKeys || []
      return !filterKeys.includes(event.type)
    })
  function renderTooltip() {
    d3.select('.d3-tooltip').remove()
    my_tooltip = d3
      .select(my_svg.node().parentElement.parentElement)
      .append('div')
      .attr('class', 'd3-tooltip')
      .style('opacity', 0)
  }
  const d3_drag = (simulation: any) => {
    function dragsubject(event: any) {
      return simulation.find(event.x, event.y)
    }

    function dragstarted(event: any) {
      if (!event.active) simulation.alphaTarget(0.3).restart()
      event.subject.fx = event.subject.x
      event.subject.fy = event.subject.y
    }

    function dragged(event: any) {
      event.subject.fx = event.x
      event.subject.fy = event.y
    }

    function dragended(event: any) {
      if (!event.active) simulation.alphaTarget(-0.1)
      // 注释以下代码，使拖动结束后固定节点
      // event.subject.fx = null
      // event.subject.fy = null
    }

    return d3.drag().subject(dragsubject).on('start', dragstarted).on('drag', dragged).on('end', dragended)
  }
  function changeGraphStyle(selectData: string | any) {
    const _val =
      typeof selectData === 'string'
        ? { clientIp: selectData, serverIp: selectData }
        : { clientIp: selectData.clientIp, serverIp: selectData.serverIp }

    const list =
      typeof selectData === 'string'
        ? d3_links.value
            .filter((link) => link.clientIp === _val.clientIp || link.serverIp === _val.serverIp)
            .map((i) => [i.clientIp, i.serverIp])
            .flat()
        : [selectData.clientIp, selectData.serverIp]

    d3_instance.value
      .select('.nodes')
      .selectAll('.node')
      .attr('class', (d: any) => {
        if (d.id === _val.clientIp) return typeof selectData === 'string' ? 'fixed' : 'active'
        if (list.includes(d.id)) return 'active'
        return ''
      })
    // 处理相邻的边line
    d3_instance.value
      .select('.links')
      .selectAll('line')
      .attr('class', (d: any) => {
        if (typeof selectData === 'string')
          return d.clientIp === _val.clientIp || d.serverIp === _val.clientIp ? 'active' : ''
        return d.clientIp === _val.clientIp && d.serverIp === _val.serverIp ? 'active' : ''
      })
  }
  function clearGraphStyle() {
    // 移除所有样式
    my_svg.select('.nodes').selectAll('circle').attr('class', '')
    my_svg.select('.texts').selectAll('text').attr('class', '')
    my_svg.select('.links').selectAll('line').attr('class', '').attr('marker-end', 'url(#posMarker)')
    my_svg.select('.linkTexts').selectAll('text').attr('class', '')
  }
  function textBreaking(d3text: any, text: string) {
    const texts = text.split('\n')
    const topY = -16
    const midY = 0
    const botY = 16
    texts.forEach((text, index) => {
      d3text
        .append('tspan')
        .attr('x', -1)
        .attr('y', 50 + index * 16)
        .text(text)
    })
  }
  function d3init() {
    const { links, nodes } = props.graphData
    if (nodes.length == 0) return
    d3_links.value = JSON.parse(JSON.stringify(links))
    d3_nodes.value = JSON.parse(JSON.stringify(nodes))
    d3_instance.value = d3.select(d3_force_ref.value)
    const svg = d3_instance.value.style('font', '12px sans-serif').call(svg_zoom).append('g').attr('class', 'group')
    my_svg = svg
    const simulation = d3
      .forceSimulation(d3_nodes.value)
      .force(
        'link',
        d3
          .forceLink(d3_links.value)
          .id((d: any) => d.id)
          .distance(10)
      )
      .force('charge', d3.forceManyBody().strength(10).distanceMin(30).distanceMax(70))
      .force('collision', d3.forceCollide(2).radius(80).iterations(0.5))
      .force(
        'center',
        d3.forceCenter(svg.node().parentElement.clientWidth / 2, svg.node().parentElement.clientHeight / 2)
      )
    renderTooltip()
    renderMarkers(svg)
    const link = svg
      .append('g')
      .attr('class', 'links')
      .attr('fill', 'transparent')
      .attr('stroke-width', 2)
      .selectAll('path')
      .data(d3_links.value)
      .join('path')
      .attr('stroke', (d: any) => {
        return d?.alarm?.length > 0 ? '#F56C6C' : '#4E7CBE'
      })
      .attr('cursor', 'pointer')
      .attr('marker-end', 'url(#posMarker)')
      .on('click', function (this: any, event: any) {
        const _node = d3.select(this)
        const nodeData: any = _node.data()[0]
        my_tooltip.style('opacity', 0)
        // my_tooltip.style('dis', 0)
        emits('click', nodeData)
      })
      .on('mouseenter', function (this: any, event: any) {
        const _link = d3.select(this)
        const nodeData: any = _link.data()[0]
        const serverPort = Array.isArray(nodeData?.serverPort)
          ? nodeData?.serverPort
          : Object.entries(nodeData?.serverPort)
        if (serverPort.length || nodeData?.alarm?.length) {
          my_tooltip
            .html(
              `<div>
                 ${serverPort
                   .map(
                     (item: any) =>
                       `${
                         typeof item === 'string'
                           ? `<p style="height:14px;margin:0 0 10px 0">目的端口：${item}</p>`
                           : `${item[0]}:${item[1]?.protocolStr} `
                       }`
                   )
                   .join('')}
                ${
                  nodeData?.alarm?.length
                    ? `<p style="height:14px;margin:0 0 10px 0">告警类型：${nodeData.alarm}</p>`
                    : ''
                }
                ${
                  props.isTraceability
                    ? nodeData.count.map((item: string) => `<p style="height:14px;margin:0 0 10px 0">${item}</p>`)
                    : ''
                }
              </div>`
            )
            .style('transform', `translate(calc(${event.layerX}px - 50%), calc(${event.layerY}px - 110%)`)
          nextTick(() => my_tooltip.style('opacity', 1))
        }
        // 遍历节点，并调整图的样式
        changeGraphStyle(nodeData)
      })
      .on('mouseleave', (event: MouseEvent) => {
        my_tooltip.style('opacity', 0)
        my_tooltip.style('display', 'none')
        clearGraphStyle()
      })
    const node = svg
      .append('g')
      .attr('class', 'nodes')
      .selectAll('.node')
      .data(d3_nodes.value)
      .enter()
      .append('g')
      .attr('class', 'node')
      .call(d3_drag(simulation))
      .on('click', function (this: any, event: any) {
        const _node = d3.select(this)
        const nodeData: any = _node.data()[0]
        my_tooltip.style('opacity', 0).style('display', 'none')
        emits('click', nodeData)
      })
      .on('mouseenter', function (this: any, event: any) {
        const _node = d3.select(this)
        const nodeData: any = _node.data()[0]
        my_tooltip
          .html(`<div style="white-space:pre-line;max-width:500px">${nodeData?.IpLocation || nodeData?.ip}</div>`)
          .style('opacity', 1)
          .style('display', 'block')
          .style('transform', `translate(calc(${event.layerX}px - 50%), calc(${event.layerY}px - 130%)`)

        // 遍历节点，并调整图的样式
        changeGraphStyle(nodeData.id)
      })
      .on('mouseleave', (event: MouseEvent) => {
        my_tooltip.style('opacity', 0).style('display', 'none')
        clearGraphStyle()
      })
    if (props.contextMenu && props.contextMenu.length > 0) node.on('contextmenu', d3.contextMenu(props.contextMenu))
    node
      .append('circle')
      .attr('r', symbolSize / 2)
      .attr('fill', (d: any) => (d.color === 'red' ? '#FFF0F0' : '#B2DEF9'))
      .attr('stroke', (d: any) => (d.color === 'red' ? '#F56C6C' : '#6C9FDA'))
    node
      .append('use')
      .attr('height', (d: any) => (d.icon !== 'client' || d.icon !== 'server' ? 32 : symbolSize))
      .attr('width', (d: any) => (d.icon !== 'client' || d.icon !== 'server' ? 32 : symbolSize))
      .attr('transform', function (this: any, d: any) {
        if (d.icon !== 'client' || d.icon !== 'server') return `translate( -16 , -16 )`
        return `translate( -${symbolSize / 2} , -${symbolSize / 2} )`
      })
      .attr('cursor', 'pointer')
      .attr('class', (d: any) => (d.color === 'red' ? 'errNode' : ''))
      .attr('xlink:href', (d: any) => `#${d.icon ? d.icon : 'client'}`)
    const texts = node.append('g').attr('class', 'text').attr('text-anchor', 'middle').attr('dominant-baseline', 'top')
    texts.append('rect').attr('class', 'host-text-bg').attr('fill', 'transparent').attr('height', '18px')
    texts
      .append('text')
      .attr('id', (d: any) => d.id)
      .attr('class', 'host-name')
      .attr('x', function (this: any, d: any) {
        return textBreaking(d3.select(this), d.name || '')
      })
    // .attr('y', 100)
    simulation.on('tick', () => {
      link
        .attr('d', function (d: any) {
          if (d.source.x < d.target.x) {
            return `M ${d.source.x} ${d.source.y} L ${d.target.x} ${d.target.y}`
          } else {
            return `M ${d.target.x} ${d.target.y} L ${d.source.x} ${d.source.y}`
          }
        })
        .attr('marker-end', function (d: any) {
          if (d.source.x < d.target.x) {
            return 'url(#posMarker)'
            // return d.alarm.length > 0 ? 'url(#errPosMarker)' : 'url(#posMarker)'
          } else {
            return null
          }
        })
        .attr('marker-start', function (d: any) {
          if (d.source.x < d.target.x) {
            return null
          } else {
            // return d.alarm.length > 0 ? 'url(#errNegativeMarker)' : 'url(#negativeMarker)'
            return 'url(#negativeMarker)'
          }
        })
      node.attr('transform', (d: any) => `translate(${d.x},${d.y})`)
    })
    simulation.force('link').links(d3_links.value).distance(190)
  }
  // 绘制关系箭头
  function renderMarkers(svg: any) {
    // 定义箭头的标识
    const defs = svg.append('defs')
    const posMarker = defs
      .append('marker')
      .attr('id', 'posMarker')
      .attr('orient', 'auto')
      .attr('stroke-width', 2)
      .attr('markerUnits', 'strokeWidth')
      .attr('markerUnits', 'userSpaceOnUse')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 35)
      .attr('refY', 0)
      .attr('markerWidth', 12)
      .attr('markerHeight', 12)
      .append('path')
      .attr('d', 'M 0 -5 L 10 0 L 0 5')
      .attr('fill', '#4E7CBE')
      .attr('stroke-opacity', 0.6)
    const errPosMarker = defs
      .append('marker')
      .attr('id', 'errPosMarker')
      .attr('orient', 'auto')
      .attr('stroke-width', 2)
      .attr('markerUnits', 'strokeWidth')
      .attr('markerUnits', 'userSpaceOnUse')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 35)
      .attr('refY', 0)
      .attr('markerWidth', 12)
      .attr('markerHeight', 12)
      .append('path')
      .attr('d', 'M 0 -5 L 10 0 L 0 5')
      .attr('fill', '#F56C6C')
      .attr('stroke-opacity', 0.6)
    const errNegativeMarker = defs
      .append('marker')
      .attr('id', 'errNegativeMarker')
      .attr('orient', 'auto')
      .attr('stroke-width', 2)
      .attr('markerUnits', 'strokeWidth')
      .attr('markerUnits', 'userSpaceOnUse')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', -25)
      .attr('refY', 0)
      .attr('markerWidth', 12)
      .attr('markerHeight', 12)
      .append('path')
      .attr('d', 'M 10 -5 L 0 0 L 10 5')
      .attr('fill', '#F56C6C')
      .attr('stroke-opacity', 0.6)

    const negativeMarker = defs
      .append('marker')
      .attr('id', 'negativeMarker')
      .attr('orient', 'auto')
      .attr('stroke-width', 2)
      .attr('markerUnits', 'strokeWidth')
      .attr('markerUnits', 'userSpaceOnUse')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', -25)
      .attr('refY', 0)
      .attr('markerWidth', 12)
      .attr('markerHeight', 12)
      .append('path')
      .attr('d', 'M 10 -5 L 0 0 L 10 5')
      .attr('fill', '#4E7CBE')
      .attr('stroke-opacity', 0.6)
  }
  // 清除事件监听防止内存泄漏
  const removeListener = () => {
    if (!d3_instance.value) return
    d3_instance.value.on('.', null)
    d3_instance.value.selectAll('*').remove().on('.', null)
  }
  const zoomHandle = (zoom: number) => {
    nextTick(() => {
      transformData.k = zoom <= 0.5 ? 0.5 : zoom >= 1.5 ? 1.5 : zoom
      d3_instance.value
        .transition()
        .duration(1000)
        .call(svg_zoom.transform, d3.zoomIdentity.translate(transformData.x, transformData.y).scale(transformData.k))
    })
  }
  watch(
    () => props.graphData,
    () => {
      removeListener()
      d3init()
    },
    {
      deep: true,
    }
  )
  onUnmounted(() => {
    removeListener()
    d3_force_ref.value = null
  })

  defineExpose({
    reset: () => {
      transformData.k = 1
      zoomHandle(transformData.k)
    },
  })
</script>

<template>
  <div class="d3_force_box">
    <div class="zoomBtn" @click.stop="zoomHandle((transformData.k += 0.1))">+</div>
    <div class="zoomBtn scale" @click.stop="zoomHandle((transformData.k -= 0.1))">-</div>
    <svg
      aria-hidden="true"
      style="position: absolute; width: 0; height: 0"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
    >
      <symbol id="server" viewBox="0 0 80 80" xmlns="http://www.w3.org/2000/svg">
        <path d="M20.5 37.833V30.5h39v7.333h-39z" fill="#fff" stroke="#154EAA" />
        <path
          d="M55 33.333h-5V35h5v-1.667zM58.333 36.667H21.667v5h36.666v-5zM58.333 46.667H21.667v5h36.666v-5z"
          fill="#154EAA"
        />
        <path d="M20.5 47.833V40.5h39v7.333h-39z" fill="#fff" stroke="#154EAA" />
        <path d="M55 43.333h-5V45h5v-1.667z" fill="#154EAA" />
        <path d="M20.5 57.833V50.5h39v7.333h-39z" fill="#fff" stroke="#154EAA" />
        <path d="M55 53.333h-5V55h5v-1.667zM25 21.667h30L60 30H20l5-8.333z" fill="#154EAA" />
      </symbol>
      <symbol id="client" viewBox="0 0 80 80" xmlns="http://www.w3.org/2000/svg">
        <path d="M20.5 37.833V30.5h39v7.333h-39z" fill="#fff" stroke="#154EAA" />
        <path
          d="M55 33.333h-5V35h5v-1.667zM58.333 36.667H21.667v5h36.666v-5zM58.333 46.667H21.667v5h36.666v-5z"
          fill="#154EAA"
        />
        <path d="M20.5 47.833V40.5h39v7.333h-39z" fill="#fff" stroke="#154EAA" />
        <path d="M55 43.333h-5V45h5v-1.667z" fill="#154EAA" />
        <path d="M20.5 57.833V50.5h39v7.333h-39z" fill="#fff" stroke="#154EAA" />
        <path d="M55 53.333h-5V55h5v-1.667zM25 21.667h30L60 30H20l5-8.333z" fill="#154EAA" />
      </symbol>

      <svg
        id="web"
        height="60px"
        version="1.1"
        viewBox="0 0 60 60"
        width="60px"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
      >
        <g id="页面优化" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
          <g id="算法" transform="translate(-651.000000, -653.000000)">
            <g id="编组" transform="translate(612.000000, 614.000000)">
              <g id="chrome-fill-" transform="translate(40.000000, 40.000000)">
                <circle id="椭圆形" cx="28.9998438" cy="28.9998438" fill="#FFFFFF" r="28.9998438" />
                <path
                  id="形状"
                  d="M14.6141088,25.3186761 L6.34371582,10.9311286 C9.02499304,7.53210525 12.3478918,4.86049465 16.3124121,2.91629679 C20.2769324,0.972098929 24.5060763,0 28.9998438,0 C34.2113573,0 39.034998,1.28384725 43.4707658,3.85154175 C47.9065336,6.41923625 51.4088689,9.8744468 53.9777717,14.2171734 L30.3012117,14.2171734 C29.8855473,14.1797153 29.4511538,14.1609862 28.9980313,14.1609862 C25.5990079,14.1609862 22.5594618,15.2086056 19.8793929,17.3038443 C17.199324,19.399083 15.4436251,22.0706936 14.6122963,25.3186761 L14.6141088,25.3186761 Z M39.3654754,18.4076508 L56.0186357,18.4076508 C57.3405452,21.8066742 58.0015,25.3374052 58.0015,28.9998438 C58.0015,32.8894478 57.2462957,36.6183444 55.7358872,40.1865335 C54.2254787,43.7547226 52.1960938,46.832331 49.6477325,49.4193587 C47.0993712,52.0063865 44.0501585,54.0738337 40.5000943,55.6217003 C36.9500301,57.169567 33.2308001,57.9628335 29.3424044,58.0015 L41.1797781,37.4406108 C42.9548102,34.9103744 43.8423263,32.0973896 43.8423263,29.0016562 C43.8423263,24.8474286 42.3506468,21.3166977 39.3672879,18.4094633 L39.3654754,18.4076508 Z M18.4656505,28.9998438 C18.4656505,26.0926094 19.4945408,23.610102 21.5523214,21.5523214 C23.610102,19.4945408 26.0926094,18.4656505 28.9998438,18.4656505 C31.9070781,18.4656505 34.3895855,19.4945408 36.4473661,21.5523214 C38.5051467,23.610102 39.534037,26.0926094 39.534037,28.9998438 C39.534037,31.9070781 38.5051467,34.3895855 36.4473661,36.4473661 C34.3895855,38.5051467 31.9070781,39.534037 28.9998438,39.534037 C26.0926094,39.534037 23.610102,38.5051467 21.5523214,36.4473661 C19.4945408,34.3895855 18.4656505,31.9070781 18.4656505,28.9998438 Z M32.9655724,43.3293916 L24.6951794,57.6589393 C20.0878293,56.9786513 15.8961435,55.2984729 12.1201222,52.618404 C8.34410088,49.9383351 5.38007518,46.5211869 3.22804511,42.3669592 C1.07601504,38.2127316 0,33.7570265 0,28.9998438 C0,23.9019129 1.24638912,19.1628551 3.73916735,14.7826704 L15.5765411,35.2873724 C16.7848679,37.8550669 18.5973581,39.9225141 21.0140118,41.489714 C23.4306654,43.0569139 26.0926094,43.8405138 28.9998438,43.8405138 C30.3592114,43.8405138 31.681121,43.6707439 32.9655724,43.331204 L32.9655724,43.3293916 Z"
                  fill="#0049AA"
                  fill-rule="nonzero"
                />
              </g>
            </g>
          </g>
        </g>
      </svg>
      <svg
        id="pc"
        height="60px"
        version="1.1"
        viewBox="0 0 60 60"
        width="60px"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
      >
        <g id="页面优化" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
          <g id="算法" transform="translate(-846.000000, -482.000000)">
            <g id="编组" transform="translate(804.000000, 441.000000)">
              <g id="computer-fill-" transform="translate(47.000000, 49.000000)">
                <rect id="矩形" fill="#FFFFFF" height="15" rx="4" width="50" x="0" y="20.9090909" />
                <path
                  id="形状"
                  d="M44.288214,37.3687338 L31.4148231,37.3687338 L31.4148231,43.1201927 L34.2644669,43.1201927 C34.7831521,43.1012941 35.2768404,43.3721734 35.5393076,43.8257386 C35.808024,44.2793038 35.808024,44.8399608 35.5393076,45.293526 C35.2705912,45.7470912 34.7831521,46.0179704 34.2644669,45.9990719 L17.1541057,45.9990719 C16.3667042,45.9990719 15.7292838,45.3565212 15.7292838,44.562782 C15.7292838,43.7690429 16.3667042,43.1264922 17.1541057,43.1264922 L20.0037495,43.1264922 L20.0037495,37.3750333 L5.70553681,37.3750333 C2.55593051,37.3750333 0,34.7985309 0,31.6235745 L0,5.74515935 C0,2.57020287 2.55593051,0 5.70553681,0 L44.2944632,0 C47.4440695,0 49.9937508,2.57020287 50,5.74515935 L50,31.617275 C49.9937508,34.7922314 47.4378203,37.3687338 44.288214,37.3687338 Z M2.84964379,28.7446953 L2.84964379,31.6235745 C2.84964379,33.2110527 4.12448444,34.4961541 5.69928759,34.4961541 L44.288214,34.4961541 C45.8630171,34.4961541 47.144107,33.2047532 47.144107,31.6235745 L47.144107,28.7446953 L2.84964379,28.7446953 Z"
                  fill="#0046B0"
                  fill-rule="nonzero"
                />
              </g>
            </g>
          </g>
        </g>
      </svg>
      <svg
        id="db"
        height="60px"
        version="1.1"
        viewBox="0 0 60 60"
        width="60px"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
      >
        <g id="页面优化" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
          <g id="算法" transform="translate(-651.000000, -480.000000)">
            <g id="编组" transform="translate(612.000000, 441.000000)">
              <g id="database-2-fill-" transform="translate(44.000000, 44.000000)">
                <polygon
                  id="矩形"
                  fill="#FFFFFF"
                  points="0 11.1111111 50 11.1111111 50 38.8888889 24.9998896 43.9463033 0 38.8888889"
                />
                <path
                  id="形状"
                  d="M24.9998896,0 C11.197019,0 0,4.96338487 0,11.1065907 C0,17.2497966 11.197019,22.2131814 24.9998896,22.2131814 C38.8027603,22.2131814 49.9998899,17.2497966 49.9998899,11.1065907 C50.0404956,4.96338487 38.8434767,0 24.9998896,0 M0,16.6395443 L0,24.9796583 C0,31.1228641 11.197019,36.086249 24.9998896,36.086249 C38.8027603,36.086249 49.9997792,31.1228641 49.9997792,24.9796583 L49.9997792,16.6395443 C49.9997792,22.7827502 38.8027603,27.7461351 24.9998896,27.7461351 C11.197019,27.7868186 0,22.7827502 0,16.6395443 M0,30.5532954 L0,38.8934093 C0,45.0366151 11.197019,50 24.9998896,50 C38.8027603,50 49.9997792,45.0366151 49.9997792,38.8934093 L49.9997792,30.5532954 C49.9997792,36.6965012 38.8027603,41.6598861 24.9998896,41.6598861 C11.197019,41.6598861 0,36.6965012 0,30.5532954 Z"
                  fill="#0049AA"
                  fill-rule="nonzero"
                />
              </g>
            </g>
          </g>
        </g>
      </svg>
      <svg
        id="linux"
        height="60px"
        version="1.1"
        viewBox="0 0 60 60"
        width="60px"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
      >
        <g id="页面优化" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
          <g id="算法" fill-rule="nonzero" transform="translate(-651.000000, -307.000000)">
            <g id="编组" transform="translate(612.000000, 268.000000)">
              <g id="ubuntu-fill" transform="translate(40.000000, 40.000000)">
                <path
                  id="路径"
                  d="M58,29 C58,45.008 45.008,58 29,58 C12.992,58 0,45.008 0,29 C0,12.992 12.992,0 29,0 C45.008,0 58,12.992 58,29"
                  fill="#FFFFFF"
                />
                <path
                  id="形状"
                  d="M58,29 C58,45.008 45.008,58 29,58 C12.992,58 0,45.008 0,29 C0,12.992 12.992,0 29,0 C45.008,0 58,12.992 58,29 M35.786,16.646 C37.468,17.603 39.585,17.023 40.6,15.37 C41.499,13.717 40.948,11.6 39.266,10.614 C37.613,9.65699999 35.467,10.15 34.51,11.89 C33.553,13.543 34.133,15.689 35.786,16.646 M28.652,39.15 C27.115,39.15 25.665,38.831 24.389,38.222 L21.953,42.572 C23.983,43.5 26.245,44.138 28.652,44.138 C30.073,44.138 31.407,43.935 32.712,43.587 C32.944,42.166 33.756,40.89 35.09,40.136 C36.424,39.353 37.932,39.295 39.266,39.788 C41.847,37.265 43.5,33.814 43.761,29.957 L38.802,29.899 C38.338,35.09 33.988,39.15 28.652,39.15 M28.652,18.85 C33.988,18.85 38.338,22.881 38.802,28.101 L43.761,28.014 C43.5,24.186 41.847,20.735 39.266,18.212 C37.932,18.705 36.395,18.618 35.09,17.864 C33.756,17.11 32.944,15.805 32.712,14.413 C31.407,14.065 30.073,13.862 28.652,13.862 C26.245,13.862 23.983,14.413 21.953,15.428 L24.389,19.778 C25.665,19.169 27.115,18.85 28.652,18.85 M18.473,29 C18.473,25.549 20.184,22.504 22.794,20.677 L20.3,16.385 C17.226,18.444 14.935,21.547 14.007,25.201 C15.109,26.1 15.805,27.463 15.805,29 C15.805,30.537 15.109,31.9 14.007,32.799 C14.935,36.424 17.226,39.556 20.3,41.586 L22.794,37.323 C20.184,35.496 18.473,32.451 18.473,29 M35.786,41.354 C34.133,42.311 33.553,44.428 34.51,46.11 C35.467,47.763 37.613,48.343 39.266,47.386 C40.948,46.4 41.499,44.283 40.6,42.63 C39.585,40.948 37.468,40.397 35.786,41.354 M10.904,25.52 C8.99000002,25.52 7.42400002,27.086 7.42400002,29 C7.42400002,30.914 8.99000002,32.48 10.904,32.48 C12.847,32.48 14.384,30.914 14.384,29 C14.384,27.086 12.847,25.52 10.904,25.52 Z"
                  fill="#154BA4"
                />
              </g>
            </g>
          </g>
        </g>
      </svg>
      <svg
        id="mail"
        height="60px"
        version="1.1"
        viewBox="0 0 60 60"
        width="60px"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
      >
        <g id="页面优化" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
          <g id="算法" transform="translate(-1021.000000, -306.000000)">
            <g id="编组" transform="translate(982.000000, 268.000000)">
              <g id="mail-fill" transform="translate(45.000000, 50.000000)">
                <rect id="矩形" fill="#FFFFFF" height="25" width="42" x="3" y="3" />
                <path
                  id="形状"
                  d="M43.2,3.55271368e-15 L4.80000004,3.55271368e-15 C2.16,3.55271368e-15 0,2.1375 0,4.74999998 L0,33.25 C0,35.8625 2.16,38 4.79999998,38 L43.2,38 C45.84,38 48,35.8625 48,33.25 L48,4.74999998 C48,2.1375 45.84,3.55271368e-15 43.2,3.55271368e-15 Z M43.2,9.49999996 L24,21.3749999 L4.80000004,9.49999996 L4.80000004,4.74999998 L24,16.625 L43.2,4.74999998 L43.2,9.49999996 Z"
                  fill="#0A45AA"
                  fill-rule="nonzero"
                />
              </g>
            </g>
          </g>
        </g>
      </svg>
      <svg
        id="windows"
        height="60px"
        version="1.1"
        viewBox="0 0 60 60"
        width="60px"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
      >
        <g id="页面优化" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
          <g id="算法" fill-rule="nonzero" transform="translate(-840.000000, -305.000000)">
            <g id="编组" transform="translate(804.000000, 268.000000)">
              <g id="windows-fill" transform="translate(41.000000, 42.000000)">
                <polygon id="路径" fill="#FFFFFF" points="0 41.8510148 50 50 50 0 0 8.29857127" />
                <path
                  id="形状"
                  d="M25.6848261,4.02919786 L25.6848261,22.9451718 L50,22.9451718 L50,0 L25.6848261,4.02919786 Z M25.6848261,46.0355929 L50,50 L50,27.0799049 L25.6848261,27.0799049 L25.6848261,46.0355929 Z M0,22.9386604 L20.5440435,22.9386604 L20.5440435,4.8889579 L0,8.30105418 L0,22.9386604 Z M0,41.8477389 L20.5440435,45.1989312 L20.5440435,27.0799049 L0,27.0799049 L0,41.8477389 Z"
                  fill="#0049AA"
                />
              </g>
            </g>
          </g>
        </g>
      </svg>
      <svg
        id="user"
        height="60px"
        version="1.1"
        viewBox="0 0 60 60"
        width="60px"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
      >
        <g id="页面优化" fill="none" fill-rule="evenodd" stroke="none" stroke-width="1">
          <g id="算法" transform="translate(-843.000000, -653.000000)">
            <g id="编组" transform="translate(804.000000, 614.000000)">
              <g id="user-2-fill-" transform="translate(44.000000, 44.000000)">
                <path
                  id="形状"
                  d="M7.02795208,33.6801731 C10.9241301,32.2541571 13.0233885,31.5523776 13.7763833,31.3165797 C14.4780376,31.1313099 15.1397604,30.8393696 15.7444381,30.4407589 C16.0867085,30.2386464 16.4004564,29.9860057 16.6742727,29.6940655 L16.6742727,29.564938 C16.5487735,29.2505408 16.4746149,28.9193009 16.4575014,28.5880609 C16.4346834,28.4140196 16.4346834,28.2343641 16.4575014,28.0547085 L15.5904164,27.2687154 C14.9857387,26.7185203 14.5065602,26.0504262 14.1814033,25.309347 C13.9646321,24.7984515 13.7478608,24.310013 13.5767256,23.7766606 L13.5767256,23.3948925 C13.1660011,23.0299672 12.7952082,22.6257422 12.4757558,22.1822175 C11.9965773,21.4916665 11.6428979,20.7169019 11.4375357,19.9028376 C11.0952653,18.8978894 10.9697661,17.8311845 11.0724472,16.7700939 C11.1409013,16.1356852 11.2892185,15.5181192 11.5059897,14.9173959 C11.6086708,14.6198414 11.7284655,14.3279011 11.8710782,14.0415751 C11.8767827,12.9074994 11.9566458,11.7790379 12.1106674,10.6561907 C12.2475756,9.49404381 12.4643468,8.33751116 12.7609812,7.20904969 L12.3274387,7.20904969 C11.779806,7.21466393 11.2378779,7.27080629 10.7073588,7.37747678 C10.1426127,7.47853303 9.58927553,7.64134588 9.06446092,7.86591532 C8.48830576,8.14101289 7.94637764,8.47786706 7.44438106,8.86524936 C6.87393041,9.32000249 6.38334284,9.86458339 5.99543639,10.4821494 C5.52196235,11.2120001 5.16828294,12.0092216 4.93439817,12.8457428 C4.70621791,13.659807 4.53508272,14.4963282 4.4381061,15.3384636 C4.32972048,16.2591984 4.28408443,17.1799331 4.30690245,18.1062821 C3.95322305,18.5778779 3.74786081,19.1393015 3.70222476,19.7231821 C3.63947519,20.0881074 3.63947519,20.4642612 3.70222476,20.8291866 C3.79920137,21.1604265 3.94181403,21.4748237 4.13576726,21.766764 C4.26697091,21.9464196 4.40958357,22.1204609 4.56930975,22.2776595 C4.72333143,22.3899442 4.88876212,22.4741578 5.06560183,22.5359144 L5.3907587,23.6419189 C5.52196235,23.9619304 5.63034797,24.2594849 5.77866515,24.5794963 C5.88705077,24.8433654 6.05248146,25.0791634 6.25213919,25.2812759 C6.61722761,25.5844446 6.96520251,25.9156845 7.29035938,26.2581529 C7.63262978,26.6230783 7.81517399,27.1059026 7.80946948,27.5999554 L7.80946948,28.5992894 C7.81517399,28.9136866 7.76383343,29.2224696 7.6554478,29.5144099 C7.52994866,29.8344214 7.34740445,30.1319759 7.11351968,30.3902307 C6.7826583,30.7439276 6.40616087,31.0527106 5.98973189,31.3053512 C5.32230462,31.726419 4.59212778,32.0408162 3.8277239,32.2429287 C3.26297775,32.3945131 1.31774102,33.0962926 0,33.5622742 L0,42.5899659 L0.193953223,42.5899659 C0.193953223,42.0566135 0.325156874,41.4615044 0.433542499,40.6755114 L0.433542499,40.4621704 C0.638904735,38.6375436 1.52880776,36.9532728 2.92070736,35.7349836 C4.06160867,34.8479342 5.36223617,34.1742259 6.74843126,33.7531582 L7.02795208,33.7531582 L7.02795208,33.6801731 Z M48.4826013,32.2541571 C47.849401,32.0801158 47.2333143,31.8723891 46.6229321,31.6365912 C46.0867085,31.4120217 45.5904164,31.1032387 45.1511694,30.7214707 C44.8146035,30.4239162 44.5407872,30.0646051 44.3525385,29.6547658 C44.2042213,29.2954547 44.1129492,28.9136866 44.0901312,28.5263043 C44.0444952,27.9031241 44.1357673,27.2799439 44.3525385,26.6960634 C44.4381061,26.4995651 44.5521962,26.3199095 44.7005134,26.1627109 C44.8488306,26.0055123 45.0142613,25.8707706 45.1968055,25.7584859 L45.7615516,25.2700474 C45.9783229,25.0679349 46.1437536,24.8209085 46.2578437,24.5458109 C46.4118654,24.2482564 46.5373645,23.9450876 46.6457501,23.6306904 C46.6457501,23.2713793 46.8625214,22.883997 46.948089,22.502229 C47.1990873,22.4236297 47.4272675,22.288888 47.6212208,22.1204609 C47.815174,21.9071199 47.9749002,21.6600935 48.0946948,21.3962244 C48.2715345,21.0144564 48.3799201,20.6046171 48.4198517,20.1835494 C48.4654877,19.8579237 48.4654877,19.5266838 48.4198517,19.2066723 C48.3571021,18.9877171 48.277239,18.7743762 48.1802624,18.5666494 C48.1289218,18.3869939 48.0148317,18.2297953 47.8551055,18.1175105 C47.8722191,17.1518619 47.815174,16.1805991 47.6839703,15.2205647 C47.5527667,14.3335154 47.3588135,13.4576945 47.1021107,12.5987164 C46.7712493,11.6835959 46.2977752,10.830232 45.6930975,10.0666959 C45.4135767,9.68492783 45.0770108,9.34807367 44.7005134,9.06736186 C43.7763833,8.36558234 42.7267541,7.83784414 41.6086708,7.51221845 C40.9982886,7.34379136 40.3650884,7.25396358 39.7261837,7.25396358 L38.8590987,7.25396358 C39.167142,8.30943998 39.3896178,9.38175908 39.5322305,10.4709209 C39.6862521,11.54324 39.7718197,12.6267876 39.7946378,13.7103352 C40.1711352,14.372815 40.4335425,15.0914372 40.5761552,15.8381307 C40.6959498,16.6297379 40.6959498,17.4325737 40.5761552,18.224181 C40.4677695,19.128073 40.2281803,20.0095081 39.8630918,20.8460293 C39.4694809,21.6657078 38.9047347,22.3955585 38.1973759,22.9738248 L38.0433542,23.4398064 L38.0433542,23.5689338 C37.8436965,24.0798293 37.6269253,24.5907248 37.3930405,25.1016203 C36.9652025,26.1234113 36.2692527,27.0216891 35.3793497,27.7010116 L35.3793497,28.6385891 C35.3793497,29.0540425 35.2880776,29.4582675 35.0998289,29.8288071 C35.4535083,30.1095189 35.8471192,30.3397026 36.2692527,30.5081297 C36.970907,30.8056842 37.6953793,31.0470964 38.4312607,31.2323662 C39.1728465,31.4569356 39.8973189,31.7320332 40.5932687,32.0632731 L42.7552767,33.0401502 L43.9874501,33.5735026 C45.0656018,34.0001846 46.069595,34.5784509 46.970907,35.2970731 C48.0433542,36.2346506 48.8362807,37.4417113 49.2641187,38.789128 C49.5094124,39.6144207 49.6691386,40.4621704 49.7375927,41.3211485 C49.8288648,41.6692312 49.8745009,42.0285423 49.8687963,42.3878534 L50,42.3878534 L50,33.107521 C49.2926412,32.7875096 48.5738734,32.5011835 47.837992,32.2541571 L48.4826013,32.2541571 Z"
                  fill="#0043B7"
                  fill-rule="nonzero"
                />
                <rect id="矩形" fill="#FFFFFF" height="14" width="11" x="24" y="32" />
                <path
                  id="形状"
                  d="M47.5527667,40.1084735 C47.2846549,39.2831808 46.7940673,38.5477159 46.12664,37.980678 C45.4256733,37.4418762 44.6543804,36.9981309 43.8334284,36.6613325 C42.9891614,36.3188641 39.2698232,34.533537 38.3399886,34.3819526 C37.4671991,34.1405405 36.6058186,33.8486002 35.7672561,33.5061318 C35.0313748,33.2029631 34.346834,32.7762811 33.7535653,32.2485429 C33.2915003,31.8387037 32.9207074,31.3390367 32.6697091,30.777613 C32.4529378,30.2947887 32.3274387,29.7726648 32.3046207,29.2449266 C32.2532801,28.37472 32.3958928,27.5045134 32.7153451,26.6904491 C32.8294353,26.4209658 32.9891614,26.1795536 33.1888192,25.9662127 C33.3827724,25.7584859 33.6052481,25.5732161 33.8391329,25.4104033 C34.1243582,25.2082908 34.3924701,24.9781071 34.6377638,24.7310807 C34.9229892,24.4503689 35.1568739,24.1303574 35.3280091,23.7710463 C35.5219624,23.3443644 35.6930975,22.9176824 35.8699373,22.4910005 C36.0410724,22.0643185 36.1494581,21.4692095 36.2806617,20.958314 C36.5715916,20.8291866 36.8282944,20.6439168 37.0393611,20.4025046 C37.3245864,20.1161786 37.5584712,19.7849387 37.7296064,19.4256276 C37.9634912,18.8978894 38.1061038,18.3364657 38.1631489,17.7638137 C38.2373075,17.320289 38.2373075,16.8655359 38.1631489,16.4220112 C38.1003993,16.1132282 37.9977182,15.8100595 37.86081,15.5293477 C37.7809469,15.276707 37.6212208,15.0633661 37.4044495,14.9117817 C37.432972,13.5812077 37.353109,12.2506338 37.1648602,10.9312883 C37.0108386,9.70738478 36.7313177,8.50593824 36.3434113,7.33256289 C35.9669139,6.06935976 35.3736452,4.88475593 34.5921278,3.81805106 C34.1985168,3.30154133 33.7307473,2.84117397 33.2059327,2.45379168 C32.6183685,1.97096737 31.9851683,1.55551389 31.3006275,1.21865972 C30.5476326,0.864962847 29.7661152,0.573022567 28.9617798,0.342838885 C28.1175128,0.112655203 27.2447233,-0.00524375639 26.3662293,0.000178953321 C25.641757,0.0115989521 24.9172847,0.084584022 24.2042213,0.213711453 C23.456931,0.354067357 22.7324586,0.584251039 22.0422133,0.893034027 C20.2224758,1.6341132 18.6651455,2.88608786 17.5641757,4.4917594 C16.9195665,5.49109343 16.4346834,6.58025524 16.1380491,7.73117366 C15.5504849,10.0779244 15.2994866,12.4920459 15.4021677,14.9061675 C15.1625784,15.2205647 14.9686252,15.5742616 14.8374216,15.9504154 C14.6833999,16.298498 14.5693098,16.6634234 14.4894467,17.033963 C14.4438106,17.5448584 14.5065602,18.0557539 14.6605819,18.5441925 C14.7689675,19.0045598 14.9572162,19.4480845 15.225328,19.841081 C15.4078722,20.0937217 15.6075299,20.3239053 15.8300057,20.5428605 C16.0296634,20.7000592 16.2578437,20.8123439 16.5031375,20.8853289 L16.9766115,22.4404724 C17.1306332,22.8671543 17.3017684,23.2938363 17.4729036,23.6980613 C17.6269253,24.0686008 17.8551055,24.3998408 18.1460354,24.6749383 C18.6537365,25.1016203 19.1272105,25.5507592 19.572162,26.0167408 C20.034227,26.522022 20.2909298,27.1732734 20.3080434,27.8469818 L20.3080434,29.2505408 C20.3023388,29.6772228 20.2281803,30.1039047 20.0912721,30.5081297 C19.9087279,30.9460401 19.6520251,31.3446509 19.3325727,31.6983478 C18.8876212,32.1755578 18.3799201,32.5910113 17.815174,32.9334797 C16.9555794,33.5275337 15.9956319,33.9666461 14.9800342,34.2303683 C13.8790645,34.5279228 9.16143754,36.2514933 8.14033086,36.6388756 C7.17056475,36.9645013 6.26925271,37.4641683 5.4763263,38.1098055 C4.6948089,38.8508846 4.19851683,39.8389902 4.06731318,40.9000808 C3.78208785,42.1576697 3.71363377,43.4545583 3.87335995,44.7346041 C3.85624643,45.6048107 4.12435824,46.4525604 4.63205933,47.1599541 C5.40216771,47.5473364 6.22932116,47.822434 7.07929264,47.9908611 C13.5653166,49.3102066 20.1711352,49.9783007 26.7883628,49.9951434 C33.2800913,50.0737427 39.7490017,49.1979218 45.9783229,47.395752 C46.4746149,47.2610104 46.93668,47.013984 47.3188819,46.6715156 C47.5698802,46.1718486 47.7353109,45.6328819 47.792356,45.0770725 C47.9007416,44.3359933 47.9463776,43.5892999 47.9463776,42.8426065 C47.9463776,42.0734562 47.9463776,42.6910221 47.837992,41.927486 C47.792356,41.2762346 47.6725613,40.6362117 47.4900171,40.0130315 L47.5527667,40.1084735 L47.5527667,40.1084735 Z M29.1842556,44.8805742 C28.9332573,45.2005857 28.659441,45.4981402 28.3628066,45.7732378 C28.0718768,46.0595638 27.7524244,46.3122045 27.410154,46.5423881 C27.1990873,46.7164295 26.948089,46.8287142 26.6742727,46.8623996 C26.4232744,46.8174857 26.1836851,46.7164295 25.9840274,46.5648451 C25.641757,46.3346614 25.3166001,46.0764065 25.0085568,45.7956947 C24.7005134,45.5149829 24.4152881,45.2118142 24.1642898,44.8805742 C23.9760411,44.6953045 23.8562464,44.4538923 23.8163149,44.2012517 C23.8277239,43.7857982 23.8790645,43.3703447 23.9703366,42.9661197 C23.9703366,42.4327673 24.2099258,41.8376583 24.335425,41.2425492 C24.4609241,40.6474402 24.6149458,40.0523312 24.7917855,39.4740648 C24.9686252,38.8957985 25.0313748,38.598244 25.1568739,38.1771763 C24.8374216,37.8964645 24.5521962,37.5708388 24.312607,37.2171419 C24.1471763,36.9869582 24.0559042,36.7118606 24.0501997,36.4311488 C24.0559042,36.1448228 24.1414718,35.8697252 24.2897889,35.6226988 C24.4324016,35.3756724 24.6092413,35.1454888 24.808899,34.9433763 C25.0541928,34.6795072 25.3223046,34.438095 25.6075299,34.2191398 L27.9463776,34.2191398 L28.7450086,34.9433763 C28.9845978,35.1454888 29.1842556,35.3925152 29.3268682,35.6676127 C29.4980034,35.9034106 29.5892755,36.1841225 29.5892755,36.4760627 C29.5778665,36.8241454 29.4637764,37.1666138 29.2641187,37.4529398 C29.0017114,37.7841798 28.705077,38.0817343 28.3742156,38.3456034 C28.3742156,38.7273714 28.5909869,39.1765103 28.6993725,39.7716194 L29.087279,41.5401037 C29.087279,42.1127558 29.3268682,42.6685652 29.4124358,43.1794607 C29.4808899,43.5219291 29.526526,43.8700117 29.5436395,44.2237086 C29.5037079,44.4875777 29.3839133,44.7346041 29.1956646,44.9254881 L29.1842556,44.8805742 L29.1842556,44.8805742 Z"
                  fill="#0043B7"
                  fill-rule="nonzero"
                />
              </g>
            </g>
          </g>
        </g>
      </svg>
    </svg>
    <svg ref="d3_force_ref" class="d3-svg" />
  </div>
</template>
<style scoped lang="scss">
  @import url(./d3-plugins/d3-context-menu.scss);
  .d3_force_box {
    position: relative;
    width: 100%;
    height: 100%;
    transition: all 0.3s;
    .zoomBtn {
      position: absolute;
      right: 52px;
      bottom: 20px;
      width: 28px;
      height: 28px;
      background-color: #25223f;
      border-radius: 6px;
      text-align: center;
      line-height: 26px;
      color: #fff;
      user-select: none;
      cursor: pointer;
      font-size: 18px;
      &.scale {
        right: 20px;
      }
    }
    .d3-svg {
      :deep() {
        width: 100%;
        height: 100%;
        user-select: none;
        $opacity: 0.15; /* 显示的不透明度 */
        $activeColor: #1e90ff; /* 激活的颜色 */
        svg {
          box-sizing: border-box;
          width: 100%;
          height: 100%;
          margin: 20px 0px;
          background-color: rgb(0, 124, 249);
        }
        .errNode {
          svg {
            fill: red;
          }
        }
      }
    }
  }
</style>
