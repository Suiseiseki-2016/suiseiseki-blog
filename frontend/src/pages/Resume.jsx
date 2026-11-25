import React, { useEffect } from 'react';
import { Link } from 'react-router-dom';

const ProfessionalBio = () => {
  useEffect(() => {
    document.title = 'Resume - Di Wu';
    return () => { document.title = 'Di Wu'; };
  }, []);

  return (
    <div className="bg-white text-gray-900 rounded-lg shadow-sm py-10 px-6 sm:px-10 md:px-12 font-sans">
      <style>{`
        .serif-font { font-family: 'Noto Serif SC', 'Songti SC', 'STSong', Georgia, serif; }
        .inter-font { font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'PingFang SC', sans-serif; }
      `}</style>

      {/* 头部 */}
      <header className="mb-12 border-b border-gray-200 pb-8">
        <h1 className="text-3xl sm:text-5xl font-bold mb-4 serif-font">吴迪 (Di Wu)</h1>
        <p className="text-lg sm:text-xl text-gray-600 italic inter-font">北京大学 · 信息与计算科学（双学位心理与认知科学）</p>
        <div className="mt-6 flex flex-wrap gap-x-6 gap-y-2 text-sm text-gray-600 inter-font">
          <span>📧 <a href="mailto:suiseiseki@stu.pku.edu.cn" className="hover:text-blue-600 underline">suiseiseki@stu.pku.edu.cn</a></span>
          <span>💻 <a href="https://github.com/Suiseiseki-2016" target="_blank" rel="noopener noreferrer" className="hover:text-blue-600 underline">GitHub</a></span>
          <span>📝 <Link to="/" className="hover:text-blue-600 underline">Technical Blog</Link></span>
        </div>
      </header>

      <main className="space-y-12 inter-font">
        {/* 教育背景 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">教育背景</h2>
          <div className="space-y-6 text-gray-800 leading-relaxed">
            <div>
              <div className="flex flex-wrap justify-between items-baseline gap-2">
                <h3 className="font-bold text-lg">北京大学</h3>
                <span className="text-sm text-gray-500">2023.09 ~ 至今</span>
              </div>
              <p className="text-gray-600 mt-1">信息与计算科学（双学位心理与认知科学）（本科）</p>
              <ul className="list-disc ml-5 mt-2 space-y-1 text-gray-700">
                <li>2023 至今，担任班级团支部书记。</li>
                <li>2024～2025 春季学期，担任北京大学 Python 程序设计与数据科学课程助教。</li>
                <li>2025～2026 春季学期，担任北京大学程序设计实习、算法设计与分析课程助教。</li>
                <li>2026～2027 秋季学期，担任北京大学计算概论（B）课程助教。</li>
                <li>2026～2027 秋季学期，担任朋辈一对一咨询师。</li>
              </ul>
            </div>
            <div>
              <div className="flex flex-wrap justify-between items-baseline gap-2">
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
              <div className="flex flex-wrap justify-between items-baseline gap-2">
                <h3 className="font-bold text-lg">LIBIHT — Cross-Platform Intel Hardware Trace Library</h3>
                <span className="text-sm text-gray-500">Blackhat 2024</span>
              </div>
              <p className="mt-2 text-gray-700">
                基于 Intel 芯片中的 LBR (Last Branch Record) 与 BTS (Branch Trace Store) 跟踪目标进程跳转行为，帮助用户更好地理解目标进程和更好地调试。主导其中 library 模块的开发，作为 Blackhat 2024 Arsenal 第二作者。
              </p>
            </div>
            <div>
              <div className="flex flex-wrap justify-between items-baseline gap-2">
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
              <div className="flex flex-wrap justify-between items-baseline gap-2">
                <h3 className="font-bold text-lg">乾象投资（Metabit）— Engineer</h3>
                <span className="text-sm text-gray-500">2026.06 ~ 至今</span>
              </div>
              <ul className="list-disc ml-5 mt-2 space-y-2">
                <li><strong>多 Agent 分布式架构与通信：</strong>设计并实现多 Agent 协同系统的网络架构与通信机制，通过网络协议调度大规模分布式 Agent 集群，实现 Agent 多平台化。</li>
                <li><strong>多机权限与凭证管理：</strong>设计并实现面向多机环境的权限与凭证管理 RPC 接口，保障分布式节点间的统一鉴权与凭证安全流转。</li>
                <li><strong>后台任务生命周期管理：</strong>实现后台任务（Background Task）的统一生命周期管理，覆盖任务调度、状态追踪与资源回收。</li>
                <li><strong>性能优化与前端平台开发：</strong>负责系统性能与缓存优化（包括 Skill 逐层按需暴露机制等），降低上下文与调度开销、提升响应速度；同时参与配套前端平台的开发与功能迭代。</li>
              </ul>
            </div>
            <div>
              <div className="flex flex-wrap justify-between items-baseline gap-2">
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
          </div>
        </section>

        {/* 荣誉证书 */}
        <section>
          <h2 className="text-2xl font-bold mb-4 border-l-4 border-black pl-4 serif-font">荣誉证书</h2>
          <ul className="grid grid-cols-1 md:grid-cols-2 gap-2 text-gray-700 list-none">
            <li>2025–2026 学年北京大学三好学生、北京大学三等奖学金</li>
            <li>2026 年阿里云 CTF 决赛第 8 名</li>
            <li>2024 年强网杯挑战赛线下决赛二等奖</li>
            <li>2024 年京麒 CTF 决赛第 6 名</li>
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
            <li>熟练使用 C、C++、Python 和 TypeScript 进行开发，具备 Go、PHP 等语言及前后端全栈开发经验。</li>
            <li>熟悉分布式多 Agent 系统架构、RPC 通信、任务生命周期管理与缓存/上下文性能优化，具备大模型安全与 AI Agent 工程落地经验。</li>
          </ul>
        </section>
      </main>
    </div>
  );
};

export default ProfessionalBio;
