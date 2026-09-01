<template>
  <div class="app-container">
    <Sidebar @tool-change="handleToolChange"/>
    <div class="main-content">
      <!-- 使用动态组件根据activeTool显示对应组件，KeepAlive缓存组件状态，切换菜单后数据不丢失 -->
      <KeepAlive>
        <component :is="currentComponent" />
      </KeepAlive>
    </div>
  </div>
</template>

<script>
import Sidebar from '../components/Sidebar.vue'
import TimeTool from '../components/TimeTool.vue'
import Md5Tool from '../components/Md5Tool.vue'
import JsonFormat from "../components/JsonFormat.vue";
import IpQuery from "../components/IpQuery.vue";
import AesTool from "../components/AesTool.vue";
import RsaTool from "../components/RsaTool.vue";
import TextTool from "../components/TextTool.vue";

export default {
  components: {Sidebar, TimeTool, Md5Tool, JsonFormat, IpQuery, AesTool, RsaTool, TextTool},
  data() {
    return {
      // 映射工具ID到组件
      toolComponents: {
        'time': TimeTool,
        'md5': Md5Tool,
        'json': JsonFormat,
        'ip': IpQuery,
        'aes': AesTool,
        'rsa': RsaTool,
        'text': TextTool,
      },
      activeTool: 'time'
    }
  },
  computed: {
    currentComponent() {
      return this.toolComponents[this.activeTool] || TimeTool;
    }
  },
  methods: {
    handleToolChange(toolId) {
      this.activeTool = toolId
    }
  }
}
</script>

<style>
.app-container {
  display: flex;
  height: 100vh;
}

.main-content {
  flex: 1;
  font-family: "思源宋体", "Noto Serif SC", serif;
  padding: 20px;
  overflow: hidden;
}
</style>