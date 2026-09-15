<template>
  <div class="page">
    <div class="grid">
      <section class="card wide">
        <div class="card-head">
          <span class="ci">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="11" cy="11" r="7"/><line x1="21" y1="21" x2="16.2" y2="16.2"/>
            </svg>
          </span>
          <div>
            <h3>查询条件</h3>
            <p>回车即可查询</p>
          </div>
        </div>
        <div class="card-body">
          <div class="row">
            <input
                v-model="ipAddress"
                class="inp mono"
                placeholder="请输入 IP 地址或域名"
                @keyup.enter="queryIp"
            >
            <button class="btn primary" @click="queryIp" :disabled="loading">查询</button>
            <button class="btn ghost" @click="clearQuery">清空</button>
          </div>
        </div>
      </section>

      <!-- 加载中用行内提示条，避免结果区高度跳动 -->
      <section class="card wide" v-if="loading">
        <div class="state-row">
          <span class="spinner"></span>
          正在查询…
        </div>
      </section>

      <section class="card wide" v-else-if="error">
        <div class="out out-error">
          <span class="val">{{ error }}</span>
        </div>
      </section>

      <section class="card wide" v-else-if="ipInfo">
        <div class="card-head">
          <span class="ci ci-3">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="9"/><path d="M3 12h18"/>
              <path d="M12 3a14 14 0 0 1 4 9 14 14 0 0 1-4 9 14 14 0 0 1-4-9 14 14 0 0 1 4-9z"/>
            </svg>
          </span>
          <div>
            <h3>归属信息</h3>
            <p>数据来源 api.ip.sb</p>
          </div>
          <div class="card-head-actions">
            <span class="badge mono">{{ ipInfo.query }}</span>
          </div>
        </div>
        <div class="card-body">
          <div class="kv"><span class="k">国家/地区</span><span class="v">{{ info('country', 'countryCode') }}</span></div>
          <div class="kv"><span class="k">省份/地区</span><span class="v">{{ info('regionName', 'region') }}</span></div>
          <div class="kv"><span class="k">城市</span><span class="v">{{ ipInfo.city }}</span></div>
          <div class="kv"><span class="k">经纬度</span><span class="v mono">{{ ipInfo.lat }}, {{ ipInfo.lon }}</span></div>
          <div class="kv"><span class="k">时区</span><span class="v mono">{{ ipInfo.timezone }}</span></div>
          <div class="kv"><span class="k">ISP</span><span class="v">{{ ipInfo.isp }}</span></div>
          <div class="kv"><span class="k">组织</span><span class="v">{{ ipInfo.org }}</span></div>
          <div class="kv"><span class="k">AS 编号</span><span class="v mono">{{ ipInfo.as }}</span></div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
  name: 'IpQuery',
  setup() {
    const ipAddress = ref('');
    const ipInfo = ref(null);
    const loading = ref(false);
    const error = ref('');

    const queryIp = async () => {
      try {
        loading.value = true;
        error.value = '';
        ipInfo.value = null;

        const query = ipAddress.value.trim() || '';
        const data = await window.go.handler.NetworkHandler.GetIpInfo(query);
        if (data.status === 'success') {
          ipInfo.value = data;
        } else {
          error.value = data.message || '查询失败，请检查 IP 地址是否正确';
        }
      } catch (err) {
        error.value = `查询出错: ${err.message}`;
      } finally {
        loading.value = false;
      }
    };

    const clearQuery = () => {
      ipAddress.value = '';
      ipInfo.value = null;
      error.value = '';
    };

    return {
      ipAddress,
      ipInfo,
      loading,
      error,
      queryIp,
      clearQuery
    };
  },
  methods: {
    // 「主字段 (编码)」的统一拼法，字段为空时不留多余的括号
    info(main, code) {
      const m = this.ipInfo[main];
      const c = this.ipInfo[code];
      if (m && c) return `${m} (${c})`;
      return m || c || '-';
    }
  }
};
</script>

<style scoped>
.badge {
  padding: 5px 9px;
  border-radius: 999px;
  font-size: 11.5px;
  color: var(--accent-strong);
  background: var(--accent-soft);
  border: 1px solid var(--accent-border);
}
</style>
