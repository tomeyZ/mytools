<template>
  <div class="ip-query-container">
    <h2>IP地址查询工具</h2>

    <div class="search-area">
      <input
          v-model="ipAddress"
          class="ip-input"
          placeholder="请输入IP地址或域名"
          @keyup.enter="queryIp"
      />
      <button @click="queryIp" class="query-btn">查询</button>
      <button @click="clearQuery" class="clear-btn">清空</button>
    </div>

    <div v-if="loading" class="loading-indicator">
      <div class="spinner"></div>
      正在查询中...
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div v-if="ipInfo" class="result-container">
      <div class="result-card">
        <div class="result-header">
          <h2>IP地址信息</h2>
          <span class="query-ip">{{ ipInfo.query }}</span>
        </div>

        <div class="result-grid">
          <div class="result-item">
            <span class="label">国家/地区</span>
            <span class="value">{{ ipInfo.country }} ({{ ipInfo.countryCode }})</span>
          </div>
          <div class="result-item">
            <span class="label">省份/地区</span>
            <span class="value">{{ ipInfo.regionName }} ({{ ipInfo.region }})</span>
          </div>
          <div class="result-item">
            <span class="label">城市</span>
            <span class="value">{{ ipInfo.city }}</span>
          </div>
          <div class="result-item">
            <span class="label">经纬度</span>
            <span class="value">{{ ipInfo.lat }}, {{ ipInfo.lon }}</span>
          </div>
          <div class="result-item">
            <span class="label">时区</span>
            <span class="value">{{ ipInfo.timezone }}</span>
          </div>
          <div class="result-item">
            <span class="label">ISP</span>
            <span class="value">{{ ipInfo.isp }}</span>
          </div>
          <div class="result-item">
            <span class="label">组织</span>
            <span class="value">{{ ipInfo.org }}</span>
          </div>
          <div class="result-item">
            <span class="label">AS编号</span>
            <span class="value">{{ ipInfo.as }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
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
        const data = await window.go.handler.NetworkHandler.GetIpInfo(
            query
        );
        if (data.status === 'success') {
          ipInfo.value = data;
        } else {
          error.value = data.message || '查询失败，请检查IP地址是否正确';
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
  }
};
</script>

<style scoped>
.ip-query-container {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
  font-family: "思源宋体", "Noto Serif SC", serif;
}
.ip-query-container {
  font-family: inherit;
}

h1 {
  color: #333;
  text-align: center;
  margin-bottom: 20px;
}

.search-area {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.search-area input,
.search-area button {
  font-family: "思源宋体", "Noto Serif SC", serif;
}

.ip-input {
  flex: 1;
  padding: 10px 15px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.query-btn {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.clear-btn {
  padding: 10px 20px;
  background-color: #f0f0f0;
  color: #333;
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
}

.clear-btn:hover {
  background-color: #e0e0e0;
}

.query-btn:hover, .clear-btn:hover {
  opacity: 0.9;
}

.loading-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 20px;
  color: #666;
}

.spinner {
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top: 3px solid #3498db;
  width: 20px;
  height: 20px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  padding: 15px;
  background-color: #ffebee;
  color: #f44336;
  border-radius: 4px;
  margin-bottom: 20px;
}

.result-container {
  margin-top: 20px;
}

.result-card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin-bottom: 20px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.result-header h2 {
  margin: 0;
  color: #333;
}

.query-ip {
  background-color: #e3f2fd;
  color: #1976d2;
  padding: 5px 10px;
  border-radius: 4px;
}

.result-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 15px;
}

.result-item {
  display: flex;
  flex-direction: column;
}

.label {
  font-size: 14px;
  color: #666;
  margin-top: 5px;
  margin-bottom: 5px;
}

.value {
  font-weight: 500;
  word-break: break-all;
}

.map-container {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.map-container h3 {
  margin-top: 0;
  color: #333;
}

.map-placeholder {
  height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
  border-radius: 4px;
  color: #666;
}

.map-placeholder a {
  margin-top: 10px;
  color: #1976d2;
  text-decoration: none;
}

.map-placeholder a:hover {
  text-decoration: underline;
}
</style>