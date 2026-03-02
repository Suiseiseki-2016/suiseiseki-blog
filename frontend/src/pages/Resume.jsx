import React from 'react';

const ProfessionalBio = () => {
  return (
    <div className="bg-white text-gray-900 py-16 px-6 md:px-20 lg:px-48 font-sans">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Noto+Serif+SC:wght@300;400;700&family=Inter:wght@300;400;600&display=swap');
        .serif-font { font-family: 'Noto Serif SC', serif; }
        .inter-font { font-family: 'Inter', sans-serif; }
      `}</style>

      {/* 头部：突出专业形象 */}
      <header className="mb-16 border-b border-gray-200 pb-10">
        <h1 className="text-5xl font-bold mb-4 serif-font">吴迪 (Di Wu)</h1>
        <p className="text-xl text-gray-600 italic inter-font">北京大学 · 信息与计算科学</p>
        <div className="mt-6 flex flex-wrap gap-6 text-sm text-gray-500 inter-font">
          <span>📧 suiseiseki@stu.pku.edu.cn</span>
          <span>📞 177-6432-9625</span>
          <a href="https://www.aeoluswu.info/resume" className="hover:text-blue-600 underline">🌐 Personal Site</a>
        </div>
      </header>

      <main className="space-y-12 inter-font">
        {/* 关于我：强调跨学科背景 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">个人简介</h2>
          <p className="text-lg leading-relaxed text-gray-800">
            北京大学信息与计算科学本科生，辅修心理与认知科学。我致力于深度挖掘“人工智能+安全”的交叉领域，专注于构建自主式安全代理（Agentic Security）与高性能网络防御体系。凭借在竞赛与工业界的深厚积累，我擅长将 LLM 能力引入工程实践，实现安全防御架构的智能化升级。
          </p>
        </section>

        {/* 核心经历：强调技术沉淀 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">研究与工作经历</h2>
          <div className="space-y-8 text-gray-800 leading-relaxed">
            <div>
              <h3 className="font-bold text-lg">腾讯安全 | 基础安全研究员</h3>
              <p className="text-sm text-gray-500 mb-2">2024.07 - 2026.03</p>
              <ul className="list-disc ml-5 space-y-2">
                <li><strong>NDR 智能化升级：</strong> 基于大模型设计并实现全流量自动研判系统，通过深度上下文解析与专家知识库协同，大幅提升告警误报识别准确度并缩短处置周期。</li>
                <li><strong>自动化渗透 Pipeline：</strong> 基于 Claude Agent SDK 与 MonoAgent 架构，成功构建具备自我推理能力的漏洞复现流水线，通过“Skills”模块化封装实现了 POC 开发与执行的自动化闭环。</li>
                <li><strong>情报平台工程：</strong> 构建全栈自动化情报平台（FastAPI/Next.js/MongoDB），引入 AI Agent 流水线实现多源数据智能归纳与趋势分析。</li>
              </ul>
            </div>
            <div>
              <h3 className="font-bold text-lg">开源与科研贡献</h3>
              <ul className="list-disc ml-5 space-y-2">
                <li><strong>LIBIHT (Blackhat 2024)：</strong> 作为核心开发者（第二作者）参与该跨平台 Intel 硬件追踪库开发，利用 LBR/BTS 跟踪技术优化目标进程跳转分析与调试效率。</li>
                <li><strong>学术产出：</strong> 第一发明人申报多项大模型安全防御专利，涉及自动化研判、白名单创建及数据清洗技术。</li>
              </ul>
            </div>
          </div>
        </section>

        {/* 竞赛与教学：突出实战力 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">实战竞赛与教学</h2>
          <p className="leading-relaxed text-gray-800 mb-4">
            拥有丰富的 CTF 实战背景，曾获 ICPC 区域赛金牌及多项全国顶级安全竞赛奖项。在腾讯星火计划中，负责微型企业网络靶场的全链路测试，并以助教身份指导学员攻防技术。在护网行动中，展现了独立攻破复杂 AI 靶标的实战响应能力。
          </p>
        </section>

        {/* 荣誉列表 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">核心荣誉</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-2 text-gray-700">
            <li>ICPC 区域赛合肥站金牌</li>
            <li>APIO 国际金牌 / NOIWC 金牌</li>
            <li>2024 强网杯决赛二等奖</li>
            <li>2024 京麟 CTF 决赛第 6 名</li>
          </div>
        </section>
      </main>

      <footer className="mt-20 pt-10 border-t border-gray-100 text-center text-gray-400 text-sm">
        © {new Date().getFullYear()} 吴迪 · 持续探索安全技术的边界
      </footer>
    </div>
  );
};

export default ProfessionalBio;