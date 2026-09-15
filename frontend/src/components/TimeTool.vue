<template>
  <div class="time-page">
    <!-- 页头 -->
    <div class="page-head">
      <p>时间戳与日期互转 · 支持 {{ timezoneOptions.length }} 个常用时区</p>
      <!-- 右侧一组：时区标识 + 时区下拉。两者说的是同一件事（当前时区），所以贴在一起放。
           顺序是「值在前、控件在后」：先看到现在是哪个时区，再决定要不要换。
           注意 .tz-badge 必须放在 .tz-select-wrap **外面**——wrap 是下拉面板的定位基准，
           徽章一旦进去，wrap 变宽，面板用 left/right 定宽就会比触发器宽（上次修过的 12px 偏移会重演） -->
      <div class="hz-right">
        <button
            type="button"
            class="tz-badge mono"
            title="点击复制时区标识"
            @click="copyResult(timezone, '时区标识')"
        >{{ timezone }}</button>

        <!-- 时区选择：原生 <select> 的展开列表由 WebView2 自己绘制，圆角 / 行高 / hover / 选中态
             一个都改不了，所以改成自绘面板。键盘：↑↓ 移动、Enter 选中、Esc 关闭；点击面板外关闭 -->
        <div ref="tzWrap" class="tz-select-wrap" :class="{ 'is-open': tzOpen }">
          <button
              type="button"
              class="tz-trigger"
              role="combobox"
              aria-haspopup="listbox"
              :aria-expanded="tzOpen ? 'true' : 'false'"
              :aria-activedescendant="tzOpen ? 'tz-opt-' + tzActive : null"
              @click="toggleTz"
              @keydown="onTzKeydown"
          >
            <svg class="tz-globe" viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="9"/><path d="M3 12h18"/><path d="M12 3a14 14 0 0 1 4 9 14 14 0 0 1-4 9 14 14 0 0 1-4-9 14 14 0 0 1 4-9z"/>
            </svg>
            <span class="tz-current">{{ currentTz.name }}</span>
            <span class="tz-offset mono">{{ currentTz.offset }}</span>
            <svg class="tz-caret" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="6 9 12 15 18 9"/>
            </svg>
          </button>

          <transition name="tz-pop">
            <ul v-if="tzOpen" ref="tzMenu" class="tz-menu" role="listbox">
              <li
                  v-for="(tz, i) in tzItems"
                  :id="'tz-opt-' + i"
                  :key="tz.value"
                  class="tz-item"
                  :class="{ act: i === tzActive, on: tz.value === timezone }"
                  role="option"
                  :aria-selected="tz.value === timezone ? 'true' : 'false'"
                  @mousedown.prevent
                  @mouseenter="tzActive = i"
                  @click="pickTz(tz)"
              >
                <span class="tz-name">{{ tz.name }}</span>
                <span class="tz-off mono">{{ tz.offset }}</span>
                <svg v-if="tz.value === timezone" class="tz-check" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
              </li>
            </ul>
          </transition>
        </div>
      </div>
    </div>

    <!-- 当前时间：一行式。
         标签用 12px 灰字退到后面，**值保持 18px 等宽数字**——这是全页最该一眼看到的东西，
         跟着标签一起缩成 13px 就变成"一行说明文字"了；暂停按钮留在行尾，是该行唯一的交互。
         时区标识（Asia/Shanghai）已上移到页头，和下拉放成一组 -->
    <div class="hero">
      <div class="h-cell">
        <span class="h-key">当前时间</span>
        <span class="h-val mono">{{ currentTime }}</span>
        <span class="h-dow">{{ currentWeekday }}</span>
      </div>
      <div class="h-cell">
        <span class="h-key">当前时间戳</span>
        <span class="h-val mono">{{ currentTimestamp }}</span>
      </div>
      <button
          type="button"
          class="icon-btn"
          :title="isPaused ? '继续刷新' : '暂停刷新'"
          @click="toggleTimeUpdate"
      >
        <svg v-if="!isPaused" viewBox="0 0 24 24" width="14" height="14" fill="currentColor">
          <rect x="6.5" y="5" width="3.6" height="14" rx="1.2"/><rect x="13.9" y="5" width="3.6" height="14" rx="1.2"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" width="14" height="14" fill="currentColor">
          <path d="M8 5l11 7-11 7z"/>
        </svg>
      </button>
    </div>

    <div class="grid">
      <!-- 时间戳 → 日期 -->
      <section class="card">
        <div class="card-head">
          <span class="ci"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="8.7"/><path d="M12 6.6V12l4 2.5"/>
          </svg></span>
          <div>
            <h3>时间戳 → 日期</h3>
            <p>秒级 Unix 时间戳</p>
          </div>
        </div>
        <div class="card-body">
          <span class="field-label">时间戳</span>
          <div class="row">
            <input
                v-model="timestampInput"
                class="inp mono"
                placeholder="1789228377"
                @input="clearResultIfEmpty('timestamp')"
                @keyup.enter="convertToDate"
            >
            <button class="btn primary" @click="convertToDate">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="m16 3 4 4-4 4"/><path d="M20 7H4"/><path d="m8 21-4-4 4-4"/><path d="M4 17h16"/>
              </svg>
              转换
            </button>
          </div>
          <div v-if="dateResult" class="out" :class="isErr(dateResult) ? 'out-error' : 'out-ok'">
            <span class="val mono">{{ dateResult }}</span>
            <button v-if="!isErr(dateResult)" class="copy" title="复制" @click="copyResult(dateResult, '日期')">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
              </svg>
            </button>
          </div>
          <div v-else class="out out-empty"><span class="val">输入时间戳后点击转换</span></div>
        </div>
      </section>

      <!-- 日期 → 时间戳 -->
      <section class="card">
        <div class="card-head">
          <span class="ci ci-2"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="4.5" width="18" height="16.5" rx="2.4"/><line x1="16" y1="2.5" x2="16" y2="6.5"/><line x1="8" y1="2.5" x2="8" y2="6.5"/><line x1="3" y1="10" x2="21" y2="10"/>
          </svg></span>
          <div>
            <h3>日期 → 时间戳</h3>
            <p>支持 YYYY-MM-DD / 斜杠 / 不补零</p>
          </div>
        </div>
        <div class="card-body">
          <span class="field-label">日期时间</span>
          <div class="row">
            <input
                v-model="dateInput"
                class="inp mono"
                placeholder="YYYY-MM-DD HH:mm:ss"
                @input="clearResultIfEmpty('date')"
                @keyup.enter="convertToTimestamp"
            >
            <button class="btn soft" @click="convertToTimestamp">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="m16 3 4 4-4 4"/><path d="M20 7H4"/><path d="m8 21-4-4 4-4"/><path d="M4 17h16"/>
              </svg>
              转换
            </button>
          </div>
          <div v-if="timestampResult" class="out" :class="isErr(timestampResult) ? 'out-error' : 'out-ok'">
            <span class="val mono">{{ timestampResult }}</span>
            <button v-if="!isErr(timestampResult)" class="copy" title="复制" @click="copyResult(timestampResult, '时间戳')">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
              </svg>
            </button>
          </div>
          <div v-else class="out out-empty"><span class="val">输入日期后点击转换</span></div>
        </div>
      </section>

      <!-- 日期计算 -->
      <section class="card">
        <div class="card-head">
          <span class="ci ci-3"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="4.5" width="18" height="16.5" rx="2.4"/><line x1="16" y1="2.5" x2="16" y2="6.5"/><line x1="8" y1="2.5" x2="8" y2="6.5"/><line x1="3" y1="10" x2="21" y2="10"/><line x1="12" y1="13.5" x2="12" y2="18.5"/><line x1="9.5" y1="16" x2="14.5" y2="16"/>
          </svg></span>
          <div>
            <h3>日期计算</h3>
            <p>负数向前，正数向后</p>
          </div>
        </div>
        <div class="card-body">
          <div class="row calc-row">
            <div class="fld">
              <span class="field-label">基准日期</span>
              <input
                  v-model="baseDateInput"
                  class="inp mono fld-date"
                  placeholder="YYYY-MM-DD"
                  @input="clearResultIfEmpty('dateCalc')"
              >
            </div>
            <div class="fld">
              <span class="field-label">相差天数</span>
              <input
                  v-model="daysOffset"
                  class="inp mono fld-days"
                  placeholder="0"
                  @keyup.enter="calculateDate"
              >
            </div>
            <button class="btn primary" @click="calculateDate">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3.6" y="5" width="16.8" height="15.4" rx="2.2"/><path d="M3.6 10.2h16.8"/><path d="M8.2 3.2v3.6M15.8 3.2v3.6"/><path d="M12 12.8v4.6M9.7 15.1h4.6"/>
              </svg>
              计算
            </button>
          </div>
          <div v-if="dateCalcResult" class="out" :class="isErr(dateCalcResult) ? 'out-error' : 'out-ok'">
            <span class="val mono">{{ dateCalcResult }}</span>
            <button v-if="!isErr(dateCalcResult)" class="copy" title="复制" @click="copyResult(dateCalcResult, '计算结果')">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
              </svg>
            </button>
          </div>
          <div v-else class="out out-empty"><span class="val">输入日期和天数后计算</span></div>
        </div>
      </section>
    </div>

    <!-- 复制提示：图标由 .message-toast::before 统一绘制（见 style.css），这里只放文案 -->
    <transition name="toast-fade">
      <div v-if="toast.show" :class="['message-toast', toast.type]">{{ toast.text }}</div>
    </transition>
  </div>
</template>

<script>
// Date.getDay() 返回 0=周日 … 6=周六
const WEEKDAY_CN = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];

export default {
  data() {
    return {
      timezone: "Asia/Shanghai",
      timezoneOptions: [
        { text: "中国标准时间 (UTC+8)", value: "Asia/Shanghai" },
        { text: "菲律宾标准时间 (UTC+8)", value: "Asia/Manila" },
        { text: "印度标准时间 (UTC+5:30)", value: "Asia/Kolkata" },
        { text: "印尼西部时间 (UTC+7)", value: "Asia/Jakarta" },
        { text: "越南标准时间 (UTC+7)", value: "Asia/Ho_Chi_Minh" },
        { text: "墨西哥中部时间 (UTC-6)", value: "America/Mexico_City" },
        { text: "美国东部时间 (UTC-5)", value: "America/New_York" },
        { text: "美国太平洋时间 (UTC-8)", value: "America/Los_Angeles" },
      ],
      currentTime: "",
      currentWeekday: "",
      currentTimestamp: "",
      timestampInput: "",
      dateInput: "",
      dateResult: "",
      timestampResult: "",
      baseDateInput: "",
      daysOffset: "",
      dateCalcResult: "",
      timer: null,
      isPaused: false,
      toast: { show: false, text: '', type: 'success' },
      toastTimer: null,
      tzOpen: false,
      tzActive: 0
    }
  },
  computed: {
    // timezoneOptions 的 text 是「名称 (UTC+8)」这种合并串，下拉里要分两列排版，
    // 所以在不改动原数据结构的前提下派生一份，顺便按 UTC 偏移从东到西排序
    tzItems() {
      return this.timezoneOptions
          .map(opt => {
            const m = opt.text.match(/^(.*?)\s*\(([^)]+)\)\s*$/);
            return {
              value: opt.value,
              name: m ? m[1] : opt.text,
              offset: m ? m[2] : '',
              sort: this.offsetRank(m ? m[2] : '')
            };
          })
          // 同偏移的（中国 / 菲律宾）保持原有先后，sort 在现代引擎里是稳定排序
          .sort((a, b) => b.sort - a.sort);
    },
    currentTz() {
      const hit = this.tzItems.find(t => t.value === this.timezone);
      // 兜底：万一以后 timezone 被设成列表外的值，至少要能显示出来
      return hit || { name: this.timezone, offset: '' };
    }
  },
  mounted() {
    this.updateCurrentTime();
    this.timer = setInterval(this.updateCurrentTime, 1000);
    document.addEventListener('mousedown', this.onDocMousedown);
  },
  beforeUnmount() {
    if (this.timer) {
      clearInterval(this.timer);
    }
    document.removeEventListener('mousedown', this.onDocMousedown);
  },
  methods: {
    // 后端返回的错误文案统一带这些前缀，用它区分「成功结果」和「错误提示」
    isErr(text) {
      if (!text) return false;
      return text.startsWith('无效') || text.startsWith('格式错误') || text.startsWith('日期格式错误');
    },
    // ---------------- 时区下拉 ----------------
    // 'UTC+5:30' → 5.5，用于排序；解析不出来时排到最后
    offsetRank(offset) {
      const m = String(offset).match(/UTC([+-])(\d{1,2})(?::(\d{2}))?/);
      if (!m) return -Infinity;
      const sign = m[1] === '-' ? -1 : 1;
      return sign * (parseInt(m[2], 10) + (m[3] ? parseInt(m[3], 10) / 60 : 0));
    },
    toggleTz() {
      this.tzOpen ? this.closeTz() : this.openTz();
    },
    openTz() {
      if (this.tzOpen) return;
      this.tzOpen = true;
      // 打开时把高亮定到当前选中项，避免每次都要从头按 ↓
      const i = this.tzItems.findIndex(t => t.value === this.timezone);
      this.tzActive = i < 0 ? 0 : i;
      this.$nextTick(this.scrollTzToActive);
    },
    closeTz() {
      this.tzOpen = false;
    },
    pickTz(tz) {
      this.timezone = tz.value;   // watch 会顺带刷新当前时间
      this.tzOpen = false;
    },
    scrollTzToActive() {
      const menu = this.$refs.tzMenu;
      const el = menu && menu.children[this.tzActive];
      if (el && el.scrollIntoView) {
        el.scrollIntoView({ block: 'nearest' });
      }
    },
    onTzKeydown(e) {
      const k = e.key;
      if (!this.tzOpen) {
        // 关闭状态下这些键直接展开（原生 select 的习惯）
        if (k === 'ArrowDown' || k === 'ArrowUp' || k === 'Enter' || k === ' ') {
          e.preventDefault();
          this.openTz();
        }
        return;
      }
      if (k === 'Escape') {
        e.preventDefault();
        this.closeTz();
      } else if (k === 'Tab') {
        this.closeTz();          // 不拦 Tab，让它正常移走焦点
      } else if (k === 'ArrowDown') {
        e.preventDefault();
        this.tzActive = Math.min(this.tzActive + 1, this.tzItems.length - 1);
        this.scrollTzToActive();
      } else if (k === 'ArrowUp') {
        e.preventDefault();
        this.tzActive = Math.max(this.tzActive - 1, 0);
        this.scrollTzToActive();
      } else if (k === 'Enter' || k === ' ') {
        e.preventDefault();
        const tz = this.tzItems[this.tzActive];
        if (tz) this.pickTz(tz);
      }
    },
    onDocMousedown(e) {
      if (!this.tzOpen) return;
      const wrap = this.$refs.tzWrap;
      if (wrap && !wrap.contains(e.target)) {
        this.closeTz();
      }
    },
    clearTimer() {
      if (this.timer) {
        clearInterval(this.timer);
        this.timer = null;
      }
    },
    toggleTimeUpdate() {
      this.isPaused = !this.isPaused;
      if (this.isPaused) {
        this.clearTimer();
      } else {
        this.timer = setInterval(this.updateCurrentTime, 1000);
      }
    },
    updateCurrentTime() {
      const now = new Date();
      this.currentTimestamp = Math.floor(now.getTime() / 1000);

      const options = {
        timeZone: this.timezone,
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      };
      const formatter = new Intl.DateTimeFormat('zh-CN', options);
      const parts = formatter.formatToParts(now);
      const year = parts.find(p => p.type === 'year').value;
      const month = parts.find(p => p.type === 'month').value;
      const day = parts.find(p => p.type === 'day').value;
      const hour = parts.find(p => p.type === 'hour').value.padStart(2, '0');
      const minute = parts.find(p => p.type === 'minute').value.padStart(2, '0');
      const second = parts.find(p => p.type === 'second').value.padStart(2, '0');
      this.currentTime = `${year}-${month}-${day} ${hour}:${minute}:${second}`;
      // 星期必须按「目标时区的日历日」算，不能用 new Date(...).getDay()——
      // 那会先按浏览器本地时区解释日期串，极东/极西时区跨零点时会差一天。
      // 这里用日历数字直接构造，Date 内部只记录一个本地午夜，getDay() 读回的还是同一个日历日
      const weekday = new Date(Number(year), Number(month) - 1, Number(day)).getDay();
      this.currentWeekday = WEEKDAY_CN[weekday];
    },
    async convertToDate() {
      if (!this.timestampInput) {
        this.dateResult = "";
        return;
      }
      const ts = parseInt(this.timestampInput.trim());
      if (isNaN(ts)) {
        this.dateResult = "无效的时间戳";
        return;
      }
      this.dateResult = await window.go.handler.TimeHandler.TimestampToDate(ts, this.timezone);
    },
    async convertToTimestamp() {
      if (!this.dateInput) {
        this.timestampResult = "";
        return;
      }
      // 归一化：/ 转 -，不补零的时间补零，只填日期时补 00:00:00
      const normalized = this.normalizeDate(this.dateInput.trim());
      if (!normalized) {
        this.timestampResult = "格式错误，应为 YYYY-MM-DD HH:mm:ss";
        return;
      }
      this.timestampResult = await window.go.handler.TimeHandler.DateToTimestamp(
          normalized,
          this.timezone
      );
    },
    // 日期输入归一化：支持 YYYY-MM-DD、YYYY/MM/DD、不补零、带或不带时分秒
    normalizeDate(input) {
      const str = input.replace(/\//g, '-').trim();
      // 无时间部分，补 00:00:00
      const dateOnly = str.match(/^(\d{4})-(\d{1,2})-(\d{1,2})$/);
      if (dateOnly) {
        return `${dateOnly[1]}-${dateOnly[2].padStart(2, '0')}-${dateOnly[3].padStart(2, '0')} 00:00:00`;
      }
      // 带时间
      const full = str.match(/^(\d{4})-(\d{1,2})-(\d{1,2})[ T](\d{1,2}):(\d{1,2})(?::(\d{1,2}))?$/);
      if (full) {
        return `${full[1]}-${full[2].padStart(2, '0')}-${full[3].padStart(2, '0')} ${full[4].padStart(2, '0')}:${full[5].padStart(2, '0')}:${(full[6] || '0').padStart(2, '0')}`;
      }
      return null;
    },
    async copyResult(text, label) {
      if (!text || text.startsWith('格式错误') || text.startsWith('无效')) {
        return;
      }
      let ok = false;
      try {
        await navigator.clipboard.writeText(text);
        ok = true;
      } catch (e) {
        // 降级
        const textarea = document.createElement('textarea');
        textarea.value = text;
        document.body.appendChild(textarea);
        textarea.select();
        ok = document.execCommand('copy');
        document.body.removeChild(textarea);
      }
      // 文案带对象：这一页有 4 个复制点，泛泛的「已复制」分不清复制了哪个
      this.showToast(ok ? label + '已复制' : '复制失败', ok ? 'success' : 'error');
    },
    showToast(text, type) {
      if (this.toastTimer) {
        clearTimeout(this.toastTimer);
      }
      this.toast = { show: true, text, type: type || 'success' };
      this.toastTimer = setTimeout(() => {
        this.toast.show = false;
      }, 3000);
    },
    clearResultIfEmpty(type) {
      if (type === 'timestamp' && !this.timestampInput.trim()) {
        this.dateResult = '';
      } else if (type === 'date' && !this.dateInput.trim()) {
        this.timestampResult = '';
      } else if (type === 'dateCalc' && !this.baseDateInput.trim() && !this.daysOffset.trim()) {
        this.dateCalcResult = '';
      }
    },
    calculateDate() {
      if (!this.baseDateInput) {
        this.dateCalcResult = '';
        return;
      }
      const days = parseInt(this.daysOffset || 0);
      const baseDate = new Date(this.baseDateInput.replace(/-/g, '/'));
      if (isNaN(baseDate.getTime())) {
        this.dateCalcResult = '日期格式错误';
        return;
      }
      const resultDate = new Date(baseDate);
      resultDate.setDate(resultDate.getDate() + days);

      const year = resultDate.getFullYear();
      const month = String(resultDate.getMonth() + 1).padStart(2, '0');
      const day = String(resultDate.getDate()).padStart(2, '0');
      const hours = String(resultDate.getHours()).padStart(2, '0');
      const minutes = String(resultDate.getMinutes()).padStart(2, '0');
      const seconds = String(resultDate.getSeconds()).padStart(2, '0');

      this.dateCalcResult = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    }
  },
  watch: {
    timezone() {
      this.updateCurrentTime();
    }
  }
}
</script>

<style scoped>
/* 本页只保留时区页特有的样式；页头 / 卡片 / 输入 / 按钮 / 结果条 / 提示
   全部来自 style.css 里的全局类，改主题不用动这里 */

.time-page {
  padding: 22px 24px 28px;
}

/* ---------------- 时区选择（自绘下拉） ----------------
   为什么不用原生 <select>：展开列表由 WebView2 自己绘制，圆角 / 行高 / hover / 选中态
   一个都改不了（Windows 上 <option> 只认 background/color，圆角和内边距直接无效），
   截图里那个灰白块就是这么来的。改成自绘面板后视觉与顶栏搜索下拉保持一套语言。 */
.tz-select-wrap {
  position: relative;
}

.tz-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  /* 固定最小宽度：时区名长短不一时按钮宽度不跳 */
  min-width: 222px;
  height: 36px;
  padding: 0 10px 0 12px;
  font-family: inherit;
  font-size: 13px;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: border-color var(--tr), box-shadow var(--tr);
}

.tz-trigger:hover {
  border-color: var(--border-strong);
}

/* 展开时给焦点环，和输入框获得焦点时的表现统一 */
.tz-trigger:focus-visible,
.tz-select-wrap.is-open .tz-trigger {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--ring);
  outline: none;
}

.tz-globe,
.tz-caret {
  flex-shrink: 0;
  color: var(--text-3);
}

.tz-caret {
  transition: transform var(--tr);
}

.tz-select-wrap.is-open .tz-caret {
  transform: rotate(180deg);
}

.tz-current {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-align: left;
  white-space: nowrap;
  text-overflow: ellipsis;
}

/* UTC 偏移单独做成小胶囊：主次分明，也让「名称 + 偏移」有个视觉分隔 */
.tz-offset {
  flex-shrink: 0;
  padding: 2px 6px;
  border-radius: 5px;
  font-size: 11px;
  line-height: 1.35;
  color: var(--text-2);
  background: var(--surface-2);
  border: 1px solid var(--border);
}

/* ---------------- 下拉面板 ---------------- */
.tz-menu {
  position: absolute;
  top: calc(100% + 6px);
  /* 用 left/right 双约束让面板宽度严格等于触发器：绝对定位下这两条会把
     padding + border 一起算进去（border-box 语义），左右边缘自然对齐。
     不能用 min-width:100% —— 面板没设 box-sizing，那是**内容盒**的下限，
     padding 10 + border 2 会额外撑出 12px，再配合 right:0 就整体向左歪 12px。 */
  left: 0;
  right: 0;
  z-index: 60;
  max-height: 320px;
  overflow-y: auto;
  margin: 0;
  padding: 5px;
  list-style: none;
  background: var(--surface);
  /* 这里用 --border-strong 而不是 --border：深色主题下 --shadow-lg 是 none，
     面板与卡片又同为 --surface，只剩这条边框能把浮层从背景里勾出来 */
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-lg);
}

/* 三列定宽：名称吸左、偏移靠右、勾选占固定列。
   第三列即使没有勾也保留宽度，切换时行内元素不会左右跳 */
.tz-item {
  display: grid;
  grid-template-columns: 1fr auto 16px;
  align-items: center;
  gap: 8px;
  height: 34px;
  padding: 0 9px;
  border-radius: var(--radius-xs);
  font-size: 13px;
  color: var(--text-2);
  cursor: pointer;
  transition: background var(--tr), color var(--tr);
}

/* 键盘 / 鼠标经过的高亮 */
.tz-item.act {
  background: var(--surface-2);
  color: var(--text);
}

/* 当前生效的时区。写在 .act 之后，两者同时命中时以它为准 */
.tz-item.on {
  background: var(--accent-soft);
  color: var(--accent-strong);
  font-weight: 600;
}

.tz-name {
  /* min-width:0 是给 grid 列的保险：面板现在是定宽，时区名再长也只会省略号截断，
     不会把 1fr 列撑开导致行内容溢出面板 */
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.tz-off {
  font-size: 11px;
  color: var(--text-3);
  transition: color var(--tr);
}

.tz-item.act .tz-off {
  color: var(--text-2);
}

.tz-item.on .tz-off {
  color: var(--accent-strong);
}

.tz-check {
  color: var(--accent-strong);
}

/* 展开 / 收起：轻微下移 + 淡入，别做位移过大的动画（列表类控件容易显得晃） */
.tz-pop-enter-active,
.tz-pop-leave-active {
  transition: opacity .14s ease, transform .14s ease;
}

.tz-pop-enter-from,
.tz-pop-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* ---------------- 当前时间（一行式） ----------------
   原来是「上下两行的标签 + 21px 数字」分两栏，占 79px 且是整页唯一的大卡片，
   抢了下面操作区的视觉重心。改成一行：标签退成 12px 灰字，值保持 18px 等宽数字。 */
.hero {
  /* 两栏：第一栏贴卡片左边（与下方卡片标题同一竖线），第二栏在剩余空间里居中。
     第二栏曾整组右对齐（贴着暂停按钮），结果是「时间戳被刻意推到行尾」，
     中间留一道 300px+ 的洞 —— 现在改成居中，让它落在分隔线与按钮的正中间，
     两侧各留一半，既不像右对齐那样贴死行尾，也不像左对齐那样把洞甩到右边。
     1.29fr 不是随手写的：解「两栏行尾富余相等」（栏1 内容 286、栏2 内容 164 + 24 分隔线内距），
     得 C1-C2 = 98 → 比例约 1.29:1；字体或内容一变就失效，但这行内容宽度是固定的 */
  display: grid;
  grid-template-columns: minmax(0, 1.29fr) minmax(0, 1fr) 28px;
  column-gap: 28px;
  align-items: center;
  padding: 13px 18px;
  margin-bottom: 16px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
}

.h-cell {
  display: flex;
  align-items: baseline;
  gap: 7px;
  min-width: 0;
}

/* 分隔线：画在第二栏的左边缘。
   第二栏内容在该栏内**居中**（justify-content:center）——不是贴分隔线，也不是贴按钮，
   而是落在两者之间，两侧留白均分；分隔线本身由第一栏的宽度决定位置，不受影响 */
.h-cell + .h-cell {
  justify-content: center;
  padding-left: 24px;
  border-left: 1px solid var(--border-strong);
}

/* 星期：紧跟在时间值后面，用比标签淡一档的字重和色阶，避免被当成第二个标签 */
.h-dow {
  flex-shrink: 0;
  padding: 2px 7px;
  border-radius: 999px;
  font-size: 11.5px;
  line-height: 1.4;
  color: var(--text-2);
  background: var(--surface-2);
  border: 1px solid var(--border);
}

/* 窄窗口（≤1020px）：比例分栏会把两侧留白压到 0，内容贴住分隔线。
   这时改成「第一栏自然宽 + 第二栏吃掉剩余」，被压缩的空间全部从行尾富余里出，
   内容不会被挤压；第二栏仍居中，于是分隔线两侧的留白自动趋于相等 */
@media (max-width: 1020px) {
  .hero {
    grid-template-columns: max-content minmax(0, 1fr) 28px;
  }
}

/* 极窄（≤880px）：两栏真的放不下了，改成两行堆叠——
   第一行「当前时间」独占整宽，第二行「时间戳 + 按钮」。比让内容互相贴住好读 */
@media (max-width: 880px) {
  .hero {
    grid-template-columns: minmax(0, 1fr) 28px;
    row-gap: 8px;
  }

  .h-cell:first-child {
    grid-column: 1 / -1;
  }

  .h-cell + .h-cell {
    grid-column: 1;
    justify-content: flex-start;
    padding-left: 0;
    border-left: 0;
  }

  .hero > .icon-btn {
    grid-column: 2;
    grid-row: 2;
  }
}

.h-key {
  font-size: 12px;
  color: var(--text-3);
  white-space: nowrap;
}

.h-val {
  font-size: 18px;
  font-weight: 620;
  letter-spacing: -.01em;
  white-space: nowrap;
  color: var(--text);
}

/* 页头右侧：时区标识 + 下拉，两件东西说的是同一件事，所以贴成一组 */
.hz-right {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

/* 时区标识（Asia/Shanghai）：从 hero 上移到页头，点击可复制——
   这种 IANA 字符串经常要贴进代码里，比显示着更实用 */
.tz-badge {
  padding: 5px 9px;
  border-radius: 999px;
  font-family: var(--font-mono);
  font-size: 11.5px;
  line-height: 1.35;
  color: var(--text-2);
  background: var(--surface-2);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: color var(--tr), border-color var(--tr), background var(--tr);
}

.tz-badge:hover {
  color: var(--accent-strong);
  border-color: var(--accent);
  background: var(--accent-soft);
}

/* 暂停 / 继续：该行唯一的交互，放在行尾。
   grid 布局下不再需要 margin-left:auto —— 它自己占据最后一列，天然在最右 */
.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: 1px solid var(--border);
  border-radius: var(--radius-xs);
  background: var(--surface);
  color: var(--text-2);
  cursor: pointer;
  transition: color var(--tr), border-color var(--tr), background var(--tr);
}

.icon-btn:hover {
  color: var(--accent-strong);
  border-color: var(--accent);
  background: var(--accent-soft);
}

/* ---------------- 日期计算 ---------------- */
/* 行内三个控件必须永远排在一行：卡片是普通格子（半宽），窗口一窄就会被挤到折行，
   折行会让这张卡比上面两张高出一大截。所以日期栏做成弹性、其余两项固定，绝不 wrap */
.calc-row {
  align-items: flex-end;
  gap: 12px;
  flex-wrap: nowrap;
}

/* 「基准日期」吃掉剩余宽度。nowrap 下这行的硬需求 = 天数框 + 按钮 + 24(gap) + min-width，
   而卡片最窄内宽 = 320(栅格列宽下限) - 32(padding) - 2(border) = 286。
   按钮加了 14px 图标后从 55 涨到 80，固定项合计 80(天数)+80(按钮)+24 = 184，
   留给日期框的只剩 102 —— min-width 一旦超过它，「2 列 + 卡片 320px」时按钮就会被挤出卡片。
   96 留 6px 余量；正常窗口下这一栏由 flex 撑到 120~212px，兜底值根本碰不到 */
.calc-row > .fld:first-child {
  flex: 1 1 auto;
  min-width: 96px;
}

.calc-row > .btn {
  flex: none;
}

/* 固定宽输入框必须显式 border-box。项目没有全局 `* { box-sizing:border-box }`，
   声明的是**内容盒**：`.inp` 带 `padding:0 12px` + `border:1px`，
   于是 width:190 实际渲染 216、90 渲染 116（各胖 26px）——比卡片内宽还宽时就折行 */
.fld-date {
  box-sizing: border-box;
  width: 100%;
}

.fld-days {
  box-sizing: border-box;
  width: 80px;
  text-align: center;
}
</style>
