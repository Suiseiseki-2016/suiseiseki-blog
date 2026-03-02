import React from 'react';

const ProfessionalBio = () => {
  return (
    <div className="bg-white text-gray-900 py-16 px-6 md:px-20 lg:px-48 font-sans">
      {/* 嵌入式样式，保持字体与学术排版 */}
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Noto+Serif+SC:wght@300;400;700&family=Inter:wght@300;400;600&display=swap');
        .serif-font { font-family: 'Noto Serif SC', serif; }
        .inter-font { font-family: 'Inter', sans-serif; }
      `}</style>

      {/* 头部 */}
      <header className="mb-16 border-b border-gray-200 pb-10">
        <h1 className="text-5xl font-bold mb-4 serif-font">吴迪 (Di Wu)</h1>
        <p className="text-xl text-gray-600 italic inter-font">北京大学 | 信息与计算科学</p>
        <div className="mt-6 flex flex-wrap gap-4 text-sm text-gray-500 inter-font">
          <span>📧 suiseiseki@stu.pku.edu.cn</span>
          <span>📞 177-6432-9625</span>
          <span>📍 北京，中国</span>
        </div>
      </header>

      {/* 内容主体 */}
      <main className="space-y-12 inter-font">
        
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">关于我</h2>
          <p className="text-lg leading-relaxed text-gray-800">
            我目前就读于北京大学，主修信息与计算科学，辅修心理与认知科学。我致力于探索“人工智能+安全”的交叉领域，专注于构建自主式安全代理（Agentic Security）与智能网络防御体系。
          </p>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">核心研究与工作经历</h2>
          <div className="space-y-6 text-gray-800 leading-relaxed">
            <p>
              <strong>在腾讯安全玄武实验室期间</strong>，我深度参与了 NDR（网络检测与响应）系统的智能化研究。我利用大语言模型（LLM）对大规模流量告警进行自动研判，通过解析专家经验与流量上下文，有效提升了误报识别的准确率。此外，我独自构建了一套基于 Claude Agent SDK 的漏洞复现流水线，引入 MiniAgent 架构与模块化“Skill”设计，实现了安全测试任务的闭环自动化。
            </p>
            <p>
              <strong>在科研与项目层面</strong>，我是 Blackhat 2024 开源项目 LIBIHT 的核心开发者（第二作者），致力于跨平台硬件跟踪技术的落地；我以第一发明人身份申请了多项基于大模型的安全防御专利，涵盖自动化研判、白名单生成及数据清洗等技术方向。
            </p>
          </div>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">实战与教学</h2>
          <p className="leading-relaxed text-gray-800">
            我热衷于在实战中验证技术架构。在护网专项行动中，我独立攻破了讯飞星火 AI 靶标，展现了复杂环境下的攻防能力。同时，我曾作为腾讯星火计划的助教与靶场负责人，参与设计并测试了微型企业网络靶场，在指导青少年探索网络安全的过程中，进一步夯实了自己的理论深度。
          </p>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">部分荣誉</h2>
          <ul className="list-disc list-inside space-y-1 text-gray-700">
            <li>2023 ICPC 区域赛（合肥）金牌</li>
            <li>2022 APIO 国际金牌 / NOIWC 全国金牌</li>
            <li>2024 强网杯决赛二等奖 / 京麟 CTF 决赛第 6 名</li>
          </ul>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">技术栈</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h3 className="font-bold uppercase tracking-wider text-sm text-gray-500">研究领域</h3>
              <p className="text-gray-800">漏洞分析、NDR 流量研判、红队实战、AI 安全</p>
            </div>
            <div>
              <h3 className="font-bold uppercase tracking-wider text-sm text-gray-500">工程能力</h3>
              <p className="text-gray-800">Python (FastAPI), TypeScript (Next.js), MongoDB, Claude Agent SDK</p>
            </div>
          </div>
        </section>

      </main>

      <footer className="mt-20 pt-10 border-t border-gray-100 text-center text-gray-400 text-sm">
        © {new Date().getFullYear()} 吴迪 · Powered by Curiosity
      </footer>
    </div>
  );
};

export default ProfessionalBio;