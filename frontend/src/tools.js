// 工具注册表：菜单结构 + 组件映射的唯一数据源。
// Sidebar 与 MainLayout（顶栏搜索）共用同一份，避免两处各维护一份清单导致不一致。
import TimeTool from './components/TimeTool.vue';
import Md5Tool from './components/Md5Tool.vue';
import JsonFormat from './components/JsonFormat.vue';
import IpQuery from './components/IpQuery.vue';
import AesTool from './components/AesTool.vue';
import RsaTool from './components/RsaTool.vue';
import TextTool from './components/TextTool.vue';
import QrcodeTool from './components/QrcodeTool.vue';

// 图标全部内联 SVG：跟随 currentColor 变色、零依赖、离线可用
const svg = (body) =>
    '<svg viewBox="0 0 24 24" width="17" height="17" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">' + body + '</svg>';

export const ICON = {
  clock: svg('<circle cx="12" cy="12" r="9"/><polyline points="12 7 12 12 15.5 14.2"/>'),
  lock: svg('<rect x="4" y="10.5" width="16" height="10.5" rx="2.2"/><path d="M8 10.5V7a4 4 0 0 1 8 0v3.5"/>'),
  shield: svg('<path d="M12 21.5s7.5-3.6 7.5-9.3V5.6L12 2.8 4.5 5.6v6.6c0 5.7 7.5 9.3 7.5 9.3z"/><path d="M9.2 12l2 2 3.6-3.6"/>'),
  key: svg('<circle cx="8" cy="15" r="4.2"/><path d="M11 11.9 20.5 2.5"/><path d="M17 6l2 2"/><path d="M19.4 3.6l1.9 1.9"/>'),
  braces: svg('<path d="M8.5 3H7.8A2.8 2.8 0 0 0 5 5.8v3.7A2.8 2.8 0 0 1 2.2 12A2.8 2.8 0 0 1 5 14.5v3.7A2.8 2.8 0 0 0 7.8 21h.7"/><path d="M15.5 3h.7A2.8 2.8 0 0 1 19 5.8v3.7A2.8 2.8 0 0 0 21.8 12A2.8 2.8 0 0 0 19 14.5v3.7A2.8 2.8 0 0 1 16.2 21h-.7"/>'),
  type: svg('<path d="M4.5 6.5V4.5h15v2"/><path d="M12 4.5v15"/><path d="M9 19.5h6"/>'),
  globe: svg('<circle cx="12" cy="12" r="9"/><path d="M3 12h18"/><path d="M12 3a14 14 0 0 1 4 9 14 14 0 0 1-4 9 14 14 0 0 1-4-9 14 14 0 0 1 4-9z"/>'),
  qrcode: svg('<rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/><rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/><path d="M21 21v.01"/><path d="M12 7v3a2 2 0 0 1-2 2H7"/><path d="M3 12h.01"/><path d="M12 3h.01"/><path d="M12 16v.01"/><path d="M16 12h1"/><path d="M21 12v.01"/><path d="M12 21v-1"/>')
};

// 侧边栏菜单。新增工具只需往这里加一条，搜索和侧栏会同步出现
export const TOOL_GROUPS = [
  { label: '时间', items: [{ id: 'time', name: '时区转换', icon: ICON.clock }] },
  {
    label: '编码加密',
    items: [
      { id: 'md5', name: 'MD5加密', icon: ICON.lock },
      { id: 'aes', name: 'AES加解密', icon: ICON.shield },
      { id: 'rsa', name: 'RSA密钥生成', icon: ICON.key }
    ]
  },
  {
    label: '数据处理',
    items: [
      { id: 'json', name: 'JSON美化', icon: ICON.braces },
      { id: 'text', name: '文本处理', icon: ICON.type },
      { id: 'qrcode', name: '二维码', icon: ICON.qrcode }
    ]
  },
  { label: '网络', items: [{ id: 'ip', name: 'IP地址查询', icon: ICON.globe }] }
];

export const TOOL_COMPONENTS = {
  time: TimeTool,
  md5: Md5Tool,
  json: JsonFormat,
  ip: IpQuery,
  aes: AesTool,
  rsa: RsaTool,
  text: TextTool,
  qrcode: QrcodeTool
};

// 展平成列表，供顶栏搜索使用（顺带带上所属分组名）
export const flatTools = () =>
    TOOL_GROUPS.flatMap((g) => g.items.map((i) => ({ ...i, group: g.label })));
