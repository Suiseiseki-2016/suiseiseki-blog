import { useEffect } from 'react'

function Resume() {
  useEffect(() => {
    document.title = '吴迪 - 个人简历'
    return () => { document.title = 'Blog' }
  }, [])

  return (
    <div className="bg-gray-50 min-h-screen py-10 px-4 print:bg-white print:py-0">
      {/* 打印样式优化 */}
      <style>
        {`
          @media print {
            nav, footer, .print-hidden { display: none !important; }
            body { background: white; }
            .resume-container { box-shadow: none !important; border: none !important; max-width: 100% !important; padding: 0 !important; }
            h2 { border-bottom-width: 2px !important; border-color: #e5e7eb !important; }
          }
        `}
      </style>

      <div className="resume-container max-w-4xl mx-auto bg-white shadow-lg rounded-lg overflow-hidden p-10 print:p-0">
        
        {/* 头部信息 */}
        <header className="border-b-2 border-blue-600 pb-6 mb-8 flex justify-between items-end">
          <div>
            <h1 className="text-5xl font-bold text-gray-900 tracking-tight">吴迪</h1>
            <p className="text-blue-600 font-medium mt-2 text-lg">北京大学 · 信息与计算科学</p>
          </div>
          <div className="text-right text-gray-600">
            <p className="flex items-center justify-end uppercase tracking-wide text-sm font-semibold">
              <span className="mr-2">📧</span> {`suiseiseki@stu.pku.edu.cn`}
            </p>
            <p className="flex items-center justify-end mt-1 uppercase tracking-wide text-sm font-semibold">
              <span className="mr-2">📞</span> 177-6432-9625
            </p>
            <p className="text-sm mt-1 text-gray-400">出生年月：2005-06</p>
          </div>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          
          {/* 左侧边栏 - 技能与荣誉 */}
          <div className="md:col-span-1 space-y-8">
            <section>
              <h2 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
                <span className="w-1 h-6 bg-blue-600 mr-2"></span>
                核心技能
              </h2>
              <div className="space-y-2">
                <p className="text-gray-700"><span className="font-semibold italic">熟练:</span> C, C++, Python</p>
                <p className="text-gray-700"><span className="font-semibold italic">经验:</span> TypeScript, PHP, FastAPI, Next.js, MongoDB</p>
                <p className="text-gray-600 text-sm mt-2">✨ 善于在工程中利用 AI 技术提升开发效率</p>
              </div>
            </section>

            <section>
              <h2 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
                <span className="w-1 h-6 bg-blue-600 mr-2"></span>
                竞赛荣誉
              </h2>
              <ul className="text-sm space-y-2 text-gray-700">
                <li>🏆 <strong>2023 ICPC 区域赛(合肥) 金牌</strong></li>
                <li>🏅 2022 APIO 国际金牌</li>
                <li>🥈 2022 NOI 全国银牌</li>
                <li>🥇 2022 NOIWC 全国金牌</li>
                <li>🔥 ICPC Challenge 2021 全球第6名</li>
                <li>🛡️ 2024 强网杯决赛二等奖</li>
                <li>🛡️ 2024 京麟 CTF 决赛第6名</li>
              </ul>
            </section>

            <section>
              <h2 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
                <span className="w-1 h-6 bg-blue-600 mr-2"></span>
                专利成果
              </h2>
              <ul className="text-xs space-y-2 text-gray-600 italic">
                <li>• 应用大模型创建白名单 (2024080026CN)</li>
                <li>• 基于大模型和NDR告警的自动研判 (2024080041CN)</li>
                <li>• 基于大模型和聚类的数据清洗技术 (2024080225CN)</li>
              </ul>
            </section>
          </div>

          {/* 右侧主栏 - 经历 */}
          <div className="md:col-span-2 space-y-8">
            
            <section>
              <h2 className="text-xl font-bold text-gray-900 mb-4 border-b-2 border-gray-100 pb-1">教育背景</h2>
              <div className="relative pl-4 border-l-2 border-blue-100">
                <div className="mb-4">
                  <div className="flex justify-between font-bold">
                    <span>北京大学</span>
                    <span className="text-gray-500">2023.09 - 至今</span>
                  </div>
                  <p className="text-gray-700 text-sm">信息与计算科学 (双学位心理与认知科学)</p>
                  <ul className="text-xs text-gray-500 mt-1 list-disc list-inside">
                    <li>班级团支部书记</li>
                    <li>Python 程序设计与数据科学课程助教 (2024-2025春)</li>
                  </ul>
                </div>
                <div className="text-sm text-gray-600">
                  <p>安徽师范大学附属中学 (2020 - 2023)</p>
                </div>
              </div>
            </section>

            <section>
              <h2 className="text-xl font-bold text-gray-900 mb-4 border-b-2 border-gray-100 pb-1">实习经验</h2>
              <div className="space-y-6">
                <div>
                  <div className="flex justify-between font-bold">
                    <span className="text-blue-700 text-lg">腾讯 · 基础安全研究员</span>
                    <span className="text-gray-500">2024.07 - 2025.06</span>
                  </div>
                  <ul className="mt-2 space-y-2 text-gray-700 text-sm leading-relaxed">
                    <li>• <strong>NDR 告警 AI 研判：</strong> 在玄武实验室实现流量告警自动分析，利用大模型对请求进行误报判定并提供防御建议。</li>
                    <li>• <strong>腾讯星火计划：</strong> 负责微型公司网络靶场测试与漏洞复现，并担任助教指导学员。</li>
                    <li>• <strong>信息搜集系统：</strong> 基于 <strong>FastAPI + Next.js + MongoDB</strong> 开发 AI 自动化搜集站，实现主题自搜索与自动分析汇总。</li>
                    <li>• <strong>实战成果：</strong> 协助护网，独立拿下讯飞星火 AI 靶标。</li>
                  </ul>
                </div>
              </div>
            </section>

            <section>
              <h2 className="text-xl font-bold text-gray-900 mb-4 border-b-2 border-gray-100 pb-1">项目/研究经验</h2>
              <div className="space-y-4 text-sm text-gray-700">
                <div className="bg-blue-50 p-3 rounded-r-md">
                  <div className="flex justify-between font-bold mb-1">
                    <span>Blackhat 2024 · LIBIHT (第二作者)</span>
                    <span className="text-blue-600">开源项目</span>
                  </div>
                  <p>跨平台硬件跟踪库，基于 Intel LBR/BTS 技术实现。主导 Library 模块开发，旨在帮助用户深度调试与理解进程跳转行为。</p>
                </div>
                <div>
                  <div className="flex justify-between font-bold mb-1">
                    <span>NUS 2025 暑期研究 (访问学者)</span>
                    <span className="text-gray-500">2025.06 - 2025.09</span>
                  </div>
                  <p>研究基于大模型的代码生成框架及代码生成难度评估体系。</p>
                </div>
              </div>
            </section>

          </div>
        </div>

        {/* 打印按钮 */}
        <div className="print-hidden mt-12 text-center border-t pt-6">
          <button
            onClick={() => window.print()}
            className="px-8 py-3 bg-blue-600 text-white font-bold rounded-full shadow-lg hover:bg-blue-700 transform hover:-translate-y-1 transition-all"
          >
            导出为 PDF / 打印简历
          </button>
        </div>
      </div>
    </div>
  )
}

export default Resume