import { useEffect } from 'react'

function Resume() {
  useEffect(() => {
    document.title = '简历 - Blog'
    return () => { document.title = 'Blog' }
  }, [])

  return (
    <div className="bg-white rounded-lg shadow-sm p-8 print:p-4">
      {/* 打印样式 */}
      <style>
        {`
          @media print {
            body {
              background: white;
            }
            nav, footer {
              display: none;
            }
            .print-hidden {
              display: none;
            }
            .resume-container {
              max-width: 100%;
              padding: 0;
            }
            .page-break {
              page-break-after: always;
            }
          }
        `}
      </style>

      <div className="resume-container max-w-3xl mx-auto">
        <header className="mb-8 text-center print:text-left">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Your Name</h1>
          <div className="text-gray-600 space-x-4">
            <span>Email: your.email@example.com</span>
            <span>|</span>
            <span>Phone: +86 138-0000-0000</span>
          </div>
          <div className="text-gray-600 mt-2">
            <span>GitHub: github.com/yourusername</span>
            <span className="mx-2">|</span>
            <span>Website: yourwebsite.com</span>
          </div>
        </header>

        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 border-b pb-2">
            个人简介
          </h2>
          <p className="text-gray-700 leading-relaxed">
            这里写你的个人简介。可以包括你的专业背景、主要技能、职业目标等。
          </p>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 border-b pb-2">
            工作经历
          </h2>
          <div className="space-y-6">
            <div>
              <div className="flex justify-between items-start mb-2">
                <h3 className="text-xl font-semibold text-gray-900">
                  职位名称
                </h3>
                <span className="text-gray-600">2020.01 - 至今</span>
              </div>
              <div className="text-gray-700 font-medium mb-1">公司名称</div>
              <ul className="text-gray-700 list-disc list-inside space-y-1">
                <li>工作内容描述 1</li>
                <li>工作内容描述 2</li>
                <li>工作内容描述 3</li>
              </ul>
            </div>
          </div>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 border-b pb-2">
            教育背景
          </h2>
          <div className="space-y-4">
            <div>
              <div className="flex justify-between items-start mb-1">
                <h3 className="text-xl font-semibold text-gray-900">
                  学位名称
                </h3>
                <span className="text-gray-600">2016.09 - 2020.06</span>
              </div>
              <div className="text-gray-700">学校名称</div>
            </div>
          </div>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 border-b pb-2">
            技能
          </h2>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <h3 className="font-semibold text-gray-900 mb-2">编程语言</h3>
              <ul className="text-gray-700 list-disc list-inside">
                <li>Go</li>
                <li>JavaScript/TypeScript</li>
                <li>Python</li>
              </ul>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 mb-2">框架/工具</h3>
              <ul className="text-gray-700 list-disc list-inside">
                <li>React</li>
                <li>Gin</li>
                <li>Docker</li>
              </ul>
            </div>
          </div>
        </section>

        <div className="print-hidden mt-8 text-center">
          <button
            onClick={() => window.print()}
            className="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
          >
            打印/导出 PDF
          </button>
        </div>
      </div>
    </div>
  )
}

export default Resume
