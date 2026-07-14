import React from 'react';

const ProfessionalBio = () => {
  return (
    <div className="bg-white text-gray-900 py-16 px-6 md:px-20 lg:px-48 font-sans">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Noto+Serif+SC:wght@300;400;700&family=Inter:wght@300;400;600&display=swap');
        .serif-font { font-family: 'Noto Serif SC', serif; }
        .inter-font { font-family: 'Inter', sans-serif; }
      `}</style>

      {/* 头部 */}
      <header className="mb-16 border-b border-gray-200 pb-10">
        <h1 className="text-5xl font-bold mb-4 serif-font">吴迪 (Di Wu)</h1>
        <p className="text-xl text-gray-600 italic inter-font">北京大学 · 信息与计算科学（双学位心理与认知科学）</p>
        <div className="mt-6 flex flex-wrap gap-6 text-sm text-gray-500 inter-font">
          <span>📧 suiseiseki@stu.pku.edu.cn</span>
          <span>📞 177-6432-9625</span>
          <a href="https://www.aeoluswu.info/resume" className="hover:text-blue-600 underline">🌐 Personal Site</a>
        </div>
      </header>

      <main className="space-y-12 inter-font">
        {/* 教育背景 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">教育背景</h2>
          <div className="space-y-6 text-gray-800 leading-relaxed">
            <div>
              <div className="flex justify-between items-baseline">
                <h3 className="font-bold text-lg">北京大学</h3>
                <span className="text-sm text-gray-500">2023.09 ~ 至今</span>
              </div>
              <p className="text-gray-600 mt-1">信息与计算科学（双学位心理与认知科学）（本科）</p>
              <ul className="list-disc ml-5 mt-2 space-y-1 text-gray-700">
                <li>2023 至今，担任班级团支部书记。</li>
                <li>2024～2025 春季学期，担任北京大学 Python 程序设计与数据科学课程助教。</li>
                <li>2025～2026 春季学期，担任北京大学程序设计实习、算法设计与分析课程助教。</li>
              </ul>
            </div>
            <div>
              <div className="flex justify-between items-baseline">
                <h3 className="font-bold text-lg">安徽师范大学附属中学</h3>
                <span className="text-sm text-gray-500">2020.09 ~ 2023.06</span>
              </div>
              <p className="text-gray-600 mt-1">高中</p>
            </div>
          </div>
        </section>

        {/* 项目经验 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">项目经验</h2>
          <div className="space-y-8 text-gray-800 leading-relaxed">
            <div>
              <div className="flex justify-between items-baseline">
                <h3 className="font-bold text-lg">LIBIHT — Cross-Platform Intel Hardware Trace Library</h3>
                <span className="text-sm text-gray-500">Blackhat 2024</span>
              </div>
              <p className="mt-2 text-gray-700">
                基于 Intel 芯片中的 LBR (Last Branch Record) 与 BTS (Branch Trace Store) 跟踪目标进程跳转行为，帮助用户更好地理解目标进程和更好地调试。主导其中 library 模块的开发，作为 Blackhat 2024 Arsenal 第二作者。
              </p>
            </div>
            <div>
              <div className="flex justify-between items-baseline">
                <h3 className="font-bold text-lg">NUS 2025 暑期研究 — 访问学者</h3>
                <span className="text-sm text-gray-500">2025.06 ~ 2025.09</span>
              </div>
              <p className="mt-2 text-gray-700">
                研究借助模型的代码生成框架、代码生成难度评估等问题。
              </p>
            </div>
          </div>
        </section>

        {/* 实习经验 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">实习经验</h2>
          <div className="space-y-8 text-gray-800 leading-relaxed">
            <div>
              <div className="flex justify-between items-baseline">
                <h3 className="font-bold text-lg">腾讯 — 基础安全研究员</h3>
                <span className="text-sm text-gray-500">2024.07 ~ 2026.03</span>
              </div>
              <ul className="list-disc ml-5 mt-2 space-y-2">
                <li><strong>NDR 告警自动研判：</strong>设计并实现基于大模型的 NDR（网络检测与响应）告警自动研判系统。通过深度解析全流量上下文，结合专家知识库实现告警误报自动识别，缩短高置信度告警的处置周期。针对确认威胁，自动生成针对性防御策略与加固建议。</li>
                <li><strong>星火计划：</strong>参与腾讯星火计划，负责微型企业网络靶场的搭建与漏洞全链路测试。作为助教，协助设计攻防教学实验，负责靶场环境稳定性维护及学员技术指导。</li>
                <li><strong>护网工作：</strong>协助护网工作，独立拿下讯飞星火 AI 靶标。</li>
                <li><strong>情报平台工程：</strong>开发并部署自动化情报收集平台（Full-Stack）。基于 FastAPI 与 Next.js 构建响应式架构，集成 MongoDB 实现情报存储。引入 AI 自动化 Agent 流水线，实现多源数据自动搜索、智能归纳与趋势分析。</li>
                <li><strong>自动化漏洞复现 Pipeline：</strong>构建基于 Claude Agent SDK 的自动化漏洞复现 Pipeline。通过 MonoAgent 架构实现任务拆解与逻辑解耦，将漏洞验证步骤封装为可复用的"Skills"，构建具备自我推理能力的攻击验证闭环。</li>
              </ul>
            </div>
            <div>
              <div className="flex justify-between items-baseline">
                <h3 className="font-bold text-lg">乾象投资 — Engineer</h3>
                <span className="text-sm text-gray-500">2026.06 ~ 2026.09</span>
              </div>
              <ul className="list-disc ml-5 mt-2 space-y-2">
                <li>设计实现多 agent 协同系统的网络架构与通信。通过网络协议调度大规模分布式 agent 集群，实现 agent 多平台化。</li>
              </ul>
            </div>
          </div>
        </section>

        {/* 荣誉证书 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">荣誉证书</h2>
          <ul className="grid grid-cols-1 md:grid-cols-2 gap-2 text-gray-700 list-none">
            <li>2026 年阿里云 CTF 决赛第 8 名</li>
            <li>2024 年强网杯挑战赛线下决赛二等奖</li>
            <li>2024 年京麟 CTF 决赛第 6 名</li>
            <li>2023 年 ICPC 区域赛合肥站金牌</li>
            <li>2021、2022 年 NOI 银牌</li>
            <li>2020、2021、2022 NOIWC 金牌</li>
            <li>2022 年 APIO 国际金牌</li>
            <li>ICPC Challenge 2021 全球第 6 名</li>
          </ul>
        </section>

        {/* 专利成果 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">专利成果</h2>
          <ul className="list-disc ml-5 space-y-1 text-gray-700">
            <li>2024080026CN 应用大模型创建白名单</li>
            <li>2024080041CN 一种基于大模型和 NDR 流量告警数据的自动研判</li>
            <li>2024080225CN 一种基于大模型和聚类的数据清洗技术</li>
          </ul>
        </section>

        {/* 技能特长 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">技能特长</h2>
          <ul className="list-disc ml-5 space-y-1 text-gray-700">
            <li>熟练使用 C、C++ 和 Python 进行开发，拥有 TypeScript、PHP 等多种语言开发经验。</li>
            <li>善于在工程中利用 AI 技术。</li>
          </ul>
        </section>
      </main>

      <footer className="mt-20 pt-10 border-t border-gray-100 text-center text-gray-400 text-sm">
        © {new Date().getFullYear()} 吴迪 · 持续探索安全技术的边界
      </footer>
    </div>
  );
};

export default ProfessionalBio;
